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

package scheduledevents

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	cmdutil "github.com/aws/amazon-ec2-metadata-mock/pkg/cmd/cmdutil"
	cfg "github.com/aws/amazon-ec2-metadata-mock/pkg/config"
	e "github.com/aws/amazon-ec2-metadata-mock/pkg/error"
	se "github.com/aws/amazon-ec2-metadata-mock/pkg/mock/scheduledevents"

	"github.com/spf13/cobra"
)

const (
	cfgPrefix = "scheduled-events."

	// local flags
	eventCodeFlagName         = "code"
	eventStateFlagName        = "state"
	notBeforeFlagName         = "not-before"
	notAfterFlagName          = "not-after"
	notBeforeDeadlineFlagName = "not-before-deadline"

	// event codes
	instanceReboot     = "instance-reboot"
	systemReboot       = "system-reboot"
	systemMaintenance  = "system-maintenance"
	instanceRetirement = "instance-retirement"
	instanceStop       = "instance-stop"

	// event states
	active    = "active"
	completed = "completed"
	canceled  = "canceled"

	// default date diffs (in days) in metadata
	notAfterDiff          = 7
	notBeforeDeadlineDiff = 9
)

var (
	c cfg.Config

	// Command represents the CLI command
	Command *cobra.Command

	// constraints
	validEventCodes  = []string{instanceReboot, systemReboot, systemMaintenance, instanceRetirement, instanceStop}
	validEventStates = []string{active, completed, canceled}
	constraints      = []string{
		"event-code can be one of the following: " + strings.Join(validEventCodes, ","),
		"state can be one of the following: " + strings.Join(validEventStates, ","),
	}

	// defaults
	defaultCfg = map[string]interface{}{
		cfgPrefix + eventCodeFlagName:         systemReboot,
		cfgPrefix + eventStateFlagName:        active,
		cfgPrefix + notBeforeFlagName:         time.Now().Format(time.RFC3339),
		cfgPrefix + notAfterFlagName:          time.Now().Add(time.Hour * 24 * notAfterDiff).Format(time.RFC3339),
		cfgPrefix + notBeforeDeadlineFlagName: time.Now().Add(time.Hour * 24 * notBeforeDeadlineDiff).Format(time.RFC3339),
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
		Use:     "scheduledevents [--code CODE] [--state STATE] [--not-after] [--not-before-deadline]",
		Aliases: []string{"se"},
		PreRunE: preRun,
		Example: fmt.Sprintf("  %s scheduledevents -h \tscheduledevents help \n  %s scheduledevents -o instance-stop --state active -d\t\tmocks an active and upcoming scheduled event for instance stop with a deadline for the event start time", cmdutil.BinName, cmdutil.BinName),
		Run:     run,
		Short:   "Mock EC2 Scheduled Events",
		Long:    "Mock EC2 Scheduled Events",
	}

	// local flags
	cmd.Flags().StringP(eventCodeFlagName, "o", "", "event code in the scheduled event (default: system-reboot)\nevent-code can be one of the following: "+strings.Join(validEventCodes, ","))
	cmd.Flags().StringP(eventStateFlagName, "t", "", "state of the scheduled event (default: active)\nstate can be one of the following: "+strings.Join(validEventStates, ","))
	cmd.Flags().StringP(notBeforeFlagName, "b", "", "the earliest start time for the scheduled event in RFC3339 format E.g. 2020-01-07T01:03:47Z (default: application start time in UTC)")
	cmd.Flags().StringP(notAfterFlagName, "a", "", "the latest end time for the scheduled event in RFC3339 format E.g. 2020-01-07T01:03:47Z default: application start time + 7 days in UTC))")
	cmd.Flags().StringP(notBeforeDeadlineFlagName, "l", "", "the deadline for starting the event in RFC3339 format E.g. 2020-01-07T01:03:47Z (default: application start time + 9 days in UTC)")

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
	c := c.SchEventsConfig

	// validate event code
	if ok := cmdutil.Contains(validEventCodes, c.EventCode); !ok {
		errStrings = append(errStrings, e.FlagValidationError{
			FlagName:     eventCodeFlagName,
			Allowed:      strings.Join(validEventCodes, ","),
			InvalidValue: c.EventCode}.Error(),
		)
	}

	// validate event status
	if ok := cmdutil.Contains(validEventStates, c.EventState); !ok {
		errStrings = append(errStrings, e.FlagValidationError{
			FlagName:     eventStateFlagName,
			Allowed:      strings.Join(validEventStates, ","),
			InvalidValue: c.EventState}.Error(),
		)
	}

	// validate time flags
	if err := cmdutil.ValidateRFC3339TimeFormat(notBeforeFlagName, c.NotBefore); err != nil {
		errStrings = append(errStrings, err.Error())
	}
	if err := cmdutil.ValidateRFC3339TimeFormat(notAfterFlagName, c.NotAfter); err != nil {
		errStrings = append(errStrings, err.Error())
	}
	if err := cmdutil.ValidateRFC3339TimeFormat(notBeforeDeadlineFlagName, c.NotBeforeDeadline); err != nil {
		errStrings = append(errStrings, err.Error())
	}
	return errStrings
}

func run(cmd *cobra.Command, args []string) {
	log.Printf("Initiating %s for EC2 Scheduled Events on port %s\n", cmdutil.BinName, c.Server.Port)
	cmdutil.PrintFlags(cmd.Flags())
	cmdutil.RegisterHandlers(cmd, c)
	se.Mock(c)
}
