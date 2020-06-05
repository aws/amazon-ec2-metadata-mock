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

// Event represents the metadata structure for mock json response parsing
type Event struct {
	Code              string `mapstructure:"code"`
	Description       string `mapstructure:"description"`
	State             string `mapstructure:"state"` // State of the scheduled event
	EventID           string `mapstructure:"event-id" json:"EventId"`
	NotBefore         string `mapstructure:"not-before"`                    // The earliest start time for the scheduled event
	NotAfter          string `mapstructure:"not-after,omitempty"`           // The latest end time for the scheduled event
	NotBeforeDeadline string `mapstructure:"not-before-deadline,omitempty"` // The deadline for starting the event
}
