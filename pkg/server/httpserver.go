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

package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
)

const notFoundResponse = `<?xml version="1.0" encoding="iso-8859-1"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
	"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
 <head>
  <title>404 - Not Found</title>
 </head>
 <body>
  <h1>404 - Not Found</h1>
 </body>
</html>`

// BadRequestResponse represents the IMDSv2 response in the event of missing or invalid parameters in the request
const BadRequestResponse = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
   <head>
      <title>400 - Bad Request</title>
   </head>
   <body>
      <h1>400 - Bad Request</h1>
   </body>
</html>`

// UnauthorizedResponse represents the IMDSv2 response in the event of unauthorized access
const UnauthorizedResponse = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
   <head>
      <title>401 - Unauthorized</title>
   </head>
   <body>
      <h1>401 - Unauthorized</h1>
   </body>
</html>`

var (
	// Routes represents the list of routes served by the http server
	Routes []string
)

// HandlerType represents the function passed as an argument to HandleFunc
type HandlerType func(http.ResponseWriter, *http.Request)

// HandleFunc registers the handler function for the given pattern
func HandleFunc(pattern string, requestHandler HandlerType) {
	http.HandleFunc(pattern, requestHandler)
}

func listRoutes() {
	routesMap := reflect.ValueOf(http.DefaultServeMux).Elem().FieldByIndex([]int{1})
	routes := routesMap.MapKeys()
	for _, r := range routes {
		Routes = append(Routes, r.String())
	}
	log.Printf("Serving the following routes: %s\n\n", strings.Join(Routes, ", "))
}

// ListenAndServe serves all patterns setup via their respective handlers
func ListenAndServe(hostname string, port string) {
	listRoutes()
	host := fmt.Sprint(hostname, ":", port)
	if err := http.ListenAndServe(host, nil); err != nil {
		panic(err)
	}
}

// FormatAndReturnJSONResponse formats the given data into JSON and returns the response
func FormatAndReturnJSONResponse(res http.ResponseWriter, data interface{}) {
	res.Header().Set("Content-Type", "application/json")

	var err error
	var metadataPrettyJSON []byte
	if metadataPrettyJSON, err = json.MarshalIndent(data, "", "\t"); err != nil {
		log.Fatalf("Error while attempting to format data %s for response: %s", data, err)
	}
	res.Write(metadataPrettyJSON)
	log.Println("Returned JSON mock response successfully.")
	return
}

// FormatAndReturnTextResponse formats the given data as plaintext and returns the response
func FormatAndReturnTextResponse(res http.ResponseWriter, data string) {
	res.Header().Set("Content-Type", "text/plain")
	res.Write([]byte(data))
	log.Println("Returned text mock response successfully.")
	return
}

// ReturnNotFoundResponse returns response with 404 Not Found
func ReturnNotFoundResponse(w http.ResponseWriter) {
	http.Error(w, notFoundResponse, http.StatusNotFound)
	return
}

// ReturnBadRequestResponse returns response with 400 Bad Request
func ReturnBadRequestResponse(w http.ResponseWriter) {
	http.Error(w, BadRequestResponse, http.StatusBadRequest)
	return
}

// ReturnUnauthorizedResponse returns response with 401 Unauthorized
func ReturnUnauthorizedResponse(w http.ResponseWriter) {
	http.Error(w, UnauthorizedResponse, http.StatusUnauthorized)
	return
}
