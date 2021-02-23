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

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockState struct {
	mock.Mock
}

func (m *mockState) Parse(pos int, char rune) error {
	args := m.MethodCalled("Parse", pos, char)

	return args.Error(0)
}

func TestParseBase(t *testing.T) {
	s := &mockState{}
	s.On("Parse", 0, 't').Return(nil).Once()
	s.On("Parse", 1, 'e').Return(nil).Once()
	s.On("Parse", 2, 's').Return(nil).Once()
	s.On("Parse", 3, 't').Return(nil).Once()

	err := Parse("test", s)

	assert.NoError(t, err)
	s.AssertExpectations(t)
}

func TestParseError(t *testing.T) {
	s := &mockState{}
	s.On("Parse", 0, 't').Return(nil).Once()
	s.On("Parse", 1, 'e').Return(nil).Once()
	s.On("Parse", 2, 's').Return(assert.AnError).Once()

	err := Parse("test", s)

	assert.Same(t, assert.AnError, err)
	s.AssertExpectations(t)
}
