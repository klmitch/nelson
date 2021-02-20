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

package injector

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

func TestIsInjectorErrorBase(t *testing.T) {
	err := fmt.Errorf("%w", Error{Message: "an error"})

	result := IsInjectorError(err)

	assert.True(t, result)
}

func TestIsInjectorErrorFalse(t *testing.T) {
	result := IsInjectorError(assert.AnError)

	assert.False(t, result)
}

type mockVivifier struct {
	mock.Mock
}

func (m *mockVivifier) Vivify(inj *Injector, typ reflect.Type) (interface{}, error) {
	args := m.MethodCalled("Vivify", inj, typ)

	return args.Get(0), args.Error(1)
}

func TestInjectorAddInternalBase(t *testing.T) {
	typ := reflect.TypeOf("")
	obj := &Injector{}

	result, err := obj.add(typ, "a string")

	assert.NoError(t, err)
	assert.Equal(t, "a string", result.Interface())
	assert.Equal(t, &Injector{
		Objects: map[reflect.Type]reflect.Value{
			typ: result,
		},
	}, obj)
}

func TestInjectorAddInternalValue(t *testing.T) {
	typ := reflect.TypeOf("")
	obj := &Injector{}

	result, err := obj.add(typ, reflect.ValueOf("a string"))

	assert.NoError(t, err)
	assert.Equal(t, "a string", result.Interface())
	assert.Equal(t, &Injector{
		Objects: map[reflect.Type]reflect.Value{
			typ: result,
		},
	}, obj)
}

func TestInjectorAddInternalNotAssignable(t *testing.T) {
	typ := reflect.TypeOf("")
	obj := &Injector{}

	result, err := obj.add(typ, 12345)

	assert.ErrorIs(t, err, ErrBadType)
	assert.False(t, result.IsValid())
	assert.Equal(t, &Injector{}, obj)
}

func TestInjectorAddInternalExists(t *testing.T) {
	typ := reflect.TypeOf("")
	obj := &Injector{
		Objects: map[reflect.Type]reflect.Value{
			typ: reflect.ValueOf("another string"),
		},
	}

	result, err := obj.add(typ, "a string")

	assert.ErrorIs(t, err, ErrDuplicate)
	assert.False(t, result.IsValid())
	assert.NotEqual(t, &Injector{
		Objects: map[reflect.Type]reflect.Value{
			typ: result,
		},
	}, obj)
}

func TestInjectorAddBase(t *testing.T) {
	typ := reflect.TypeOf("")
	obj := &Injector{}

	err := obj.Add("a string")

	assert.NoError(t, err)
	assert.Contains(t, obj.Objects, typ)
}

func TestInjectorAddValue(t *testing.T) {
	typ := reflect.TypeOf("")
	obj := &Injector{}

	err := obj.Add(reflect.ValueOf("a string"))

	assert.NoError(t, err)
	assert.Contains(t, obj.Objects, typ)
}

func TestInjectorAddNil(t *testing.T) {
	obj := &Injector{}

	err := obj.Add(nil)

	assert.ErrorIs(t, err, ErrNil)
	assert.Equal(t, &Injector{}, obj)
}

func TestInjectorAddError(t *testing.T) {
	typ := reflect.TypeOf("")
	obj := &Injector{
		Objects: map[reflect.Type]reflect.Value{
			typ: reflect.ValueOf("another string"),
		},
	}

	err := obj.Add("another string")

	assert.ErrorIs(t, err, ErrDuplicate)
}

type iTrialInterface interface {
	Method() int
}

type trialInterface int

func (i trialInterface) Method() int {
	return int(i)
}

func TestInjectorAddInterfaceBase(t *testing.T) {
	typ := reflect.TypeOf((*iTrialInterface)(nil)).Elem()
	obj := &Injector{}

	err := obj.AddInterface((*iTrialInterface)(nil), trialInterface(0))

	assert.NoError(t, err)
	assert.Contains(t, obj.Objects, typ)
}

func TestInjectorAddInterfaceReflected(t *testing.T) {
	typ := reflect.TypeOf((*iTrialInterface)(nil)).Elem()
	obj := &Injector{}

	err := obj.AddInterface(typ, trialInterface(0))

	assert.NoError(t, err)
	assert.Contains(t, obj.Objects, typ)
}

func TestInjectorAddInterfaceValue(t *testing.T) {
	typ := reflect.TypeOf((*iTrialInterface)(nil)).Elem()
	obj := &Injector{}

	err := obj.AddInterface((*iTrialInterface)(nil), reflect.ValueOf(trialInterface(0)))

	assert.NoError(t, err)
	assert.Contains(t, obj.Objects, typ)
}

