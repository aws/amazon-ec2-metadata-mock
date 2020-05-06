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

package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/amazon-ec2-metadata-mock/pkg/config/defaults"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultCfgFileName = "aemm-config"
	defaultCfgFileExt  = "json"
	finalCfgDir        = ".amazon-ec2-metadata-mock"
	finalCfgFile       = ".aemm-config-used.json"
	cfgFileEnvKey      = "AEMM_CONFIG_FILE"
)

var (
	savedConfigFilePath         string
	home                        = getHomeDir()
	isCfgWriteSuccessMsgEmitted = false // required to emit the success message once
)

func getHomeDir() string {
	var err error
	var h string
	if h, err = homedir.Dir(); err != nil {
		log.Printf("Warning: Failed to find home directory due to error: %s\n", err)
	}
	return h
}

// WriteConfigToFile attempts to save the configuration to a local file
func WriteConfigToFile() {
	var dir string
	if home != "" {
		dir = home + "/" + finalCfgDir // save config in home dir
	} else {
		dir = "./" + finalCfgDir // save config in working dir
	}

	// on failure, print error message(s) and carry on
	errMsg := "Warning: Failed to save the final configuration to local file"
	if err := createDir(dir); err != nil {
		log.Println(errMsg, "-", err)
		return
	}

	savedConfigFilePath = dir + "/" + finalCfgFile
	if err := viper.WriteConfigAs(savedConfigFilePath); err != nil {
		log.Printf(errMsg, " %s: %s\n", savedConfigFilePath, err)
	} else {
		fmt.Println("Successfully saved final configuration to local file ", savedConfigFilePath) // the file will be overwritten, if it exists
	}
}

// creates the directory used to save final configuration
func createDir(dir string) error {
	f, err := os.Stat(dir)

	// create dir if not exists
	if os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			return fmt.Errorf("Failed to create directory for final configuration at %s: %s", dir, err)
		}
	} else {
		// make sure the dir exists
		if !f.IsDir() {
			return fmt.Errorf("The destination '%s' for saving the configuration already exists, but is not a directory", dir)
		}
	}
	return nil
}

// LoadConfigForRoot is initiated by the root command. It loads config from various input sources
func LoadConfigForRoot(configFileFlagName string, cmdDefaults map[string]interface{}) {

	// Set config file following Viper's precedence, in order to allow Viper to apply correct precedence for values in config file
	// 1) config file name in flag
	// 2) config file name in env variable
	// 3) default config file in user's HOME dir
	configFile := viper.GetString(configFileFlagName)
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else if val, ok := os.LookupEnv(cfgFileEnvKey); ok {
		viper.SetConfigFile(val)
	} else if home != "" {
		// Search for config file in home directory by name, without including extension
		viper.SetConfigName(defaultCfgFileName)
		viper.AddConfigPath(home)
	}

	// set up env
	viper.AutomaticEnv()
	viper.SetEnvPrefix("aemm") // AEMM stands for amazon-ec2-metadata-mock, set prefix to avoid name clashes
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	// set up defaults
	LoadConfigFromDefaults(cmdDefaults)
	SetMetadataDefaults(defaults.GetDefaultValues())
	SetServerCfgDefaults()

	// read in config using viper
	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Println("Warning: ", err)
		default:
			log.Printf("Error while attempting to read config from %s: %s\n", viper.ConfigFileUsed(), err)
		}
	} else {
		fmt.Println("Using configuration from file: ", viper.ConfigFileUsed())
	}

	// if config overrides update placeholder values, then paths with placeholders need to be updated as well
	overrideMetadataPathsWithPlaceholders()
}

// LoadConfigFromDefaults loads the given defaults into the config
func LoadConfigFromDefaults(cmdDefaults map[string]interface{}) {
	for key, value := range cmdDefaults {
		viper.SetDefault(key, value)
	}
}

// GetDefaultCfgFileName returns the default CLI config file name with path
func GetDefaultCfgFileName() string {
	return "$HOME/" + defaultCfgFileName + "." + defaultCfgFileExt
}

// GetSavedCfgFileName returns the default CLI config file name with path
func GetSavedCfgFileName() string {
	return finalCfgDir + "/" + finalCfgFile
}

// BindFlagSet binds the given flag set to the config, using each flag's long name as the config key.
func BindFlagSet(flagSet *pflag.FlagSet) {
	if err := viper.BindPFlags(flagSet); err != nil {
		panic(fmt.Errorf("Error binding CLI flags %#v: %s", flagSet, err.Error()))
	}
}

// BindTopLevelGFlags binds the global flags that need to be at top level of the config with Viper.
func BindTopLevelGFlags(flags []*pflag.Flag) {
	for _, f := range flags {
		if err := viper.BindPFlag(f.Name, f); err != nil {
			panic(fmt.Errorf("Error binding CLI flag %#v: %s", f, err.Error()))
		}
	}
}

// BindFlagSetWithKeyPrefix binds the given flag set to the config, using prefix + flag's long name as the config key.
func BindFlagSetWithKeyPrefix(flagSet *pflag.FlagSet, keyPrefix string) {
	flagSet.VisitAll(func(flag *pflag.Flag) {
		if err := viper.BindPFlag(keyPrefix+flag.Name, flag); err != nil {
			panic(fmt.Errorf("Error binding CLI flag %s: %s", flag.Name, err.Error()))
		}
	})
}

func bindFlagWithKeyPrefix(flag *pflag.Flag, keyPrefix string) {
	if err := viper.BindPFlag(keyPrefix+flag.Name, flag); err != nil {
		panic(fmt.Errorf("Error binding CLI flag %s: %s", flag.Name, err.Error()))
	}
}

func overrideMetadataPathsWithPlaceholders() {
	valueToPlaceholderPathsKeyMap := GetMetadataValueToPlaceholderPathsKeyMap()
	_, valueDefaults := GetMetadataDefaults()
	// override placeholders in paths, precedence is applied in the following order:
	// (1) placeholder value in metadata path itself is overridden with a non-default value
	// (2) placeholder value represented as metadata value is overridden
	// (3) default value for the placeholder
	for mdValueKey, listOfPaths := range valueToPlaceholderPathsKeyMap {
		defaultValue := valueDefaults[mdValueKey].(string)
		overriddenValue := viper.Get(mdValueKey).(string)
		// apply (1) and (2) from above
		if defaultValue != overriddenValue {
			for _, pathKey := range listOfPaths {
				// get current path value
				switch currentPath := viper.Get(pathKey).(type) {
				case string:
					// find default value and replace with overridden value
					updatedPath := strings.Replace(currentPath, defaultValue, overriddenValue, -1)
					// set new path value
					log.Printf("Updating path from %s to %s\n", currentPath, updatedPath)
					viper.Set(pathKey, updatedPath)
				default:
					log.Printf("Something went wrong trying to update path. Expected type 'string' for %v\n", currentPath)
				}
			}
		}
	}
}
