// Copyright (c) 2021 Kevin L. Mitchell
//
// Licensed under the Apache License, Version 2.0 (the "License"); you
// may not use this file except in compliance with the License.  You
// may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.  See the License for the specific language governing
// permissions and limitations under the License.

package nelson

import "errors"

// CommandError is an implementation of the error interface that wraps
// another error and associates with it an error code to return.
type CommandError struct {
	Err   error // The wrapped error (if any)
	Code  int   // The exit code for the program
	Usage bool  // If true, emit a usage message
}

// Error returns the error message.
func (e *CommandError) Error() string {
	return e.Err.Error()
}

// Unwrap returns the wrapped error, if any.
func (e *CommandError) Unwrap() error {
	return e.Err
}

// ExitControl is a helper that determines the exit type.  If the
// error is not a CommandError, a default exit code of 1 and usage
// emission of false will be returned.
func ExitControl(err error) (int, bool) {
	var tmp *CommandError

	// Is it a CommandError?
	if errors.As(err, &tmp) {
		return tmp.Code, tmp.Usage
	}

	return 1, false
}
