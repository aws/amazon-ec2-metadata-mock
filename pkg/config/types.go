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
	dynamic "github.com/aws/amazon-ec2-metadata-mock/pkg/mock/dynamic/types"
	events "github.com/aws/amazon-ec2-metadata-mock/pkg/mock/events/config"
	spot "github.com/aws/amazon-ec2-metadata-mock/pkg/mock/spot/config"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/static/types"
)

// Config set via various input sources (Sample CLI config in cli-config.json)
type Config struct {
	// ----- metadata config ----- //
	Metadata Metadata `mapstructure:"metadata"`

	// ----- CLI config ----- //
	// config keys that are also cli flags
	CfgFile          string `mapstructure:"config-file"`
	MockDelayInSec   int64  `mapstructure:"mock-delay-sec"`
	MockTriggerTime  string `mapstructure:"mock-trigger-time"`
	SaveConfigToFile bool   `mapstructure:"save-config-to-file"`
	Server           Server `mapstructure:"server"`
	Imdsv2Required   bool   `mapstructure:"imdsv2"`

	// config keys for subcommands
	SpotConfig   spot.Config   `mapstructure:"spot"`
	EventsConfig events.Config `mapstructure:"events"`

	// ----- dynamic config ----- //
	Dynamic Dynamic `mapstructure:"dynamic"`
}

// Server represents server config
type Server struct {
	HostName string `mapstructure:"hostname"`
	Port     string `mapstructure:"port"`
}

// Metadata represents metadata config used by the mock (Json values in metadata-config.json)
type Metadata struct {
	Paths  Paths  `mapstructure:"paths"`
	Values Values `mapstructure:"values"`
}

// Dynamic represents metadata config used by the mock (Json values in metadata-config.json)
type Dynamic struct {
	Paths  DynamicPaths  `mapstructure:"paths"`
	Values DynamicValues `mapstructure:"values"`
}

// Paths represents EC2 metadata paths
type Paths struct {
	AmiID                        string `mapstructure:"ami-id"`
	AmiLaunchIndex               string `mapstructure:"ami-launch-index"`
	AmiManifestPath              string `mapstructure:"ami-manifest-path"`
	BlockDeviceMappingAmi        string `mapstructure:"block-device-mapping-ami"`
	BlockDeviceMappingEbs        string `mapstructure:"block-device-mapping-ebs"`
	BlockDeviceMappingEphemeral  string `mapstructure:"block-device-mapping-ephemeral"`
	BlockDeviceMappingRoot       string `mapstructure:"block-device-mapping-root"`
	BlockDeviceMappingSwap       string `mapstructure:"block-device-mapping-swap"`
	ElasticInferenceAccelerator  string `mapstructure:"elastic-inference-accelerator"`
	ElasticInferenceAssociations string `mapstructure:"elastic-inference-associations"`
	Events                       string `mapstructure:"events"`
	Hostname                     string `mapstructure:"hostname"`
	IamInformation               string `mapstructure:"iam-info"`
	IamSecurityCredentialsRole   string `mapstructure:"iam-security-credentials-role"`
	IamSecurityCredentials       string `mapstructure:"iam-security-credentials"`
	InstanceAction               string `mapstructure:"instance-action"`
	InstanceID                   string `mapstructure:"instance-id"`
	InstanceLifecycle            string `mapstructure:"instance-life-cycle"`
	InstanceType                 string `mapstructure:"instance-type"`
	LocalHostName                string `mapstructure:"local-hostname"`
	LocalIpv4                    string `mapstructure:"local-ipv4"`
	Mac                          string `mapstructure:"mac"`
	MacDeviceNumber              string `mapstructure:"mac-device-number"`
	MacNetworkInterfaceID        string `mapstructure:"mac-network-interface-id"`
	MacIpv4Associations          string `mapstructure:"mac-ipv4-associations"`
	MacIpv6Associations          string `mapstructure:"mac-ipv6-associations"`
	MacLocalHostname             string `mapstructure:"mac-local-hostname"`
	MacLocalIpv4s                string `mapstructure:"mac-local-ipv4s"`
	MacMac                       string `mapstructure:"mac-mac"`
	MacOwnerID                   string `mapstructure:"mac-owner-id"`
	MacPublicHostname            string `mapstructure:"mac-public-hostname"`
	MacPublicIpv4s               string `mapstructure:"mac-public-ipv4s"`
	MacSecurityGroups            string `mapstructure:"mac-security-groups"`
	MacSecurityGroupIds          string `mapstructure:"mac-security-group-ids"`
	MacSubnetID                  string `mapstructure:"mac-subnet-id"`
	MacSubnetIpv4CidrBlock       string `mapstructure:"mac-subnet-ipv4-cidr-block"`
	MacSubnetIpv6CidrBlocks      string `mapstructure:"mac-subnet-ipv6-cidr-blocks"`
	MacVpcID                     string `mapstructure:"mac-vpc-id"`
	MacVpcIpv4CidrBlock          string `mapstructure:"mac-vpc-ipv4-cidr-block"`
	MacVpcIpv4CidrBlocks         string `mapstructure:"mac-vpc-ipv4-cidr-blocks"`
	MacVpcIpv6CidrBlocks         string `mapstructure:"mac-vpc-ipv6-cidr-blocks"`
	PlacementAvailabilityZone    string `mapstructure:"placement-availability-zone"`
	ProductCodes                 string `mapstructure:"product-codes"`
	PublicHostName               string `mapstructure:"public-hostname"`
	PublicIpv4                   string `mapstructure:"public-ipv4"`
	PublicKey                    string `mapstructure:"public-key"`
	ReservationID                string `mapstructure:"reservation-id"`
	SecurityGroups               string `mapstructure:"security-groups"`
	ServicesDomain               string `mapstructure:"services-domain"`
	ServicesPartition            string `mapstructure:"services-partition"`
	Spot                         string `mapstructure:"spot"`
	SpotTerminationTime          string `mapstructure:"spot-termination-time"`
}

