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
	"bytes"
	"testing"

	h "github.com/aws/amazon-ec2-metadata-mock/test"
)

func TestNewCmdName(t *testing.T) {
	expected := "target-lifecycle-state"
	actual := newCmd().Name()

	h.Assert(t, expected == actual, "expected %q, got %q", expected, actual)
}

func TestNewCmdHasExample(t *testing.T) {
	h.Assert(t, newCmd().HasExample(), "expected command to have an example")
}

func TestExecuteHelpExists(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := newCmd()
	cmd.SetOutput(buf)
	cmd.SetArgs([]string{"-h"})

	err := cmd.Execute()
	h.Ok(t, err)

	h.Assert(t, buf.Len() > 0, "expected command to have help documentation")
}
