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
	// Blank import else compiler complains it's unused
	_ "embed"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/aws/amazon-ec2-metadata-mock/pkg/cmd/asglifecycle"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/cmd/cmdutil"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/cmd/events"
	gf "github.com/aws/amazon-ec2-metadata-mock/pkg/cmd/root/globalflags"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/cmd/spot"
	cfg "github.com/aws/amazon-ec2-metadata-mock/pkg/config"
	r "github.com/aws/amazon-ec2-metadata-mock/pkg/mock/root"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/server"
)

var (
	c       cfg.Config
	command *cobra.Command
	//go:embed version.txt
	version string

	// defaults
	cfgMdPrefix = cfg.GetCfgMdValPrefix()
	cfgDnPrefix = cfg.GetCfgDnValPrefix()
	defaultCfg  = map[string]interface{}{
		gf.ConfigFileFlag:                cfg.GetDefaultCfgFileName(),
		gf.MockDelayInSecFlag:            0,
		gf.MockTriggerTimeFlag:           "",
		gf.MockIPCountFlag:               2,
		gf.SaveConfigToFileFlag:          false,
		gf.Imdsv2Flag:                    false,
		gf.RebalanceDelayInSecFlag:       0,
		gf.RebalanceTriggerTimeFlag:      "",
		gf.ASGTerminationDelayInSecFlag:  0,
		gf.ASGTerminationTriggerTimeFlag: "",
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
	cmd.PersistentFlags().BoolP(gf.WatchConfigFileFlag, "s", false, "whether to watch the config file "+cfg.GetSavedCfgFileName()+" in $HOME or working dir, if homedir is not found (default: false)")
	cmd.PersistentFlags().Int64P(gf.MockDelayInSecFlag, "d", 0, "spot itn delay in seconds, relative to the application start time (default: 0 seconds)")
	cmd.PersistentFlags().String(gf.MockTriggerTimeFlag, "", "spot itn trigger time in RFC3339 format. This takes priority over "+gf.MockDelayInSecFlag+" (default: none)")
	cmd.PersistentFlags().Int64P(gf.MockIPCountFlag, "x", 2, "number of IPs in a cluster that can receive a Spot Interrupt Notice and/or Scheduled Event")
	cmd.PersistentFlags().BoolP(gf.Imdsv2Flag, "I", false, "whether to enable IMDSv2 only, requiring a session token when submitting requests (default: false, meaning both IMDS v1 and v2 are enabled)")
	cmd.PersistentFlags().Int64(gf.RebalanceDelayInSecFlag, 0, "rebalance rec delay in seconds, relative to the application start time (default: 0 seconds)")
	cmd.PersistentFlags().String(gf.RebalanceTriggerTimeFlag, "", "rebalance rec trigger time in RFC3339 format. This takes priority over "+gf.RebalanceDelayInSecFlag+" (default: none)")
	cmd.PersistentFlags().Int64P(gf.ASGTerminationDelayInSecFlag, "", 0, "asg termination delay in seconds, relative to the application start time (default: 0 seconds)")
	cmd.PersistentFlags().Int64P(gf.ASGTerminationTriggerTimeFlag, "", 0, "asg termination trigger time in RFC3339 format. This takes priority over "+gf.ASGTerminationDelayInSecFlag+" (default: none)")

	// add subcommands
	cmd.AddCommand(spot.Command, events.Command, asglifecycle.Command)

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

	if watchCfg := viper.GetBool(gf.WatchConfigFileFlag); watchCfg {
		viper.OnConfigChange(func(_ fsnotify.Event) {
			if err := injectViperConfig(); err != nil {
				log.Printf("Failed to reset config on config change: %v\n", err)
				return
			}
			server.Reset()
			cmdutil.RegisterHandlers(cmd, c)
		})
		viper.WatchConfig()
	}

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
	asglifecycle.SetConfig(config)
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
	errStrings = append(errStrings, asglifecycle.ValidateLocalConfig()...)

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
