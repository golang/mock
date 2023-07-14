// Code generated by MockGen. DO NOT EDIT.
// Source: bugreport.go
//
// Generated by this command:
//
//	mockgen -destination bugreport_mock.go -package bugreport -source=bugreport.go
//
// Package bugreport is a generated GoMock package.
package bugreport

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockExample is a mock of Example interface.
type MockExample struct {
	ctrl     *gomock.Controller
	recorder *MockExampleMockRecorder
}

// MockExampleMockRecorder is the mock recorder for MockExample.
type MockExampleMockRecorder struct {
	mock *MockExample
}

// NewMockExample creates a new mock instance.
func NewMockExample(ctrl *gomock.Controller) *MockExample {
	mock := &MockExample{ctrl: ctrl}
	mock.recorder = &MockExampleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExample) EXPECT() *MockExampleMockRecorder {
	return m.recorder
}

// Method mocks base method.
func (m_2 *MockExample) Method(_m, _mr, m, mr int) {
	m_2.ctrl.T.Helper()
	m_2.ctrl.Call(m_2, "Method", _m, _mr, m, mr)
}

// Method indicates an expected call of Method.
func (mr_2 *MockExampleMockRecorder) Method(_m, _mr, m, mr interface{}) *gomock.Call {
	mr_2.mock.ctrl.T.Helper()
	return mr_2.mock.ctrl.RecordCallWithMethodType(mr_2.mock, "Method", reflect.TypeOf((*MockExample)(nil).Method), _m, _mr, m, mr)
}

// VarargMethod mocks base method.
func (m *MockExample) VarargMethod(_s, _x, a, ret int, varargs ...int) {
	m.ctrl.T.Helper()
	varargs_2 := []interface{}{_s, _x, a, ret}
	for _, a_2 := range varargs {
		varargs_2 = append(varargs_2, a_2)
	}
	m.ctrl.Call(m, "VarargMethod", varargs_2...)
}

// VarargMethod indicates an expected call of VarargMethod.
func (mr *MockExampleMockRecorder) VarargMethod(_s, _x, a, ret interface{}, varargs ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs_2 := append([]interface{}{_s, _x, a, ret}, varargs...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VarargMethod", reflect.TypeOf((*MockExample)(nil).VarargMethod), varargs_2...)
}
