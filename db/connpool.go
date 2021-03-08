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
	"fmt"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/PluginCiscoACI/config"
	redisSentinel "github.com/go-redis/redis"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// ConnPool is the established connection
type ConnPool struct {
	ReadPool        *redis.Pool
	WritePool       *redis.Pool
	MasterIP        string
	PoolUpdatedTime time.Time
	Mux             sync.Mutex
}

var connPool *ConnPool

// RedisExternalCalls containes the methods to make calls to external client libraries of Redis DB
type RedisExternalCalls interface {
	newSentinelClient(opt *redisSentinel.Options) *redisSentinel.SentinelClient
	getMasterAddrByName(snlClient *redisSentinel.SentinelClient) []string
}

type redisExtCallsImp struct{}

var redisExtCalls RedisExternalCalls

func init() {
	redisExtCalls = redisExtCallsImp{}
}

func (r redisExtCallsImp) newSentinelClient(opt *redisSentinel.Options) *redisSentinel.SentinelClient {
	return redisSentinel.NewSentinelClient(opt)
}

func (r redisExtCallsImp) getMasterAddrByName(snlClient *redisSentinel.SentinelClient) []string {
	return snlClient.GetMasterAddrByName(config.Data.DBConf.MasterSet).Val()
}

// GetConnPool is used to get the new Connection Pool for DB
func GetConnPool() (*ConnPool, *errors.Error) {
	var err *errors.Error
	if connPool == nil || connPool.ReadPool == nil {
		log.Info("GetConnPool : connPool OR connPool.ReadPool is nil")
		connPool, err = createConnPool()
		if err != nil {
			log.Error("error while trying to get Readpool connection : " + err.Error())
		}
		connPool.PoolUpdatedTime = time.Now()
	}
	if connPool.WritePool == nil {
		log.Info("GetDBConnection : connPool.WritePool is nil, invoking resetDBWriteConection ")
		resetDBWriteConection()
	}

	return connPool, err
}

func createConnPool() (*ConnPool, *errors.Error) {
	var err error
	var masterIP string
	var masterPort string
	connPools := &ConnPool{}
	masterIP = config.Data.DBConf.Host
	masterPort = config.Data.DBConf.Port
	if config.Data.DBConf.RedisHAEnabled {
		masterIP, masterPort = getCurrentMasterHostPort()
	}

	connPools.ReadPool, err = getPool(config.Data.DBConf.Host, config.Data.DBConf.Port)
	//Check if any connection error occured
	if err != nil {
		if errs, aye := isDbConnectError(err); aye {
			log.Error("error while trying to get Readpool connection : " + errs.Error())
			return nil, errs
		}
		return nil, errors.PackError(errors.UndefinedErrorType, err)
	}
	connPools.WritePool, err = getPool(masterIP, masterPort)
	//Check if any connection error occured
	if err != nil {
		if errs, aye := isDbConnectError(err); aye {
			log.Error("error while trying to get Writepool connection : " + errs.Error())
			return nil, errs
		}
		return nil, errors.PackError(errors.UndefinedErrorType, err)
	}
	connPools.MasterIP = masterIP

	return connPools, nil
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

func sentinelNewClient() *redisSentinel.SentinelClient {
	rdb := redisExtCalls.newSentinelClient(&redisSentinel.Options{
		Addr:     config.Data.DBConf.Host + ":" + config.Data.DBConf.SentinelPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb
}

// isDbConnectError is for checking if error is dial connection error
func isDbConnectError(err error) (*errors.Error, bool) {
	if strings.HasSuffix(err.Error(), "connect: connection refused") || err.Error() == "EOF" {
		return errors.PackError(errors.DBConnFailed, err), true
	}
	return nil, false
}

//resetDBWriteConection is used to reset the WriteConnection Pool
func resetDBWriteConection() {
	if config.Data.DBConf.RedisHAEnabled {
		connPool.Mux.Lock()
		defer connPool.Mux.Unlock()
		if connPool.WritePool != nil {
			return
		}
		err := connPool.setWritePool()
		if err != nil {
			log.Error("Reset of inMemory write pool failed: " + err.Error())
			return
		}
		log.Info("New inMemory connection pool created")
	}
}

func (p *ConnPool) setWritePool() error {
	currentMasterIP, currentMasterPort := retryForMasterIP(p)
	if currentMasterIP == "" {
		return fmt.Errorf("unable to retrieve master ip from sentinel master election")
	}
	log.Info("new write pool master IP found: " + currentMasterIP)
	writePool, _ := getPool(currentMasterIP, currentMasterPort)
	if writePool == nil {
		return fmt.Errorf("write pool creation failed")
	}

	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool)), unsafe.Pointer(writePool))
	p.MasterIP = currentMasterIP
	p.PoolUpdatedTime = time.Now()
	return nil
}

func retryForMasterIP(pool *ConnPool) (currentMasterIP, currentMasterPort string) {
	for i := 0; i < 120; i++ {
		currentMasterIP, currentMasterPort = getCurrentMasterHostPort()
		if currentMasterIP != "" && pool.MasterIP != currentMasterIP {
			break
		}
		time.Sleep(1 * time.Second)
	}
	return
}

//getPool is used is utility function to get the Connection Pool from DB.
func getPool(host, port string) (*redis.Pool, error) {
	protocol := config.Data.DBConf.Protocol
	p := &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: config.Data.DBConf.MaxIdleConns,
		// max number of connections
		MaxActive: config.Data.DBConf.MaxActiveConns,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(protocol, host+":"+port)
			return c, err
		},
		/*TestOnBorrow is an optional application supplied function to
		  check the health of an idle connection before the connection is
		  used again by the application. Argument t is the time that the
		  connection was returned to the pool.This function PINGs
		  connections that have been idle more than a minute.
		  If the function returns an error, then the connection is closed.
		*/
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	return p, nil
}
