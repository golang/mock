// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/golang/mock/gomock (interfaces: Matcher)

package mock_gomock

import (
	gomock "github.com/golang/mock/gomock"
)

// MockMatcher is a Mock of Matcher interface
type MockMatcher struct {
	ctrl     *gomock.Controller
	recorder *MockMatcherRecorder
}

// MockMatcherRecorder is a Recorder for MockMatcher (not exported)
type MockMatcherRecorder struct {
	mock *MockMatcher
}

// NewMockMatcher creates a new instance of the mock
func NewMockMatcher(ctrl *gomock.Controller) *MockMatcher {
	mock := &MockMatcher{ctrl: ctrl}
	mock.recorder = &MockMatcherRecorder{mock}
	return mock
}

// EXPECT allows declaring expectations on the mock
func (_m *MockMatcher) EXPECT() *MockMatcherRecorder {
	return _m.recorder
}

// Matches is a mock version of the original Matches function
func (_m *MockMatcher) Matches(_param0 interface{}) bool {
	ret := _m.ctrl.Call(_m, "Matches", _param0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Matches is a recorder version of the mocked Matches method
func (_mr *MockMatcherRecorder) Matches(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Matches", arg0)
}

// String is a mock version of the original String function
func (_m *MockMatcher) String() string {
	ret := _m.ctrl.Call(_m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String is a recorder version of the mocked String method
func (_mr *MockMatcherRecorder) String() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "String")
}
