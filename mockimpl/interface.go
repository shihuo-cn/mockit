// Code generated by MockGen. DO NOT EDIT.
// Source: ./iface/interface.go

// Package mockimpl is a generated GoMock package.
package mockimpl

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	iface "github.com/shihuo-cn/mockit/iface"
)

// MockDemoInterface is a mock of DemoInterface interface.
type MockDemoInterface struct {
	ctrl     *gomock.Controller
	recorder *MockDemoInterfaceMockRecorder
}

// MockDemoInterfaceMockRecorder is the mock recorder for MockDemoInterface.
type MockDemoInterfaceMockRecorder struct {
	mock *MockDemoInterface
}

// NewMockDemoInterface creates a new mock instance.
//func NewMockDemoInterface(ctrl *gomock.Controller) *MockDemoInterface
func NewMockDemoInterface(ctrl *gomock.Controller) interface{} {
	mock := &MockDemoInterface{ctrl: ctrl}
	mock.recorder = &MockDemoInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDemoInterface) EXPECT() *MockDemoInterfaceMockRecorder {
	return m.recorder
}

// Aggregate mocks base method.
func (m *MockDemoInterface) Aggregate(ctx context.Context, dm *iface.DemoInterfaceModel, time time.Time) ([]*iface.DemoInterfaceModel, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Aggregate", ctx, dm, time)
	ret0, _ := ret[0].([]*iface.DemoInterfaceModel)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Aggregate indicates an expected call of Aggregate.
func (mr *MockDemoInterfaceMockRecorder) Aggregate(ctx, dm, time interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Aggregate", reflect.TypeOf((*MockDemoInterface)(nil).Aggregate), ctx, dm, time)
}

// First mocks base method.
func (m *MockDemoInterface) First(ctx context.Context, key int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "First", ctx, key)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// First indicates an expected call of First.
func (mr *MockDemoInterfaceMockRecorder) First(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "First", reflect.TypeOf((*MockDemoInterface)(nil).First), ctx, key)
}

// List mocks base method.
func (m *MockDemoInterface) List(ctx context.Context, relationId int64, pageIndex, pageSize int) ([]*iface.KV, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, relationId, pageIndex, pageSize)
	ret0, _ := ret[0].([]*iface.KV)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockDemoInterfaceMockRecorder) List(ctx, relationId, pageIndex, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockDemoInterface)(nil).List), ctx, relationId, pageIndex, pageSize)
}

// Put mocks base method.
func (m *MockDemoInterface) Put(ctx context.Context, key, val int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", ctx, key, val)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockDemoInterfaceMockRecorder) Put(ctx, key, val interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockDemoInterface)(nil).Put), ctx, key, val)
}
