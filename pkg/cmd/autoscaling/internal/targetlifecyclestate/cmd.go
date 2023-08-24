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
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/amazon-ec2-metadata-mock/pkg/cmd/cmdutil"
	cfg "github.com/aws/amazon-ec2-metadata-mock/pkg/config"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/autoscaling/targetlifecyclestate"

	"github.com/spf13/cobra"
)

var (
	Command *cobra.Command

	config cfg.Config
)

func init() {
	Command = newCmd()
}

func newCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "target-lifecycle-state -h | --help | [STATE [TIME|DELAY STATE]...]",
		Example: fmt.Sprintf("  target-lifecycle-state %[1]s 2m %[2]s",
			targetlifecyclestate.InService, targetlifecyclestate.Terminated),
		Args:  validateArgs,
		RunE:  run,
		Short: "Mock autoscaling lifecycle state transitions",
		Long: "Mock autoscaling lifecycle state transitions.\n\n" +
			"Target states may be changed at a specified time or after a delay, e.g. \"InService 2m Terminated\" will initially " +
			"return the InService target lifecycle state then after 2 minutes will return Terminated.\n" +
			"Valid states are " + strings.Join(targetlifecyclestate.Names(), ", ") + "\n" +
			"Times must be in RFC 3339 format, e.g. \"2018-09-01T12:00:00Z\".\n" +
			"Delays are specified as one or more numbers and suffixes, e.g. \"2m30s\". " +
			"Valid suffixes are \"h\", \"m\", \"s\", \"ms\", \"us\", \"ns\".",
	}

	return cmd
}

// SetConfig sets the config for the command.
func SetConfig(c cfg.Config) {
	config = c
}

func validateArgs(_ *cobra.Command, args []string) error {
	_, err := parse(args)
	return err
}

func run(cmd *cobra.Command, args []string) error {
	log.Printf(
		"Initiating %s for autoscaling target lifecycle state on port %s\n",
		cmdutil.BinName,
		config.Server.Port,
	)

	cmdutil.PrintFlags(cmd.Flags())
	cmdutil.RegisterHandlers(cmd, config)

	if len(args) == 0 {
		args = config.AutoscalingConfig.TargetLifecycleState
	}
	s, err := parse(args)
	if err != nil {
		return fmt.Errorf("parsing state transitions: %w", err)
	}

	setTargetState(s.Initial)()
	for _, t := range s.Transitions {
		go schedule(cmd.Context(), t.Time, setTargetState(t.State))
	}

	targetlifecyclestate.Mock(config)
	return nil
}

func setTargetState(s targetlifecyclestate.TargetState) func() {
	return func() {
		log.Printf("Setting autoscaling target lifecycle state to %s\n", s)
		targetlifecyclestate.Set(s)
	}
}

func schedule(ctx context.Context, t time.Time, f func()) {
	select {
	case <-time.After(t.Sub(time.Now())):
		f()
	case <-ctx.Done():
		return
	}
}
