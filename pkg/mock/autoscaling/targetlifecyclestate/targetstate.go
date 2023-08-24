// Copyright 2023 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

package targetlifecyclestate

import (
	"fmt"
	"sort"
)

// A TargetState represents a valid autoscaling lifecyle state that is being
// transitioned to. The zero value is "InService".
//
// Valid states are from https://docs.aws.amazon.com/autoscaling/ec2/userguide/retrieving-target-lifecycle-state-through-imds.html
type TargetState int

const (
	InService TargetState = iota
	Detached
	Standby
	Terminated
	WarmedHibernated
	WarmedStopped
	WarmedRunning
	WarmedTerminated
)

var (
	stateToName = map[TargetState]string{
		InService:        "InService",
		Detached:         "Detached",
		Standby:          "Standby",
		Terminated:       "Terminated",
		WarmedHibernated: "Warmed:Hibernated",
		WarmedStopped:    "Warmed:Stopped",
		WarmedRunning:    "Warmed:Running",
		WarmedTerminated: "Warmed:Terminated",
	}

	nameToState map[string]TargetState

	names []string
)

func init() {
	nameToState = make(map[string]TargetState, len(stateToName))
	names = make([]string, 0, len(stateToName))
	for state, name := range stateToName {
		nameToState[name] = state
		names = append(names, name)
	}

	sort.Strings(names)
}

// String returns the name of the target state.
func (s TargetState) String() string {
	return stateToName[s]
}

// Parse returns the target state for the given name. Name matching is case-sensitive.
func Parse(targetStateName string) (TargetState, error) {
	state, ok := nameToState[targetStateName]
	if !ok {
		return InService, fmt.Errorf("invalid autoscaling lifecycle target state name: %s", targetStateName)
	}
	return state, nil
}

// Names returns all valid target state names in alphabetical order.
func Names() []string {
	return append([]string(nil), names...)
}
