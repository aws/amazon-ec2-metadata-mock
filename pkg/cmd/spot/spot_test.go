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
	"bytes"
	"fmt"
	"testing"

	h "github.com/aws/amazon-ec2-metadata-mock/test"

	"github.com/spf13/pflag"
)

func TestNewCmdName(t *testing.T) {
	expected := "spot"
	actual := newCmd().Name()
	h.Assert(t, expected == actual, fmt.Sprintf("Expected the name for spot command to be %s, but was %s", expected, actual))
}
func TestNewCmdLocalFlags(t *testing.T) {
	expectedFlags := []string{"action", "time", "noticeTime"}

	cmd := newCmd()
	actualFlagSet := cmd.LocalFlags()

	var actualFlags []string
	actualFlagSet.VisitAll(func(flag *pflag.Flag) {
		actualFlags = append(actualFlags, flag.Name)
	})

	h.ItemsMatch(t, expectedFlags, actualFlags)
}
func TestNewCmdHasPreRunE(t *testing.T) {
	pre := newCmd().PreRunE
	h.Assert(t, pre != nil, "Expected a non nil PreRunE for the spot command")
}
func TestNewCmdHasRun(t *testing.T) {
	run := newCmd().Run
	h.Assert(t, run != nil, "Expected a non nil Run for the spot command")
}
func TestNewCmdHasExample(t *testing.T) {
	hasExample := newCmd().HasExample()
	h.Assert(t, hasExample, "Expected spot command to have an example, but wasn't found")
}
func TestExecuteHelpExists(t *testing.T) {
	cmd := newCmd()
	buf := new(bytes.Buffer)
	cmd.SetOutput(buf)
	cmd.SetArgs([]string{"-h"})
	err := cmd.Execute()
	h.Ok(t, err)

	output := buf.String()
	h.Assert(t, output != "", "Expected help subcommand for spot, but wasn't found")
}
