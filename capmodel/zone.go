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

	"github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/PluginCiscoACI/capdata"
	"github.com/ODIM-Project/PluginCiscoACI/db"
)

// GetZone collects the zone data from the DB
func GetZone(fabricID, zoneURI string) (model.Zone, error) {
	var zone model.Zone
	key := fmt.Sprintf("%s:%s", fabricID, zoneURI)
	data, err := db.Connector.Get(db.TableZone, key)
	if err != nil {
		return zone, fmt.Errorf("while trying to collect zone data, got: %w", err)
	}
	if err = json.Unmarshal([]byte(data), &zone); err != nil {
		return zone, fmt.Errorf("while trying to unmarshal zone data, got: %v", err)
	}
	return zone, nil
}

// GetZoneDomain collects the ZoneToDomainDN data from the DB
func GetZoneDomain(zoneURI string) (capdata.ACIDomainData, error) {
	var domainData capdata.ACIDomainData
	data, err := db.Connector.Get(db.TableZoneDomain, zoneURI)
	if err != nil {
		return domainData, fmt.Errorf("while trying to collect zone domain data, got: %w", err)
	}
	if err = json.Unmarshal([]byte(data), &domainData); err != nil {
		return domainData, fmt.Errorf("while trying to unmarshal zone domain data, got: %v", err)
	}
	return domainData, nil
}

// GetAllZones collects the zone data from the DB
func GetAllZones(fabricID string) (map[string]model.Zone, error) {
	allZones := make(map[string]model.Zone)
	keySet := fmt.Sprintf("%s:%s", db.TableZone, fabricID)
	zoneURIs, err := db.Connector.GetKeySetMembers(keySet)
	if err != nil {
		return nil, fmt.Errorf("while trying to collect all zone keys, got: %w", err)
	}
	for _, zoneURI := range zoneURIs {
		zone, err := GetZone(fabricID, zoneURI)
		if err != nil {
			return nil, fmt.Errorf("while trying collect individual zone data, got: %w", err)
		}
		allZones[zoneURI] = zone
	}
	return allZones, nil
}

// SaveZone stores the zone data in the DB
func SaveZone(fabricID, zoneURI string, data *model.Zone) error {
	key := fmt.Sprintf("%s:%s", fabricID, zoneURI)
	if err := SaveToDB(db.TableZone, key, *data); err != nil {
		return fmt.Errorf("while trying to store zone data, got: %v", err)
	}
	keySet := fmt.Sprintf("%s:%s", db.TableZone, fabricID)
	if err := db.Connector.UpdateKeySet(keySet, zoneURI); err != nil {
		return fmt.Errorf("while trying to update zone key set members, got: %v", err)
	}
	return nil
}

// SaveZoneDomain stores the ZoneToDomainDN data in the DB
func SaveZoneDomain(zoneURI string, data *capdata.ACIDomainData) error {
	if err := SaveToDB(db.TableZoneDomain, zoneURI, *data); err != nil {
		return fmt.Errorf("while trying to store zone domain data, got: %v", err)
	}
	return nil
}

// UpdateZone updates the zone data stored in the DB
func UpdateZone(fabricID, zoneURI string, data *model.Zone) error {
	key := fmt.Sprintf("%s:%s", fabricID, zoneURI)
	return UpdateDbData(db.TableZone, key, *data)
}

// DeleteZone deletes the zone data stored in the DB
func DeleteZone(fabricID, zoneURI string) error {
	key := fmt.Sprintf("%s:%s", fabricID, zoneURI)
	if err := db.Connector.Delete(db.TableZone, key); err != nil {
		return fmt.Errorf("while trying to remove zone data, got: %v", err)
	}
	keySet := fmt.Sprintf("%s:%s", db.TableZone, fabricID)
	if err := db.Connector.DeleteKeySetMembers(keySet, zoneURI); err != nil {
		return fmt.Errorf("while trying to remove member from zone key set, got: %v", err)
	}
	return nil
}

// DeleteZoneDomain deletes the ZoneToDomainDN data stored in the DB
func DeleteZoneDomain(zoneURI string) error {
	if err := db.Connector.Delete(db.TableZoneDomain, zoneURI); err != nil {
		return fmt.Errorf("while trying to remove zone domain data, got: %v", err)
	}
	return nil
}
