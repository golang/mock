// Copyright 2011 Google Inc.
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

package gomock_test

import (
	"fmt"
	"reflect"
	"testing"

	"gomock.googlecode.com/hg/gomock"
)

type ErrorReporter struct {
	t          *testing.T
	log        []string
	failed     bool
	fatalToken struct{}
}

func NewErrorReporter(t *testing.T) *ErrorReporter {
	return &ErrorReporter{t: t}
}

func (e *ErrorReporter) reportLog() {
	for _, entry := range e.log {
		e.t.Log(entry)
	}
}

func (e *ErrorReporter) assertPass(msg string) {
	if e.failed {
		e.t.Errorf("Expected pass, but got failure(s): %s", msg)
		e.reportLog()
	}
}

func (e *ErrorReporter) assertFail(msg string) {
	if !e.failed {
		e.t.Error("Expected failure, but got pass: %s", msg)
	}
}

// Use to check that code triggers a fatal test failure.
func (e *ErrorReporter) assertFatal(fn func()) {
	defer func() {
		err := recover()
		if err == nil {
			var actual string
			if e.failed {
				actual = "non-fatal failure"
			} else {
				actual = "pass"
			}
			e.t.Error("Expected fatal failure, but got a", actual)
		} else if token, ok := err.(*struct{}); ok && token == &e.fatalToken {
			// This is okay - the panic is from FailNow().
			return
		} else {
			// Some other panic.
			panic(err)
		}
	}()

	fn()
}

// recoverUnexpectedFatal can be used as a deferred call in test cases to
// recover from and display a call to ErrorReporter.Fatalf().
func (e *ErrorReporter) recoverUnexpectedFatal() {
	err := recover()
	if err == nil {
		// No panic.
	} else if token, ok := err.(*struct{}); ok && token == &e.fatalToken {
		// Unexpected fatal error happened.
		e.t.Error("Got unexpected fatal error(s). All errors up to this point:")
		e.reportLog()
		return
	} else {
		// Some other panic.
		panic(err)
	}
}

func (e *ErrorReporter) Logf(format string, args ...interface{}) {
	e.log = append(e.log, fmt.Sprintf(format, args...))
}

func (e *ErrorReporter) Errorf(format string, args ...interface{}) {
	e.Logf(format, args...)
	e.failed = true
}

func (e *ErrorReporter) Fatalf(format string, args ...interface{}) {
	e.Logf(format, args...)
	e.failed = true
	panic(&e.fatalToken)
}

// A type purely for use as a receiver in testing the Controller.
type Subject struct{}

func (s *Subject) FooMethod(arg string) int {
	return 0
}

func (s *Subject) BarMethod(arg string) int {
	return 0
}

func assertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected %+v, but got %+v", expected, actual)
	}
}

func createFixtures(t *testing.T) (reporter *ErrorReporter, ctrl *gomock.Controller) {
	// reporter acts as a testing.T-like object that we pass to the
	// Controller. We use it to test that the mock considered tests
	// successful or failed.
	reporter = NewErrorReporter(t)
	ctrl = gomock.NewController(reporter)
	return
}

func TestNoCalls(t *testing.T) {
	reporter, ctrl := createFixtures(t)
	ctrl.Finish()
	reporter.assertPass("No calls expected or made.")
}

func TestExpectedMethodCall(t *testing.T) {
	reporter, ctrl := createFixtures(t)
	subject := new(Subject)

	ctrl.RecordCall(subject, "FooMethod", "argument")
	ctrl.Call(subject, "FooMethod", "argument")
	ctrl.Finish()

	reporter.assertPass("Expected method call made.")
}

func TestUnexpectedMethodCall(t *testing.T) {
	reporter, ctrl := createFixtures(t)
	subject := new(Subject)

	reporter.assertFatal(func() {
		ctrl.Call(subject, "FooMethod", "argument")
	})

	ctrl.Finish()
}

func TestRepeatedCall(t *testing.T) {
	reporter, ctrl := createFixtures(t)
	subject := new(Subject)

	ctrl.RecordCall(subject, "FooMethod", "argument").Times(3)
	ctrl.Call(subject, "FooMethod", "argument")
	ctrl.Call(subject, "FooMethod", "argument")
	ctrl.Call(subject, "FooMethod", "argument")
	reporter.assertPass("After expected repeated method calls.")
	reporter.assertFatal(func() {
		ctrl.Call(subject, "FooMethod", "argument")
	})
	ctrl.Finish()
	reporter.assertFail("After calling one too many times.")
}

func TestUnexpectedArgCount(t *testing.T) {
	reporter, ctrl := createFixtures(t)
	defer reporter.recoverUnexpectedFatal()
	subject := new(Subject)

	ctrl.RecordCall(subject, "FooMethod", "argument")
	reporter.assertFatal(func() {
		// This call is made with the wrong number of arguments...
		ctrl.Call(subject, "FooMethod", "argument", "extra_argument")
	})
	reporter.assertFatal(func() {
		// ... so is this.
		ctrl.Call(subject, "FooMethod")
	})
	reporter.assertFatal(func() {
		// The expected call wasn't made.
		ctrl.Finish()
	})
}

func TestAnyTimes(t *testing.T) {
	reporter, ctrl := createFixtures(t)
	subject := new(Subject)

	ctrl.RecordCall(subject, "FooMethod", "argument").AnyTimes()
	for i := 0; i < 100; i++ {
		ctrl.Call(subject, "FooMethod", "argument")
	}
	reporter.assertPass("After 100 method calls.")
	ctrl.Finish()
}

func TestDo(t *testing.T) {
	_, ctrl := createFixtures(t)
	subject := new(Subject)

	doCalled := false
	var argument string
	ctrl.RecordCall(subject, "FooMethod", "argument").Do(
		func(arg string) {
			doCalled = true
			argument = arg
		})
	if doCalled {
		t.Error("Do() callback called too early.")
	}

	ctrl.Call(subject, "FooMethod", "argument")

	if !doCalled {
		t.Error("Do() callback not called.")
	}
	if "argument" != argument {
		t.Error("Do callback received wrong argument.")
	}

	ctrl.Finish()
}

func TestReturn(t *testing.T) {
	_, ctrl := createFixtures(t)
	subject := new(Subject)

	// Unspecified return should produce "zero" result.
	ctrl.RecordCall(subject, "FooMethod", "zero")
	ctrl.RecordCall(subject, "FooMethod", "five").Return(5)

	assertEqual(
		t,
		[]interface{}{0},
		ctrl.Call(subject, "FooMethod", "zero"))

	assertEqual(
		t,
		[]interface{}{5},
		ctrl.Call(subject, "FooMethod", "five"))
	ctrl.Finish()
}

func TestUnorderedCalls(t *testing.T) {
	reporter, ctrl := createFixtures(t)
	defer reporter.recoverUnexpectedFatal()
	subjectTwo := new(Subject)
	subjectOne := new(Subject)

	ctrl.RecordCall(subjectOne, "FooMethod", "1")
	ctrl.RecordCall(subjectOne, "BarMethod", "2")
	ctrl.RecordCall(subjectTwo, "FooMethod", "3")
	ctrl.RecordCall(subjectTwo, "BarMethod", "4")

	// Make the calls in a different order, which should be fine.
	ctrl.Call(subjectOne, "BarMethod", "2")
	ctrl.Call(subjectTwo, "FooMethod", "3")
	ctrl.Call(subjectTwo, "BarMethod", "4")
	ctrl.Call(subjectOne, "FooMethod", "1")

	reporter.assertPass("After making all calls in different order")

	ctrl.Finish()

	reporter.assertPass("After finish")
}
