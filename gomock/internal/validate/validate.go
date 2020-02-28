// Copyright 2020 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validate

import (
	"fmt"
	"reflect"
)

// InputAndOutputSig compares the argument and return signatures of actualFunc
// against expectedFunc. It returns an error unless everything matches.
func InputAndOutputSig(actualFunc, expectedFunc reflect.Type) error {
	if err := InputSig(actualFunc, expectedFunc); err != nil {
		return err
	}

	if err := outputSig(actualFunc, expectedFunc); err != nil {
		return err
	}

	return nil
}

// InputSig compares the argument signatures of actualFunc
// against expectedFunc. It returns an error unless everything matches.
func InputSig(actualFunc, expectedFunc reflect.Type) error {
	// check number of arguments and type of each argument
	if actualFunc.NumIn() != expectedFunc.NumIn() {
		return fmt.Errorf(
			"expected function to have %d arguments not %d",
			expectedFunc.NumIn(), actualFunc.NumIn())
	}

	lastIdx := expectedFunc.NumIn()

	// If the function has a variadic argument validate that one first so that
	// we aren't checking for it while we iterate over the other args
	if expectedFunc.IsVariadic() {
		if ok := variadicArg(lastIdx, actualFunc, expectedFunc); !ok {
			i := lastIdx - 1
			return fmt.Errorf(
				"expected function to have"+
					" arg of type %v at position %d"+
					" not type %v",
				expectedFunc.In(i), i, actualFunc.In(i),
			)
		}

		lastIdx--
	}

	for i := 0; i < lastIdx; i++ {
		expectedArg := expectedFunc.In(i)
		actualArg := actualFunc.In(i)

		if err := arg(actualArg, expectedArg); err != nil {
			return fmt.Errorf("input argument at %d: %s", i, err)
		}
	}

	return nil
}

func outputSig(actualFunc, expectedFunc reflect.Type) error {
	// check number of return vals and type of each val
	if actualFunc.NumOut() != expectedFunc.NumOut() {
		return fmt.Errorf(
			"expected function to have %d return vals not %d",
			expectedFunc.NumOut(), actualFunc.NumOut())
	}

	for i := 0; i < expectedFunc.NumOut(); i++ {
		expectedArg := expectedFunc.Out(i)
		actualArg := actualFunc.Out(i)

		if err := arg(actualArg, expectedArg); err != nil {
			return fmt.Errorf("return argument at %d: %s", i, err)
		}
	}

	return nil
}

func variadicArg(lastIdx int, actualFunc, expectedFunc reflect.Type) bool {
	if actualFunc.In(lastIdx-1) != expectedFunc.In(lastIdx-1) {
		if actualFunc.In(lastIdx-1).Kind() != reflect.Slice {
			return false
		}

		expectedArgT := expectedFunc.In(lastIdx - 1)
		expectedElem := expectedArgT.Elem()
		if expectedElem.Kind() != reflect.Interface {
			return false
		}

		actualArgT := actualFunc.In(lastIdx - 1)
		actualElem := actualArgT.Elem()

		if ok := actualElem.ConvertibleTo(expectedElem); !ok {
			return false
		}

	}

	return true
}

func interfaceArg(actualArg, expectedArg reflect.Type) error {
	if !actualArg.ConvertibleTo(expectedArg) {
		return fmt.Errorf(
			"expected arg convertible to type %v not type %v",
			expectedArg, actualArg,
		)
	}

	return nil
}

func mapArg(actualArg, expectedArg reflect.Type) error {
	expectedKey := expectedArg.Key()
	actualKey := actualArg.Key()

	switch expectedKey.Kind() {
	case reflect.Interface:
		if err := interfaceArg(actualKey, expectedKey); err != nil {
			return fmt.Errorf("map key: %s", err)
		}
	default:
		if actualKey != expectedKey {
			return fmt.Errorf("expected map key of type %v not type %v",
				expectedKey, actualKey)
		}
	}

	expectedElem := expectedArg.Elem()
	actualElem := actualArg.Elem()

	switch expectedElem.Kind() {
	case reflect.Interface:
		if err := interfaceArg(actualElem, expectedElem); err != nil {
			return fmt.Errorf("map element: %s", err)
		}
	default:
		if actualElem != expectedElem {
			return fmt.Errorf("expected map element of type %v not type %v",
				expectedElem, actualElem)
		}
	}

	return nil
}

func arg(actualArg, expectedArg reflect.Type) error {
	switch expectedArg.Kind() {
	// If the expected arg is an interface we only care if the actual arg is convertible
	// to that interface
	case reflect.Interface:
		if err := interfaceArg(actualArg, expectedArg); err != nil {
			return err
		}
	default:
		// If the expected arg is not an interface then first check to see if
		// the actual arg is even the same reflect.Kind
		if expectedArg.Kind() != actualArg.Kind() {
			return fmt.Errorf("expected arg of kind %v not %v",
				expectedArg.Kind(), actualArg.Kind())
		}

		switch expectedArg.Kind() {
		// If the expected arg is a map then we need to handle the case where
		// the map key or element type is an interface
		case reflect.Map:
			if err := mapArg(actualArg, expectedArg); err != nil {
				return err
			}
		default:
			if actualArg != expectedArg {
				return fmt.Errorf(
					"Expected arg of type %v not type %v",
					expectedArg, actualArg,
				)
			}
		}
	}

	return nil
}
