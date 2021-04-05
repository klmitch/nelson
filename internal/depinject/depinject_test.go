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

package depinject

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestErrorImplementsError(t *testing.T) {
	assert.Implements(t, (*error)(nil), Error{})
}

func TestErrorError(t *testing.T) {
	obj := Error{Message: "an error"}

	result := obj.Error()

	assert.Equal(t, "an error", result)
}

func TestIsErrorBase(t *testing.T) {
	err := fmt.Errorf("%w", Error{Message: "an error"})

	result := IsError(err)

	assert.True(t, result)
}

func TestIsErrorFalse(t *testing.T) {
	result := IsError(assert.AnError)

	assert.False(t, result)
}

func TestDepsMerge(t *testing.T) {
	obj := Deps{
		reflect.TypeOf(""): reflect.Value{},
	}
	other := Deps{
		reflect.TypeOf(0): reflect.ValueOf(0),
	}

	obj.Merge(other)

	assert.Equal(t, Deps{
		reflect.TypeOf(""): reflect.Value{},
		reflect.TypeOf(0):  reflect.Value{},
	}, obj)
}

func TestDepsCopy(t *testing.T) {
	obj := Deps{
		reflect.TypeOf(""): reflect.ValueOf("test"),
		reflect.TypeOf(0):  reflect.ValueOf(5),
	}

	result := obj.Copy()

	assert.Equal(t, Deps{
		reflect.TypeOf(""): reflect.Value{},
		reflect.TypeOf(0):  reflect.Value{},
	}, result)
	assert.NotEqual(t, result, obj)
}

type methods struct {
	mock.Mock
}

func (m *methods) Niladic() {
	m.MethodCalled("Niladic")
}

func (m *methods) NiladicErr() error {
	args := m.MethodCalled("NiladicErr")

	return args.Error(0)
}

func (m *methods) Basic(i int, s string) {
	m.MethodCalled("Basic", i, s)
}

func (m *methods) Variadic(i int, s ...string) {
	m.MethodCalled("Variadic", i, s)
}

func (m *methods) TwoReturn() (interface{}, error) {
	args := m.MethodCalled("TwoReturn")

	return args.Get(0), args.Error(1)
}

func (m *methods) NonError() interface{} {
	args := m.MethodCalled("NonError")

	return args.Get(0)
}

func (m *methods) DuplicatedInput(a int, b int) {
	m.MethodCalled("DuplicatedInput", a, b)
}

func TestNewNiladic(t *testing.T) {
	val := &methods{}
	val.On("Niladic").Once()

	result, err := New(val, "Niladic")

	assert.NoError(t, err)
	assert.Equal(t, "Niladic", result.Name)
	assert.Equal(t, Deps{}, result.Deps)
	assert.Nil(t, result.Args)
	callResult := result.Method.Call([]reflect.Value{})
	assert.Len(t, callResult, 0)
	val.AssertExpectations(t)
}

func TestNewNiladicValue(t *testing.T) {
	val := &methods{}
	val.On("Niladic").Once()

	result, err := New(reflect.ValueOf(val), "Niladic")

	assert.NoError(t, err)
	assert.Equal(t, "Niladic", result.Name)
	assert.Equal(t, Deps{}, result.Deps)
	assert.Nil(t, result.Args)
	callResult := result.Method.Call([]reflect.Value{})
	assert.Len(t, callResult, 0)
	val.AssertExpectations(t)
}

func TestNewNiladicErr(t *testing.T) {
	val := &methods{}
	val.On("NiladicErr").Return(assert.AnError).Once()

	result, err := New(val, "NiladicErr")

	assert.NoError(t, err)
	assert.Equal(t, "NiladicErr", result.Name)
	assert.Equal(t, Deps{}, result.Deps)
	assert.Nil(t, result.Args)
	callResult := result.Method.Call([]reflect.Value{})
	assert.Len(t, callResult, 1)
	assert.Same(t, assert.AnError, callResult[0].Interface())
	val.AssertExpectations(t)
}

