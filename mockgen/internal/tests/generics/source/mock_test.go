// Code generated by MockGen. DO NOT EDIT.
// Source: generics.go

// Package source is a generated GoMock package.
package source

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	generics "github.com/golang/mock/mockgen/internal/tests/generics"
	other "github.com/golang/mock/mockgen/internal/tests/generics/other"
)

// MockBar is a mock of Bar interface.
type MockBar[T any, R any] struct {
	ctrl     *gomock.Controller
	recorder *MockBarMockRecorder[T, R]
}

// MockBarMockRecorder is the mock recorder for MockBar.
type MockBarMockRecorder[T any, R any] struct {
	mock *MockBar[T, R]
}

// NewMockBar creates a new mock instance.
func NewMockBar[T any, R any](ctrl *gomock.Controller) *MockBar[T, R] {
	mock := &MockBar[T, R]{ctrl: ctrl}
	mock.recorder = &MockBarMockRecorder[T, R]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBar[T, R]) EXPECT() *MockBarMockRecorder[T, R] {
	return m.recorder
}

// Eight mocks base method.
func (m *MockBar[T, R]) Eight(arg0 T) other.Two[T, R] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Eight", arg0)
	ret0, _ := ret[0].(other.Two[T, R])
	return ret0
}

// Eight indicates an expected call of Eight.
func (mr *MockBarMockRecorder[T, R]) Eight(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Eight", reflect.TypeOf((*MockBar[T, R])(nil).Eight), arg0)
}

// Five mocks base method.
func (m *MockBar[T, R]) Five(arg0 T) generics.Baz[T] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Five", arg0)
	ret0, _ := ret[0].(generics.Baz[T])
	return ret0
}

// Five indicates an expected call of Five.
func (mr *MockBarMockRecorder[T, R]) Five(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Five", reflect.TypeOf((*MockBar[T, R])(nil).Five), arg0)
}

// Four mocks base method.
func (m *MockBar[T, R]) Four(arg0 T) generics.Foo[T, R] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Four", arg0)
	ret0, _ := ret[0].(generics.Foo[T, R])
	return ret0
}

// Four indicates an expected call of Four.
func (mr *MockBarMockRecorder[T, R]) Four(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Four", reflect.TypeOf((*MockBar[T, R])(nil).Four), arg0)
}

// Nine mocks base method.
func (m *MockBar[T, R]) Nine(arg0 generics.Iface[T]) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Nine", arg0)
}

// Nine indicates an expected call of Nine.
func (mr *MockBarMockRecorder[T, R]) Nine(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Nine", reflect.TypeOf((*MockBar[T, R])(nil).Nine), arg0)
}

// One mocks base method.
func (m *MockBar[T, R]) One(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "One", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// One indicates an expected call of One.
func (mr *MockBarMockRecorder[T, R]) One(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "One", reflect.TypeOf((*MockBar[T, R])(nil).One), arg0)
}

// Seven mocks base method.
func (m *MockBar[T, R]) Seven(arg0 T) other.One[T] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Seven", arg0)
	ret0, _ := ret[0].(other.One[T])
	return ret0
}

// Seven indicates an expected call of Seven.
func (mr *MockBarMockRecorder[T, R]) Seven(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Seven", reflect.TypeOf((*MockBar[T, R])(nil).Seven), arg0)
}

// Six mocks base method.
func (m *MockBar[T, R]) Six(arg0 T) *generics.Baz[T] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Six", arg0)
	ret0, _ := ret[0].(*generics.Baz[T])
	return ret0
}

// Six indicates an expected call of Six.
func (mr *MockBarMockRecorder[T, R]) Six(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Six", reflect.TypeOf((*MockBar[T, R])(nil).Six), arg0)
}

// Ten mocks base method.
func (m *MockBar[T, R]) Ten(arg0 *T) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Ten", arg0)
}

// Ten indicates an expected call of Ten.
func (mr *MockBarMockRecorder[T, R]) Ten(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ten", reflect.TypeOf((*MockBar[T, R])(nil).Ten), arg0)
}

// Three mocks base method.
func (m *MockBar[T, R]) Three(arg0 T) R {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Three", arg0)
	ret0, _ := ret[0].(R)
	return ret0
}

// Three indicates an expected call of Three.
func (mr *MockBarMockRecorder[T, R]) Three(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Three", reflect.TypeOf((*MockBar[T, R])(nil).Three), arg0)
}

// Two mocks base method.
func (m *MockBar[T, R]) Two(arg0 T) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Two", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// Two indicates an expected call of Two.
func (mr *MockBarMockRecorder[T, R]) Two(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Two", reflect.TypeOf((*MockBar[T, R])(nil).Two), arg0)
}

// MockIface is a mock of Iface interface.
type MockIface[T any] struct {
	ctrl     *gomock.Controller
	recorder *MockIfaceMockRecorder[T]
}

// MockIfaceMockRecorder is the mock recorder for MockIface.
type MockIfaceMockRecorder[T any] struct {
	mock *MockIface[T]
}

// NewMockIface creates a new mock instance.
func NewMockIface[T any](ctrl *gomock.Controller) *MockIface[T] {
	mock := &MockIface[T]{ctrl: ctrl}
	mock.recorder = &MockIfaceMockRecorder[T]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIface[T]) EXPECT() *MockIfaceMockRecorder[T] {
	return m.recorder
}