// Values represents config used in the mock responses
type Values struct {
	AmiID                        string                            `mapstructure:"ami-id"`
	AmiLaunchIndex               string                            `mapstructure:"ami-launch-index"`
	AmiManifestPath              string                            `mapstructure:"ami-manifest-path"`
	BlockDeviceMappingAmi        string                            `mapstructure:"block-device-mapping-ami"`
	BlockDeviceMappingEbs        string                            `mapstructure:"block-device-mapping-ebs"`
	BlockDeviceMappingEphemeral  string                            `mapstructure:"block-device-mapping-ephemeral"`
	BlockDeviceMappingRoot       string                            `mapstructure:"block-device-mapping-root"`
	BlockDeviceMappingSwap       string                            `mapstructure:"block-device-mapping-swap"`
	ElasticInferenceAccelerator  types.ElasticInferenceAccelerator `mapstructure:"elastic-inference-accelerator"`
	ElasticInferenceAssociations string                            `mapstructure:"elastic-inference-associations"`
	EventID                      string                            `mapstructure:"event-id" json:"EventId"`
	Hostname                     string                            `mapstructure:"hostname"`
	IamInformation               types.IamInformation              `mapstructure:"iam-info"`
	IamSecurityCredentialsRole   string                            `mapstructure:"iam-security-credentials-role"`
	IamSecurityCredentials       types.IamSecurityCredentials      `mapstructure:"iam-security-credentials"`
	InstanceAction               string                            `mapstructure:"instance-action"`
	InstanceID                   string                            `mapstructure:"instance-id"`
	InstanceLifecycle            string                            `mapstructure:"instance-life-cycle"`
	InstanceType                 string                            `mapstructure:"instance-type"`
	LocalHostName                string                            `mapstructure:"local-hostname"`
	LocalIpv4                    string                            `mapstructure:"local-ipv4"`
	Mac                          string                            `mapstructure:"mac"`
	MacDeviceNumber              string                            `mapstructure:"mac-device-number"`
	MacNetworkInterfaceID        string                            `mapstructure:"mac-network-interface-id"`
	MacIpv4Associations          string                            `mapstructure:"mac-ipv4-associations"`
	MacIpv6Associations          string                            `mapstructure:"mac-ipv6-associations"`
	MacLocalHostname             string                            `mapstructure:"mac-local-hostname"`
	MacLocalIpv4s                string                            `mapstructure:"mac-local-ipv4s"`
	MacMac                       string                            `mapstructure:"mac-mac"`
	MacOwnerID                   string                            `mapstructure:"mac-owner-id"`
	MacPublicHostname            string                            `mapstructure:"mac-public-hostname"`
	MacPublicIpv4s               string                            `mapstructure:"mac-public-ipv4s"`
	MacSecurityGroups            string                            `mapstructure:"mac-security-groups"`
	MacSecurityGroupIds          string                            `mapstructure:"mac-security-group-ids"`
	MacSubnetID                  string                            `mapstructure:"mac-subnet-id"`
	MacSubnetIpv4CidrBlock       string                            `mapstructure:"mac-subnet-ipv4-cidr-block"`
	MacSubnetIpv6CidrBlocks      string                            `mapstructure:"mac-subnet-ipv6-cidr-blocks"`
	MacVpcID                     string                            `mapstructure:"mac-vpc-id"`
	MacVpcIpv4CidrBlock          string                            `mapstructure:"mac-vpc-ipv4-cidr-block"`
	MacVpcIpv4CidrBlocks         string                            `mapstructure:"mac-vpc-ipv4-cidr-blocks"`
	MacVpcIpv6CidrBlocks         string                            `mapstructure:"mac-vpc-ipv6-cidr-blocks"`
	PlacementAvailabilityZone    string                            `mapstructure:"placement-availability-zone"`
	ProductCodes                 string                            `mapstructure:"product-codes"`
	PublicHostName               string                            `mapstructure:"public-hostname"`
	PublicIpv4                   string                            `mapstructure:"public-ipv4"`
	PublicKey                    string                            `mapstructure:"public-key"`
	ReservationID                string                            `mapstructure:"reservation-id"`
	SecurityGroups               string                            `mapstructure:"security-groups"`
	ServicesDomain               string                            `mapstructure:"services-domain"`
	ServicesPartition            string                            `mapstructure:"services-partition"`
}

// DynamicPaths represents EC2 dynamic paths
type DynamicPaths struct {
	InstanceIdentityDocument string `mapstructure:"instance-identity-document"`
}

// DynamicValues represents EC2 dynamic paths
type DynamicValues struct {
	InstanceIdentityDocument dynamic.InstanceIdentityDocument `mapstructure:"instance-identity-document"`
}
