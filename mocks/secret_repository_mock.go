// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/nalawade41/secret-server/internal/domain (interfaces: SecretRepository)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/nalawade41/secret-server/internal/domain"
)

// MockSecretRepository is a mock of SecretRepository interface.
type MockSecretRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSecretRepositoryMockRecorder
}

// MockSecretRepositoryMockRecorder is the mock recorder for MockSecretRepository.
type MockSecretRepositoryMockRecorder struct {
	mock *MockSecretRepository
}

// NewMockSecretRepository creates a new mock instance.
func NewMockSecretRepository(ctrl *gomock.Controller) *MockSecretRepository {
	mock := &MockSecretRepository{ctrl: ctrl}
	mock.recorder = &MockSecretRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretRepository) EXPECT() *MockSecretRepositoryMockRecorder {
	return m.recorder
}

// DeleteSecret mocks base method.
func (m *MockSecretRepository) DeleteSecret(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSecret", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSecret indicates an expected call of DeleteSecret.
func (mr *MockSecretRepositoryMockRecorder) DeleteSecret(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSecret", reflect.TypeOf((*MockSecretRepository)(nil).DeleteSecret), arg0, arg1)
}

// GetByHash mocks base method.
func (m *MockSecretRepository) GetByHash(arg0 context.Context, arg1 string) (domain.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByHash", arg0, arg1)
	ret0, _ := ret[0].(domain.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByHash indicates an expected call of GetByHash.
func (mr *MockSecretRepositoryMockRecorder) GetByHash(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByHash", reflect.TypeOf((*MockSecretRepository)(nil).GetByHash), arg0, arg1)
}

// Save mocks base method.
func (m *MockSecretRepository) Save(arg0 context.Context, arg1 domain.Secret) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockSecretRepositoryMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockSecretRepository)(nil).Save), arg0, arg1)
}

// UpdateSecretViews mocks base method.
func (m *MockSecretRepository) UpdateSecretViews(arg0 context.Context, arg1 string, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSecretViews", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSecretViews indicates an expected call of UpdateSecretViews.
func (mr *MockSecretRepositoryMockRecorder) UpdateSecretViews(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSecretViews", reflect.TypeOf((*MockSecretRepository)(nil).UpdateSecretViews), arg0, arg1, arg2)
}
