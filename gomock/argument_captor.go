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

package gomock

// ArgumentCaptor is a struct that composes a Matcher, and extends it by storing given arguments in the values slice.
type ArgumentCaptor struct {
	m      Matcher
	values []interface{}
}

// Matches method overrides the Matcher.Matches method
// First it appends any argument(s) used to the values slice.
// Then the parent Matches method is called.
func (ac *ArgumentCaptor) Matches(x interface{}) bool {
	ac.values = append(ac.values, x)
	return ac.m.Matches(x)
}

// String simply calls the String method for the composed Matcher
func (ac *ArgumentCaptor) String() string {
	return ac.m.String()
}

// LastValue returns the last argument the matcher was called with as an interface{}.
// If the matcher was never called, nil is returned.
func (ac *ArgumentCaptor) LastValue() interface{} {
	if len(ac.values) < 1 {
		return nil
	}
	return ac.values[len(ac.values)-1]
}

// Values returns the all arguments the matcher was called with as a []interface{}.
// The values are ordered from first called to last called.
func (ac *ArgumentCaptor) Values() []interface{} {
	return ac.values
}

// Captor is a helper method that returns a new *ArgumentCaptor struct with Matcher set to the given matcher m
func Captor(m Matcher) *ArgumentCaptor {
	return &ArgumentCaptor{m: m}
}

// AnyCaptor is a helper method that returns a new *ArgumentCaptor struct with the matcher set to an anyMatcher
func AnyCaptor() *ArgumentCaptor {
	return &ArgumentCaptor{m: Any()}
}
