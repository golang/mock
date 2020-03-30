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

package gomock_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golang/mock/gomock/internal/mock_gomock"
)

const (
	intArg = 45982
	item1  = "item1"
	item2  = "item2"
)

var sliceArg = []string{item1, item2}

func TestLastValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMatcher := mock_gomock.NewMockMatcher(ctrl)
	captor := gomock.Captor(mockMatcher)

	mockMatcher.EXPECT().Matches(intArg).Times(1)
	mockMatcher.EXPECT().Matches(sliceArg).Times(1)

	captor.Matches(intArg)
	captor.Matches(sliceArg)

	actualValue := captor.LastValue().([]string)

	if len(sliceArg) != len(actualValue) {
		t.Errorf("expected length %d, but was %d", len(sliceArg), len(actualValue))
	}
	if item1 != actualValue[0] {
		t.Errorf("expected %s, but was %s", item1, actualValue[0])
	}
	if item2 != actualValue[1] {
		t.Errorf("expected %s, but was %s", item2, actualValue[1])
	}
}

func TestLastValueWithNoElements(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMatcher := mock_gomock.NewMockMatcher(ctrl)
	captor := gomock.Captor(mockMatcher)

	mockMatcher.EXPECT().Matches(gomock.Any()).Times(0)

	actualValue := captor.LastValue()

	if actualValue != nil {
		t.Errorf("expected nil, but was %s", actualValue)
	}
}

func TestValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMatcher := mock_gomock.NewMockMatcher(ctrl)
	captor := gomock.Captor(mockMatcher)

	mockMatcher.EXPECT().Matches(intArg).Times(1)
	mockMatcher.EXPECT().Matches(sliceArg).Times(1)

	captor.Matches(intArg)
	captor.Matches(sliceArg)

	actualValues := captor.Values()

	if len(actualValues) != 2 {
		t.Errorf("expected 2 values, but got %s", actualValues)
	}

	actualVal1 := actualValues[0].(int)
	if intArg != actualVal1 {
		t.Errorf("expected %d, but got %d", intArg, actualVal1)
	}

	actualVal2 := actualValues[1].([]string)
	if len(sliceArg) != len(actualVal2) {
		t.Errorf("expected length %d, but was %d", len(sliceArg), len(actualVal2))
	}
	if item1 != actualVal2[0] {
		t.Errorf("expected %s, but was %s", item1, actualVal2[0])
	}
	if item2 != actualVal2[1] {
		t.Errorf("expected %s, but was %s", item2, actualVal2[1])
	}
}

func TestValuesWithNoElements(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMatcher := mock_gomock.NewMockMatcher(ctrl)
	captor := gomock.Captor(mockMatcher)

	mockMatcher.EXPECT().Matches(gomock.Any()).Times(0)

	actualValue := captor.Values()

	if len(actualValue) != 0 {
		t.Errorf("expected length 0, but slice had elements: %s", actualValue)
	}
}

func TestAnyCaptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	anyCaptor := gomock.AnyCaptor()

	if !anyCaptor.Matches(192394) {
		t.Errorf("expected Matches to be true for %d", 192394)
	}
	if !anyCaptor.Matches([]interface{}{make([]int, 3)}) {
		t.Errorf("expected Matches to be true for %s", []interface{}{make([]int, 3)})
	}
	if fmt.Sprintf("%s", gomock.Any()) != fmt.Sprintf("%s", anyCaptor) {
		t.Errorf("expected string representation to be '%s', but was '%s'", gomock.Any(), anyCaptor)
	}
}
