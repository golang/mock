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
	doFunc *reflect.FuncValue
}

func (c *Call) AnyTimes() *Call {
	c.minCalls, c.maxCalls = 0, 1e8 // close enough to infinity
	return c
}

// Do declares the action to run when the call is matched.
// It takes an interface{} argument to support n-arity functions.
func (c *Call) Do(f interface{}) *Call {
	// TODO: Check arity and types here, rather than dying badly elsewhere.
	c.doFunc = reflect.NewValue(f).(*reflect.FuncValue)
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