func TestInjectorAddInterfaceNil(t *testing.T) {
	obj := &Injector{}

	err := obj.AddInterface((*iTrialInterface)(nil), nil)

	assert.ErrorIs(t, err, ErrNil)
	assert.Equal(t, &Injector{}, obj)
}

func TestInjectorAddInterfaceNotPtr(t *testing.T) {
	obj := &Injector{}

	err := obj.AddInterface(5, trialInterface(0))

	assert.ErrorIs(t, err, ErrBadInterface)
	assert.Equal(t, &Injector{}, obj)
}

func TestInjectorAddInterfaceNotInterface(t *testing.T) {
	obj := &Injector{}

	err := obj.AddInterface(&mock.Mock{}, trialInterface(0))

	assert.ErrorIs(t, err, ErrBadInterface)
	assert.Equal(t, &Injector{}, obj)
}

func TestInjectorAddInterfaceReflectedNotInterface(t *testing.T) {
	obj := &Injector{}

	err := obj.AddInterface(reflect.TypeOf(mock.Mock{}), trialInterface(0))

	assert.ErrorIs(t, err, ErrBadInterface)
	assert.Equal(t, &Injector{}, obj)
}

func TestInjectorAddInterfaceError(t *testing.T) {
	typ := reflect.TypeOf((*iTrialInterface)(nil)).Elem()
	obj := &Injector{
		Objects: map[reflect.Type]reflect.Value{
			typ: reflect.ValueOf(trialInterface(5)),
		},
	}

	err := obj.AddInterface((*iTrialInterface)(nil), trialInterface(0))

	assert.ErrorIs(t, err, ErrDuplicate)
}

func TestInjectorAddVivifierBase(t *testing.T) {
	typ := reflect.TypeOf("")
	viv := &mockVivifier{}
	obj := &Injector{}

	err := obj.AddVivifier("a string", viv)

	assert.NoError(t, err)
	assert.Equal(t, &Injector{
		Vivifiers: map[reflect.Type]Vivifier{
			typ: viv,
		},
	}, obj)
}

func TestInjectorAddVivifierType(t *testing.T) {
	typ := reflect.TypeOf("")
	viv := &mockVivifier{}
	obj := &Injector{}

	err := obj.AddVivifier(typ, viv)

	assert.NoError(t, err)
	assert.Equal(t, &Injector{
		Vivifiers: map[reflect.Type]Vivifier{
			typ: viv,
		},
	}, obj)
}

func TestInjectorAddVivifierInterface(t *testing.T) {
	typ := reflect.TypeOf((*iTrialInterface)(nil)).Elem()
	viv := &mockVivifier{}
	obj := &Injector{}

	err := obj.AddVivifier((*iTrialInterface)(nil), viv)

	assert.NoError(t, err)
	assert.Equal(t, &Injector{
		Vivifiers: map[reflect.Type]Vivifier{
			typ: viv,
		},
	}, obj)
}

func TestInjectorAddVivifierInterfaceType(t *testing.T) {
	typ := reflect.TypeOf((*iTrialInterface)(nil))
	viv := &mockVivifier{}
	obj := &Injector{}

	err := obj.AddVivifier(typ, viv)

	assert.NoError(t, err)
	assert.Equal(t, &Injector{
		Vivifiers: map[reflect.Type]Vivifier{
			typ: viv,
		},
	}, obj)
}

func TestInjectorAddVivifierNil(t *testing.T) {
	obj := &Injector{}

	err := obj.AddVivifier("a string", nil)

	assert.ErrorIs(t, err, ErrNil)
	assert.Equal(t, &Injector{}, obj)
}

func TestInjectorAddVivifierDuplicate(t *testing.T) {
	typ := reflect.TypeOf("")
	viv := &mockVivifier{}
	obj := &Injector{
		Vivifiers: map[reflect.Type]Vivifier{
			typ: &mockVivifier{},
		},
	}

	err := obj.AddVivifier("a string", viv)

	assert.ErrorIs(t, err, ErrDuplicate)
}

func TestInjectorGetBase(t *testing.T) {
	val := reflect.ValueOf("a string")
	typ := val.Type()
	obj := &Injector{
		Objects: map[reflect.Type]reflect.Value{
			typ: val,
		},
	}

	result, err := obj.Get(typ)

	assert.NoError(t, err)
	assert.Equal(t, val, result)
}

