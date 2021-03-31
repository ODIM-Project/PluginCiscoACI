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

	"github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/PluginCiscoACI/db"
)

func TestGetZone(t *testing.T) {
	db.Connector = db.MockConnector{}
	type args struct {
		fabricID string
		zoneID   string
	}
	tests := []struct {
		name    string
		args    args
		want    model.Zone
		wantErr bool
	}{
		{
			name: "successful get on zone",
			args: args{
				fabricID: "validID",
				zoneID:   "zoneID",
			},
			want:    model.Zone{ID: "zoneID"},
			wantErr: false,
		},
		{
			name: "failed get on zone",
			args: args{
				fabricID: "invalidID",
				zoneID:   "invalidID",
			},
			want:    model.Zone{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetZone(tt.args.fabricID, tt.args.zoneID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetZone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetZone() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllZones(t *testing.T) {
	db.Connector = db.MockConnector{}
	type args struct {
		fabricID string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]model.Zone
		wantErr bool
	}{
		{
			name: "successful get on zone collection",
			args: args{
				fabricID: "validID",
			},
			want: map[string]model.Zone{
				"zoneID": model.Zone{ID: "zoneID"},
			},
			wantErr: false,
		},
		{
			name: "failed get on zone collection",
			args: args{
				fabricID: "invalidID",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllZones(tt.args.fabricID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllZones() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllZones() got = %v, want %v", got, tt.want)
			}
		})
	}
}
