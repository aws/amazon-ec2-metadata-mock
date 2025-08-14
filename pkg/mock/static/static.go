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

package static

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
	jsonTextResponse = map[string]bool{"/latest/meta-data/elastic-inference/associations/eia-bfa21c7904f64a82a21b9f4540169ce1": true}

	// ServicePath defines the static service path
	ServicePath = "/latest/meta-data"
)

// Handler processes http requests
func Handler(res http.ResponseWriter, req *http.Request) {
	log.Println("Received request to mock static metadata:", req.URL.Path)

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

// RegisterHandlers registers handlers for static paths
func RegisterHandlers(config cfg.Config) {
	// Register all tags/instance/<TAGNAME> from the tags- map in the config struct
	if config.Metadata.Values.TagsInstance != nil {
		for tag, value := range config.Metadata.Values.TagsInstance {
			tagPath := "/latest/meta-data/tags/instance/" + tag
			supportedPaths[tagPath] = value
			if config.Imdsv2Required {
				server.HandleFunc(tagPath, imdsv2.ValidateToken(Handler))
			} else {
				server.HandleFunc(tagPath, Handler)
			}
		}
	}
	server.HandleFunc("/latest/api/token", imdsv2.GenerateToken)

	pathValues := reflect.ValueOf(config.Metadata.Paths)
	mdValues := reflect.ValueOf(config.Metadata.Values)

	// Iterate over fields in config.Metadata.Paths to
	// determine intersections with config.Metadata.Values.
	// Intersections represent which paths and values to bind.
	for i := 0; i < pathValues.NumField(); i++ {
		pathFieldName := pathValues.Type().Field(i).Name
		mdValueFieldName := mdValues.FieldByName(pathFieldName)
		if mdValueFieldName.IsValid() {
			path := pathValues.Field(i).Interface().(string)
			value := mdValueFieldName.Interface()
			if path != "" && value != nil {
				supportedPaths[path] = value
				if config.Imdsv2Required {
					server.HandleFunc(path, imdsv2.ValidateToken(Handler))
				} else {
					server.HandleFunc(path, Handler)
				}
			} else {
				log.Printf("There was an issue registering path %v with mdValue: %v", path, value)
			}
		}
	}

	if config.Metadata.Values.TagsInstance != nil {
		for tag, value := range config.Metadata.Values.TagsInstance {
			tagPath := "/latest/meta-data/tags/instance/" + tag
			supportedPaths[tagPath] = value
			if config.Imdsv2Required {
				server.HandleFunc(tagPath, imdsv2.ValidateToken(Handler))
			} else {
				server.HandleFunc(tagPath, Handler)
			}
		}
	}
}
