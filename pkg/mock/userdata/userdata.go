// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
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

package userdata

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"reflect"

	cfg "github.com/aws/amazon-ec2-metadata-mock/pkg/config"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/mock/imdsv2"
	"github.com/aws/amazon-ec2-metadata-mock/pkg/server"
)

var (
	supportedPaths = make(map[string]interface{})
	response       interface{}
	// ServicePath defines the userdata service path
	ServicePath = "/latest/user-data"
)

// Handler processes http requests
func Handler(res http.ResponseWriter, req *http.Request) {
	log.Println("Received request to mock userdata:", req.URL.Path)

	if val, ok := supportedPaths[req.URL.Path]; ok {

		response = val
	} else {
		response = "Something went wrong with: " + req.URL.Path
	}
	server.FormatAndReturnOctetResponse(res, response.(string))

}

// RegisterHandlers registers handlers for userdata paths
func RegisterHandlers(config cfg.Config) {
	pathValues := reflect.ValueOf(config.Userdata.Paths)
	udValues := reflect.ValueOf(config.Userdata.Values)
	// Iterate over fields in config.Userdata.Paths to
	// determine intersections with config.Userdata.Values.
	// Intersections represent which paths and values to bind.
	for i := 0; i < pathValues.NumField(); i++ {
		pathFieldName := pathValues.Type().Field(i).Name
		udValueFieldName := udValues.FieldByName(pathFieldName)
		if udValueFieldName.IsValid() {
			path := pathValues.Field(i).Interface().(string)
			value := udValueFieldName.Interface()
			if path != "" && value != nil {
				// Ex: "/latest/meta-data/instance-id" : "i-1234567890abcdef0"
				bvalue, err := base64.StdEncoding.DecodeString(config.Userdata.Values.Userdata)
				value = string(bvalue)

				if err != nil {
					fmt.Println("There was an issue decoding base64 data from config")
					panic(err)
				}
				supportedPaths[path] = value
				if config.Imdsv2Required {
					server.HandleFunc(path, imdsv2.ValidateToken(Handler))
				} else {

					server.HandleFunc(path, Handler)
				}
			} else {
				log.Printf("There was an issue registering path %v with udValue: %v", path, value)
			}
		}
	}
}
