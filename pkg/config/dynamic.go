// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package config

import (
	"encoding/json"

	"github.com/spf13/pflag"
)

const (
	// dynamic keys
	instanceIdentityDocument = "instance-identity-document"
)

var (
	dyCfgPrefix       = "dynamic."
	dyPathsCfgPrefix  = dyCfgPrefix + "paths."
	dyValuesCfgPrefix = dyCfgPrefix + "values."

	// mapping of dyValue KEYS to dyPath KEYS requiring path substitutions on override
	dyValueToPlaceholderPathsKeyMap = map[string][]string{}

	// supported URL paths to run a mock
	dyPathsDefaults = map[string]interface{}{}

	// values in mock responses
	dyValuesDefaults = map[string]interface{}{}

	// mapping of dynamic value keys to its nested struct
	dyNestedValues = map[string]interface{}{}
)

// GetCfgDnValPrefix returns the prefix to use to access dynamic values in config
func GetCfgDnValPrefix() string {
	return dyCfgPrefix + "." + dyValuesCfgPrefix + "."
}

// GetCfgDnPathsPrefix returns the prefix to use to access dynamic values in config
func GetCfgDnPathsPrefix() string {
	return dyCfgPrefix + "." + dyPathsCfgPrefix + "."
}

// BindDynamicCfg binds a flag that represents a dynamic value to configuration
func BindDynamicCfg(flag *pflag.Flag) {
	bindFlagWithKeyPrefix(flag, dyValuesCfgPrefix)
}

// SetDynamicDefaults sets config defaults for dynamic paths and values
func SetDynamicDefaults(jsonWithDefaults []byte) {
	// Unmarshal to map to preserve keys for Paths and Values
	var defaultsMap map[string]interface{}
	json.Unmarshal(jsonWithDefaults, &defaultsMap)

	dyPaths := defaultsMap["dynamic"].(map[string]interface{})["paths"].(map[string]interface{})
	dyValues := defaultsMap["dynamic"].(map[string]interface{})["values"].(map[string]interface{})

	for k, v := range dyPaths {
		newKey := dyPathsCfgPrefix + k
		// ex: "dynamic.paths.instance-identity-document": "/latest/dynamic/instance-identity/document"
		dyPathsDefaults[newKey] = v
	}

	for k, v := range dyValues {
		newKey := dyValuesCfgPrefix + k
		// ex: "dynamic.values.instance-identity-document": {"accountId" ...}
		dyValuesDefaults[newKey] = v

		// if dyvalue is a nested struct, then re-unmarshal json data to correct type
		if nestedStruct, ok := dyNestedValues[newKey]; ok {
			updatedVal, err := unmarshalToNestedStruct(v, nestedStruct)
			if err == nil {
				dyValuesDefaults[newKey] = updatedVal
			}
		}
	}

	LoadConfigFromDefaults(dyPathsDefaults)
	LoadConfigFromDefaults(dyValuesDefaults)
}

// GetDynamicDefaults returns config defaults for dynamic paths and values
func GetDynamicDefaults() (map[string]interface{}, map[string]interface{}) {
	return dyPathsDefaults, dyValuesDefaults
}

// GetDynamicValueToPlaceholderPathsKeyMap returns collection of dynamic values that are substituted into paths
func GetDynamicValueToPlaceholderPathsKeyMap() map[string][]string {
	return dyValueToPlaceholderPathsKeyMap
}
