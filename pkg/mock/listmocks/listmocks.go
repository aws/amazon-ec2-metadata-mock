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

package listmocks

import (
	"net/http"
	"sort"
	"strings"

	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/dynamic"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/static"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/server"
)

const (
	// shortestRouteLength represents smallest metadata path without prefix + "/" ex: "mac/"
	shortestRouteLength = 4
)

var (
	// trimmedRoutes represents the list of routes served by the http server without "latest/meta-data/" prefix
	trimmedRoutes        []string
	trimmedRoutesDynamic []string
)

// Handler handles http requests
func Handler(res http.ResponseWriter, req *http.Request) {
	// Routes are not available until runtime; only want to do this ONCE
	if len(trimmedRoutes) == 0 {
		formatRoutes()
	}

	// return 404 for unsupported paths; this is needed due to DefaultServeMux path-pattern matching
	if req.URL.Path == static.ServicePath || req.URL.Path == static.ServicePath2 {
		server.FormatAndReturnTextResponse(res, strings.Join(trimmedRoutes, "\n")+"\n")
	} else if req.URL.Path == dynamic.ServicePath || req.URL.Path == dynamic.ServicePath2 {
		server.FormatAndReturnTextResponse(res, strings.Join(trimmedRoutesDynamic, "\n")+"\n")
	} else {
		server.ReturnNotFoundResponse(res)
	}
	return
}

func formatRoutes() {
	var trimmedRoute string
	for _, route := range server.Routes {

		if strings.HasPrefix(route, dynamic.ServicePath) {
			// Omit /latest/meta-data
			trimmedRoute = strings.TrimPrefix(route, dynamic.ServicePath)
			// Omit empty paths and "/"
			if len(trimmedRoute) >= shortestRouteLength {
				trimmedRoute = strings.TrimPrefix(trimmedRoute, "/")
				trimmedRoutesDynamic = append(trimmedRoutesDynamic, trimmedRoute)
			}
		} else {
			// Omit /latest/meta-data
			trimmedRoute = strings.TrimPrefix(route, static.ServicePath)
			// Omit empty paths and "/"
			if len(trimmedRoute) >= shortestRouteLength {
				trimmedRoute = strings.TrimPrefix(trimmedRoute, "/")
				trimmedRoutes = append(trimmedRoutes, trimmedRoute)
			}
		}
	}
	sort.Sort(sort.StringSlice(trimmedRoutes))
	sort.Sort(sort.StringSlice(trimmedRoutesDynamic))
}
