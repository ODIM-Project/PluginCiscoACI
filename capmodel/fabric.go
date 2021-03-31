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

	"github.com/ODIM-Project/PluginCiscoACI/capdata"
	"github.com/ODIM-Project/PluginCiscoACI/db"
)

// GetFabric collects the fabric data from the DB
func GetFabric(fabricID string) (capdata.Fabric, error) {
	var fabric capdata.Fabric
	data, err := db.Connector.Get(db.TableFabric, fabricID)
	if err != nil {
		return fabric, err
	}
	if err = json.Unmarshal([]byte(data), &fabric); err != nil {
		return fabric, fmt.Errorf("while trying to unmarshal fabric data, got: %v", err)
	}
	return fabric, nil
}

// GetAllFabric collects the fabric data from the DB
func GetAllFabric(pattern string) (map[string]capdata.Fabric, error) {
	allFabricData := make(map[string]capdata.Fabric)
	fabricIDs, err := db.Connector.GetAllMatchingKeys(db.TableFabric, pattern)
	if err != nil {
		return nil, fmt.Errorf("while trying to collect all fabric data, got: %w", err)
	}
	for _, fabricID := range fabricIDs {
		fabric, err := GetFabric(fabricID)
		if err != nil {
			return nil, fmt.Errorf("while trying collect individual fabric data, got: %w", err)
		}
		allFabricData[fabricID] = fabric
	}
	return allFabricData, nil
}

// SaveFabric stores the fabric data in the DB
func SaveFabric(fabricID string, data *capdata.Fabric) error {
	return SaveToDB(db.TableFabric, fabricID, *data)
}

// UpdateFabric updates the fabric data stored in the DB
func UpdateFabric(fabricID string, data *capdata.Fabric) error {
	return UpdateDbData(db.TableFabric, fabricID, *data)
}
