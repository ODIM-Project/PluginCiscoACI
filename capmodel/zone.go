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

// GetZone collects the zone data from the DB
func GetZone(zoneID string) (*capdata.ZoneData, error) {
	var zone capdata.ZoneData
	data, err := db.Connector.Get(TableZone, zoneID)
	if err != nil {
		return nil, fmt.Errorf("while trying to collect zone data, got: %w", err)
	}
	err = json.Unmarshal([]byte(data), &zone)
	if err != nil {
		return nil, fmt.Errorf("while trying to unmarshal zone data, got: %v", err)
	}
	return &zone, nil
}

// GetAllZones collects the zone data from the DB
func GetAllZones(pattern string) ([]capdata.ZoneData, error) {
	var allZones []capdata.ZoneData
	allKeys, err := db.Connector.GetAllMatchingKeys(TableZone, pattern)
	if err != nil {
		return nil, fmt.Errorf("while trying to collect all zone keys, got: %w", err)
	}
	for _, key := range allKeys {
		zone, err := GetZone(key)
		if err != nil {
			return nil, fmt.Errorf("while trying collect individual zone data, got: %w", err)
		}
		allZones = append(allZones, *zone)
	}
	return allZones, nil
}
