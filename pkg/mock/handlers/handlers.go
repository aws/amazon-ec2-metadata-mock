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

package handlers

import (
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/dynamic"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/static"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/userdata"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/server"
)

const (
	// shortestRouteLength represents smallest metadata path without prefix + "/" ex: "mac/"
	shortestRouteLength = 4
	versionsPath        = "/"
	latestPath          = "/latest"
)

var (
	routeLookupTable = make(map[string][]string)

	supportedVersions   = []string{"latest"}
	supportedCategories = []string{"dynamic", "meta-data", "user-data"}

	// trimmedRoutes represents the list of routes served by the http server without "latest/meta-data/" prefix
	trimmedRoutes         []string
	trimmedRoutesDynamic  []string
	trimmedRoutesUserdata []string
)

// CatchAllHandler returns subpath listings, if available; 404 status code otherwise
func CatchAllHandler(res http.ResponseWriter, req *http.Request) {
	log.Println("Received request to CatchAllHandler: ", req.URL.Path)
	// CatchAllHandler may be invoked before ListRoutesHandler
	if len(trimmedRoutes) == 0 {
		formatRoutes()
	}

	var routes []string

	// Clean request path and determine which route list to search
	trimmedRoute := req.URL.Path
	if strings.HasPrefix(trimmedRoute, static.ServicePath) {
		trimmedRoute = strings.TrimPrefix(trimmedRoute, static.ServicePath+"/")
		log.Println("static prefix detected..trimming: ", trimmedRoute)
		routes = trimmedRoutes
	} else if strings.HasPrefix(trimmedRoute, dynamic.ServicePath) {
		trimmedRoute = strings.TrimPrefix(trimmedRoute, dynamic.ServicePath+"/")
		log.Println("dynamic prefix detected..trimming: ", trimmedRoute)
		routes = trimmedRoutesDynamic
	} else if strings.HasPrefix(trimmedRoute, userdata.ServicePath) {
		trimmedRoute = strings.TrimPrefix(trimmedRoute, userdata.ServicePath+"/")
		log.Println("userdata prefix detected..trimming: ", trimmedRoute)
		routes = trimmedRoutesUserdata
	} else {
		server.ReturnNotFoundResponse(res)
		return
	}

	if paths, ok := routeLookupTable[trimmedRoute]; ok {
		log.Printf("CatchAllHandler entry %s already in map: %v \n", trimmedRoute, routeLookupTable[trimmedRoute])
		server.FormatAndReturnTextResponse(res, strings.Join(paths, "\n"))
		return
	}

	/*
		The request /latest/meta-data/iam will populate results as [info, security-credentials/]
		Note: not every path for which iam is a prefix needs to be appended to results.
		ex: iam/security-credentials/baskinc-role should not be added because security-credentials/ already exists
		results are cached in routeLookupTable
	*/
	resultSet := map[string]bool{}
	for _, route := range routes {
		if strings.Contains(route, trimmedRoute) {
			// ex: iam/security-credentials contains iam
			route = strings.TrimPrefix(route, trimmedRoute+"/")
			resultSet[trimRoute(route)] = true
		}
	}

	if len(resultSet) <= 0 {
		server.ReturnNotFoundResponse(res)
		return
	}

	results := make([]string, 0, len(resultSet))
	for key := range resultSet {
		results = append(results, key)
	}
	sort.Strings(results)

	routeLookupTable[trimmedRoute] = results
	log.Printf("CatchAllHandler: adding  %s  and its routes: %v to the map\n", trimmedRoute, results)
	server.FormatAndReturnTextResponse(res, strings.Join(results, "\n"))
	return
}

// ListRoutesHandler returns the list of supported paths
func ListRoutesHandler(res http.ResponseWriter, req *http.Request) {
	log.Println("Received request to display paths: ", req.URL.Path)
	// Routes are not available until runtime; only want to do this ONCE
	if len(trimmedRoutes) == 0 {
		formatRoutes()
	}

	// these paths do not use routeLookupTable due to inconsistency of trailing "/" with IMDS
	switch req.URL.Path {
	case userdata.ServicePath:
		server.FormatAndReturnOctetResponse(res, strings.Join(trimmedRoutesUserdata, "\n")+"\n")
	case static.ServicePath:
		server.FormatAndReturnTextResponse(res, strings.Join(trimAndSortRoutes(trimmedRoutes), "\n")+"\n")
	case dynamic.ServicePath:
		server.FormatAndReturnTextResponse(res, strings.Join(trimAndSortRoutes(trimmedRoutesDynamic), "\n")+"\n")
	case latestPath:
		server.FormatAndReturnTextResponse(res, strings.Join(supportedCategories, "\n")+"\n")
	case versionsPath:
		server.FormatAndReturnTextResponse(res, strings.Join(supportedVersions, "\n")+"\n")
	default:
		server.ReturnNotFoundResponse(res)
	}
	return
}

func formatRoutes() {
	var trimmedRoute string
	for _, route := range server.Routes {
		if strings.HasPrefix(route, dynamic.ServicePath) {
			// Omit /latest/dynamic and /latest/user-data
			trimmedRoute = strings.TrimPrefix(route, dynamic.ServicePath)
			// Omit empty paths and "/"
			if len(trimmedRoute) >= shortestRouteLength {
				trimmedRoute = strings.TrimPrefix(trimmedRoute, "/")
				trimmedRoutesDynamic = append(trimmedRoutesDynamic, trimmedRoute)
			}
		} else if strings.HasPrefix(route, userdata.ServicePath) {
			// Omit /latest/dynamic and /latest/meta-data
			trimmedRoute = strings.TrimPrefix(route, userdata.ServicePath)
			// Omit empty paths and "/"
			if len(trimmedRoute) >= shortestRouteLength {
				trimmedRoute = strings.TrimPrefix(trimmedRoute, "/")
				trimmedRoutesUserdata = append(trimmedRoutesUserdata, trimmedRoute)
			}

		} else if strings.HasPrefix(route, static.ServicePath) {
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
	sort.Sort(sort.StringSlice(trimmedRoutesUserdata))
}

func trimRoute(route string) string {
	// Remove trailing path elements, e.g. "0e:49:61:0f:c3:11/device-number" => "0e:49:61:0f:c3:11"
	route, _, foundSlash := stringCut(route, "/")
	if !foundSlash {
		return route
	}
	return route + "/"
}

// stringCut is a backfill for `strings.Cut` in Go <1.18.
func stringCut(s, sep string) (before, after string, found bool) {
	if i := strings.Index(s, sep); i > -1 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}

func trimAndSortRoutes(notTrimmedRoutes []string) []string {
	trimmedRouteSet := map[string]bool{}
	for _, route := range notTrimmedRoutes {
		trimmedRouteSet[trimRoute(route)] = true
	}

	trimmedRoutes := make([]string, 0, len(trimmedRouteSet))
	for key := range trimmedRouteSet {
		trimmedRoutes = append(trimmedRoutes, key)
	}
	sort.Strings(trimmedRoutes)

	return trimmedRoutes
}
