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
	"github.com/gomodule/redigo/redis"
	redisSentinel "github.com/go-redis/redis"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"sync"
	"time"
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
		inMemDBConnPool.PoolUpdatedTime = time.Now()
	}
	if inMemDBConnPool.WritePool == nil {
		log.Info("GetDBConnection : connPool.WritePool is nil, invoking resetDBWriteConection ")
		resetDBWriteConection(InMemory)
	}

	return inMemDBConnPool, err
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
	log.Info("GetCurrentMasterHostPort masterIP : "+masterIP, ", masterPort : "+masterPort)
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
