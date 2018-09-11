// Code generated by MockGen. DO NOT EDIT.
// Source: source.go
// Hash: 57f355b8f035d301b9518cbdcb09eb4a1c625fd60475bcffb44bc29a71d2c3e5

// Package mock_source is a generated GoMock package.
package mock_source

import (
	gomock "github.com/golang/mock/gomock"
	definition "github.com/golang/mock/mockgen/tests/import_source/definition"
	reflect "reflect"
)

// MockS is a mock of S interface
type MockS struct {
	ctrl     *gomock.Controller
	recorder *MockSMockRecorder
}

// MockSMockRecorder is the mock recorder for MockS
type MockSMockRecorder struct {
	mock *MockS
}

// NewMockS creates a new mock instance
func NewMockS(ctrl *gomock.Controller) *MockS {
	mock := &MockS{ctrl: ctrl}
	mock.recorder = &MockSMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockS) EXPECT() *MockSMockRecorder {
	return m.recorder
}

// F mocks base method
func (m *MockS) F(arg0 definition.X) {
	m.ctrl.Call(m, "F", arg0)
}

// F indicates an expected call of F
func (mr *MockSMockRecorder) F(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "F", reflect.TypeOf((*MockS)(nil).F), arg0)
}
