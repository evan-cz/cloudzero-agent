// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cloudzero/cloudzero-insights-controller/pkg/types (interfaces: Runnable)
//
// Generated by this command:
//
//	mockgen -destination=mocks/runnable_mock.go -package=mocks . Runnable
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockRunnable is a mock of Runnable interface.
type MockRunnable struct {
	ctrl     *gomock.Controller
	recorder *MockRunnableMockRecorder
	isgomock struct{}
}

// MockRunnableMockRecorder is the mock recorder for MockRunnable.
type MockRunnableMockRecorder struct {
	mock *MockRunnable
}

// NewMockRunnable creates a new mock instance.
func NewMockRunnable(ctrl *gomock.Controller) *MockRunnable {
	mock := &MockRunnable{ctrl: ctrl}
	mock.recorder = &MockRunnableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRunnable) EXPECT() *MockRunnableMockRecorder {
	return m.recorder
}

// IsRunning mocks base method.
func (m *MockRunnable) IsRunning() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRunning")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsRunning indicates an expected call of IsRunning.
func (mr *MockRunnableMockRecorder) IsRunning() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRunning", reflect.TypeOf((*MockRunnable)(nil).IsRunning))
}

// Run mocks base method.
func (m *MockRunnable) Run() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run")
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockRunnableMockRecorder) Run() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockRunnable)(nil).Run))
}

// Shutdown mocks base method.
func (m *MockRunnable) Shutdown() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shutdown")
	ret0, _ := ret[0].(error)
	return ret0
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockRunnableMockRecorder) Shutdown() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockRunnable)(nil).Shutdown))
}
