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
	"strings"
)

// Call represents an expected call to a mock.
type Call struct {
	receiver interface{}   // the receiver of the method call
	method   string        // the name of the method
	args     []Matcher     // the args
	rets     []interface{} // the return values (if any)

	preReqs []*Call // prerequisite calls

	// Expectations
	minCalls, maxCalls int

	numCalls int // actual number made

	// Actions
	doFunc  reflect.Value
	setArgs map[int]reflect.Value
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

// SetArg declares an action that will set the nth argument's value,
// indirected through a pointer.
func (c *Call) SetArg(n int, value interface{}) *Call {
	if c.setArgs == nil {
		c.setArgs = make(map[int]reflect.Value)
	}
	// TODO: Check func arity and value type.
	c.setArgs[n] = reflect.ValueOf(value)
	return c
}

// isPreReq returns true if other is a direct or indirect prerequisite to c.
func (c *Call) isPreReq(other *Call) bool {
	for _, preReq := range c.preReqs {
		if other == preReq || preReq.isPreReq(other) {
			return true
		}
	}
	return false
}

// After declares that the call may only match after preReq has been exhausted.
func (c *Call) After(preReq *Call) *Call {
	if preReq.isPreReq(c) {
		msg := fmt.Sprintf(
			"Loop in call order: %v is a prerequisite to %v (possibly indirectly).",
			c, preReq,
		)
		panic(msg)
	}

	c.preReqs = append(c.preReqs, preReq)
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
	args := make([]string, len(c.args))
	for i, arg := range c.args {
		args[i] = arg.String()
	}
	arguments := strings.Join(args, ", ")
	return fmt.Sprintf("%T.%v(%s)", c.receiver, c.method, arguments)
}

// Tests if the given call matches the expected call.
func (c *Call) matches(args []interface{}) bool {
	if len(args) != len(c.args) {
		return false
	}
	for i, m := range c.args {
		if !m.Matches(args[i]) {
			return false
		}
	}

	// Check that all prerequisite calls have been satisfied.
	for _, preReqCall := range c.preReqs {
		if !preReqCall.satisfied() {
			return false
		}
	}

	return true
}

// dropPrereqs tells the expected Call to not re-check prerequite calls any
// longer, and to return its current set.
func (c *Call) dropPrereqs() (preReqs []*Call) {
	preReqs = c.preReqs
	c.preReqs = nil
	return
}

func (c *Call) call(args []interface{}) []interface{} {
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
	for n, v := range c.setArgs {
		reflect.ValueOf(args[n]).Elem().Set(v)
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

// InOrder declares that the given calls should occur in order.
func InOrder(calls ...*Call) {
	for i := 1; i < len(calls); i++ {
		calls[i].After(calls[i-1])
	}
}