func TestInjectorGetVivifyBase(t *testing.T) {
	typ := reflect.TypeOf("")
	viv := &mockVivifier{}
	obj := &Injector{
		Vivifiers: map[reflect.Type]Vivifier{
			typ: viv,
		},
	}
	viv.On("Vivify", obj, typ).Return("a string", nil)

	result, err := obj.Get(typ)

	assert.NoError(t, err)
	assert.Equal(t, "a string", result.Interface())
	assert.Equal(t, &Injector{
		Objects: map[reflect.Type]reflect.Value{
			typ: result,
		},
		Vivifiers: map[reflect.Type]Vivifier{
			typ: viv,
		},
	}, obj)
	viv.AssertExpectations(t)
}

func TestInjectorGetVivifyObjectMissing(t *testing.T) {
	typ := reflect.TypeOf("")
	viv := &mockVivifier{}
	obj := &Injector{
		Objects: map[reflect.Type]reflect.Value{},
		Vivifiers: map[reflect.Type]Vivifier{
			typ: viv,
		},
	}
	viv.On("Vivify", obj, typ).Return("a string", nil)

	result, err := obj.Get(typ)

	assert.NoError(t, err)
	assert.Equal(t, "a string", result.Interface())
	assert.Equal(t, &Injector{
		Objects: map[reflect.Type]reflect.Value{
			typ: result,
		},
		Vivifiers: map[reflect.Type]Vivifier{
			typ: viv,
		},
	}, obj)
	viv.AssertExpectations(t)
}

func TestInjectorGetFallbackBase(t *testing.T) {
	typ := reflect.TypeOf("")
	viv := &mockVivifier{}
	obj := &Injector{
		Fallback: viv,
	}
	viv.On("Vivify", obj, typ).Return("a string", nil)

	result, err := obj.Get(typ)

	assert.NoError(t, err)
	assert.Equal(t, "a string", result.Interface())
	assert.Equal(t, &Injector{
		Objects: map[reflect.Type]reflect.Value{
			typ: result,
		},
		Fallback: viv,
	}, obj)
	viv.AssertExpectations(t)
}

func TestInjectorGetFallbackVivifierMissing(t *testing.T) {
	typ := reflect.TypeOf("")
	viv := &mockVivifier{}
	obj := &Injector{
		Vivifiers: map[reflect.Type]Vivifier{},
		Fallback:  viv,
	}
	viv.On("Vivify", obj, typ).Return("a string", nil)

	result, err := obj.Get(typ)

	assert.NoError(t, err)
	assert.Equal(t, "a string", result.Interface())
	assert.Equal(t, &Injector{
		Objects: map[reflect.Type]reflect.Value{
			typ: result,
		},
		Vivifiers: map[reflect.Type]Vivifier{},
		Fallback:  viv,
	}, obj)
	viv.AssertExpectations(t)
}

func TestInjectorGetVivifyError(t *testing.T) {
	typ := reflect.TypeOf("")
	viv := &mockVivifier{}
	obj := &Injector{
		Vivifiers: map[reflect.Type]Vivifier{
			typ: viv,
		},
	}
	viv.On("Vivify", obj, typ).Return(nil, assert.AnError)

	result, err := obj.Get(typ)

	assert.ErrorIs(t, err, assert.AnError)
	assert.False(t, result.IsValid())
	assert.Equal(t, &Injector{
		Vivifiers: map[reflect.Type]Vivifier{
			typ: viv,
		},
	}, obj)
	viv.AssertExpectations(t)
}

func TestInjectorGetFallbackError(t *testing.T) {
	typ := reflect.TypeOf("")
	viv := &mockVivifier{}
	obj := &Injector{
		Fallback: viv,
	}
	viv.On("Vivify", obj, typ).Return(nil, assert.AnError)

	result, err := obj.Get(typ)

	assert.ErrorIs(t, err, assert.AnError)
	assert.False(t, result.IsValid())
	assert.Equal(t, &Injector{
		Fallback: viv,
	}, obj)
	viv.AssertExpectations(t)
}

