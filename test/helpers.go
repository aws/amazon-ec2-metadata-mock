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

package test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

// ItemsMatch fails the test if the items in exp and act slices dont match.
// A nil argument is equivalent to an empty slice.
func ItemsMatch(tb testing.TB, exp, act []string) {
	if len(exp) != len(act) {
		tb.Errorf(fmt.Sprintf("Expected %d items in slice, but was %d", len(exp), len(act)))
	}
	for _, v := range exp {
		if !Contains(act, v) {
			tb.Errorf(fmt.Sprintf("Expected to find item %s in slice, but was not found", v))
		}
	}
}

// Contains returns a bool indicating whether the slice contains the given val
func Contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// Assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// Ok fails the test if an err is not nil.
func Ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}
