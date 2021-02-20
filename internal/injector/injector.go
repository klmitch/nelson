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
func IsInjectorError(e error) bool {
	var tmp Error

	return errors.As(e, &tmp)
}

// Standard errors that may occur within the package.
var (
	ErrNil          = Error{Message: "cannot inject nil"}
	ErrDuplicate    = Error{Message: "object of type already available"}
	ErrBadInterface = Error{Message: "bad interface type"}
	ErrBadType      = Error{Message: "cannot be assigned to type"}
	ErrNoMethod     = Error{Message: "no such method"}
	ErrBadMethod    = Error{Message: "method is not a function"}
	ErrMissingValue = Error{Message: "injector missing value for type"}
)

// errType is the type of the error interface.
var errType = reflect.TypeOf((*error)(nil)).Elem()

// Vivifier is an interface representing vivifiers.  A vivifier, if
// found in the Injector, creates a desired type.  A Vivifier is
// passed the Injector it is a part of, and must return the object; it
// MUST NOT modify the Injector directly as part of its action.
type Vivifier interface {
	// Vivify must construct the desired object and return it, or
	// an error if an error occurs.  It is passed the Injector,
	// but MUST NOT modify it.
	Vivify(inj *Injector, typ reflect.Type) (interface{}, error)
}

// Injector is a type that allows for dependency injection.  It
// contains a number of things that may be injected, or special types
// that automatically vivify such a type, and can then invoke a
// specified method injecting the correct arguments.
type Injector struct {
	Objects   map[reflect.Type]reflect.Value // Injectible objects
	Vivifiers map[reflect.Type]Vivifier      // Vivifiers
	Fallback  Vivifier                       // Fallback vivifier
}

// add adds an object associated with a specific type to the injector.
// It verifies that there are no duplicate types and that the value
// can be assigned to that type.
func (i *Injector) add(typ reflect.Type, obj interface{}) (reflect.Value, error) {
	// Reflect the value
	val, ok := obj.(reflect.Value)
	if !ok {
		val = reflect.ValueOf(obj)
	}

	// Make sure it's assignable
	if !val.Type().AssignableTo(typ) {
		return reflect.Value{}, fmt.Errorf("%#v %w %s", val, ErrBadType, typ.String())
	}

	// Add object to the injector
	if i.Objects == nil {
		i.Objects = map[reflect.Type]reflect.Value{}
	} else if _, ok = i.Objects[typ]; ok {
		return reflect.Value{}, fmt.Errorf("type %s: %w", typ.String(), ErrDuplicate)
	}
	i.Objects[typ] = val

	return val, nil
}

// Add adds an object to an Injector.  It returns an error if another
// object with the same type is already present.
func (i *Injector) Add(obj interface{}) error {
	// Make sure we're not trying to inject a nil...
	if obj == nil {
		return ErrNil
	}

	// Induct the object
	val, ok := obj.(reflect.Value)
	if !ok {
		val = reflect.ValueOf(obj)
	}

	// Add to the Injector
	if _, err := i.add(val.Type(), val); err != nil {
		return err
	}

	return nil
}

// AddInterface is a variant of Add which allows explicitly specifying
// an interface type.  This would typically be used when the type to
// be injected is an interface.  The type passed should be a nil
// pointer to the interface type.
func (i *Injector) AddInterface(iface interface{}, obj interface{}) error {
	// Make sure we're not trying to inject a nil...
	if obj == nil {
		return ErrNil
	}

	// Determine the type
	typ, ok := iface.(reflect.Type)
	if !ok {
		typ = reflect.TypeOf(iface)
		if typ.Kind() != reflect.Ptr {
			return ErrBadInterface
		}
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Interface {
		return ErrBadInterface
	}

	// Add to the Injector
	if _, err := i.add(typ, obj); err != nil {
		return err
	}

	return nil
}

// AddVivifier is a variant of Add which adds a Vivifier, an object
// that constructs the desired object to the Injector on demand.
func (i *Injector) AddVivifier(obj interface{}, viv Vivifier) error {
	// Make sure we have a vivifier
	if viv == nil {
		return ErrNil
	}

	// Determine the type
	typ, ok := obj.(reflect.Type)
	if !ok {
		typ = reflect.TypeOf(obj)
		if typ.Kind() == reflect.Ptr {
			tmp := typ.Elem()
			if tmp.Kind() == reflect.Interface {
				typ = tmp
			}
		}
	}

	// Add the vivifier
	if i.Vivifiers == nil {
		i.Vivifiers = map[reflect.Type]Vivifier{}
	} else if _, ok = i.Vivifiers[typ]; ok {
		return fmt.Errorf("type %s: %w", typ.String(), ErrDuplicate)
	}
	i.Vivifiers[typ] = viv

	return nil
}

// Get retrieves the object matching the specified type from the
// Injector.
func (i *Injector) Get(typ reflect.Type) (reflect.Value, error) {
	// First, look in Objects
	if i.Objects != nil {
		if val, ok := i.Objects[typ]; ok {
			return val, nil
		}
	}

	// OK, maybe we can vivify it?
	if i.Vivifiers != nil {
		if viv, ok := i.Vivifiers[typ]; ok {
			obj, err := viv.Vivify(i, typ)
			if err != nil {
				return reflect.Value{}, err
			}

			// Add to the Injector
			return i.add(typ, obj)
		}
	}

	// OK, try the fallback
	if i.Fallback != nil {
		obj, err := i.Fallback.Vivify(i, typ)
		if err != nil {
			return reflect.Value{}, err
		}

		// Add to the Injector
		return i.add(typ, obj)
	}

	return reflect.Value{}, fmt.Errorf("%w %s", ErrMissingValue, typ.String())
}

// Call calls a specified method on a specified object.  The method
// must either return nothing or return an error.
func (i *Injector) Call(obj interface{}, method string) error {
	// Get a value for the object
	if obj == nil {
		return fmt.Errorf("%w %q", ErrNoMethod, method)
	}
	val, ok := obj.(reflect.Value)
	if !ok {
		val = reflect.ValueOf(obj)
	}

	// Look up the method
	meth := val.MethodByName(method)
	if !meth.IsValid() {
		return fmt.Errorf("%w %q", ErrNoMethod, method)
	}

	return i.CallMethod(meth)
}

// CallMethod calls the specified method.  The method must either
// return nothing or return an error.
func (i *Injector) CallMethod(meth reflect.Value) error {
	// Make sure we have a method to call
	if !meth.IsValid() {
		return ErrNoMethod
	} else if meth.Kind() != reflect.Func {
		return ErrBadMethod
	}

	// Get the method type information
	mTyp := meth.Type()
	if mTyp.IsVariadic() || mTyp.NumOut() > 1 || (mTyp.NumOut() == 1 && !mTyp.Out(0).AssignableTo(errType)) {
		return ErrBadMethod
	}

	// Put together the list of input values
	values := []reflect.Value{}
	for j := 0; j < mTyp.NumIn(); j++ {
		val, err := i.Get(mTyp.In(j))
		if err != nil {
			return err
		}
		values = append(values, val)
	}

	// Call the method
	result := meth.Call(values)
	if len(result) > 0 {
		return result[0].Interface().(error)
	}

	return nil
}
