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

package gomock_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	mock_matcher "github.com/golang/mock/gomock/mock_matcher"
	my_gomock "github.com/noibar/mock/gomock"
)

func TestMatchers(t *testing.T) {
	type e interface{}
	type testCase struct {
		matcher gomock.Matcher
		yes, no []e
	}
	tests := []testCase{
		testCase{gomock.Any(), []e{3, nil, "foo"}, nil},
		testCase{gomock.Eq(4), []e{4}, []e{3, "blah", nil, int64(4)}},
		testCase{my_gomock.Str(4), []e{4, "4", int64(4), float64(4.0)}, []e{3, float64(4.1)}},
		testCase{gomock.Nil(),
			[]e{nil, (error)(nil), (chan bool)(nil), (*int)(nil)},
			[]e{"", 0, make(chan bool), errors.New("err"), new(int)}},
		testCase{gomock.Not(gomock.Eq(4)), []e{3, "blah", nil, int64(4)}, []e{4}},
	}
	for i, test := range tests {
		for _, x := range test.yes {
			if !test.matcher.Matches(x) {
				t.Errorf(`test %d: "%v %s" should be true.`, i, x, test.matcher)
			}
		}
		for _, x := range test.no {
			if test.matcher.Matches(x) {
				t.Errorf(`test %d: "%v %s" should be false.`, i, x, test.matcher)
			}
		}
	}
}

// A more thorough test of notMatcher
func TestNotMatcher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMatcher := mock_matcher.NewMockMatcher(ctrl)
	notMatcher := gomock.Not(mockMatcher)

	mockMatcher.EXPECT().Matches(4).Return(true)
	if match := notMatcher.Matches(4); match {
		t.Errorf("notMatcher should not match 4")
	}

	mockMatcher.EXPECT().Matches(5).Return(false)
	if match := notMatcher.Matches(5); !match {
		t.Errorf("notMatcher should match 5")
	}
}

type testStruct struct {
	id                string
	interestingValues []string
}

func (t testStruct) String() string {
	return strings.Join(t.interestingValues, ";")
}

// A more thorough test of strMatcher
func TestStrMatcher(t *testing.T) {
	for _, test := range []struct {
		title    string
		first    testStruct
		second   testStruct
		expected bool
	}{
		{
			title:    "Empty structs",
			expected: true,
		}, {
			title: "Identical structs",
			first: testStruct{
				id:                "id",
				interestingValues: []string{"1", "2", "3"},
			},
			second: testStruct{
				id:                "id",
				interestingValues: []string{"1", "2", "3"},
			},
			expected: true,
		}, {
			title: "Identical interesting values",
			first: testStruct{
				id:                "first",
				interestingValues: []string{"1", "2", "3"},
			},
			second: testStruct{
				id:                "second",
				interestingValues: []string{"1", "2", "3"},
			},
			expected: true,
		}, {
			title: "Different interesting values",
			first: testStruct{
				id:                "first",
				interestingValues: []string{"1", "2", "3"},
			},
			second: testStruct{
				id:                "second",
				interestingValues: []string{"1", "2", "3", "4"},
			},
			expected: false,
		},
	} {
		matcher := my_gomock.Str(test.first)
		if match := matcher.Matches(test.second); match != test.expected {
			t.Errorf("Test %v failed. expected %v, recieved %v", test.title, test.expected, match)
		}
	}
}
