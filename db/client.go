//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.

package db

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/PluginCiscoACI/config"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

// Client is the established connection
type Client struct {
	readPool        *redis.Client
	writePool       *redis.Client
	masterIP        string
	poolUpdatedTime time.Time
	mux             sync.Mutex
}

var client *Client

// RedisExternalCalls containes the methods to make calls to external client libraries of Redis DB
type RedisExternalCalls interface {
	getNewClient(host, port string) *redis.Client
	newSentinelClient(opt *redis.Options) *redis.SentinelClient
	getMasterAddrByName(snlClient *redis.SentinelClient) []string
}

type redisExtCallsImp struct{}

var redisExtCalls RedisExternalCalls

func init() {
	redisExtCalls = redisExtCallsImp{}
}

func (r redisExtCallsImp) newSentinelClient(opt *redis.Options) *redis.SentinelClient {
	return redis.NewSentinelClient(opt)
}

func (r redisExtCallsImp) getMasterAddrByName(snlClient *redis.SentinelClient) []string {
	ctx := context.Background()
	return snlClient.GetMasterAddrByName(ctx, config.Data.DBConf.MasterSet).Val()
}

// GetClient retrieves client
func GetClient() (*Client, error) {
	if client == nil || client.readPool == nil {
		log.Info("GetClient : client OR client.readPool is nil")
		client = createClient()
	}
	if client.writePool == nil {
		log.Info("GetDBConnection : client.writePool is nil, invoking resetDBWriteConection ")
		err := resetDBWriteConection()
		if err != nil {
			log.Error("GetClient: unable to create new DB connection: " + err.Error())
			return nil, err
		}
	}
	return client, nil
}

func createClient() *Client {
	var masterIP string
	var masterPort string
	connPool := &Client{}
	masterIP = config.Data.DBConf.Host
	masterPort = config.Data.DBConf.Port

	if config.Data.DBConf.RedisHAEnabled {
		masterIP, masterPort = getCurrentMasterHostPort()
	}

	connPool.readPool = redisExtCalls.getNewClient(config.Data.DBConf.Host, config.Data.DBConf.Port)
	connPool.writePool = redisExtCalls.getNewClient(masterIP, masterPort)
	connPool.masterIP = masterIP
	connPool.poolUpdatedTime = time.Now()

	return connPool
}

func getCurrentMasterHostPort() (string, string) {
	sentinelClient := sentinelNewClient()
	stringSlice := redisExtCalls.getMasterAddrByName(sentinelClient)
	var masterIP string
	var masterPort string
	if len(stringSlice) == 2 {
		masterIP = stringSlice[0]
		masterPort = stringSlice[1]
	}
	log.Info("getCurrentMasterHostPort masterIP : "+masterIP, ", masterPort : "+masterPort)
	return masterIP, masterPort
}

// isDbConnectError is for checking if error is dial connection error
func isDbConnectError(err error) (*errors.Error, bool) {
	if strings.HasSuffix(err.Error(), "connect: connection refused") || err.Error() == "EOF" {
		return errors.PackError(errors.DBConnFailed, err), true
	}
	return nil, false
}

// resetDBWriteConection is used to reset the WriteConnection Pool
func resetDBWriteConection() (err error) {
	if config.Data.DBConf.RedisHAEnabled {
		client.mux.Lock()
		defer client.mux.Unlock()
		if client.writePool != nil {
			return
		}
		err = client.setWritePool()
		if err != nil {
			return fmt.Errorf("Reset of inMemory write pool failed: %s", err.Error())
		}
		log.Info("New inMemory connection pool created")
	}
	return
}

func (p *Client) setWritePool() (err error) {
	currentMasterIP, currentMasterPort := retryForMasterIP(p)
	if currentMasterIP == "" {
		return fmt.Errorf("unable to retrieve master ip from sentinel master election")
	}
	log.Info("new write pool master IP found: " + currentMasterIP)
	writePool := redisExtCalls.getNewClient(currentMasterIP, currentMasterPort)
	if writePool == nil {
		return fmt.Errorf("write pool creation failed")
	}

	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&p.writePool)), unsafe.Pointer(writePool))
	p.masterIP = currentMasterIP
	p.poolUpdatedTime = time.Now()
	return
}

func retryForMasterIP(pool *Client) (currentMasterIP, currentMasterPort string) {
	for i := 0; i < 120; i++ {
		currentMasterIP, currentMasterPort = getCurrentMasterHostPort()
		if currentMasterIP != "" && pool.masterIP != currentMasterIP {
			break
		}
		time.Sleep(1 * time.Second)
	}
	return
}

// getNewClient is used is utility function to create new connection pool for DB.
func (r redisExtCallsImp) getNewClient(host, port string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Network:      config.Data.DBConf.Protocol,
		Addr:         fmt.Sprintf("%s:%s", host, port),
		PoolSize:     config.Data.DBConf.PoolSize,
		MinIdleConns: config.Data.DBConf.MinIdleConns,
	})
}

func sentinelNewClient() *redis.SentinelClient {
	return redisExtCalls.newSentinelClient(&redis.Options{
		Addr:     config.Data.DBConf.Host + ":" + config.Data.DBConf.SentinelPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
