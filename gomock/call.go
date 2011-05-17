// Copyright 2010 Google Inc.
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

package gomock

import (
	"fmt"
	"reflect"
)

// Call represents an expected call to a mock.
type Call struct {
	receiver interface{}   // the receiver of the method call
	method   string        // the name of the method
	args     []Matcher     // the args
	rets     []interface{} // the return values (if any)

	// Expectations
	minCalls, maxCalls int

	numCalls int // actual number made

	// Actions
	doFunc reflect.Value
}

func (c *Call) AnyTimes() *Call {
	c.minCalls, c.maxCalls = 0, 1e8 // close enough to infinity
	return c
}

// Do declares the action to run when the call is matched.
// It takes an interface{} argument to support n-arity functions.
func (c *Call) Do(f interface{}) *Call {
	// TODO: Check arity and types here, rather than dying badly elsewhere.
	c.doFunc = reflect.ValueOf(f)
	return c
}

func (c *Call) Return(rets ...interface{}) *Call {
	// TODO: Check return-arity and types here, rather than dying badly elsewhere.
	c.rets = rets
	return c
}

func (c *Call) Times(n int) *Call {
	c.minCalls, c.maxCalls = n, n
	return c
}

// Returns true iff the minimum number of calls have been made.
func (c *Call) satisfied() bool {
	return c.numCalls >= c.minCalls
}

// Returns true iff the maximum number of calls have been made.
func (c *Call) exhausted() bool {
	return c.numCalls >= c.maxCalls
}

func (c *Call) String() string {
	return fmt.Sprintf("%T.%v", c.receiver, c.method)
}

// Tests if the given call matches the expected call.
func (c *Call) matches(receiver interface{}, method string, args ...interface{}) (match bool, failure string) {
	if receiver != c.receiver || method != c.method {
		return false, fmt.Sprintf("got a %T.%v method call, expected %v", receiver, method, c)
	}
	if len(args) != len(c.args) {
		return false, fmt.Sprintf("got %d args to %v, expected %d args", len(args), c, len(c.args))
	}
	for i, m := range c.args {
		if !m.Matches(args[i]) {
			// TODO: Tune this error message.
			return false, fmt.Sprintf("arg #%d to %v was %v, expected: %v", i, c, args[i], m)
		}
	}

	return true, ""
}

func (c *Call) call(args ...interface{}) []interface{} {
	c.numCalls++

	// Actions
	if c.doFunc.IsValid() {
		doArgs := make([]reflect.Value, len(args))
		ft := c.doFunc.Type()
		for i := 0; i < ft.NumIn(); i++ {
			doArgs[i] = reflect.ValueOf(args[i])
		}
		c.doFunc.Call(doArgs)
	}

	rets := c.rets
	if rets == nil {
		// Synthesize the zero value for each of the return args' types.
		recv := reflect.ValueOf(c.receiver)
		var mt reflect.Type
		for i := 0; i < recv.Type().NumMethod(); i++ {
			if recv.Type().Method(i).Name == c.method {
				mt = recv.Method(i).Type()
				break
			}
		}
		rets = make([]interface{}, mt.NumOut())
		for i := 0; i < mt.NumOut(); i++ {
			rets[i] = reflect.Zero(mt.Out(i)).Interface()
		}
	}

	return rets
}