func TestInjectorGetMissing(t *testing.T) {
	typ := reflect.TypeOf("")
	obj := &Injector{}

	result, err := obj.Get(typ)

	assert.ErrorIs(t, err, ErrMissingValue)
	assert.False(t, result.IsValid())
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

func TestInjectorCallBase(t *testing.T) {
	val := &methods{}
	val.On("Niladic")
	obj := &Injector{}

	err := obj.Call(val, "Niladic")

	assert.NoError(t, err)
	val.AssertExpectations(t)
}

func TestInjectorCallValue(t *testing.T) {
	val := &methods{}
	val.On("Niladic")
	obj := &Injector{}

	err := obj.Call(reflect.ValueOf(val), "Niladic")

	assert.NoError(t, err)
	val.AssertExpectations(t)
}

func TestInjectorCallNoObject(t *testing.T) {
	obj := &Injector{}

	err := obj.Call(nil, "Niladic")

	assert.ErrorIs(t, err, ErrNoMethod)
}

func TestInjectorCallNoMethod(t *testing.T) {
	val := &methods{}
	obj := &Injector{}

	err := obj.Call(val, "NoMethod")

	assert.ErrorIs(t, err, ErrNoMethod)
	val.AssertExpectations(t)
}

func TestInjectorCallError(t *testing.T) {
	val := &methods{}
	val.On("NiladicErr").Return(assert.AnError)
	obj := &Injector{}

	err := obj.Call(val, "NiladicErr")

	assert.ErrorIs(t, err, assert.AnError)
	val.AssertExpectations(t)
}

func TestInjectorCallMethodBase(t *testing.T) {
	val := &methods{}
	val.On("Niladic")
	meth := reflect.ValueOf(val).MethodByName("Niladic")
	obj := &Injector{}

	err := obj.CallMethod(meth)

	assert.NoError(t, err)
	val.AssertExpectations(t)
}

func TestInjectorCallMethodBaseError(t *testing.T) {
	val := &methods{}
	val.On("NiladicErr").Return(assert.AnError)
	meth := reflect.ValueOf(val).MethodByName("NiladicErr")
	obj := &Injector{}

	err := obj.CallMethod(meth)

	assert.ErrorIs(t, err, assert.AnError)
	val.AssertExpectations(t)
}

func TestInjectorCallMethodBasic(t *testing.T) {
	val := &methods{}
	val.On("Basic", 5, "a string")
	meth := reflect.ValueOf(val).MethodByName("Basic")
	obj := &Injector{
		Objects: map[reflect.Type]reflect.Value{
			reflect.TypeOf(0):  reflect.ValueOf(5),
			reflect.TypeOf(""): reflect.ValueOf("a string"),
		},
	}

	err := obj.CallMethod(meth)

	assert.NoError(t, err)
	val.AssertExpectations(t)
}

func TestInjectorCallMethodNoMethod(t *testing.T) {
	obj := &Injector{}

	err := obj.CallMethod(reflect.Value{})

	assert.ErrorIs(t, err, ErrNoMethod)
}

func TestInjectorCallMethodNotFunc(t *testing.T) {
	obj := &Injector{}

	err := obj.CallMethod(reflect.ValueOf(5))

	assert.ErrorIs(t, err, ErrBadMethod)
}

func TestInjectorCallMethodVariadic(t *testing.T) {
	val := &methods{}
	meth := reflect.ValueOf(val).MethodByName("Variadic")
	obj := &Injector{
		Objects: map[reflect.Type]reflect.Value{
			reflect.TypeOf(0):  reflect.ValueOf(5),
			reflect.TypeOf(""): reflect.ValueOf("a string"),
		},
	}

	err := obj.CallMethod(meth)

	assert.ErrorIs(t, err, ErrBadMethod)
	val.AssertExpectations(t)
}

func TestInjectorCallMethodTwoReturn(t *testing.T) {
	val := &methods{}
	meth := reflect.ValueOf(val).MethodByName("TwoReturn")
	obj := &Injector{}

	err := obj.CallMethod(meth)

	assert.ErrorIs(t, err, ErrBadMethod)
	val.AssertExpectations(t)
}

func TestInjectorCallMethodNonError(t *testing.T) {
	val := &methods{}
	meth := reflect.ValueOf(val).MethodByName("NonError")
	obj := &Injector{}

	err := obj.CallMethod(meth)

	assert.ErrorIs(t, err, ErrBadMethod)
	val.AssertExpectations(t)
}

func TestInjectorCallMethodMissing(t *testing.T) {
	val := &methods{}
	meth := reflect.ValueOf(val).MethodByName("Basic")
	obj := &Injector{
		Objects: map[reflect.Type]reflect.Value{
			reflect.TypeOf(0): reflect.ValueOf(5),
		},
	}

	err := obj.CallMethod(meth)

	assert.ErrorIs(t, err, ErrMissingValue)
	val.AssertExpectations(t)
}
