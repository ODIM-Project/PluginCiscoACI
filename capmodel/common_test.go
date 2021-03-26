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
	"fmt"
	"testing"

	"github.com/ODIM-Project/PluginCiscoACI/db"
)

type MockConnector struct{}

func TestSaveToDB(t *testing.T) {
	db.Connector = MockConnector{}
	type args struct {
		table      string
		resourceID string
		data       interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "saving data successfully",
			args: args{
				table:      "someTable",
				resourceID: "someResource",
				data:       args{data: "someData"},
			},
			wantErr: false,
		},
		{
			name: "invalid data",
			args: args{
				table:      "someTable",
				resourceID: "resourceAlreadyPresent",
				data:       func() {},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveToDB(tt.args.table, tt.args.resourceID, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("SaveToDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (d MockConnector) Create(table, resourceID, data string) error {
	return nil
}

func (d MockConnector) Update(table, resourceID, data string) error {
	return nil
}

func (d MockConnector) GetAllMatchingKeys(table, pattern string) ([]string, error) {
	return []string{"validID"}, nil
}

func (d MockConnector) Get(table, resourceID string) (string, error) {
	if resourceID == "validID" {
		switch table {
		case TableFabric:
			return `{"SwitchData": ["test"], "PodID": "test"}`, nil
		case TableSwitch:
			return `{"Id": "validID", "FabricID": "validID"}`, nil
		case TableSwitchPorts:
			return `{"Id": "validID", "FabricID": "validID"}`, nil
		case TablePort:
			return `{"Id": "validID", "FabricID": "validID"}`, nil
		case TableZone:
			return `{"Id": "validID", "FabricID": "validID"}`, nil
		default:
		}
	}
	return "", fmt.Errorf("not found")
}

func (d MockConnector) UpdateKeySet(key string, members ...interface{}) (err error) {
	return nil
}

func (d MockConnector) GetKeySetMembers(key string) (list []string, err error) {
	return []string{"validID"}, nil
}

func (d MockConnector) Delete(table, resourceID string) (err error) {
	return nil
}

func (d MockConnector) DeleteKeySetMembers(key string, members ...interface{}) (err error) {
	return nil
}
