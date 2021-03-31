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

// GetEndpoints collects the endpoint data belonging to a fabric from the DB
func GetEndpoints(fabricID, oid string) (capdata.EndpointData, error) {
	var endpoint capdata.EndpointData
	key := fmt.Sprintf("%s:%s", fabricID, oid)
	data, err := db.Connector.Get(db.TableEndPoint, key)
	if err != nil {
		return endpoint, err
	}
	if err = json.Unmarshal([]byte(data), &endpoint); err != nil {
		return endpoint, fmt.Errorf("while trying to unmarshal endpoint data, got: %v", err)
	}
	return endpoint, nil
}

// GetAllEndpoints collects all the endpoint data belonging to a fabric from the DB
func GetAllEndpoints(fabricID string) (map[string]capdata.EndpointData, error) {
	allEndpointData := make(map[string]capdata.EndpointData)
	keySet := fmt.Sprintf("%s:%s", db.TableEndPoint, fabricID)
	endpointOids, err := db.Connector.GetKeySetMembers(keySet)
	if err != nil {
		return nil, fmt.Errorf("while trying to collect all endpoint data, got: %v", err)
	}
	for _, oid := range endpointOids {
		endpoint, err := GetEndpoints(fabricID, oid)
		if err != nil {
			return nil, fmt.Errorf("while trying collect individual endpoint data, got: %v", err)
		}
		allEndpointData[oid] = endpoint
	}
	return allEndpointData, nil
}

// SaveEndpoint stores the endpoint data in the DB
func SaveEndpoint(fabricID, oid string, data *capdata.EndpointData) error {
	key := fmt.Sprintf("%s:%s", fabricID, oid)
	if err := SaveToDB(db.TableEndPoint, key, *data); err != nil {
		return fmt.Errorf("while trying to store endpoint data, got: %v", err)
	}
	keySet := fmt.Sprintf("%s:%s", db.TableEndPoint, fabricID)
	if err := db.Connector.UpdateKeySet(keySet, oid); err != nil {
		return fmt.Errorf("while trying to update endpoint key set members, got: %v", err)
	}
	return nil
}

// UpdateEndpoint updates the endpoint data stored in the DB
func UpdateEndpoint(fabricID, oid string, data *capdata.EndpointData) error {
	key := fmt.Sprintf("%s:%s", fabricID, oid)
	return UpdateDbData(db.TableEndPoint, key, *data)
}

// DeleteEndpoint deletes the endpoint data stored in the DB
func DeleteEndpoint(fabricID, oid string) error {
	key := fmt.Sprintf("%s:%s", fabricID, oid)
	if err := db.Connector.Delete(db.TableEndPoint, key); err != nil {
		return fmt.Errorf("while trying to remove endpoint data, got: %v", err)
	}
	keySet := fmt.Sprintf("%s:%s", db.TableEndPoint, fabricID)
	if err := db.Connector.DeleteKeySetMembers(keySet, oid); err != nil {
		return fmt.Errorf("while trying to remove member from endpoint key set, got: %v", err)
	}
	return nil
}
