// Code generated by MockGen. DO NOT EDIT.
// Source: ./backend.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockBackend is a mock of Backend interface
type MockBackend struct {
	ctrl     *gomock.Controller
	recorder *MockBackendMockRecorder
}

// MockBackendMockRecorder is the mock recorder for MockBackend
type MockBackendMockRecorder struct {
	mock *MockBackend
}

// NewMockBackend creates a new mock instance
func NewMockBackend(ctrl *gomock.Controller) *MockBackend {
	mock := &MockBackend{ctrl: ctrl}
	mock.recorder = &MockBackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBackend) EXPECT() *MockBackendMockRecorder {
	return m.recorder
}

// Start mocks base method
func (m *MockBackend) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (mr *MockBackendMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockBackend)(nil).Start))
}

// UnmarshalGetConfig mocks base method
func (m *MockBackend) UnmarshalGetConfig(ctx context.Context, ptr interface{}, keys ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, ptr}
	for _, a := range keys {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UnmarshalGetConfig", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnmarshalGetConfig indicates an expected call of UnmarshalGetConfig
func (mr *MockBackendMockRecorder) UnmarshalGetConfig(ctx, ptr interface{}, keys ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, ptr}, keys...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnmarshalGetConfig", reflect.TypeOf((*MockBackend)(nil).UnmarshalGetConfig), varargs...)
}
