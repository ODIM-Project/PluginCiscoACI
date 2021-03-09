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
	"reflect"
	"strings"
	"testing"

	"github.com/ODIM-Project/PluginCiscoACI/config"
	"github.com/go-redis/redis/v8"
)

type redisExtCallsImpMock struct{}

func (r redisExtCallsImpMock) newSentinelClient(opt *redis.Options) *redis.SentinelClient {
	strSlice := strings.Split(opt.Addr, ":")
	sentinelHost := strSlice[0]
	sentinelPort := strSlice[1]
	if sentinelHost == "ValidHost" && sentinelPort == "ValidPort" {
		return &redis.SentinelClient{}
	}
	return nil
}

func (r redisExtCallsImpMock) getMasterAddrByName(snlClient *redis.SentinelClient) []string {
	if config.Data.DBConf.MasterSet == "ValidMasterSet" {
		return []string{"ValidMasterIP", "ValidMasterPort"}
	}
	return []string{"", ""}
}

func (r redisExtCallsImpMock) getNewClient(host, port string) *redis.Client {
	return &redis.Client{}
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
			name: "successfully created client",
			want: &Client{
				masterIP: "ValidMasterIP",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.masterIP, tt.want.masterIP) {
				t.Errorf("GetClient() = %v, want %v", got.masterIP, tt.want.masterIP)
			}
			if got.readPool == nil {
				t.Error("readPool is nil")
			}
			if got.writePool == nil {
				t.Error("writePool is nil")
			}
		})
	}
}
