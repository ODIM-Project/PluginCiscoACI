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
	"github.com/ODIM-Project/PluginCiscoACI/db"
)

// GetAddressPool collects the AddressPool data belonging to a fabric from the DB
func GetAddressPool(fabricID, oid string) (model.AddressPool, error) {
	var addressPool model.AddressPool
	key := fmt.Sprintf("%s:%s", fabricID, oid)
	data, err := db.Connector.Get(db.TableAddressPool, key)
	if err != nil {
		return addressPool, err
	}
	if err = json.Unmarshal([]byte(data), &addressPool); err != nil {
		return addressPool, fmt.Errorf("while trying to unmarshal addressPool data, got: %v", err)
	}
	return addressPool, nil
}

// GetAllAddressPools collects all the AddressPool data belonging to a fabric from the DB
func GetAllAddressPools(fabricID string) (map[string]model.AddressPool, error) {
	allAddressPoolData := make(map[string]model.AddressPool)
	keySet := fmt.Sprintf("%s:%s", db.TableAddressPool, fabricID)
	addressPoolOids, err := db.Connector.GetKeySetMembers(keySet)
	if err != nil {
		return nil, fmt.Errorf("while trying to collect all addressPool data, got: %v", err)
	}
	for _, oid := range addressPoolOids {
		addressPool, err := GetAddressPool(fabricID, oid)
		if err != nil {
			return nil, fmt.Errorf("while trying collect individual addressPool data, got: %v", err)
		}
		allAddressPoolData[oid] = addressPool
	}
	return allAddressPoolData, nil
}

// SaveAddressPool stores the AddressPool data in the DB
func SaveAddressPool(fabricID, oid string, data *model.AddressPool) error {
	key := fmt.Sprintf("%s:%s", fabricID, oid)
	if err := SaveToDB(db.TableAddressPool, key, *data); err != nil {
		return fmt.Errorf("while trying to store AddressPool data, got: %v", err)
	}
	keySet := fmt.Sprintf("%s:%s", db.TableAddressPool, fabricID)
	if err := db.Connector.UpdateKeySet(keySet, oid); err != nil {
		return fmt.Errorf("while trying to update AddressPool key set members, got: %v", err)
	}
	return nil
}

// UpdateAddressPool updates the AddressPool data stored in the DB
func UpdateAddressPool(fabricID, oid string, data *model.AddressPool) error {
	key := fmt.Sprintf("%s:%s", fabricID, oid)
	return UpdateDbData(db.TableAddressPool, key, *data)
}

// DeleteAddressPool deletes the AddressPool data stored in the DB
func DeleteAddressPool(fabricID, oid string) error {
	key := fmt.Sprintf("%s:%s", fabricID, oid)
	if err := db.Connector.Delete(db.TableAddressPool, key); err != nil {
		return fmt.Errorf("while trying to remove AddressPool data, got: %v", err)
	}
	keySet := fmt.Sprintf("%s:%s", db.TableAddressPool, fabricID)
	if err := db.Connector.DeleteKeySetMembers(keySet, oid); err != nil {
		return fmt.Errorf("while trying to remove member from AddressPool key set, got: %v", err)
	}
	return nil
}
