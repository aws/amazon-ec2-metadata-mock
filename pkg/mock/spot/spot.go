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

package spot

import (
	"log"
	"net/http"
	"time"

	cfg "github.com/aws/amazon-ec2-metadata-mock/pkg/config"
	t "github.com/aws/amazon-ec2-metadata-mock/pkg/mock/spot/internal/types"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/server"
)

const (
	instanceActionPath  = "/latest/meta-data/spot/instance-action"
	terminationTimePath = "/latest/meta-data/spot/termination-time"
)

var spotItnStartTime int64 = time.Now().Unix()
var c cfg.Config

// Mock starts spot itn mock
func Mock(config cfg.Config) {
	SetConfig(config)

	server.ListenAndServe(config.Server.HostName, config.Server.Port)
}

// SetConfig sets the local config
func SetConfig(config cfg.Config) {
	c = config
}

// Handler processes http requests
func Handler(res http.ResponseWriter, req *http.Request) {
	log.Println("Received request to mock spot interruption:", req.URL.Path)

	requestTime := time.Now().Unix()

	if c.MockTriggerTime != "" {
		triggerTime, _ := time.Parse(time.RFC3339, c.MockTriggerTime)

		delayRemaining := triggerTime.Unix() - requestTime
		if delayRemaining > 0 {
			log.Printf("MockTriggerTime %s was not reached yet. The mock response will be available in %ds. Returning `notFoundResponse` for now", triggerTime, delayRemaining)
			server.ReturnNotFoundResponse(res)
			return
		}
	} else {
		delayInSeconds := c.MockDelayInSec
		delayRemaining := delayInSeconds - (requestTime - spotItnStartTime)
		if delayRemaining > 0 {
			log.Printf("Delaying the response by %ds as requested. The mock response will be available in %ds. Returning `notFoundResponse` for now", delayInSeconds, delayRemaining)
			server.ReturnNotFoundResponse(res)
			return
		}
	}

	// default time to requestTime + 2min, unless overridden
	timePlus2Min := time.Now().UTC().Add(time.Minute * time.Duration(2)).Format(time.RFC3339)
	mockResponseTime := timePlus2Min
	if c.SpotConfig.TerminationTime != "" {
		mockResponseTime = c.SpotConfig.TerminationTime
	}

	// return mock response after the delay or trigger time has elapsed
	switch req.URL.Path {
	case instanceActionPath:
		server.FormatAndReturnJSONResponse(res, getInstanceActionResponse(mockResponseTime))
	case terminationTimePath:
		server.FormatAndReturnTextResponse(res, mockResponseTime)
	}
}

func getInstanceActionResponse(time string) t.InstanceActionResponse {
	return t.InstanceActionResponse{
		Action: c.SpotConfig.InstanceAction,
		Time:   time,
	}
}
