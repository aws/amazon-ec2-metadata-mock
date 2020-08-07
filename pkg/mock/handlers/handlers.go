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
	supportedCategories = []string{"dynamic", "meta-data"}

	// trimmedRoutes represents the list of routes served by the http server without "latest/meta-data/" prefix
	trimmedRoutes        []string
	trimmedRoutesDynamic []string
)

// CatchAllHandler returns subpath listings, if available; 404 status code otherwise
func CatchAllHandler(res http.ResponseWriter, req *http.Request) {
	log.Println("Received request to CatchAllHandler: ", req.URL.Path)
	// CatchAllHandler may be invoked before ListRoutesHandler
	if len(trimmedRoutes) == 0 {
		formatRoutes()
	}

	var routes, results []string

	// Clean request path and determine which route list to search
	trimmedRoute := req.URL.Path
	log.Println("removing suffix slash: ", trimmedRoute)
	if strings.HasPrefix(trimmedRoute, static.ServicePath) {
		trimmedRoute = strings.TrimPrefix(trimmedRoute, static.ServicePath+"/")
		log.Println("static prefix detected..trimming: ", trimmedRoute)
		routes = trimmedRoutes
	} else if strings.HasPrefix(trimmedRoute, dynamic.ServicePath) {
		trimmedRoute = strings.TrimPrefix(trimmedRoute, dynamic.ServicePath+"/")
		log.Println("dynamic prefix detected..trimming: ", trimmedRoute)
		routes = trimmedRoutesDynamic
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
	appendToResults := true
	for _, route := range routes {
		if strings.Contains(route, trimmedRoute) {
			// ex: iam/security-credentials contains iam
			route = strings.TrimPrefix(route, trimmedRoute+"/")
			// route is now security-credentials
			for i, existingRoute := range results {
				if strings.Contains(route, existingRoute) {
					// security-credentials/baskinc-role contains security-credentials
					results[i] = existingRoute + "/"
					// display as security-credentials/
					appendToResults = false
					// do not add security-credentials/baskinc-role to results
					break
				}
			}
			if appendToResults {
				log.Printf("adding route: %s to results\n", route)
				results = append(results, route)
			}
			appendToResults = true
		}
	}

	if len(results) <= 0 {
		server.ReturnNotFoundResponse(res)
		return
	}

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
	case static.ServicePath:
		server.FormatAndReturnTextResponse(res, strings.Join(trimmedRoutes, "\n")+"\n")
	case dynamic.ServicePath:
		server.FormatAndReturnTextResponse(res, strings.Join(trimmedRoutesDynamic, "\n")+"\n")
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
			// Omit /latest/dynamic
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
