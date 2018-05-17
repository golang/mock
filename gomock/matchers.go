//go:generate mockgen -destination mock_matcher/mock_matcher.go github.com/golang/mock/gomock Matcher

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

// A Matcher is a representation of a class of values.
// It is used to represent the valid or expected arguments to a mocked method.
type Matcher interface {
	// Matches returns whether x is a match.
	Matches(x interface{}) bool

	// String describes what the matcher matches.
	String() string
}

type anyMatcher struct{}

func (anyMatcher) Matches(x interface{}) bool {
	return true
}

func (anyMatcher) String() string {
	return "is anything"
}

type eqMatcher struct {
	x interface{}
}

func (e eqMatcher) Matches(x interface{}) bool {
	return reflect.DeepEqual(e.x, x)
}

func (e eqMatcher) String() string {
	return fmt.Sprintf("is equal to %v", e.x)
}

type nilMatcher struct{}

func (nilMatcher) Matches(x interface{}) bool {
	if x == nil {
		return true
	}

	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}

	return false
}

func (nilMatcher) String() string {
	return "is nil"
}

type notMatcher struct {
	m Matcher
}

func (n notMatcher) Matches(x interface{}) bool {
	return !n.m.Matches(x)
}

func (n notMatcher) String() string {
	// TODO: Improve this if we add a NotString method to the Matcher interface.
	return "not(" + n.m.String() + ")"
}

type assignableToTypeOfMatcher struct {
	targetType reflect.Type
}

func (m assignableToTypeOfMatcher) Matches(x interface{}) bool {
	return reflect.TypeOf(x).AssignableTo(m.targetType)
}

func (m assignableToTypeOfMatcher) String() string {
	return "is assignable to " + m.targetType.Name()
}

// Constructors
func Any() Matcher             { return anyMatcher{} }
func Eq(x interface{}) Matcher { return eqMatcher{x} }
func Nil() Matcher             { return nilMatcher{} }
func Not(x interface{}) Matcher {
	if m, ok := x.(Matcher); ok {
		return notMatcher{m}
	}
	return notMatcher{Eq(x)}
}

// AssignableToTypeOf is a Matcher that matches if the parameter to the mock
// function is assignable to the type of the parameter to this function.
//
// Example usage:
//
// 		dbMock.EXPECT().
// 			Insert(gomock.AssignableToTypeOf(&EmployeeRecord{})).
// 			Return(errors.New("DB error"))
//
func AssignableToTypeOf(x interface{}) Matcher {
	return assignableToTypeOfMatcher{reflect.TypeOf(x)}
}

type compositeMatcher struct {
	Matchers []Matcher
}

func (m compositeMatcher) Matches(x interface{}) bool {
	for _, _m := range m.Matchers {
		if !_m.Matches(x) {
			return false
		}
	}
	return true
}

func (m compositeMatcher) String() string {
	ss := make([]string, 0, len(m.Matchers))
	for _, matcher := range m.Matchers {
		ss = append(ss, matcher.String())
	}
	return strings.Join(ss, "; ")
}

// All returns a composite Matcher that returns true if and only if all element
// Matchers are satisfied for a given parameter to the mock function.
func All(ms ...Matcher) Matcher {
	return compositeMatcher{ms}
}

type lenMatcher struct {
	n int
}

func (m lenMatcher) Matches(x interface{}) bool {
	return reflect.TypeOf(x).Kind() == reflect.Slice &&
		m.n == reflect.ValueOf(x).Cap()
}

func (m lenMatcher) String() string {
	return fmt.Sprintf("is of length %v", m.n)
}

// Len returns a Matcher that matches on length. The returned Matcher returns
// false if its input is not a slice.
func Len(l int) Matcher {
	return lenMatcher{l}
}

// Contains returns a Matcher that asserts that its input contains the requisite
// elements, in any order. The returned Matcher returns false if its input is
// not a slice. It is useful for variadic parameters where you only care that
// specific values are represented, not their specific ordering.
func Contains(vs ...interface{}) Matcher {
	return containsElementsMatcher{vs}
}

type containsElementsMatcher struct {
	Elements []interface{}
}

func (m containsElementsMatcher) Matches(x interface{}) bool {
	if reflect.TypeOf(x).Kind() != reflect.Slice {
		return false
	}

	xv := reflect.ValueOf(x)

	eq := true
Elements:
	for i := 0; i < len(m.Elements) && eq; i++ {
		e := m.Elements[i]
		for j := 0; j < xv.Cap(); j++ {
			if reflect.DeepEqual(xv.Index(j).Interface(), e) {
				continue Elements
			}
		}
		eq = false
	}
	return eq
}

func (m containsElementsMatcher) String() string {
	return fmt.Sprintf("contains elements: %v", m.Elements)
}
