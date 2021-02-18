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

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandErrorImplementsError(t *testing.T) {
	assert.Implements(t, (*error)(nil), &CommandError{})
}

func TestCommandErrorError(t *testing.T) {
	obj := &CommandError{
		Err: errors.New("some random error"), //nolint:goerr113
	}

	result := obj.Error()

	assert.Equal(t, "some random error", result)
}

func TestCommandErrorUnwrap(t *testing.T) {
	obj := &CommandError{
		Err: assert.AnError,
	}

	result := obj.Unwrap()

	assert.Same(t, assert.AnError, result)
}

func TestExitControlBase(t *testing.T) {
	code, usage := ExitControl(assert.AnError)

	assert.Equal(t, 1, code)
	assert.False(t, usage)
}

func TestExitControlSpecific(t *testing.T) {
	err := fmt.Errorf("Test error %w", &CommandError{
		Code: 5,
	})

	code, usage := ExitControl(err)

	assert.Equal(t, 5, code)
	assert.False(t, usage)
}

func TestExitControlUsage(t *testing.T) {
	err := fmt.Errorf("Test error %w", &CommandError{
		Usage: true,
	})

	code, usage := ExitControl(err)

	assert.Equal(t, 0, code)
	assert.True(t, usage)
}
