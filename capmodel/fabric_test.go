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

package capmodel

import (
	"reflect"
	"testing"

	"github.com/ODIM-Project/PluginCiscoACI/capdata"
	"github.com/ODIM-Project/PluginCiscoACI/db"
)

func TestGetFabric(t *testing.T) {
	db.Connector = db.MockConnector{}
	type args struct {
		fabricID string
	}
	tests := []struct {
		name    string
		args    args
		want    capdata.Fabric
		wantErr bool
	}{
		{
			name: "successful get on fabric",
			args: args{
				fabricID: "validID",
			},
			want: capdata.Fabric{
				SwitchData: []string{"test"},
				PodID:      "test",
			},
			wantErr: false,
		},
		{
			name: "failed get on fabric",
			args: args{
				fabricID: "invalidID",
			},
			want:    capdata.Fabric{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFabric(tt.args.fabricID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFabric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFabric() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllFabric(t *testing.T) {
	db.Connector = db.MockConnector{}
	tests := []struct {
		name    string
		want    map[string]capdata.Fabric
		wantErr bool
	}{
		{
			name: "successful get on fabric collection",
			want: map[string]capdata.Fabric{
				"validID": capdata.Fabric{
					SwitchData: []string{"test"},
					PodID:      "test",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllFabric("")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllFabric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllFabric() got = %v, want %v", got, tt.want)
			}
		})
	}
}
