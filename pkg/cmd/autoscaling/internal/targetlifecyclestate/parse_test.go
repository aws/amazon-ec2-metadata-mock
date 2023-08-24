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
	"testing"
	"time"

	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/autoscaling/targetlifecyclestate"
	h "github.com/aws/amazon-ec2-metadata-mock/test"
)

func TestParse(t *testing.T) {
	t.Run("empty args", func(t *testing.T) {
		expectedState := targetlifecyclestate.InService

		p, err := parse([]string{})

		h.Ok(t, err)
		h.Assert(t, p.Initial == expectedState, "expected initial state to be %s, got %s", expectedState, p.Initial)
		h.Assert(t, len(p.Transitions) == 0, "expected no transitions, got %d", len(p.Transitions))
	})

	t.Run("initial state only", func(t *testing.T) {
		expectedState := targetlifecyclestate.Standby

		p, err := parse([]string{expectedState.String()})

		h.Ok(t, err)
		h.Assert(t, p.Initial == expectedState, "expected initial state to be %s, got %s", expectedState, p.Initial)
		h.Assert(t, len(p.Transitions) == 0, "expected no transitions, got %d", len(p.Transitions))
	})

	t.Run("initial state and transitions using RFC3339 times", func(t *testing.T) {
		firstTransitionTime := time.Now().Add(1 * time.Minute)
		secondTransitionTime := firstTransitionTime.Add(30 * time.Second)

		p, err := parse([]string{
			targetlifecyclestate.Standby.String(),
			firstTransitionTime.Format(time.RFC3339),
			targetlifecyclestate.InService.String(),
			secondTransitionTime.Format(time.RFC3339),
			targetlifecyclestate.Terminated.String(),
		})

		h.Ok(t, err)

		h.Assert(t, p.Initial == targetlifecyclestate.Standby,
			"expected initial state to be %s, got %s", targetlifecyclestate.Standby, p.Initial)

		h.Assert(t, len(p.Transitions) == 2, "expected 2 transitions, got %d", len(p.Transitions))

		h.Assert(t, equalTime(p.Transitions[0].Time, firstTransitionTime),
			"expected first transition time to be %s, got %s",
			firstTransitionTime.Format(time.RFC3339), p.Transitions[0].Time.Format(time.RFC3339))
		h.Assert(t, p.Transitions[0].State == targetlifecyclestate.InService,
			"expected first transition state to be %s, got %s",
			targetlifecyclestate.InService, p.Transitions[0].State)

		h.Assert(t, equalTime(p.Transitions[1].Time, secondTransitionTime),
			"expected second transition time to be %s, got %s",
			secondTransitionTime.Format(time.RFC3339), p.Transitions[1].Time.Format(time.RFC3339))
		h.Assert(t, p.Transitions[1].State == targetlifecyclestate.Terminated,
			"expected second transition state to be %s, got %s",
			targetlifecyclestate.Terminated, p.Transitions[1].State)
	})

	t.Run("initial state and transitions using durations", func(t *testing.T) {
		firstTransitionTime := time.Now().Add(1 * time.Minute)
		secondTransitionDuration := 30 * time.Second
		secondTransitionTime := firstTransitionTime.Add(secondTransitionDuration)
		p, err := parse([]string{
			targetlifecyclestate.Standby.String(),
			firstTransitionTime.Format(time.RFC3339),
			targetlifecyclestate.InService.String(),
			secondTransitionDuration.String(),
			targetlifecyclestate.Terminated.String(),
		})

		h.Ok(t, err)

		h.Assert(t, p.Initial == targetlifecyclestate.Standby,
			"expected initial state to be %s, got %s", targetlifecyclestate.Standby, p.Initial)

		h.Assert(t, len(p.Transitions) == 2, "expected 2 transitions, got %d", len(p.Transitions))

		h.Assert(t, equalTime(p.Transitions[0].Time, firstTransitionTime),
			"expected first transition time to be %s, got %s",
			firstTransitionTime.Format(time.RFC3339), p.Transitions[0].Time.Format(time.RFC3339))
		h.Assert(t, p.Transitions[0].State == targetlifecyclestate.InService,
			"expected first transition state to be %s, got %s",
			targetlifecyclestate.InService, p.Transitions[0].State)

		h.Assert(t, equalTime(p.Transitions[1].Time, secondTransitionTime),
			"expected second transition time to be %s, got %s",
			secondTransitionTime.Format(time.RFC3339), p.Transitions[1].Time.Format(time.RFC3339))
		h.Assert(t, p.Transitions[1].State == targetlifecyclestate.Terminated,
			"expected second transition state to be %s, got %s",
			targetlifecyclestate.Terminated, p.Transitions[1].State)
	})
}

func equalTime(t1, t2 time.Time) bool {
	return t1.Format(time.RFC3339) == t2.Format(time.RFC3339)
}
