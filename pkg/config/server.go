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
	"github.com/spf13/pflag"
)

var (
	serverCfgPrefix   = "server."
	serverCfgDefaults = map[string]interface{}{
		serverCfgPrefix + "hostname": "localhost",
		serverCfgPrefix + "port":     "1338",
	}
)

// BindServerCfg binds a flag that represents a server config to configuration
func BindServerCfg(flag *pflag.Flag) {
	bindFlagWithKeyPrefix(flag, serverCfgPrefix)
}

// SetServerCfgDefaults sets config defaults for server config
func SetServerCfgDefaults() {
	LoadConfigFromDefaults(serverCfgDefaults)
}
