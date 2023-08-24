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
	"time"

	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/autoscaling/targetlifecyclestate"
)

type (
	targetLifecycleStateSchedule struct {
		Initial     targetlifecyclestate.TargetState
		Transitions []scheduledTransition
	}

	scheduledTransition struct {
		State targetlifecyclestate.TargetState
		Time  time.Time
	}
)

// parse parses the target lifecycle state transitions from the given arguments.
// Args should be empty or have an odd length. Every even indexed item must be
// a valid lifecycle target state, odd indexed items must be a RFC3339 time or
// duration.
func parse(args []string) (p targetLifecycleStateSchedule, err error) {
	if len(args) == 0 {
		return p, nil
	}

	p.Initial, err = targetlifecyclestate.Parse(args[0])
	if err != nil {
		return p, fmt.Errorf("parsing initial state: %w", err)
	}

	args = args[1:]
	if len(args)%2 != 0 {
		return p, fmt.Errorf("invalid number of arguments")
	}

	for i := 0; i < len(args); i += 2 {
		t, tErr := time.Parse(time.RFC3339, args[i])
		if tErr != nil {
			d, dErr := time.ParseDuration(args[i])
			if dErr != nil {
				return p, fmt.Errorf("%q is not a duration or RFC3339 time: %w: %w", args[i], dErr, tErr)
			}

			t = time.Now()
			if l := len(p.Transitions); l > 0 {
				t = p.Transitions[l-1].Time
			}
			t = t.Add(d)
		}
		if t.Before(time.Now()) {
			return p, fmt.Errorf("%s is in the past", t)
		}

		st, err := targetlifecyclestate.Parse(args[i+1])
		if err != nil {
			return p, fmt.Errorf("parsing state: %w", err)
		}

		p.Transitions = append(p.Transitions, scheduledTransition{
			State: st,
			Time:  t,
		})
	}

	return p, nil
}
