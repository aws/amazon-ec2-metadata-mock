// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

package server

import (
	"net/http/httptest"
	"testing"

	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/dynamic/types"
	h "github.com/aws/amazon-ec2-metadata-mock/test"
)

func TestFormatAndReturnJSONResponse(t *testing.T) {
	expected := `{
	"accountId": "123456789012",
	"imageId": "ami-02f471c4f805553d3",
	"availabilityZone": "us-east-1a",
	"ramdiskId": null,
	"kernelId": null,
	"devpayProductCodes": null,
	"marketplaceProductCodes": ["4i20ezfza3p7xx2kt2g8weu2u","entry"],
	"version": "2017-09-30",
	"privateIp": "172.31.85.190",
	"billingProducts": null,
	"instanceId": "i-048bcb15d2686eec7",
	"pendingTime": "2022-06-23T06:21:55Z",
	"architecture": "x86_64",
	"instanceType": "t2.nano",
	"region": "us-east-1"
}`

	testDoc := types.InstanceIdentityDocument{
		AccountId:               "123456789012",
		ImageId:                 "ami-02f471c4f805553d3",
		AvailabilityZone:        "us-east-1a",
		RamdiskId:               nil,
		KernelId:                nil,
		DevpayProductCodes:      nil,
		MarketplaceProductCodes: []string{"4i20ezfza3p7xx2kt2g8weu2u", "entry"},
		Version:                 "2017-09-30",
		PrivateIp:               "172.31.85.190",
		BillingProducts:         nil,
		InstanceId:              "i-048bcb15d2686eec7",
		PendingTime:             "2022-06-23T06:21:55Z",
		Architecture:            "x86_64",
		InstanceType:            "t2.nano",
		Region:                  "us-east-1",
	}
	rr := httptest.NewRecorder()
	FormatAndReturnJSONResponse(rr, testDoc)
	actual := rr.Body.String()
	h.Assert(t, expected == actual, "FormatAndReturnJSONResponse did not format InstanceIdentityDocument as expected.")
}
