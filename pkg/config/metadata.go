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
	// metadata keys
	mac                          = "mac"
	publicIpv4                   = "public-ipv4"
	iamSecurityCredentialsRole   = "iam-security-credentials-role"
	iamSecurityCredentials       = "iam-security-credentials"
	elasticInferenceAssociations = "elastic-inference-associations"
	elasticInferenceAccelerator  = "elastic-inference-accelerator"
	// prefixed with "mac" because query requires mac address and to avoid collisions with existing metadata keys
	// ex: network/interfaces/macs/<mac_address>/device-number
	macDeviceNumber         = "mac-device-number"
	macNetworkInterfaceId   = "mac-network-interface-id"
	macIpv4Associations     = "mac-ipv4-associations"
	macIpv6Associations     = "mac-ipv6-associations"
	macLocalHostname        = "mac-local-hostname"
	macLocalIpv4s           = "mac-local-ipv4s"
	macMac                  = "mac-mac"
	macOwnerId              = "mac-owner-id"
	macPublicHostname       = "mac-public-hostname"
	macPublicIpv4s          = "mac-public-ipv4s"
	macSecurityGroups       = "mac-security-groups"
	macSecurityGroupIds     = "mac-security-group-ids"
	macSubnetId             = "mac-subnet-id"
	macSubnetIpv4CidrBlock  = "mac-subnet-ipv4-cidr-block"
	macSubnetIpv6CidrBlocks = "mac-subnet-ipv6-cidr-blocks"
	macVpcId                = "mac-vpc-id"
	macVpcIpv4CidrBlock     = "mac-vpc-ipv4-cidr-block"
	macVpcIpv4CidrBlocks    = "mac-vpc-ipv4-cidr-blocks"
	macVpcIpv6CidrBlocks    = "mac-vpc-ipv6-cidr-blocks"
)

var (
	mdCfgPrefix       = "metadata."
	mdPathsCfgPrefix  = mdCfgPrefix + "paths."
	mdValuesCfgPrefix = mdCfgPrefix + "values."

	// mapping of mdValue KEYS to mdPath KEYS requiring path substitutions on override
	mdValueToPlaceholderPathsKeyMap = map[string][]string{
		// ex: "metadata.values.mac-address" : {"metadata.paths.network-interface-id", "metadata.paths.device-index", ...}
		mdValuesCfgPrefix + mac: {
			mdPathsCfgPrefix + macDeviceNumber, mdPathsCfgPrefix + macNetworkInterfaceId, mdPathsCfgPrefix + macIpv4Associations, mdPathsCfgPrefix + macLocalHostname,
			mdPathsCfgPrefix + macLocalIpv4s, mdPathsCfgPrefix + macMac, mdPathsCfgPrefix + macOwnerId, mdPathsCfgPrefix + macPublicHostname, mdPathsCfgPrefix + macPublicIpv4s,
			mdPathsCfgPrefix + macSecurityGroups, mdPathsCfgPrefix + macSecurityGroupIds, mdPathsCfgPrefix + macSubnetId, mdPathsCfgPrefix + macSubnetIpv4CidrBlock, mdPathsCfgPrefix + macVpcId,
			mdPathsCfgPrefix + macVpcIpv4CidrBlock, mdPathsCfgPrefix + macVpcIpv4CidrBlocks, mdPathsCfgPrefix + macIpv6Associations, mdPathsCfgPrefix + macSubnetIpv6CidrBlocks, mdPathsCfgPrefix + macVpcIpv6CidrBlocks},
		mdValuesCfgPrefix + publicIpv4: {
			mdPathsCfgPrefix + macIpv4Associations},
		mdValuesCfgPrefix + iamSecurityCredentialsRole: {
			mdPathsCfgPrefix + iamSecurityCredentials},
		mdValuesCfgPrefix + elasticInferenceAssociations: {
			mdPathsCfgPrefix + elasticInferenceAccelerator},
	}

	// supported URL paths to run a mock
	mdPathsDefaults = map[string]interface{}{}

	// values in mock responses
	mdValuesDefaults = map[string]interface{}{}
)

// GetCfgMdValPrefix returns the prefix to use to access metadata values in config
func GetCfgMdValPrefix() string {
	return mdCfgPrefix + "." + mdValuesCfgPrefix + "."
}

// GetCfgMdPathsPrefix returns the prefix to use to access metadata values in config
func GetCfgMdPathsPrefix() string {
	return mdCfgPrefix + "." + mdPathsCfgPrefix + "."
}

// BindMetadataCfg binds a flag that represents a metadata value to configuration
func BindMetadataCfg(flag *pflag.Flag) {
	bindFlagWithKeyPrefix(flag, mdValuesCfgPrefix)
}

// SetMetadataDefaults sets config defaults for metadata paths and values
func SetMetadataDefaults(jsonWithDefaults []byte) {
	// Unmarshal to map to preserve keys for Paths and Values
	var defaultsMap map[string]interface{}
	json.Unmarshal(jsonWithDefaults, &defaultsMap)

	mdPaths := defaultsMap["metadata"].(map[string]interface{})["paths"].(map[string]interface{})
	mdValues := defaultsMap["metadata"].(map[string]interface{})["values"].(map[string]interface{})

	for k, v := range mdPaths {
		newKey := mdPathsCfgPrefix + k
		// ex: "metadata.paths.ami-id": "/latest/meta-data/ami-id"
		mdPathsDefaults[newKey] = v
	}

	for k, v := range mdValues {
		newKey := mdValuesCfgPrefix + k
		// ex: "metadata.values.ami-id": "ami-0a887e401f7654935"
		mdValuesDefaults[newKey] = v
	}

	LoadConfigFromDefaults(mdPathsDefaults)
	LoadConfigFromDefaults(mdValuesDefaults)
}

// GetMetadataDefaults returns config defaults for metadata paths and values
func GetMetadataDefaults() (map[string]interface{}, map[string]interface{}) {
	return mdPathsDefaults, mdValuesDefaults
}

// GetMetadataValueToPlaceholderPathsKeyMap returns collection of metadata values that are substituted into paths
func GetMetadataValueToPlaceholderPathsKeyMap() map[string][]string {
	return mdValueToPlaceholderPathsKeyMap
}
