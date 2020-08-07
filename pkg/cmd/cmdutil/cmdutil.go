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

package cmdutil

import (
	"fmt"
	"strings"
	"time"

	cfg "github.com/aws/amazon-ec2-metadata-mock/pkg/config"
	e "github.com/aws/amazon-ec2-metadata-mock/pkg/error"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/dynamic"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/events"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/handlers"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/imdsv2"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/spot"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/static"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/server"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// BinName is the name of this tool's binary
const BinName = "ec2-metadata-mock"

// handlerPair holds a tuple of a path and its associated handler
type handlerPair struct {
	path    string
	handler server.HandlerType
}

// Contains finds a string in the given array
func Contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// PrintFlags prints all flags of a command, if set
func PrintFlags(flags *pflag.FlagSet) {
	f := make(map[string]interface{})
	flags.Visit(func(flag *pflag.Flag) {
		f[flag.Name] = flag.Value
	})

	if len(f) > 0 {
		fmt.Println("\nFlags:")
		for key, value := range f {
			fmt.Printf("%s: %s\n", key, value)
		}
		fmt.Println()
	}
}

// ValidateRFC3339TimeFormat validates an input time matches RFC3339 format
func ValidateRFC3339TimeFormat(flagName string, input string) error {
	if _, err := time.Parse(time.RFC3339, input); err != nil {
		return e.FlagValidationError{
			FlagName:     flagName,
			Allowed:      "time in RFC3339 format, e.g. 2020-01-07T01:03:47Z",
			InvalidValue: input}
	}
	return nil
}

// RegisterHandlers binds paths to handlers for ALL commands
func RegisterHandlers(cmd *cobra.Command, config cfg.Config) {
	handlerPairsToRegister := getHandlerPairs(cmd, config)
	for _, handlerPair := range handlerPairsToRegister {
		if config.Imdsv2Required {
			server.HandleFunc(handlerPair.path, imdsv2.ValidateToken(handlerPair.handler))
		} else {
			server.HandleFunc(handlerPair.path, handlerPair.handler)
		}
	}

	static.RegisterHandlers(config)
	dynamic.RegisterHandlers(config)

	// paths without explicit handler bindings will fallback to CatchAllHandler
	server.HandleFuncPrefix("/", handlers.CatchAllHandler)
}

// getHandlerPairs returns a slice of {paths, handlers} to register
func getHandlerPairs(cmd *cobra.Command, config cfg.Config) []handlerPair {
	// always register these paths
	handlerPairs := []handlerPair{
		{path: "/", handler: handlers.ListRoutesHandler},
		{path: "/latest", handler: handlers.ListRoutesHandler},
		{path: static.ServicePath, handler: handlers.ListRoutesHandler},
		{path: dynamic.ServicePath, handler: handlers.ListRoutesHandler},
	}

	isSpot := strings.Contains(cmd.Name(), "spot")
	isEvents := strings.Contains(cmd.Name(), "events")

	subCommandHandlers := map[string][]handlerPair{
		"spot": {{path: config.Metadata.Paths.Spot, handler: spot.Handler},
			{path: config.Metadata.Paths.SpotTerminationTime, handler: spot.Handler}},
		"events": {{path: config.Metadata.Paths.Events, handler: events.Handler}},
	}

	if isSpot {
		handlerPairs = append(handlerPairs, subCommandHandlers["spot"]...)
	} else if isEvents {
		handlerPairs = append(handlerPairs, subCommandHandlers["events"]...)
	} else {
		// root registers all subcommands
		for k := range subCommandHandlers {
			handlerPairs = append(handlerPairs, subCommandHandlers[k]...)
		}
	}

	return handlerPairs
}
