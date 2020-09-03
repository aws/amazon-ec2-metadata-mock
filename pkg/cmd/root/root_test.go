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

package root

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	h "github.com/aws/amazon-ec2-metadata-mock/test"

	"github.com/spf13/pflag"
)

func TestNewCmdName(t *testing.T) {
	expected := "ec2-metadata-mock"
	actual := NewCmd().Name()

	h.Assert(t, expected == actual, fmt.Sprintf("Expected the name for root command to be %s, but was %s", expected, actual))
}
func TestNewCmdFlags(t *testing.T) {
	expectedFlags := []string{"config-file", "save-config-to-file", "mock-delay-sec", "mock-trigger-time", "termination-nodes", "hostname", "port", "imdsv2"}

	cmd := NewCmd()
	actualFlagSet := cmd.PersistentFlags()
	var actualFlags []string
	actualFlagSet.VisitAll(func(flag *pflag.Flag) {
		actualFlags = append(actualFlags, flag.Name)
	})

	h.ItemsMatch(t, expectedFlags, actualFlags)
}
func TestNewCmdHasSubcommands(t *testing.T) {
	expSubcommandNames := []string{"spot", "events"}

	cmd := NewCmd()
	actSubcommands := cmd.Commands()
	var actSubcommandsNames []string
	for _, cmd := range actSubcommands {
		n := strings.Split(cmd.Use, " ")[0]
		actSubcommandsNames = append(actSubcommandsNames, n)
	}

	h.ItemsMatch(t, expSubcommandNames, actSubcommandsNames)
}
func TestNewCmdHasPersistentRunE(t *testing.T) {
	ppe := NewCmd().PersistentPreRunE
	h.Assert(t, ppe != nil, "Expected a non nil PersistentPreRunE for the root command")
}
func TestNewCmdHasPreRunE(t *testing.T) {
	pe := NewCmd().PreRunE
	h.Assert(t, pe != nil, "Expected a non nil PreRunE for the root command")
}
func TestNewCmdHasRun(t *testing.T) {
	run := NewCmd().Run
	h.Assert(t, run != nil, "Expected a non nil Run for the root command")
}
func TestNewCmdHasExample(t *testing.T) {
	hasExample := NewCmd().HasExample()
	h.Assert(t, hasExample, "Expected root command to have an example, but wasn't found")
}
func TestExecuteHelpExists(t *testing.T) {
	cmd := NewCmd()
	buf := new(bytes.Buffer)
	cmd.SetOutput(buf)
	cmd.SetArgs([]string{"-h"})
	err := cmd.Execute()
	h.Ok(t, err)

	output := buf.String()
	h.Assert(t, output != "", "Expected help subcommand for root, but wasn't found")
}
