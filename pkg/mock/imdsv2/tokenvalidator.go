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

package imdsv2

import (
	"log"
	"net/http"
	"time"

	"github.com/aws/amazon-ec2-metadata-mock/pkg/server"
)

const (
	tokenRequestHeader = "X-aws-ec2-metadata-token"
)

// ValidateToken is a wrapper to validate token before passing request to provided handler
func ValidateToken(pathHandler server.HandlerType) server.HandlerType {
	return func(res http.ResponseWriter, req *http.Request) {
		providedToken := req.Header.Get(tokenRequestHeader)
		if providedToken == "" {
			log.Println("Token required; No token provided.")
			server.ReturnUnauthorizedResponse(res)
			return
		}

		if actualToken, ok := generatedTokens[providedToken]; ok {
			accessTime := time.Now()
			duration := accessTime.Sub(actualToken.CreatedAt)
			if int(duration.Seconds()) >= actualToken.TTL {
				log.Println("Token has expired")
				delete(generatedTokens, providedToken)
				server.ReturnUnauthorizedResponse(res)
				return
			}
			log.Println("Token validated!")
			pathHandler(res, req)
		} else {
			log.Printf("Invalid token provided: %v", providedToken)
			server.ReturnUnauthorizedResponse(res)
			return
		}
	}
}
