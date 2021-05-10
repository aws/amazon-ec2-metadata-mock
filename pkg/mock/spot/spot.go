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
	"strings"
	"time"

	cfg "github.com/aws/amazon-ec2-metadata-mock/pkg/config"
	t "github.com/aws/amazon-ec2-metadata-mock/pkg/mock/spot/internal/types"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/server"
)

const (
	instanceActionPath  = "/latest/meta-data/spot/instance-action"
	terminationTimePath = "/latest/meta-data/spot/termination-time"
	rebalanceRecPath    = "/latest/meta-data/events/recommendations/rebalance"
)

var (
	eligibleIPs            = make(map[string]bool)
	spotItnStartTime int64 = time.Now().Unix()
	c                cfg.Config
)

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
	// specify negative value to disable this feature
	if c.MockIPCount >= 0 {
		// req.RemoteAddr is formatted as IP:port
		requestIP := strings.Split(req.RemoteAddr, ":")[0]
		if !eligibleIPs[requestIP] {
			if len(eligibleIPs) < c.MockIPCount {
				eligibleIPs[requestIP] = true
			} else {
				log.Printf("Requesting IP %s is not eligible for Spot ITN or Rebalance Recommendation because the max number of IPs configured (%d) has been reached.\n", requestIP, c.MockIPCount)
				server.ReturnNotFoundResponse(res)
				return
			}
		}
	}
	switch req.URL.Path {
	case instanceActionPath, terminationTimePath:
		handleSpotITN(res, req)
	case rebalanceRecPath:
		handleRebalance(res, req)
	}
}

func handleSpotITN(res http.ResponseWriter, req *http.Request) {
	requestTime := time.Now().Unix()
	if c.MockTriggerTime != "" {
		triggerTime, _ := time.Parse(time.RFC3339, c.MockTriggerTime)
		delayRemaining := triggerTime.Unix() - requestTime
		if delayRemaining > 0 {
			log.Printf("MockTriggerTime %s was not reached yet. The spot itn will be available in %ds. Returning `notFoundResponse` for now", triggerTime, delayRemaining)
			server.ReturnNotFoundResponse(res)
			return
		}
	} else {
		delayInSeconds := c.MockDelayInSec
		delayRemaining := delayInSeconds - (requestTime - spotItnStartTime)
		if delayRemaining > 0 {
			log.Printf("Delaying the response by %ds as requested. The spot itn will be available in %ds. Returning `notFoundResponse` for now", delayInSeconds, delayRemaining)
			server.ReturnNotFoundResponse(res)
			return
		}
	}
	// default time to requestTime + 2min, unless overridden
	mockResponseTime := time.Now().UTC().Add(time.Minute * time.Duration(2)).Format(time.RFC3339)
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

func handleRebalance(res http.ResponseWriter, req *http.Request) {
	requestTime := time.Now().Unix()
	if c.RebalanceTriggerTime != "" {
		triggerTime, _ := time.Parse(time.RFC3339, c.RebalanceTriggerTime)
		delayRemaining := triggerTime.Unix() - requestTime
		if delayRemaining > 0 {
			log.Printf("RebalanceTriggerTime %s was not reached yet. The rebalance rec will be available in %ds. Returning `notFoundResponse` for now", triggerTime, delayRemaining)
			server.ReturnNotFoundResponse(res)
			return
		}
	} else {
		delayInSeconds := c.RebalanceDelayInSec
		delayRemaining := delayInSeconds - (requestTime - spotItnStartTime)
		if delayRemaining > 0 {
			log.Printf("Delaying the response by %ds as requested. The rebalance rec will be available in %ds. Returning `notFoundResponse` for now", delayInSeconds, delayRemaining)
			server.ReturnNotFoundResponse(res)
			return
		}
	}
	// default time to requestTime, unless overridden
	mockResponseTime := time.Now().UTC().Format(time.RFC3339)
	if c.SpotConfig.RebalanceRecTime != "" {
		mockResponseTime = c.SpotConfig.RebalanceRecTime
	}
	server.FormatAndReturnJSONResponse(res, t.RebalanceRecommendationResponse{NoticeTime: mockResponseTime})
}

func getInstanceActionResponse(time string) t.InstanceActionResponse {
	return t.InstanceActionResponse{
		Action: c.SpotConfig.InstanceAction,
		Time:   time,
	}
}
