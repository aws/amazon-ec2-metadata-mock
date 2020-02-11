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
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/aws/amazon-ec2-metadata-mock/pkg/server"
	h "github.com/aws/amazon-ec2-metadata-mock/test"
)

const (
	testURL             = "http://test-example.com/"
	tokenRegex          = "^[a-zA-Z0-9]{43}="
	successMockResponse = "Success!"
)

func MockHandler(res http.ResponseWriter, req *http.Request) {
	server.FormatAndReturnTextResponse(res, successMockResponse)
}

// Token Generator Tests
func TestGenerateToken(t *testing.T) {
	req := httptest.NewRequest("PUT", testURL, nil)
	req.Header.Set(tokenTTLHeader, "21500")
	generateTokenResp := executeTestHTTPRequest(req, GenerateToken)
	h.Assert(t, isTokenValid(generateTokenResp), fmt.Sprintf("Expected valid token generation, but was %s", generateTokenResp))
}
func TestInvalidGenerateTokenRequestGet(t *testing.T) {
	req := httptest.NewRequest("GET", testURL, nil)
	req.Header.Set(tokenTTLHeader, "21500")
	generateTokenResp := executeTestHTTPRequest(req, GenerateToken)
	h.Assert(t, generateTokenResp == "", fmt.Sprintf("Expected no token, but was %s", generateTokenResp))
}
func TestInvalidGenerateTokenInvalidTTL(t *testing.T) {
	req := httptest.NewRequest("PUT", testURL, nil)
	req.Header.Set(tokenTTLHeader, "0")
	generateTokenResp := executeTestHTTPRequest(req, GenerateToken)
	h.Assert(t, strings.TrimSpace(generateTokenResp) == server.BadRequestResponse, fmt.Sprintf("Expected 400 -- Bad Request, but was %s", generateTokenResp))
}
func TestInvalidGenerateTokenNoTTL(t *testing.T) {
	req := httptest.NewRequest("PUT", testURL, nil)
	generateTokenResp := executeTestHTTPRequest(req, GenerateToken)
	h.Assert(t, generateTokenResp == server.BadRequestResponse, fmt.Sprintf("Expected 400 -- Bad Request, but was %s", generateTokenResp))
}

// Token Validator Tests
func TestValidateToken(t *testing.T) {
	req := httptest.NewRequest("PUT", testURL, nil)
	req.Header.Set(tokenTTLHeader, "21500")
	validTestToken := executeTestHTTPRequest(req, GenerateToken)

	req = httptest.NewRequest("GET", testURL, nil)
	req.Header.Set(tokenRequestHeader, validTestToken)
	validateTokenResp := executeTestHTTPRequest(req, http.HandlerFunc(ValidateToken(MockHandler)))
	h.Assert(t, validateTokenResp == successMockResponse, fmt.Sprintf("Expected successful token validation, but was %s", validateTokenResp))
}
func TestValidateTokenNoToken(t *testing.T) {
	req := httptest.NewRequest("GET", testURL, nil)
	validateTokenResp := executeTestHTTPRequest(req, http.HandlerFunc(ValidateToken(MockHandler)))
	h.Assert(t, validateTokenResp == server.UnauthorizedResponse, fmt.Sprintf("Expected 401 -- Unauthorized for no token, but was %s", validateTokenResp))
}
func TestValidateTokenInvalidToken(t *testing.T) {
	invalidTestToken := "ThisTokenIsNotValid!"
	req := httptest.NewRequest("GET", testURL, nil)
	req.Header.Set(tokenRequestHeader, invalidTestToken)
	validateTokenResp := executeTestHTTPRequest(req, http.HandlerFunc(ValidateToken(MockHandler)))
	h.Assert(t, validateTokenResp == server.UnauthorizedResponse, fmt.Sprintf("401 -- Unauthorized for invalid token, but was %s", validateTokenResp))
}
func TestValidateTokenExpiredToken(t *testing.T) {
	req := httptest.NewRequest("PUT", testURL, nil)
	req.Header.Set(tokenTTLHeader, "1")
	expiredTestToken := executeTestHTTPRequest(req, GenerateToken)

	time.Sleep(1 * time.Second)

	req = httptest.NewRequest("GET", testURL, nil)
	req.Header.Set(tokenRequestHeader, expiredTestToken)
	validateTokenResp := executeTestHTTPRequest(req, http.HandlerFunc(ValidateToken(MockHandler)))
	h.Assert(t, validateTokenResp == server.UnauthorizedResponse, fmt.Sprintf("401 -- Unauthorized for expired token, but was %s", validateTokenResp))
}

// Test Helpers
func isTokenValid(token string) bool {
	matched, _ := regexp.Match(tokenRegex, []byte(token))
	return matched
}
func executeTestHTTPRequest(req *http.Request, handler http.HandlerFunc) string {
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	return strings.TrimSpace(string(body))
}
