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

package errors

import (
	"fmt"
)

// FlagValidationError denotes an error encountered while validating CLI flags
type FlagValidationError struct {
	FlagName     string
	Allowed      string
	InvalidValue string
}

// Error returns the formatted flag validation error.
func (ve FlagValidationError) Error() string {
	return fmt.Sprintf("Invalid CLI input \"%s\" for flag %s. Allowed value(s): %s.\n", ve.InvalidValue, ve.FlagName, ve.Allowed)
}
