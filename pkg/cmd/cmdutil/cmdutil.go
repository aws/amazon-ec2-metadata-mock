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
	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/imdsv2"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/listmocks"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/scheduledevents"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/spotitn"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/static"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/versions"
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

// GetKeyValueSlicesFromMap converts map to slices of keys and values
func GetKeyValueSlicesFromMap(m map[string]string) ([]string, []string) {
	keys := make([]string, len(m))
	values := make([]string, len(m))
	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}

	return keys, values
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

	// static handles its own registration
	static.RegisterHandlers(config)
}

// getHandlerPairs returns a slice of {paths, handlers} to register
func getHandlerPairs(cmd *cobra.Command, config cfg.Config) []handlerPair {
	// always register these paths
	handlerPairs := []handlerPair{
		{path: "/", handler: versions.Handler},
		{path: "/latest/meta-data", handler: listmocks.Handler},
		{path: "/latest/meta-data/", handler: listmocks.Handler},
	}

	isSpot := strings.Contains(cmd.Name(), "spotitn")
	isSchedEvents := strings.Contains(cmd.Name(), "scheduledevents")

	subCommandHandlers := map[string][]handlerPair{
		"spotitn": {{path: config.Metadata.Paths.SpotItn, handler: spotitn.Handler},
			{path: config.Metadata.Paths.SpotItnTerminationTime, handler: spotitn.Handler}},
		"scheduledevents": {{path: config.Metadata.Paths.ScheduledEvents, handler: scheduledevents.Handler}},
	}

	if isSpot {
		handlerPairs = append(handlerPairs, subCommandHandlers["spotitn"]...)
	} else if isSchedEvents {
		handlerPairs = append(handlerPairs, subCommandHandlers["scheduledevents"]...)
	} else {
		// root registers all subcommands
		for k := range subCommandHandlers {
			handlerPairs = append(handlerPairs, subCommandHandlers[k]...)
		}
	}

	return handlerPairs
}
