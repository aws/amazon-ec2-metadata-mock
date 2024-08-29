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

package asglifecycle

import (
	"log"
	"net/http"
	"strings"
	"time"

	cfg "github.com/aws/amazon-ec2-metadata-mock/pkg/config"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/server"
)

const (
	targetLifeCycleStatePath = "/latest/meta-data/autoscaling/target-lifecycle-state"
)

var (
	eligibleIPs          = make(map[string]bool)
	c                    cfg.Config
	autoscalingStartTime int64 = time.Now().Unix()
	state                      = "InService"
)

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
	case targetLifeCycleStatePath:
		handleASGTargetLifecycleState(res, req)
	}
}

func handleASGTargetLifecycleState(res http.ResponseWriter, req *http.Request) {
	requestTime := time.Now().Unix()
	if c.ASGTerminationTriggerTime != "" {
		triggerTime, _ := time.Parse(time.RFC3339, c.ASGTerminationTriggerTime)
		delayRemaining := triggerTime.Unix() - requestTime
		if delayRemaining <= 0 {
			state = "Terminated"
		}
	} else {
		delayInSeconds := c.ASGTerminationDelayInSec
		delayRemaining := delayInSeconds - (requestTime - autoscalingStartTime)
		if delayRemaining <= 0 {
			state = "Terminated"
		}
	}

	switch req.Method {
	case "GET":
		server.FormatAndReturnTextResponse(res, state)
	}
}
