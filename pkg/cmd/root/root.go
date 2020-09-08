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
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/aws/amazon-ec2-metadata-mock/pkg/cmd/cmdutil"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/cmd/events"
	gf "github.com/aws/amazon-ec2-metadata-mock/pkg/cmd/root/globalflags"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/cmd/spot"
	cfg "github.com/aws/amazon-ec2-metadata-mock/pkg/config"
	r "github.com/aws/amazon-ec2-metadata-mock/pkg/mock/root"
)

var (
	c       cfg.Config
	command *cobra.Command
	version = "dev"

	// defaults
	cfgMdPrefix = cfg.GetCfgMdValPrefix()
	cfgDnPrefix = cfg.GetCfgDnValPrefix()
	defaultCfg  = map[string]interface{}{
		gf.ConfigFileFlag:       cfg.GetDefaultCfgFileName(),
		gf.MockDelayInSecFlag:   0,
		gf.MockTriggerTimeFlag:  "",
		gf.MockIPCountFlag:      2,
		gf.SaveConfigToFileFlag: false,
		gf.Imdsv2Flag:           false,
	}
)

func init() {
	cobra.OnInitialize(initializeConfig)
	command = NewCmd()
}

func initializeConfig() {
	cfg.LoadConfigForRoot(gf.ConfigFileFlag, defaultCfg)
}

// NewCmd returns a new root command after setting it up
func NewCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:               cmdutil.BinName + " <command> [arguments]",
		SuggestFor:        []string{"mock", "ec2-mock", "ec2-metadata-mock"},
		Version:           version,
		Example:           fmt.Sprintf("  %s --mock-delay-sec 10\tmocks all metadata paths\n  %s spot --action terminate\tmocks spot ITN only", cmdutil.BinName, cmdutil.BinName),
		PersistentPreRunE: setupAndSaveConfig, // persistentPreRun runs before PreRun
		PreRunE:           preRun,
		Run:               run,
		Short:             "Tool to mock Amazon EC2 instance metadata",
		Long:              cmdutil.BinName + " is a tool to mock Amazon EC2 instance metadata.",
	}
	cmd.SetVersionTemplate(`{{.Version}}`)

	// global flags
	cmd.PersistentFlags().StringP(gf.HostNameFlag, "n", "", "the HTTP hostname for the mock url (default: 0.0.0.0)")
	cmd.PersistentFlags().StringP(gf.PortFlag, "p", "", "the HTTP port where the mock runs (default: 1338)")
	cmd.PersistentFlags().StringP(gf.ConfigFileFlag, "c", "", "config file for cli input parameters in json format (default: "+cfg.GetDefaultCfgFileName()+")")
	cmd.PersistentFlags().BoolP(gf.SaveConfigToFileFlag, "s", false, "whether to save processed config from all input sources in "+cfg.GetSavedCfgFileName()+" in $HOME or working dir, if homedir is not found (default: false)")
	cmd.PersistentFlags().Int64P(gf.MockDelayInSecFlag, "d", 0, "mock delay in seconds, relative to the application start time (default: 0 seconds)")
	cmd.PersistentFlags().String(gf.MockTriggerTimeFlag, "", "mock trigger time in RFC3339 format. This takes priority over "+gf.MockDelayInSecFlag+" (default: none)")
	cmd.PersistentFlags().Int64P(gf.MockIPCountFlag, "x", 2, "number of IPs in a cluster that can receive a Spot Interrupt Notice and/or Scheduled Event")
	cmd.PersistentFlags().BoolP(gf.Imdsv2Flag, "I", false, "whether to enable IMDSv2 only, requiring a session token when submitting requests (default: false, meaning both IMDS v1 and v2 are enabled)")

	// add subcommands
	cmd.AddCommand(spot.Command, events.Command)

	// bind all non-metadata flags at top level
	var topLevelGFlags []*pflag.Flag
	for _, n := range gf.GetTopLevelFlags() {
		topLevelGFlags = append(topLevelGFlags, cmd.PersistentFlags().Lookup(n))
	}
	cfg.BindTopLevelGFlags(topLevelGFlags)

	// bind second level flags
	cfg.BindServerCfg(cmd.PersistentFlags().Lookup(gf.HostNameFlag))
	cfg.BindServerCfg(cmd.PersistentFlags().Lookup(gf.PortFlag))

	return cmd
}

func setupAndSaveConfig(cmd *cobra.Command, args []string) error {
	if err := injectViperConfig(); err != nil {
		return err
	}
	saveConfigToFile()

	return nil
}

// injectViperConfig after the cobra initializations and before prerun
func injectViperConfig() error {
	var aemmConfig = cfg.Config{}
	err := viper.Unmarshal(&aemmConfig)

	if err != nil {
		return fmt.Errorf("Fatal error while attempting to load viper config: %s", err)
	}

	setConfig(aemmConfig)
	return nil
}

// saveConfigToFile saves the config used by the tool to a local file, errors are logged as warnings
func saveConfigToFile() {
	if saveCfg := viper.GetBool(gf.SaveConfigToFileFlag); saveCfg {
		cfg.WriteConfigToFile()
	}
}

func setConfig(config cfg.Config) {
	c = config
	spot.SetConfig(config)
	events.SetConfig(config)
}

func preRun(cmd *cobra.Command, args []string) error {
	if cfgErrors := validateConfig(); cfgErrors != nil {
		return errors.New(strings.Join(cfgErrors, ""))
	}
	return nil
}

func validateConfig() []string {
	var errStrings []string

	// validate subcommands' config
	errStrings = append(errStrings, spot.ValidateLocalConfig()...)
	errStrings = append(errStrings, events.ValidateLocalConfig()...)

	if c.MockTriggerTime != "" {
		if err := cmdutil.ValidateRFC3339TimeFormat(gf.MockTriggerTimeFlag, c.MockTriggerTime); err != nil {
			errStrings = append(errStrings, err.Error())
		}
	}

	return errStrings
}

func run(cmd *cobra.Command, args []string) {
	log.Printf("Initiating %s for all mocks on port %s\n", cmdutil.BinName, c.Server.Port)
	cmdutil.PrintFlags(cmd.Flags())
	cmdutil.RegisterHandlers(cmd, c)
	r.Mock(c)
}
