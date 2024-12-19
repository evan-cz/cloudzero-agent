// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/types/resource_store.go
//
// Generated by this command:
//
//	mockgen -source=pkg/types/resource_store.go -destination=pkg/types/mocks/resource_store_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	types "github.com/cloudzero/cloudzero-insights-controller/pkg/types"
	gomock "go.uber.org/mock/gomock"
)

// MockResourceStore is a mock of ResourceStore interface.
type MockResourceStore struct {
	ctrl     *gomock.Controller
	recorder *MockResourceStoreMockRecorder
	isgomock struct{}
}

// MockResourceStoreMockRecorder is the mock recorder for MockResourceStore.
type MockResourceStoreMockRecorder struct {
	mock *MockResourceStore
}

// NewMockResourceStore creates a new mock instance.
func NewMockResourceStore(ctrl *gomock.Controller) *MockResourceStore {
	mock := &MockResourceStore{ctrl: ctrl}
	mock.recorder = &MockResourceStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResourceStore) EXPECT() *MockResourceStoreMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockResourceStore) Count(ctx context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockResourceStoreMockRecorder) Count(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockResourceStore)(nil).Count), ctx)
}

// Create mocks base method.
func (m *MockResourceStore) Create(ctx context.Context, it *types.ResourceTags) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, it)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockResourceStoreMockRecorder) Create(ctx, it any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockResourceStore)(nil).Create), ctx, it)
}

// Delete mocks base method.
func (m *MockResourceStore) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockResourceStoreMockRecorder) Delete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockResourceStore)(nil).Delete), ctx, id)
}

// DeleteAll mocks base method.
func (m *MockResourceStore) DeleteAll(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAll", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAll indicates an expected call of DeleteAll.
func (mr *MockResourceStoreMockRecorder) DeleteAll(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAll", reflect.TypeOf((*MockResourceStore)(nil).DeleteAll), ctx)
}

// FindAllBy mocks base method.
func (m *MockResourceStore) FindAllBy(ctx context.Context, conds ...any) ([]*types.ResourceTags, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx}
	for _, a := range conds {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindAllBy", varargs...)
	ret0, _ := ret[0].([]*types.ResourceTags)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllBy indicates an expected call of FindAllBy.
func (mr *MockResourceStoreMockRecorder) FindAllBy(ctx any, conds ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx}, conds...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllBy", reflect.TypeOf((*MockResourceStore)(nil).FindAllBy), varargs...)
}

// FindFirstBy mocks base method.
func (m *MockResourceStore) FindFirstBy(ctx context.Context, conds ...any) (*types.ResourceTags, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx}
	for _, a := range conds {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindFirstBy", varargs...)
	ret0, _ := ret[0].(*types.ResourceTags)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindFirstBy indicates an expected call of FindFirstBy.
func (mr *MockResourceStoreMockRecorder) FindFirstBy(ctx any, conds ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx}, conds...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindFirstBy", reflect.TypeOf((*MockResourceStore)(nil).FindFirstBy), varargs...)
}

// Get mocks base method.
func (m *MockResourceStore) Get(ctx context.Context, id string) (*types.ResourceTags, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*types.ResourceTags)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockResourceStoreMockRecorder) Get(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockResourceStore)(nil).Get), ctx, id)
}

// Tx mocks base method.
func (m *MockResourceStore) Tx(ctx context.Context, block func(context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tx", ctx, block)
	ret0, _ := ret[0].(error)
	if ret0 != nil {
		return ret0
	}
	// NOTE: Manual changed from generated code to ensure we invoke the block function
	return block(ctx)
}

// Tx indicates an expected call of Tx.
func (mr *MockResourceStoreMockRecorder) Tx(ctx, block any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tx", reflect.TypeOf((*MockResourceStore)(nil).Tx), ctx, block)
}

// Update mocks base method.
func (m *MockResourceStore) Update(ctx context.Context, it *types.ResourceTags) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, it)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockResourceStoreMockRecorder) Update(ctx, it any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockResourceStore)(nil).Update), ctx, it)
}
