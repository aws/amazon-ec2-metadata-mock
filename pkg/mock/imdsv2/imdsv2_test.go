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
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/amazon-ec2-metadata-mock/pkg/server"
	h "github.com/aws/amazon-ec2-metadata-mock/test"
)

const (
	testURL             = "http://test-example.com/"
	tokenRegex          = "^[a-zA-Z0-9+/]{43}="
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
	token, ttl := parseGenTokenResp(generateTokenResp)
	h.Assert(t, isTokenValid(token), fmt.Sprintf("Expected valid token generation, but was %s", token))
	h.Assert(t, ttl == 21500, fmt.Sprintf("Expected Token TTL header to equal requested value, but was %d", ttl))
}
func TestInvalidGenerateTokenRequestGet(t *testing.T) {
	req := httptest.NewRequest("GET", testURL, nil)
	req.Header.Set(tokenTTLHeader, "21500")
	generateTokenResp := executeTestHTTPRequest(req, GenerateToken)
	token, _ := parseGenTokenResp(generateTokenResp)
	h.Assert(t, token == "", fmt.Sprintf("Expected no token, but was %s", token))
}
func TestInvalidGenerateTokenInvalidTTL(t *testing.T) {
	req := httptest.NewRequest("PUT", testURL, nil)
	req.Header.Set(tokenTTLHeader, "0")
	generateTokenResp := executeTestHTTPRequest(req, GenerateToken)
	respContent, _ := parseGenTokenResp(generateTokenResp)
	h.Assert(t, respContent == server.BadRequestResponse, fmt.Sprintf("Expected 400 -- Bad Request, but was %s", respContent))
}
func TestInvalidGenerateTokenNoTTL(t *testing.T) {
	req := httptest.NewRequest("PUT", testURL, nil)
	generateTokenResp := executeTestHTTPRequest(req, GenerateToken)
	respContent, _ := parseGenTokenResp(generateTokenResp)
	h.Assert(t, respContent == server.BadRequestResponse, fmt.Sprintf("Expected 400 -- Bad Request, but was %s", respContent))
}

// Token Validator Tests
func TestValidateToken(t *testing.T) {
	req := httptest.NewRequest("PUT", testURL, nil)
	req.Header.Set(tokenTTLHeader, "21500")
	generateTokenResp := executeTestHTTPRequest(req, GenerateToken)
	validTestToken, _ := parseGenTokenResp(generateTokenResp)

	req = httptest.NewRequest("GET", testURL, nil)
	req.Header.Set(tokenRequestHeader, validTestToken)
	validateTokenResp := executeTestHTTPRequest(req, http.HandlerFunc(ValidateToken(MockHandler)))
	respContent, _ := ioutil.ReadAll(validateTokenResp.Body)
	h.Assert(t, strings.TrimSpace(string(respContent)) == successMockResponse, fmt.Sprintf("Expected successful token validation, but was %s", respContent))
}
func TestValidateTokenNoToken(t *testing.T) {
	req := httptest.NewRequest("GET", testURL, nil)
	validateTokenResp := executeTestHTTPRequest(req, http.HandlerFunc(ValidateToken(MockHandler)))
	respContent, _ := ioutil.ReadAll(validateTokenResp.Body)
	h.Assert(t, strings.TrimSpace(string(respContent)) == server.UnauthorizedResponse, fmt.Sprintf("Expected 401 -- Unauthorized for no token, but was %s", respContent))
}
func TestValidateTokenInvalidToken(t *testing.T) {
	invalidTestToken := "ThisTokenIsNotValid!"
	req := httptest.NewRequest("GET", testURL, nil)
	req.Header.Set(tokenRequestHeader, invalidTestToken)
	validateTokenResp := executeTestHTTPRequest(req, http.HandlerFunc(ValidateToken(MockHandler)))
	respContent, _ := ioutil.ReadAll(validateTokenResp.Body)
	h.Assert(t, strings.TrimSpace(string(respContent)) == server.UnauthorizedResponse, fmt.Sprintf("401 -- Unauthorized for invalid token, but was %s", respContent))
}
func TestValidateTokenExpiredToken(t *testing.T) {
	req := httptest.NewRequest("PUT", testURL, nil)
	req.Header.Set(tokenTTLHeader, "1")
	generateTokenResp := executeTestHTTPRequest(req, GenerateToken)
	expiredTestToken, _ := parseGenTokenResp(generateTokenResp)

	time.Sleep(1 * time.Second)

	req = httptest.NewRequest("GET", testURL, nil)
	req.Header.Set(tokenRequestHeader, expiredTestToken)
	validateTokenResp := executeTestHTTPRequest(req, http.HandlerFunc(ValidateToken(MockHandler)))
	respContent, _ := ioutil.ReadAll(validateTokenResp.Body)
	h.Assert(t, strings.TrimSpace(string(respContent)) == server.UnauthorizedResponse, fmt.Sprintf("401 -- Unauthorized for expired token, but was %s", respContent))
}

// Test Helpers
func isTokenValid(token string) bool {
	matched, _ := regexp.Match(tokenRegex, []byte(token))
	return matched
}

func executeTestHTTPRequest(req *http.Request, handler http.HandlerFunc) *http.Response {
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	resp := w.Result()
	return resp
}

func parseGenTokenResp(resp *http.Response) (string, int) {
	token, _ := ioutil.ReadAll(resp.Body)
	tokenString := strings.TrimSpace(string(token))
	ttl := resp.Header.Get(tokenTTLHeader)
	ttlInt := -1
	if ttl != "" {
		ttlInt, _ = strconv.Atoi(ttl)
	}
	return tokenString, ttlInt
}
