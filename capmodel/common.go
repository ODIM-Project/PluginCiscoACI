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
	"encoding/json"
	"fmt"

	"github.com/ODIM-Project/PluginCiscoACI/db"
)

const (
	// TableFabric is the table for storing switch and pod ids
	TableFabric = "ACI-Fabric"
	// TableSwitch is the table for storing switch information
	TableSwitch = "ACI-Switch"
	// TableSwitchChassis is the table for storing switch chassis information
	TableSwitchChassis = "ACI-SwitchChassis"
	// TableSwitchPorts is the table for storing ports of each switch
	TableSwitchPorts = "ACI-SwitchPorts"
	// TablePort is the table for storing port information
	TablePort = "ACI-Port"
	// TableZone is the table for storing zone information
	TableZone = "ACI-Zone"
	// TableAddressPool is the table for storing addresspool information
	TableAddressPool = "ACI-AddressPool"
	// TableEndPoint is the table for storing fabric endpoint information
	TableEndPoint = "ACI-EndPoint"
	// TableZoneDomain is the table for storing ZoneToDomainDN information
	TableZoneDomain = "ACI-ZoneDomain"
)

// SaveToDB is for adding port data to the DB
func SaveToDB(table, resourceID string, data interface{}) error {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("while marshalling data, got: %v", err)
	}
	return db.Connector.Create(table, resourceID, string(dataByte))
}

// UpdateDbData is for updating data in the DB
func UpdateDbData(table, resourceID string, data interface{}) error {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("while marshalling data, got: %v", err)
	}
	return db.Connector.Update(table, resourceID, string(dataByte))
}
