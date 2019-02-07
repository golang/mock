package gomock_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	mock_matcher "github.com/golang/mock/gomock/internal/mock_matcher"
)

const (
	intArg = 45982
	item1  = "item1"
	item2  = "item2"
)

var sliceArg = []string{item1, item2}

func TestValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMatcher := mock_matcher.NewMockMatcher(ctrl)
	captor := gomock.Captor(mockMatcher)

	mockMatcher.EXPECT().Matches(intArg).Times(1)
	mockMatcher.EXPECT().Matches(sliceArg).Times(1)

	captor.Matches(intArg)
	captor.Matches(sliceArg)

	actualValue := captor.Value().([]string)

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

func TestValueWithNoElements(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMatcher := mock_matcher.NewMockMatcher(ctrl)
	captor := gomock.Captor(mockMatcher)

	mockMatcher.EXPECT().Matches(gomock.Any()).Times(0)

	actualValue := captor.Value()

	if actualValue != nil {
		t.Errorf("expected nil, but was %s", actualValue)
	}
}

func TestAllValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMatcher := mock_matcher.NewMockMatcher(ctrl)
	captor := gomock.Captor(mockMatcher)

	mockMatcher.EXPECT().Matches(intArg).Times(1)
	mockMatcher.EXPECT().Matches(sliceArg).Times(1)

	captor.Matches(intArg)
	captor.Matches(sliceArg)

	actualValues := captor.AllValues()

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

func TestAllValuesWithNoElements(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMatcher := mock_matcher.NewMockMatcher(ctrl)
	captor := gomock.Captor(mockMatcher)

	mockMatcher.EXPECT().Matches(gomock.Any()).Times(0)

	actualValue := captor.AllValues()

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
	if fmt.Sprintf("%s", gomock.Any()) != fmt.Sprintf("%s", anyCaptor){
		t.Errorf("expected string representation to be '%s', but was '%s'", gomock.Any(), anyCaptor)
	}
}