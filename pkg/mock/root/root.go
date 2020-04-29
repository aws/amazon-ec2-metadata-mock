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
	cfg "github.com/aws/amazon-ec2-metadata-mock/pkg/config"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/scheduledevents"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/spotitn"

	"github.com/aws/amazon-ec2-metadata-mock/pkg/server"
)

// Mock serves all subcommand handlers
func Mock(config cfg.Config) {

	// set configs for subcommands
	spotitn.SetConfig(config)
	scheduledevents.SetConfig(config)

	server.ListenAndServe(config.Server.HostName, config.Server.Port)
}