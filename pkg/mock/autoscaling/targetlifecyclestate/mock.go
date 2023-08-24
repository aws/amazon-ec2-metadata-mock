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
	"net/http"

	"github.com/aws/amazon-ec2-metadata-mock/pkg/server"
)

// Path is the mocked URL endpoint.
const Path = "/latest/meta-data/autoscaling/target-lifecycle-state"

type mock = TargetState

func (m *mock) handler(res http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case Path:
		server.FormatAndReturnTextResponse(res, m.String())
	default:
		server.ReturnNotFoundResponse(res)
	}
}
