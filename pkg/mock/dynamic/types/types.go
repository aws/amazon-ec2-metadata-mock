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

// InstanceIdentityDocument structure for mock json response parsing
type InstanceIdentityDocument struct {
	AccountId               string `json:"accountId"`
	ImageId                 string `json:"imageId"`
	AvailabilityZone        string `json:"availabilityZone"`
	RamdiskId               string `json:"ramdiskId"`
	KernelId                string `json:"kernelId"`
	DevpayProductCodes      string `json:"devpayProductCodes"`
	MarketplaceProductCodes string `json:"marketplaceProductCodes"`
	Version                 string `json:"version"`
	PrivateIp               string `json:"privateIp"`
	BillingProducts         string `json:"billingProducts"`
	InstanceId              string `json:"instanceId"`
	PendingTime             string `json:"pendingTime"`
	Architecture            string `json:"architecture"`
	InstanceType            string `json:"instanceType"`
	Region                  string `json:"region"`
}
