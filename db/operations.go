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
	"errors"
	"fmt"
	"strings"

	"github.com/go-redis/redis"
)

// Using below variables as part of errors will enabling errors.Is() function
var (
	// ErrorServiceUnavailable is for identifing DB connection errors
	ErrorServiceUnavailable = errors.New("Failed to establish connection with the DB")
	// ErrorKeyAlreadyExist is for identifing already exist error
	ErrorKeyAlreadyExist = errors.New("Key already exist in DB")
	// ErrorKeyNotFound is for identifing not found error
	ErrorKeyNotFound = errors.New("Key not Found in DB")
)

type dbCalls interface {
	Create(table, resourceID, data string) (err error)
	Update(table, resourceID, data string) (err error)
	GetAllMatchingKeys(table, pattern string) ([]string, error)
	Get(table, resourceID string) (string, error)
	UpdateKeySet(key string, member string) (err error)
	GetKeySetMembers(key string) (list []string, err error)
	Delete(table, resourceID string) (err error)
	DeleteKeySetMembers(key string, member string) (err error)
}

// Connector is the interface which connects the DB functions
var Connector dbCalls

// connector is used as a receiver for DB communication functions
type connector struct{}

// Create will create a new entry in DB for the value with the given table and resourceID
func (d connector) Create(table, resourceID, data string) (err error) {
	c, err := getClient()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrorServiceUnavailable, err)
	}
	ok, err := c.pool.SetNX(generateKey(table, resourceID), data, 0).Result()
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

// Update will update an entry in DB with the value for the given table and resourceID
func (d connector) Update(table, resourceID, data string) (err error) {
	c, err := getClient()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrorServiceUnavailable, err)
	}
	if err = c.pool.Set(generateKey(table, resourceID), data, 0).Err(); err != nil {
		return fmt.Errorf(
			"Updating entry for resource id %s with value %v in table %s failed: %v",
			data, table, resourceID, err,
		)
	}
	return nil
}

// GetAllMatchingKeys will collect all the keys of provided table and pattern
func (d connector) GetAllMatchingKeys(table, pattern string) ([]string, error) {
	var allKeys []string
	c, err := getClient()
	if err != nil {
		return allKeys, fmt.Errorf("%w: %v", ErrorServiceUnavailable, err)
	}
	var cursor uint64 = 0
	for {
		keys, c, err := c.pool.Scan(cursor, generateKey(table, pattern+"*"), scanPaginationSize).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch all keys of table %s: %s", table, err.Error())
		}
		cursor = c
		allKeys = append(allKeys, keys...)
		if cursor == 0 {
			break
		}
	}
	return trimTableFromKeys(table, allKeys), nil
}

// Get will collect the data associated with the given key from the given table
func (d connector) Get(table, resourceID string) (string, error) {
	c, err := getClient()
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrorServiceUnavailable, err)
	}
	val, err := c.pool.Get(generateKey(table, resourceID)).Result()
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

// UpdateKeySet will add passed member to the particular key set.
func (d connector) UpdateKeySet(key string, member string) (err error) {
	c, err := getClient()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrorServiceUnavailable, err)
	}
	if err = c.pool.SAdd(key, member).Err(); err != nil {
		return fmt.Errorf(
			"Updating key set %s with member %v failed: %v",
			key, member, err,
		)
	}
	return nil
}

// GetKeySetMembers will get the list of member in the particular key set.
func (d connector) GetKeySetMembers(key string) (list []string, err error) {
	c, err := getClient()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrorServiceUnavailable, err)
	}
	if list, err = c.pool.SMembers(key).Result(); err != nil {
		return nil, fmt.Errorf(
			"Getting list of member in the key set %s failed: %v",
			key, err,
		)
	}
	return list, nil
}

// Delete will delete the data associated with the given key from the given table
func (d connector) Delete(table, resourceID string) (err error) {
	c, err := getClient()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrorServiceUnavailable, err)
	}
	err = c.pool.Del(generateKey(table, resourceID)).Err()
	switch err {
	case redis.Nil:
		return fmt.Errorf(
			"%w: %s",
			ErrorKeyNotFound,
			fmt.Sprintf("Data with resource ID %s not found in table %s", resourceID, table),
		)
	case nil:
		return nil
	default:
		return fmt.Errorf("unable to complete the operation: %s", err.Error())
	}
}

// DeleteKeySetMembers will delete the list of member in the particular key set.
func (d connector) DeleteKeySetMembers(key string, member string) (err error) {
	c, err := getClient()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrorServiceUnavailable, err)
	}
	if err = c.pool.SRem(key, member).Err(); err != nil {
		return fmt.Errorf("Deleting member from the key set %s failed: %v", key, err)
	}
	return nil
}
