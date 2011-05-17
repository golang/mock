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
//	- Support strict mocks (calls in a particular order).
//	- Handle different argument/return types (e.g. ..., chan, map, interface).
package gomock

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
	expectedCalls callSet
}

func NewController(t TestReporter) *Controller {
	return &Controller{
		t:             t,
		expectedCalls: make(callSet),
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

	ctrl.expectedCalls.Add(call)
	return call
}

func (ctrl *Controller) Call(receiver interface{}, method string, args ...interface{}) []interface{} {
	expected := ctrl.expectedCalls.FindMatch(receiver, method, args)
	if expected == nil {
		ctrl.t.Fatalf("no matching expected call: %T.%v", receiver, method)
	}

	rets := expected.call(args)
	if expected.exhausted() {
		ctrl.expectedCalls.Remove(expected)
	}

	return rets
}

func (ctrl *Controller) Finish() {
	// Check that all remaining expected calls are satisfied.
	failures := false
	for _, methodMap := range ctrl.expectedCalls {
		for _, calls := range methodMap {
			for _, call := range calls {
				if !call.satisfied() {
					ctrl.t.Errorf("missing call(s) to %v", call)
					failures = true
				}
			}
		}
	}
	if failures {
		ctrl.t.Fatalf("aborting test due to missing call(s)")
	}
}
