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
	"encoding/base64"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/amazon-ec2-metadata-mock/pkg/server"
)

const (
	tokenTTLHeader = "X-aws-ec2-metadata-token-ttl-seconds"
	maxTTL         = 21600 // 6 hours
)

var (
	generatedTokens = make(map[string]v2Token)
)

type v2Token struct {
	Value     string    // actual token value
	TTL       int       // token ttl
	CreatedAt time.Time // time the token was created
}

// GenerateToken returns a token with the specified TTL used for IMDSv2 requests
func GenerateToken(res http.ResponseWriter, req *http.Request) {
	// only valid with PUT
	if req.Method != http.MethodPut {
		return
	}
	// check that header contains valid ttl
	requestedTTL := req.Header.Get(tokenTTLHeader)
	validTTL, err := extractValidTTL(requestedTTL)
	if err != nil {
		log.Printf("Something went wrong with ttl validation: %v with requested TTL: %v", err.Error(), requestedTTL)
		server.ReturnBadRequestResponse(res)
		return
	}

	key := make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		server.FormatAndReturnTextResponse(res, "Something went wrong with token creation")
		return
	}

	tokenValue := base64.StdEncoding.EncodeToString(key)
	token := v2Token{
		Value:     tokenValue,
		TTL:       validTTL,
		CreatedAt: time.Now(),
	}
	generatedTokens[token.Value] = token
	server.FormatAndReturnTextResponse(res, token.Value)
}

func extractValidTTL(reqTTL string) (int, error) {
	if reqTTL == "" {
		log.Printf("TTL is required. requested TTL: %v", reqTTL)
		return 0, errors.New("TTL is nil")
	}

	intTTL, err := strconv.Atoi(reqTTL)

	if err != nil {
		log.Printf("Something went wrong with ttl conversion. requested TTL: %v", reqTTL)
		return 0, err
	}
	if intTTL <= 0 || intTTL > maxTTL {
		return 0, errors.New("TTL needs to be between 0-21600 seconds")
	}

	return intTTL, nil
}
