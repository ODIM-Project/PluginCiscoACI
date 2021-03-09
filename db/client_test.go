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
	"testing"

	"github.com/ODIM-Project/PluginCiscoACI/config"
)

type redisExtCallsImpMock struct{}

func (r redisExtCallsImpMock) getNewClient() *Client {
	return &Client{}
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetClient()
			if got == nil {
				t.Error("client is nil")
			}
		})
	}
}
