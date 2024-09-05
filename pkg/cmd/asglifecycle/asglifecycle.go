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
	"errors"
	"fmt"
	"log"
	"strings"

	cmdutil "github.com/aws/amazon-ec2-metadata-mock/pkg/cmd/cmdutil"
	cfg "github.com/aws/amazon-ec2-metadata-mock/pkg/config"
	se "github.com/aws/amazon-ec2-metadata-mock/pkg/mock/asglifecycle"

	"github.com/spf13/cobra"
)

const (
	cfgPrefix = "asglifecycle."
)

var (
	c cfg.Config

	// Command represents the CLI command
	Command *cobra.Command

	// defaults
	defaultCfg = map[string]interface{}{}
)

func init() {
	cobra.OnInitialize(initConfig)
	Command = newCmd()
}

func initConfig() {
	cfg.LoadConfigFromDefaults(defaultCfg)
}

func newCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "asglifecycle [--code CODE] [--state STATE] [--not-after] [--not-before-deadline]",
		Aliases: []string{"asglifecycle", "autoscaling", "asg"},
		PreRunE: preRun,
		Example: fmt.Sprintf("  %s asglifecycle -h \tasglifecycle help \n  %s asglifecycle -t target-lifecycle-state\t\tmocks asg lifecycle target lifecycle states", cmdutil.BinName, cmdutil.BinName),
		Run:     run,
		Short:   "Mock EC2 ASG Lifecycle target-lifecycle-state",
		Long:    "Mock EC2 ASG Lifecycle target-lifecycle-state",
	}

	// bind local flags to config
	cfg.BindFlagSetWithKeyPrefix(cmd.Flags(), cfgPrefix)
	return cmd
}

// SetConfig sets the local config
func SetConfig(config cfg.Config) {
	c = config
}

func preRun(cmd *cobra.Command, args []string) error {
	if cfgErrors := ValidateLocalConfig(); cfgErrors != nil {
		return errors.New(strings.Join(cfgErrors, ""))
	}
	return nil
}

// ValidateLocalConfig validates all local config and returns a slice of error messages
func ValidateLocalConfig() []string {
	// no-op
	return nil
}

func run(cmd *cobra.Command, args []string) {
	log.Printf("Initiating %s for EC2 ASG Lifecycle on port %s\n", cmdutil.BinName, c.Server.Port)
	cmdutil.PrintFlags(cmd.Flags())
	cmdutil.RegisterHandlers(cmd, c)
	se.Mock(c)
}
