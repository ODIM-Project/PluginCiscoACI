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
	"github.com/ODIM-Project/PluginCiscoACI/config"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

// Client is the established connection
type Client struct {
	*redis.Client
}

var client *Client

// RedisExternalCalls containes the methods to make calls to external client libraries of Redis DB
type RedisExternalCalls interface {
	getNewClient() *Client
}

type redisExtCallsImp struct{}

var redisExtCalls RedisExternalCalls

func init() {
	redisExtCalls = redisExtCallsImp{}
}

// GetClient retrieves client
func GetClient() *Client {
	if client == nil {
		log.Info("GetClient : DB client is nil, creating new client.")
		client = redisExtCalls.getNewClient()
	}
	return client
}

// getNewClient is used is utility function to create new connection pool for DB
func (r redisExtCallsImp) getNewClient() *Client {
	if config.Data.DBConf.SentinelMasterName == "" {
		return &Client{
			redis.NewClient(&redis.Options{
				Network:      config.Data.DBConf.Protocol,
				Addr:         config.Data.DBConf.Address,
				PoolSize:     config.Data.DBConf.PoolSize,
				MinIdleConns: config.Data.DBConf.MinIdleConns,
			}),
		}
	}

	return &Client{
		redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    config.Data.DBConf.SentinelMasterName,
			SentinelAddrs: []string{config.Data.DBConf.Address},
			PoolSize:      config.Data.DBConf.PoolSize,
			MinIdleConns:  config.Data.DBConf.MinIdleConns,
		}),
	}
}
