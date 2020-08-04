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

package events

import (
	"log"
	"net/http"
	"time"

	cfg "github.com/aws/amazon-ec2-metadata-mock/pkg/config"
	t "github.com/aws/amazon-ec2-metadata-mock/pkg/mock/events/internal/types"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/server"
)

const (
	descriptionPrefix = "The instance is scheduled for "
	timeLayout        = "2 Jan 2006 15:04:05 GMT"
)

var appStartTime int64 = time.Now().Unix()
var c cfg.Config

// Mock starts scheduled events mock
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
	log.Println("Received request to mock EC2 event:", req.URL.Path)

	delayInSeconds := c.MockDelayInSec
	requestTime := time.Now().Unix()
	delayRemaining := delayInSeconds - (requestTime - appStartTime)
	if delayRemaining > 0 {
		log.Printf("Delaying the response by %ds as requested. The mock response will be available in %ds. Returning `notFoundResponse` for now", delayInSeconds, delayRemaining)
		server.ReturnNotFoundResponse(res)
		return
	}

	// return mock response after the delay has elapsed
	server.FormatAndReturnJSONResponse(res, getMetadata())
}

func getMetadata() []t.Event {
	md := c.Metadata.Values
	se := c.EventsConfig

	b, _ := time.Parse(time.RFC3339, se.NotBefore)
	a, _ := time.Parse(time.RFC3339, se.NotAfter)
	bd, _ := time.Parse(time.RFC3339, se.NotBeforeDeadline)
	eventResp := t.Event{
		Code:              se.EventCode,
		Description:       descriptionPrefix + se.EventCode,
		EventID:           md.EventID,
		State:             se.EventState,
		NotBefore:         b.Format(timeLayout),
		NotAfter:          a.Format(timeLayout),
		NotBeforeDeadline: bd.Format(timeLayout),
	}
	// supports 1 scheduled event for now
	return []t.Event{eventResp}
}
