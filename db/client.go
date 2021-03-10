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
	"sync"
	"time"

	"github.com/ODIM-Project/PluginCiscoACI/config"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

// Client is the established connection
type Client struct {
	pool            *redis.Client
	poolUpdatedTime time.Time
	mux             sync.Mutex
}

var client *Client

// RedisExternalCalls containes the methods to make calls to external client libraries of Redis DB
type RedisExternalCalls interface {
	getNewClient() *redis.Client
}

type redisExtCallsImp struct{}

var redisExtCalls RedisExternalCalls

func init() {
	redisExtCalls = redisExtCallsImp{}
}

// GetClient retrieves client
func GetClient() (*Client, error) {
	if client == nil || client.pool == nil {
		log.Info("GetClient : DB connection pool is nil, creating a new pool.")
		client = &Client{}
		client.poolUpdatedTime = time.Now()
		if config.Data.DBConf.RedisHAEnabled {
			err := resetDBConection()
			if err != nil {
				log.Error("GetClient: unable to create new DB connection: " + err.Error())
				return nil, err
			}
		} else {
			client.pool = redisExtCalls.getNewClient()
		}
		log.Info("GetClient: new pool DB connection pool created.")
	}
	return client, nil
}

// resetDBConection is used to reset the WriteConnection Pool
func resetDBConection() (err error) {
	client.mux.Lock()
	defer client.mux.Unlock()
	if client.pool != nil {
		return
	}
	err = client.setPool()
	if err != nil {
		return fmt.Errorf("Reset of DB pool failed: %s", err.Error())
	}
	return
}

func (p *Client) setPool() (err error) {
	pool := redisExtCalls.getNewClient()
	if pool == nil {
		return fmt.Errorf("sentinel DB pool creation failed")
	}
	p.pool = pool
	p.poolUpdatedTime = time.Now()
	return
}

func retryForSentinelClient() *redis.Client {
	for i := 0; i < 120; i++ {
		pool := redisExtCalls.getNewClient()
		if pool != nil {
			return pool
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

// getNewClient is used is utility function to create new connection pool for DB.
func (r redisExtCallsImp) getNewClient() *redis.Client {
	if config.Data.DBConf.RedisHAEnabled {
		return redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    config.Data.DBConf.MasterSet,
			SentinelAddrs: []string{fmt.Sprintf("%s:%s", config.Data.DBConf.Host, config.Data.DBConf.SentinelPort)},
			PoolSize:      config.Data.DBConf.PoolSize,
			MinIdleConns:  config.Data.DBConf.MinIdleConns,
		})
	}
	return redis.NewClient(&redis.Options{
		Network:      config.Data.DBConf.Protocol,
		Addr:         fmt.Sprintf("%s:%s", config.Data.DBConf.Host, config.Data.DBConf.Port),
		PoolSize:     config.Data.DBConf.PoolSize,
		MinIdleConns: config.Data.DBConf.MinIdleConns,
	})
}
