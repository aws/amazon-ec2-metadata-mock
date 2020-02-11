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

package spotitn

import (
	"errors"
	"log"
	"strings"

	cmdutil "github.com/aws/amazon-ec2-metadata-mock/pkg/cmd/cmdutil"
	cfg "github.com/aws/amazon-ec2-metadata-mock/pkg/config"
	e "github.com/aws/amazon-ec2-metadata-mock/pkg/error"
	s "github.com/aws/amazon-ec2-metadata-mock/pkg/mock/spotitn"

	"github.com/spf13/cobra"
)

const (
	cfgPrefix = "spot-itn."

	// local flags
	instanceActionFlagName  = "instance-action"
	terminationTimeFlagName = "termination-time"

	// instance actions
	terminate = "terminate"
	hibernate = "hibernate"
	stop      = "stop"
)

var (
	c cfg.Config

	// Command represents the CLI command
	Command *cobra.Command

	// constraints
	validInstanceActions = []string{terminate, hibernate, stop}

	// defaults
	defaultCfg = map[string]interface{}{
		cfgPrefix + instanceActionFlagName: terminate,
	}
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
		Use:     "spotitn [--instance-action ACTION]",
		Aliases: []string{"spot", "spot-itn", "spotItn"},
		PreRunE: preRun,
		Example: "  amazon-ec2-metadata-mock spotitn -h \tspotitn help \n  amazon-ec2-metadata-mock spotitn -d 5 --instance-action terminate\t\tmocks spot interruption only",
		Run:     run,
		Short:   "Mock EC2 Spot interruption notice",
		Long:    "Mock EC2 Spot interruption notice",
	}

	// local flags
	cmd.Flags().StringP(instanceActionFlagName, "a", "", "instance action in the spot interruption notice (default: terminate)\ninstance-action can be one of the following: "+strings.Join(validInstanceActions, ","))
	cmd.Flags().StringP(terminationTimeFlagName, "t", "", "termination time specifies the approximate time when the spot instance will receive the shutdown signal in RFC3339 format to execute instance action E.g. 2020-01-07T01:03:47Z (default: request time + 2 minutes in UTC)")

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
	var errStrings []string
	c := c.SpotItnConfig

	// validate instance-action
	if ok := cmdutil.Contains(validInstanceActions, c.InstanceAction); !ok {
		errStrings = append(errStrings, e.FlagValidationError{
			FlagName:     instanceActionFlagName,
			Allowed:      strings.Join(validInstanceActions, ","),
			InvalidValue: c.InstanceAction}.Error(),
		)
	}

	// validate termination-time, if override provided
	if c.TerminationTime != "" {
		if err := cmdutil.ValidateRFC3339TimeFormat(terminationTimeFlagName, c.TerminationTime); err != nil {
			errStrings = append(errStrings, err.Error())
		}
	}
	return errStrings
}

func run(cmd *cobra.Command, args []string) {
	log.Println("Initiating amazon-ec2-metadata-mock for EC2 Spot interruption notice on port", c.Server.Port)
	cmdutil.PrintFlags(cmd.Flags())
	cmdutil.RegisterHandlers(cmd, c)
	s.Mock(c)
}
