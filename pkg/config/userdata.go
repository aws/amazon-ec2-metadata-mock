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

var (
	udCfgPrefix       = "userdata."
	udPathsCfgPrefix  = udCfgPrefix + "paths."
	udValuesCfgPrefix = udCfgPrefix + "values."

	// mapping of udValue KEYS to udPath KEYS requiring path substitutions on override
	udValueToPlaceholderPathsKeyMap = map[string][]string{
		// ex: "metadata.values.mac-address" : {"metadata.paths.network-interface-id", "metadata.paths.device-index", ...}
		//	udValuesCfgPrefix + publicIpv4: {
		//	udPathsCfgPrefix + macIpv4Associations}
	}

	// supported URL paths to run a mock
	udPathsDefaults = map[string]interface{}{}

	// values in mock responses
	udValuesDefaults = map[string]interface{}{}
)

// GetCfgUdValPrefix returns the prefix to use to access metadata values in config
func GetCfgUdValPrefix() string {
	return udCfgPrefix + "." + udValuesCfgPrefix + "."
}

// GetCfgUdPathsPrefix returns the prefix to use to access metadata values in config
func GetCfgUdPathsPrefix() string {
	return udCfgPrefix + "." + udPathsCfgPrefix + "."
}

// BindUserdataCfg binds a flag that represents a metadata value to configuration
func BindUserdataCfg(flag *pflag.Flag) {
	bindFlagWithKeyPrefix(flag, udValuesCfgPrefix)
}

// SetUserdataDefaults sets config defaults for metadata paths and values
func SetUserdataDefaults(jsonWithDefaults []byte) {
	// Unmarshal to map to preserve keys for Paths and Values
	var defaultsMap map[string]interface{}
	json.Unmarshal(jsonWithDefaults, &defaultsMap)
	udPaths := defaultsMap["userdata"].(map[string]interface{})["paths"].(map[string]interface{})

	udValues := defaultsMap["userdata"].(map[string]interface{})["values"].(map[string]interface{})

	for k, v := range udPaths {
		newKey := udPathsCfgPrefix + k
		// ex: "userdata": "/latest/user-data"
		udPathsDefaults[newKey] = v
	}

	for k, v := range udValues {
		newKey := udValuesCfgPrefix + k
		// ex: "userdata": "#!/bin/bash"
		udValuesDefaults[newKey] = v
	}

	LoadConfigFromDefaults(udPathsDefaults)
	LoadConfigFromDefaults(udValuesDefaults)
}

// GetUserdataDefaults returns config defaults for metadata paths and values
func GetUserdataDefaults() (map[string]interface{}, map[string]interface{}) {
	return udPathsDefaults, udValuesDefaults
}

// GetUserdataValueToPlaceholderPathsKeyMap returns collection of userdata values that are substituted into paths
func GetUserdataValueToPlaceholderPathsKeyMap() map[string][]string {
	return udValueToPlaceholderPathsKeyMap
}
