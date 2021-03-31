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

//Package capmodel ...
package capmodel

import (
	"encoding/json"
	"fmt"

	"github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/PluginCiscoACI/db"
)

//SwitchChassis ...
type SwitchChassis struct {
	TotalCount string                `json:"totalCount"`
	IMData     []SwitchChassisIMData `json:"imdata"`
}

// SwitchChassisIMData ...
type SwitchChassisIMData struct {
	SwitchChassisData SwitchChassisData `json:"eqptCh"`
}

// SwitchChassisData ...
type SwitchChassisData struct {
	Attributes map[string]interface{} `json:"attributes"`
}

// GetSwitch collects the switch data from the DB
func GetSwitch(switchID string) (model.Switch, error) {
	var switchData model.Switch
	data, err := db.Connector.Get(db.TableSwitch, switchID)
	if err != nil {
		return switchData, err
	}
	if err = json.Unmarshal([]byte(data), &switchData); err != nil {
		return switchData, fmt.Errorf("while trying to unmarshal switch data, got: %v", err)
	}
	return switchData, nil
}

// SaveSwitch stores the switch data in the DB
func SaveSwitch(switchID string, data *model.Switch) error {
	return SaveToDB(db.TableSwitch, switchID, *data)
}

// GetSwitchChassis collects the switch chassis data from the DB
func GetSwitchChassis(chassisID string) (model.Chassis, error) {
	var chassis model.Chassis
	data, err := db.Connector.Get(db.TableSwitchChassis, chassisID)
	if err != nil {
		return chassis, err
	}
	if err = json.Unmarshal([]byte(data), &chassis); err != nil {
		return chassis, fmt.Errorf("while trying to unmarshal chassis data, got: %v", err)
	}
	return chassis, nil
}

// GetAllSwitchChassis collects all the switch chassis data from the DB
func GetAllSwitchChassis(pattern string) (map[string]model.Chassis, error) {
	allChassisData := make(map[string]model.Chassis)
	chassisIDs, err := db.Connector.GetAllMatchingKeys(db.TableSwitchChassis, pattern)
	if err != nil {
		return nil, fmt.Errorf("while trying to collect all switch chassis data, got: %w", err)
	}
	for _, chassisID := range chassisIDs {
		chassis, err := GetSwitchChassis(chassisID)
		if err != nil {
			return nil, fmt.Errorf("while trying collect individual switch chassis data, got: %w", err)
		}
		allChassisData[chassisID] = chassis
	}
	return allChassisData, nil
}

// SaveSwitchChassis stores the switch chassis data in the DB
func SaveSwitchChassis(chassisID string, data *model.Chassis) error {
	return SaveToDB(db.TableSwitchChassis, chassisID, *data)
}
