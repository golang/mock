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
//   (2) Use mockgen to generate a mock from the interface.
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
// By default, expected calls are not enforced to run in any particular order.
// Call order dependency can be enforced by use of InOrder and/or Call.After.
// Call.After can create more varied call order dependencies, but InOrder is
// often more convenient.
//
// The following examples create equivalent call order dependencies.
//
// Example of using Call.After to chain expected call order:
//
//     firstCall := mockObj.EXPECT().SomeMethod(1, "first")
//     secondCall := mockObj.EXPECT().SomeMethod(2, "second").After(firstCall)
//     mockObj.EXPECT().SomeMethod(3, "third").After(secondCall)
//
// Example of using InOrder to declare expected call order:
//
//     gomock.InOrder(
//         mockObj.EXPECT().SomeMethod(1, "first"),
//         mockObj.EXPECT().SomeMethod(2, "second"),
//         mockObj.EXPECT().SomeMethod(3, "third"),
//     )
//
// Using custom mock execution.
// Suppose the defined interface:
//    type MyInterface interface {
//        ChangeString(str *string)
//        GetString() string
//        DoSomething() error
//    }
//
// Generate mock using mockgen.
//
// Use the mock in a test:
//    func TestMyInterface1(t *testing.T){
//        mockCtrl := gomock.NewController(t)
//        defer mockCtrl.Finish()
//
//        hello := "hello"
//        mockObj := something.NewMockMyInterface(mockCtrl)
//        // override ChangeString execution
//        mockObj.EXPECT().ChangeString(hello).Do(func(str *string){
//             str += " world"
//        })
//
//        mockObj.ChangeString(hello)
//        fmt.Println(hello) // print hello world
//    }
//
//    func TestMyInterface2(t *testing.T){
//        mockCtrl := gomock.NewController(t)
//        defer mockCtrl.Finish()
//
//        mockObj := something.NewMockMyInterface(mockCtrl)
//        // override GetString execution
//        mockObj.EXPECT().GetString().Return("hello world")
//
//        hello := mockObj.GetString()
//        fmt.Println(hello) // print hello world
//    }
//
//    func TestMyInterface3(t *testing.T){
//        mockCtrl := gomock.NewController(t)
//        defer mockCtrl.Finish()
//
//        mockObj := something.NewMockMyInterface(mockCtrl)
//        // override GetString execution
//        mockObj.EXPECT().DoSomething().DoAndReturn(func()error{
//             fmt.Println("Doing something")
//             return fmt.Errorf("Error on doing something")
//        })
//
//        err := mockObj.DoSomething() // print Doing something
//        fmt.Println(err) // print Error on doing something
//    }
//
// In some cases we can't be sure the exact param would be passed into the function we want to mock.
// For this case we need to implement our own logic to check if the mock is valid.
// Suppose we have code like this:
//    type Metric interface {
//        PushMetric(duration float64)
//    }
//
//    var metrics Metric
//
//    func DoSomething(){
//        start := time.Now()
//        // do something
//        duration := time.Since(start)
//        metrics.PushMetric(duration.Seconds())
//    }
//
// And here's the test.
//
//    // Since we cant be sure how much duration will be passed into PushMetric
//    // we create our own Matcher
//    type durationMatcher struct{}
//
//    func (durationMatcher) Matches(x interface{}) bool {
//        dur, ok := x.(float64)
//        return ok && dur > float64(0)
//    }
//
//    func (durationMatcher) String() string {
//        return "is float64 and more than 0"
//    }
//
//    func TestMyInterface(t *testing.T){
//        mockCtrl := gomock.NewController(t)
//        defer mockCtrl.Finish()
//        mockObj := something.NewMockMetric(mockCtrl)
//        mockObj.EXPECT().PushMetric(durationMatcher{})
//
//        metrics = mockObj
//        DoSomething()
//    }
package gomock

import (
	"fmt"
	"golang.org/x/net/context"
	"reflect"
	"runtime"
	"sync"
)

// A TestReporter is something that can be used to report test failures.
// It is satisfied by the standard library's *testing.T.
type TestReporter interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// A Controller represents the top-level control of a mock ecosystem.
// It defines the scope and lifetime of mock objects, as well as their expectations.
// It is safe to call Controller's methods from multiple goroutines.
type Controller struct {
	mu            sync.Mutex
	t             TestReporter
	expectedCalls *callSet
	finished      bool
}