func TestNewBasic(t *testing.T) {
	val := &methods{}
	val.On("Basic", 5, "test").Once()

	result, err := New(val, "Basic")

	assert.NoError(t, err)
	assert.Equal(t, "Basic", result.Name)
	assert.Equal(t, Deps{
		reflect.TypeOf(5):      reflect.Value{},
		reflect.TypeOf("test"): reflect.Value{},
	}, result.Deps)
	assert.Equal(t, []reflect.Type{
		reflect.TypeOf(5),
		reflect.TypeOf("test"),
	}, result.Args)
	values := []reflect.Value{
		reflect.ValueOf(5),
		reflect.ValueOf("test"),
	}
	callResult := result.Method.Call(values)
	assert.Len(t, callResult, 0)
	val.AssertExpectations(t)
}

func TestNewNil(t *testing.T) {
	result, err := New(nil, "Niladic")

	assert.ErrorIs(t, err, ErrNoMethod)
	assert.Nil(t, result)
}

func TestNewNoMethod(t *testing.T) {
	val := &methods{}

	result, err := New(val, "NoMethod")

	assert.ErrorIs(t, err, ErrNoMethod)
	assert.Nil(t, result)
	val.AssertExpectations(t)
}

func TestNewVariadic(t *testing.T) {
	val := &methods{}

	result, err := New(val, "Variadic")

	assert.ErrorIs(t, err, ErrBadMethod)
	assert.Nil(t, result)
	val.AssertExpectations(t)
}

func TestNewTwoReturn(t *testing.T) {
	val := &methods{}

	result, err := New(val, "TwoReturn")

	assert.ErrorIs(t, err, ErrBadMethod)
	assert.Nil(t, result)
	val.AssertExpectations(t)
}

func TestNewNonError(t *testing.T) {
	val := &methods{}

	result, err := New(val, "NonError")

	assert.ErrorIs(t, err, ErrBadMethod)
	assert.Nil(t, result)
	val.AssertExpectations(t)
}

func TestNewDuplicatedInput(t *testing.T) {
	val := &methods{}

	result, err := New(val, "DuplicatedInput")

	assert.ErrorIs(t, err, ErrBadMethod)
	assert.Nil(t, result)
	val.AssertExpectations(t)
}

func TestMethodCallNiladic(t *testing.T) {
	val := &methods{}
	val.On("Niladic")
	args := Deps{}
	obj := &Method{
		Name:   "Niladic",
		Method: reflect.ValueOf(val).MethodByName("Niladic"),
		Deps:   Deps{},
	}

	result := obj.Call(args)

	assert.NoError(t, result)
	val.AssertExpectations(t)
}

func TestMethodCallNiladicErr(t *testing.T) {
	val := &methods{}
	val.On("NiladicErr").Return(assert.AnError)
	args := Deps{}
	obj := &Method{
		Name:   "NiladicErr",
		Method: reflect.ValueOf(val).MethodByName("NiladicErr"),
		Deps:   Deps{},
	}

	result := obj.Call(args)

	assert.Same(t, assert.AnError, result)
	val.AssertExpectations(t)
}

func TestMethodCallBasic(t *testing.T) {
	val := &methods{}
	val.On("Basic", 5, "test")
	args := Deps{
		reflect.TypeOf(5):      reflect.ValueOf(5),
		reflect.TypeOf("test"): reflect.ValueOf("test"),
	}
	obj := &Method{
		Name:   "Basic",
		Method: reflect.ValueOf(val).MethodByName("Basic"),
		Deps: Deps{
			reflect.TypeOf(5):      reflect.Value{},
			reflect.TypeOf("test"): reflect.Value{},
		},
		Args: []reflect.Type{
			reflect.TypeOf(5),
			reflect.TypeOf("test"),
		},
	}

	result := obj.Call(args)

	assert.NoError(t, result)
	val.AssertExpectations(t)
}

func TestMethodCallMissingValue(t *testing.T) {
	val := &methods{}
	args := Deps{
		reflect.TypeOf(5): reflect.ValueOf(5),
	}
	obj := &Method{
		Name:   "Basic",
		Method: reflect.ValueOf(val).MethodByName("Basic"),
		Deps: Deps{
			reflect.TypeOf(5):      reflect.Value{},
			reflect.TypeOf("test"): reflect.Value{},
		},
		Args: []reflect.Type{
			reflect.TypeOf(5),
			reflect.TypeOf("test"),
		},
	}

	result := obj.Call(args)

	assert.ErrorIs(t, result, ErrMissingValue)
	val.AssertExpectations(t)
}
