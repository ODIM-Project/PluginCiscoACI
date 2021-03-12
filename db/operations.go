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
)

// Using below variables as part of errors will enabling errors.Is() function
var (
	// ErrorKeyAlreadyExist is for identifing already exist error
	ErrorKeyAlreadyExist = errors.New("Key already exist in DB")
	// ErrorKeyNotFound is for identifing not found error
	ErrorKeyNotFound = errors.New("Key not Found in DB")
)

const (
	scanPaginationSize = 100
)

// Create will create a new entry in DB for the value with the given table and key
func (c *Client) Create(table, key string, data interface{}) (err error) {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("while marshalling data, got: %v", err)
	}
	ok, err := c.pool.SetNX(generateKey(table, key), string(dataByte), 0).Result()
	switch {
	case !ok:
		return fmt.Errorf(
			"%w: %s",
			ErrorKeyAlreadyExist,
			fmt.Sprintf("An entry with key %s is already present in table %s", key, table),
		)
	case err != nil:
		return fmt.Errorf(
			"Creating new entry for value %v in table %s with key %s failed: %v",
			data, table, key, err,
		)
	default:
		return nil
	}
}

// GetAllKeys will collect all the keys of provided table
func (c *Client) GetAllKeys(table string) ([]string, error) {
	var allKeys []string
	var cursor uint64
	for {
		keys, c, err := c.pool.Scan(cursor, generateKey(table, "*"), scanPaginationSize).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch all keys of table %s: %s", table, err.Error())
		}
		allKeys = append(allKeys, keys...)
		if cursor == 0 {
			break
		}
		cursor = c
	}
	return allKeys, nil
}

// Get will collect the data associated with the given key from the given table
func (c *Client) Get(table, key string) (val string, err error) {
	val, err = c.pool.Get(generateKey(table, key)).Result()
	switch err {
	case redis.Nil:
		return "", fmt.Errorf(
			"%w: %s",
			ErrorKeyNotFound,
			fmt.Sprintf("Data with key %s not found in table %s", key, table),
		)
	case nil:
		return val, nil
	default:
		return "", fmt.Errorf("unable to complete the operation: %s", err.Error())
	}
}

func generateKey(table, key string) string {
	return fmt.Sprintf("%s:%s", table, key)
}
