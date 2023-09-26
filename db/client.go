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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"sync"
	"time"

	"github.com/ODIM-Project/PluginCiscoACI/config"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

// Client is the established connection
type Client struct {
	pool            *redis.Client
	poolUpdatedTime time.Time
	mux             sync.Mutex
}

// Config is the configuration for db which is set by the wrapper package.
/*
Port is the port number for the database connection
Protocol is the type of protocol with which the connection takes place
Host is hostname/IP on which the database is running
*/
type Config struct {
	Port         string
	Protocol     string
	Host         string
	SentinelHost string
	SentinelPort string
	MasterSet    string
	Password     string
}

var client *Client

// RedisExternalCalls containes the methods to make calls to external client libraries of Redis DB
type RedisExternalCalls interface {
	getNewClient() *redis.Client
	newSentinelClient(opt *redis.Options) *redis.SentinelClient
	getMasterAddrByName(mset string, snlClient *redis.SentinelClient) []string
}

type redisExtCallsImp struct {
}

var redisExtCalls RedisExternalCalls

func init() {
	redisExtCalls = redisExtCallsImp{}
	Connector = connector{}
}

// getClient retrieves client
func getClient() (*Client, error) {
	if client == nil || client.pool == nil {
		log.Info("getClient : DB connection pool is nil, creating a new pool.")
		client = &Client{}
		client.poolUpdatedTime = time.Now()
		if config.Data.DBConf.RedisHAEnabled {
			err := resetDBConection()
			if err != nil {
				log.Error("getClient: unable to create new DB connection: " + err.Error())
				return nil, err
			}
		} else {
			client.pool = redisExtCalls.getNewClient()
		}
		if client.pool == nil {
			log.Error("getClient: unable to create new DB connection pool")
			return nil, fmt.Errorf("getClient: unable to create new DB connection pool")
		}
		log.Info("getClient: new pool DB connection pool created.")
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
	pool := retryForSentinelClient()
	if pool == nil {
		return fmt.Errorf("sentinel DB pool creation failed")
	}
	p.pool = pool
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
	tlsConfig, e := getTLSConfig(config.Data.KeyCertConf.CertificatePath, config.Data.KeyCertConf.PrivateKeyPath, config.Data.KeyCertConf.RootCACertificatePath)
	if e != nil {
		log.Error(e.Error())
		return nil
	}

	if config.Data.DBConf.RedisHAEnabled {

		dbConfig := &Config{
			Port:         config.Data.DBConf.Port,
			Protocol:     config.Data.DBConf.Protocol,
			Host:         config.Data.DBConf.Host,
			SentinelHost: config.Data.DBConf.SentinelHost,
			SentinelPort: config.Data.DBConf.SentinelPort,
			MasterSet:    config.Data.DBConf.MasterSet,
			Password:     string(config.Data.DBConf.RedisOnDiskPassword),
		}
		masterIP, masterPort, err := getCurrentMasterHostPort(dbConfig)

		if err != nil {
			log.Error(err.Error())
			return nil
		}
		return redis.NewClient(&redis.Options{
			Network:      config.Data.DBConf.Protocol,
			Addr:         net.JoinHostPort(masterIP, masterPort),
			PoolSize:     config.Data.DBConf.PoolSize,
			MinIdleConns: config.Data.DBConf.MinIdleConns,
			TLSConfig:    tlsConfig,
			Password:     string(config.Data.DBConf.RedisOnDiskPassword),
		})
	}
	return redis.NewClient(&redis.Options{
		Network:      config.Data.DBConf.Protocol,
		Addr:         net.JoinHostPort(config.Data.DBConf.Host, config.Data.DBConf.Port),
		PoolSize:     config.Data.DBConf.PoolSize,
		MinIdleConns: config.Data.DBConf.MinIdleConns,
		TLSConfig:    tlsConfig,
		Password:     string(config.Data.DBConf.RedisOnDiskPassword),
	})
}

func getTLSConfig(cCert, cKey, caCert string) (*tls.Config, error) {

	tlsConfig := tls.Config{}

	// Load client cert
	cert, e1 := tls.LoadX509KeyPair(cCert, cKey)
	if e1 != nil {
		return &tlsConfig, e1
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	// Load CA cert
	caCertR, e2 := ioutil.ReadFile(caCert)
	if e2 != nil {
		return &tlsConfig, e2
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertR)
	tlsConfig.RootCAs = caCertPool

	tlsConfig.BuildNameToCertificate()
	return &tlsConfig, e2
}

// getCurrentMasterHostPort is to get the current Redis Master IP and Port from Sentinel.
func getCurrentMasterHostPort(dbConfig *Config) (string, string, error) {
	sentinelClient, err := sentinelNewClient(dbConfig)
	if err != nil {
		return "", "", err
	}
	stringSlice := redisExtCalls.getMasterAddrByName(dbConfig.MasterSet, sentinelClient)
	var masterIP string
	var masterPort string
	if len(stringSlice) == 2 {
		masterIP = stringSlice[0]
		masterPort = stringSlice[1]
	}

	return masterIP, masterPort, nil
}

func sentinelNewClient(dbConfig *Config) (*redis.SentinelClient, error) {
	tlsConfig, err := getTLSConfig(config.Data.KeyCertConf.CertificatePath, config.Data.KeyCertConf.PrivateKeyPath, config.Data.KeyCertConf.RootCACertificatePath)
	if err != nil {
		return nil, fmt.Errorf("error while trying to get tls configuration : %s", err.Error())
	}
	rdb := redisExtCalls.newSentinelClient(&redis.Options{
		Addr:      dbConfig.SentinelHost + ":" + dbConfig.SentinelPort,
		DB:        0, // use default DB
		TLSConfig: tlsConfig,
		Password:  dbConfig.Password,
	})
	return rdb, nil
}

func (r redisExtCallsImp) newSentinelClient(opt *redis.Options) *redis.SentinelClient {
	return redis.NewSentinelClient(opt)
}

func (r redisExtCallsImp) getMasterAddrByName(masterSet string, snlClient *redis.SentinelClient) []string {
	return snlClient.GetMasterAddrByName(masterSet).Val()
}
