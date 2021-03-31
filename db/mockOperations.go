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
	"fmt"
)

// MockConnector is for mocking DB connector interface
type MockConnector struct{}

// Create is for mocking DB Create operation
func (d MockConnector) Create(table, resourceID, data string) error {
	return nil
}

// Update is for mocking DB Update operation
func (d MockConnector) Update(table, resourceID, data string) error {
	return nil
}

// GetAllMatchingKeys is for mocking GetAllMatchingKeys operation
func (d MockConnector) GetAllMatchingKeys(table, pattern string) ([]string, error) {
	return []string{"validID"}, nil
}

// Get is for mocking DB Get operation
func (d MockConnector) Get(table, resourceID string) (string, error) {
	if resourceID == "validID" || resourceID == "validID:zoneID" {
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
			return `{"ID": "zoneID"}`, nil
		default:
		}
	}
	return "", fmt.Errorf("not found")
}

// UpdateKeySet is for mocking DB SADD operation
func (d MockConnector) UpdateKeySet(key string, member string) (err error) {
	return nil
}

// GetKeySetMembers is for mocking DB SMEMBERS operation
func (d MockConnector) GetKeySetMembers(key string) (list []string, err error) {
	return []string{"zoneID"}, nil
}

// Delete is for mocking DB Delete operation
func (d MockConnector) Delete(table, resourceID string) (err error) {
	return nil
}

// DeleteKeySetMembers is for mocking DB SREM operation
func (d MockConnector) DeleteKeySetMembers(key string, member string) (err error) {
	return nil
}
