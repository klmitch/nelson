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
	"errors"
	"fmt"
	"reflect"
)

// Error is a wrapper for errors that identifies an error as coming
// from the Injector, as opposed to having some other source.
type Error struct {
	Message string // The error message
}

// Error returns the error message.
func (e Error) Error() string {
	return e.Message
}

// IsError is a test to see if an error is an Error.
func IsError(e error) bool {
	var tmp Error

	return errors.As(e, &tmp)
}

// Standard errors that may occur within the package.
var (
	ErrNoMethod     = Error{Message: "no such method"}
	ErrBadMethod    = Error{Message: "method is not a function"}
	ErrMissingValue = Error{Message: "missing input for type"}
)

// errType is the type of the error interface.
var errType = reflect.TypeOf((*error)(nil)).Elem()

// Deps is a collection of the dependencies of a method.  It is simply
// a set of types.
type Deps map[reflect.Type]reflect.Value

// Merge merges another Deps map into this one.  Values are ignored;
// only the types are relevant.
func (d Deps) Merge(other Deps) {
	for typ := range other {
		if _, ok := d[typ]; !ok {
			d[typ] = reflect.Value{}
		}
	}
}

// Copy constructs a new Deps from an existing one.  Values are
// ignored; only the types are relevant.
func (d Deps) Copy() Deps {
	result := Deps{}
	for typ := range d {
		result[typ] = reflect.Value{}
	}
	return result
}

// Method is a type that identifies a specific method.  It collects
// together its dependencies, and can be used to call that method on a
// specific object.
type Method struct {
	Name   string         // Name of the method
	Method reflect.Value  // The actual method
	Deps   Deps           // The dependencies of the method
	Args   []reflect.Type // Ordered list of arguments
}

// New constructs a new Method object for a specific method.
func New(obj interface{}, method string) (*Method, error) {
	// Get the Value of the object
	if obj == nil {
		return nil, fmt.Errorf("%w %q", ErrNoMethod, method)
	}
	val, ok := obj.(reflect.Value)
	if !ok {
		val = reflect.ValueOf(obj)
	}

	// Look up the method
	meth := val.MethodByName(method)
	if !meth.IsValid() {
		return nil, fmt.Errorf("%w %q", ErrNoMethod, method)
	}

	// Check the method type information
	mType := meth.Type()
	if mType.IsVariadic() || mType.NumOut() > 1 || (mType.NumOut() == 1 && !mType.Out(0).AssignableTo(errType)) {
		return nil, fmt.Errorf("%q: %w", method, ErrBadMethod)
	}

	// Begin constructing the result
	result := &Method{
		Name:   method,
		Method: meth,
		Deps:   Deps{},
	}

	// Account for inputs
	for i := 0; i < mType.NumIn(); i++ {
		vType := mType.In(i)
		if _, ok := result.Deps[vType]; ok {
			return nil, fmt.Errorf("%q: %w", method, ErrBadMethod)
		}
		result.Deps[vType] = reflect.Value{}
		result.Args = append(result.Args, vType)
	}

	return result, nil
}

// Call calls the method.  Inputs are a completed Deps.
func (m *Method) Call(inputs Deps) error {
	// Assemble inputs
	values := []reflect.Value{}
	for _, typ := range m.Args {
		tmp := inputs[typ]
		if !tmp.IsValid() {
			return fmt.Errorf("%q: %w %s", m.Name, ErrMissingValue, typ.String())
		}

		values = append(values, tmp)
	}

	// Call the method
	result := m.Method.Call(values)

	// Return the result
	if len(result) > 0 {
		return result[0].Interface().(error)
	}
	return nil
}
