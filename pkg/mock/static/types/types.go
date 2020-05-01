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

package types

// IamInformation metadata structure for mock json response parsing
type IamInformation struct {
	Code               string `json:"code"`
	LastUpdated        string `json:"last-updated"`
	InstanceProfileArn string `json:"instance-profile-arn"`
	InstanceProfileId  string `json:"instance-profile-id"`
}

// IamSecurityCredentials metadata structure for mock json response parsing
type IamSecurityCredentials struct {
	Code            string `json:"code"`
	LastUpdated     string `json:"last-updated"`
	Type            string `json:"type"`
	AccessKeyId     string `json:"access-key-id"`
	SecretAccessKey string `json:"secret-access-key"`
	Token           string `json:"token"`
	Expiration      string `json:"expiration"`
}

// ElasticInferenceAccelerator metadata structure for mock json response parsing
type ElasticInferenceAccelerator struct {
	Version elasticInferenceAcceleratorMetadata `json:"version_2018_04_12"`
}

type elasticInferenceAcceleratorMetadata struct {
	ElasticInferenceAcceleratorId   string `json:"elastic-inference-accelerator-id"`
	ElasticInferenceAcceleratorType string `json:"elastic-inference-accelerator-type"`
}
