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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"strings"
)

// Using below variables as part of errors will enabling errors.Is() function
var (
	// ErrorKeyAlreadyExist is for identifing already exist error
	ErrorKeyAlreadyExist = errors.New("Key already exist in DB")
	// ErrorKeyNotFound is for identifing not found error
	ErrorKeyNotFound = errors.New("Key not Found in DB")
)

const (
	// scanPaginationSize defines the size of DB keys to be scanned on single query
	scanPaginationSize = 100
)

// Create will create a new entry in DB for the value with the given table and resourceID
func (c *Client) Create(table, resourceID string, data interface{}) (err error) {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("while marshalling data, got: %v", err)
	}
	ok, err := c.pool.SetNX(generateKey(table, resourceID), string(dataByte), 0).Result()
	switch {
	case !ok:
		return fmt.Errorf(
			"%w: %s",
			ErrorKeyAlreadyExist,
			fmt.Sprintf("An entry with resource id %s is already present in table %s", resourceID, table),
		)
	case err != nil:
		return fmt.Errorf(
			"Creating new entry for value %v in table %s with resource id %s failed: %v",
			data, table, resourceID, err,
		)
	default:
		return nil
	}
}

// GetAllMatchingKeys will collect all the keys of provided table and pattern
func (c *Client) GetAllMatchingKeys(table, pattern string) ([]string, error) {
	var allKeys []string
	var cursor uint64
	for {
		keys, c, err := c.pool.Scan(cursor, generateKey(table, pattern+"*"), scanPaginationSize).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch all keys of table %s: %s", table, err.Error())
		}
		allKeys = append(allKeys, keys...)
		if cursor == 0 {
			break
		}
		cursor = c
	}
	return trimTableFromKeys(table, allKeys), nil
}

// Get will collect the data associated with the given key from the given table
func (c *Client) Get(table, resourceID string) (val string, err error) {
	val, err = c.pool.Get(generateKey(table, resourceID)).Result()
	switch err {
	case redis.Nil:
		return "", fmt.Errorf(
			"%w: %s",
			ErrorKeyNotFound,
			fmt.Sprintf("Data with resource ID %s not found in table %s", resourceID, table),
		)
	case nil:
		return val, nil
	default:
		return "", fmt.Errorf("unable to complete the operation: %s", err.Error())
	}
}

// generateKey is for concatinating table and resourceID to for a key
func generateKey(table, resourceID string) string {
	return fmt.Sprintf("%s:%s", table, resourceID)
}

// trimTableFromKeys trims <table>: from the slice of keys in the form of <table>:<resourceID>
func trimTableFromKeys(table string, fullKeys []string) []string {
	var keys []string
	for _, fullKey := range fullKeys {
		keys = append(keys, strings.TrimPrefix(fullKey, generateKey(table, "")))
	}
	return keys
}
