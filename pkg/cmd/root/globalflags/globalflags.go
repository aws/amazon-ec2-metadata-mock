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

// Package globalflags represents the global flags (to be used with all sub-commands)
package globalflags

const (
	// ConfigFileFlag - config file for cli input parameters in json format
	ConfigFileFlag = "config-file"

	// SaveConfigToFileFlag - whether to save processed config from all input sources
	SaveConfigToFileFlag = "save-config-to-file"

	// MockDelayInSecFlag - spot itn delay in seconds, relative to the application start time
	MockDelayInSecFlag = "mock-delay-sec"

	// MockTriggerTimeFlag - spot itn trigger time in RFC3339
	MockTriggerTimeFlag = "mock-trigger-time"

	// MockIPCountFlag - the number of nodes in a cluster that can receive Spot ITNs
	MockIPCountFlag = "mock-ip-count"

	// HostNameFlag - the HTTP hostname for the mock url
	HostNameFlag = "hostname"

	// PortFlag - the HTTP port where the mock runs
	PortFlag = "port"

	// Imdsv2Flag - whether to enable IMDSv2 only requiring a session token when submitting requests
	Imdsv2Flag = "imdsv2"

	// RebalanceDelayInSecFlag - rebalance rec delay in seconds, relative to the application start time
	RebalanceDelayInSecFlag = "rebalance-delay-sec"

	// RebalanceTriggerTimeFlag - rebalance rec trigger time in RFC3339
	RebalanceTriggerTimeFlag = "rebalance-trigger-time"

	ASGTerminationDelayInSecFlag = "asg-termination-delay-sec"

	ASGTerminationTriggerTimeFlag = "asg-termination-trigger-time"
)

// GetTopLevelFlags returns the top level global flags
func GetTopLevelFlags() []string {
	return []string{ConfigFileFlag, SaveConfigToFileFlag, MockDelayInSecFlag, MockTriggerTimeFlag, MockIPCountFlag, Imdsv2Flag, RebalanceDelayInSecFlag, RebalanceTriggerTimeFlag, ASGTerminationDelayInSecFlag, ASGTerminationTriggerTimeFlag}
}
