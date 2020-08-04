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

package dynamic

import (
	"log"
	"net/http"
	"reflect"

	cfg "github.com/aws/amazon-ec2-metadata-mock/pkg/config"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/imdsv2"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/server"
)

var (
	supportedPaths   = make(map[string]interface{})
	response         interface{}
	jsonTextResponse = map[string]bool{}

	// ServicePath defines the static service path
	ServicePath = "/latest/dynamic"
)

// Handler processes http requests
func Handler(res http.ResponseWriter, req *http.Request) {
	log.Println("Received request to mock dynamic metadata:", req.URL.Path)

	if val, ok := supportedPaths[req.URL.Path]; ok {
		response = val
	} else {
		response = "Something went wrong with: " + req.URL.Path
	}

	switch response.(type) {
	// static metadata values are either string or JSON EXCEPT FOR elastic-inference associations
	case string:
		server.FormatAndReturnTextResponse(res, response.(string))
	default:
		if jsonTextResponse[req.URL.Path] {
			server.FormatAndReturnJSONTextResponse(res, response)
		} else {
			server.FormatAndReturnJSONResponse(res, response)
		}
	}
}

// RegisterHandlers registers handlers for dynamic paths
func RegisterHandlers(config cfg.Config) {
	pathValues := reflect.ValueOf(config.Dynamic.Paths)
	dyValues := reflect.ValueOf(config.Dynamic.Values)

	// Iterate over fields in config.Dynamic.Paths to
	// determine intersections with config.Dynamic.Values.
	// Intersections represent which paths and values to bind.
	for i := 0; i < pathValues.NumField(); i++ {
		pathFieldName := pathValues.Type().Field(i).Name
		dyValueFieldName := dyValues.FieldByName(pathFieldName)
		if dyValueFieldName.IsValid() {
			path := pathValues.Field(i).Interface().(string)
			value := dyValueFieldName.Interface()
			if path != "" && value != nil {
				// Ex: "/latest/dynamic/instance-identity/document"
				supportedPaths[path] = value
				if config.Imdsv2Required {
					server.HandleFunc(path, imdsv2.ValidateToken(Handler))
				} else {
					server.HandleFunc(path, Handler)
				}
			} else {
				log.Printf("There was an issue registering path %v with dyValue: %v", path, value)
			}
		}
	}
}