func NewController(t TestReporter) *Controller {
	return &Controller{
		t:             t,
		expectedCalls: newCallSet(),
	}
}

type cancelReporter struct {
	t      TestReporter
	cancel func()
}

func (r *cancelReporter) Errorf(format string, args ...interface{}) { r.t.Errorf(format, args...) }
func (r *cancelReporter) Fatalf(format string, args ...interface{}) {
	defer r.cancel()
	r.t.Fatalf(format, args...)
}

// WithContext returns a new Controller and a Context, which is cancelled on any
// fatal failure.
func WithContext(ctx context.Context, t TestReporter) (*Controller, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return NewController(&cancelReporter{t, cancel}), ctx
}

func (ctrl *Controller) RecordCall(receiver interface{}, method string, args ...interface{}) *Call {
	if h, ok := ctrl.t.(testHelper); ok {
		h.Helper()
	}

	recv := reflect.ValueOf(receiver)
	for i := 0; i < recv.Type().NumMethod(); i++ {
		if recv.Type().Method(i).Name == method {
			return ctrl.RecordCallWithMethodType(receiver, method, recv.Method(i).Type(), args...)
		}
	}
	ctrl.t.Fatalf("gomock: failed finding method %s on %T", method, receiver)
	panic("unreachable")
}

func (ctrl *Controller) RecordCallWithMethodType(receiver interface{}, method string, methodType reflect.Type, args ...interface{}) *Call {
	if h, ok := ctrl.t.(testHelper); ok {
		h.Helper()
	}

	call := newCall(ctrl.t, receiver, method, methodType, args...)

	ctrl.mu.Lock()
	defer ctrl.mu.Unlock()
	ctrl.expectedCalls.Add(call)

	return call
}

func (ctrl *Controller) Call(receiver interface{}, method string, args ...interface{}) []interface{} {
	if h, ok := ctrl.t.(testHelper); ok {
		h.Helper()
	}

	// Nest this code so we can use defer to make sure the lock is released.
	actions := func() []func([]interface{}) []interface{} {
		ctrl.mu.Lock()
		defer ctrl.mu.Unlock()

		expected, err := ctrl.expectedCalls.FindMatch(receiver, method, args)
		if err != nil {
			origin := callerInfo(2)
			ctrl.t.Fatalf("Unexpected call to %T.%v(%v) at %s because: %s", receiver, method, args, origin, err)
		}

		// Two things happen here:
		// * the matching call no longer needs to check prerequite calls,
		// * and the prerequite calls are no longer expected, so remove them.
		preReqCalls := expected.dropPrereqs()
		for _, preReqCall := range preReqCalls {
			ctrl.expectedCalls.Remove(preReqCall)
		}

		actions := expected.call(args)
		if expected.exhausted() {
			ctrl.expectedCalls.Remove(expected)
		}
		return actions
	}()

	var rets []interface{}
	for _, action := range actions {
		if r := action(args); r != nil {
			rets = r
		}
	}

	return rets
}

func (ctrl *Controller) Finish() {
	if h, ok := ctrl.t.(testHelper); ok {
		h.Helper()
	}

	ctrl.mu.Lock()
	defer ctrl.mu.Unlock()

	if ctrl.finished {
		ctrl.t.Fatalf("Controller.Finish was called more than once. It has to be called exactly once.")
	}
	ctrl.finished = true

	// If we're currently panicking, probably because this is a deferred call,
	// pass through the panic.
	if err := recover(); err != nil {
		panic(err)
	}

	// Check that all remaining expected calls are satisfied.
	failures := ctrl.expectedCalls.Failures()
	for _, call := range failures {
		ctrl.t.Errorf("missing call(s) to %v", call)
	}
	if len(failures) != 0 {
		ctrl.t.Fatalf("aborting test due to missing call(s)")
	}
}

func callerInfo(skip int) string {
	if _, file, line, ok := runtime.Caller(skip + 1); ok {
		return fmt.Sprintf("%s:%d", file, line)
	}
	return "unknown file"
}

type testHelper interface {
	TestReporter
	Helper()
}
