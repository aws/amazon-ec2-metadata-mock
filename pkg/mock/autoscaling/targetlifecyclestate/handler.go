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
	cfg "github.com/aws/amazon-ec2-metadata-mock/pkg/config"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/server"
)

var (
	instance mock

	// Handler writes the currently mocked target lifecycle state to the response.
	Handler = instance.handler
)

// Mock starts serving the mock endpoint. The hostname and port should be set in
// the given config.
func Mock(c cfg.Config) {
	server.ListenAndServe(c.Server.HostName, c.Server.Port)
}

// Set sets the current mock target lifecycle state value.
func Set(s TargetState) {
	instance = s
}
