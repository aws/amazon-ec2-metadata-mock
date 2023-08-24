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

// Package autoscaling provides the 'autoscaling' command.
package autoscaling

import (
	"github.com/aws/amazon-ec2-metadata-mock/pkg/cmd/autoscaling/internal/targetlifecyclestate"
	cfg "github.com/aws/amazon-ec2-metadata-mock/pkg/config"

	"github.com/spf13/cobra"
)

var Command *cobra.Command

func init() {
	Command = newCmd()
}

func newCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "autoscaling",
		Short: "Mock autoscaling information",
		Run:   func(cmd *cobra.Command, _ []string) { cmd.Help() },
	}

	cmd.AddCommand(targetlifecyclestate.Command)
	return cmd
}

// SetConfig sets the config for the command and subcommands.
func SetConfig(c cfg.Config) {
	targetlifecyclestate.SetConfig(c)
}
