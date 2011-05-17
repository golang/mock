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

// GoMock - a mock framework for Go.
//
// Standard usage:
//   (1) Define an interface that you wish to mock.
//         type MyInterface interface {
//           SomeMethod(x int64, y string)
//         }
//   (2) Use mockgen to automatically generate a mock from the interface.
//   (3) Use the mock in a test:
//         func TestMyThing(t *testing.T) {
//           mockCtrl := gomock.NewController(t)
//           defer mockCtrl.Finish()
//
//           mockObj := something.NewMockMyInterface(mockCtrl)
//           mockObj.EXPECT().SomeMethod(4, "blah")
//           // pass mockObj to a real object and play with it.
//         }
//
// TODO:
//	- Support loose mocks (calls in any order).
//	- Handle different argument/return types (e.g. ..., chan, map, interface).
package gomock

import (
	"container/list"
	"fmt"
	"reflect"
)

// A TestReporter is something that can be used to report test failures.
// It is satisfied by the standard library's *testing.T.
type TestReporter interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// A Controller represents the top-level control of a mock ecosystem.
// It defines the scope and lifetime of mock objects, as well as their expectations.
type Controller struct {
	t             TestReporter
	expectedCalls *list.List
}

func NewController(t TestReporter) *Controller {
	return &Controller{
		t:             t,
		expectedCalls: list.New(),
	}
}

func (ctrl *Controller) RecordCall(receiver interface{}, method string, args ...interface{}) *Call {
	margs := make([]Matcher, len(args))
	for i, arg := range args {
		if m, ok := arg.(Matcher); ok {
			margs[i] = m
		} else {
			margs[i] = Eq(arg)
		}
	}

	call := &Call{receiver: receiver, method: method, args: margs, minCalls: 1, maxCalls: 1}
	ctrl.expectedCalls.PushBack(call)
	return call
}

func matches(expected *Call, receiver interface{}, method string, args ...interface{}) (match bool, failure string) {
	callStr := fmt.Sprintf("%T.%v", expected.receiver, expected.method)

	if receiver != expected.receiver || method != expected.method {
		return false, fmt.Sprintf("got a %T.%v method call, expected %v", receiver, method, callStr)
	}
	if len(args) != len(expected.args) {
		return false, fmt.Sprintf("got %d args to %v, expected %d args", len(args), callStr, len(expected.args))
	}
	for i, m := range expected.args {
		if !m.Matches(args[i]) {
			// TODO: Tune this error message.
			return false, fmt.Sprintf("arg #%d to %v was %v, expected: %v", i, callStr, args[i], m)
		}
	}

	return true, ""
}

func (ctrl *Controller) Call(receiver interface{}, method string, args ...interface{}) []interface{} {
	var expected *Call

	e := ctrl.expectedCalls.Front()
	for e != nil {
		expected = e.Value.(*Call)
		mustMatch := expected.numCalls < expected.minCalls
		if ok, msg := matches(expected, receiver, method, args...); !ok {
			if mustMatch {
				ctrl.t.Fatalf("%s", msg)
			}
			// discard and advance
			ne := e.Next()
			ctrl.expectedCalls.Remove(e)
			e = ne
			continue
		}
		// match!
		break
	}
	if e == nil {
		ctrl.t.Fatalf("unexpected %T.%v method call; no more expected", receiver, method)
	}

	expected.numCalls++
	if expected.numCalls >= expected.maxCalls {
		ctrl.expectedCalls.Remove(ctrl.expectedCalls.Front())
	}

	// Actions
	if expected.doFunc != nil {
		doArgs := make([]reflect.Value, len(args))
		ft := expected.doFunc.Type().(*reflect.FuncType)
		for i := 0; i < ft.NumIn(); i++ {
			doArgs[i] = reflect.MakeZero(ft.In(i))
			doArgs[i].SetValue(reflect.NewValue(args[i])) // assignment compatibility assumed
		}
		expected.doFunc.Call(doArgs)
	}

	rets := expected.rets
	if rets == nil {
		// Synthesize the zero value for each of the return args' types.
		recv := reflect.NewValue(receiver)
		var mt *reflect.FuncType
		for i := 0; i < recv.Type().NumMethod(); i++ {
			if recv.Type().Method(i).Name == method {
				mt = recv.Method(i).Type().(*reflect.FuncType)
				break
			}
		}
		rets = make([]interface{}, mt.NumOut())
		for i := 0; i < mt.NumOut(); i++ {
			rets[i] = reflect.MakeZero(mt.Out(i)).Interface()
		}
	}

	return rets
}

func (ctrl *Controller) Finish() {
	// Check that all remaining expected calls are satisfied.
	failures := false
	for e := ctrl.expectedCalls.Front(); e != nil; e = e.Next() {
		exp := e.Value.(*Call)
		if exp.numCalls <= exp.minCalls {
			ctrl.t.Errorf("missing call(s) to %T.%v", exp.receiver, exp.method)
			failures = true
		}
	}
	if failures {
		ctrl.t.Fatalf("aborting test due to missing call(s)")
	}
}
