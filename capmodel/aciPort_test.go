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

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/PluginCiscoACI/db"
)

func TestGetPort(t *testing.T) {
	db.Connector = db.MockConnector{}
	type args struct {
		portID string
	}
	tests := []struct {
		name    string
		args    args
		want    *dmtf.Port
		wantErr bool
	}{
		{
			name: "successful get on port",
			args: args{
				portID: "validID",
			},
			want: &dmtf.Port{
				ID: "validID",
			},
			wantErr: false,
		},
		{
			name: "failed get on port",
			args: args{
				portID: "invalidID",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPort(tt.args.portID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPort() = %v, want %v", got, tt.want)
			}
		})
	}
}
