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
	"strings"
	"testing"

	"github.com/ODIM-Project/PluginCiscoACI/config"
	"github.com/go-redis/redis"
)

type redisExtCallsImpMock struct{}

func (r redisExtCallsImpMock) newSentinelClient(opt *redis.Options) *redis.SentinelClient {
	return newSentinelClientMock(opt)
}

func newSentinelClientMock(opt *redis.Options) *redis.SentinelClient {
	strSlice := strings.Split(opt.Addr, ":")
	sentinelHost := strSlice[0]
	sentinelPort := strSlice[1]
	if sentinelHost == "ValidHost" && sentinelPort == "ValidSentinelPort" {
		return &redis.SentinelClient{}
	}
	return nil
}

func (r redisExtCallsImpMock) getNewClient() *redis.Client {
	return &redis.Client{}
}
func (r redisExtCallsImpMock) getMasterAddrByName(masterSet string, snlClient *redis.SentinelClient) []string {
	return getMasterAddbyNameMock(masterSet, snlClient)
}
func getMasterAddbyNameMock(masterSet string, snlClient *redis.SentinelClient) []string {
	if masterSet == "ValidMasterSet" && snlClient != nil {
		return []string{"ValidMasterIP", "ValidMasterPort"}
	}
	return []string{"", ""}
}

func TestGetClient(t *testing.T) {
	redisExtCalls = redisExtCallsImpMock{}
	config.SetUpMockConfig(t)

	tests := []struct {
		name    string
		want    *Client
		wantErr bool
	}{
		{
			name:    "successfully created client",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("getClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.pool == nil {
				t.Error("pool is nil")
			}
		})
	}
}
