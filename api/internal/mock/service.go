// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/pkulik0/autocc/api/internal/service (interfaces: Service)
//
// Generated by this command:
//
//	mockgen -destination=../mock/service.go -package=mock . Service
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	model "github.com/pkulik0/autocc/api/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// AddCredentialsDeepL mocks base method.
func (m *MockService) AddCredentialsDeepL(arg0 context.Context, arg1 string) (*model.CredentialsDeepL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCredentialsDeepL", arg0, arg1)
	ret0, _ := ret[0].(*model.CredentialsDeepL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCredentialsDeepL indicates an expected call of AddCredentialsDeepL.
func (mr *MockServiceMockRecorder) AddCredentialsDeepL(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCredentialsDeepL", reflect.TypeOf((*MockService)(nil).AddCredentialsDeepL), arg0, arg1)
}

// AddCredentialsGoogle mocks base method.
func (m *MockService) AddCredentialsGoogle(arg0 context.Context, arg1, arg2 string) (*model.CredentialsGoogle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCredentialsGoogle", arg0, arg1, arg2)
	ret0, _ := ret[0].(*model.CredentialsGoogle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCredentialsGoogle indicates an expected call of AddCredentialsGoogle.
func (mr *MockServiceMockRecorder) AddCredentialsGoogle(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCredentialsGoogle", reflect.TypeOf((*MockService)(nil).AddCredentialsGoogle), arg0, arg1, arg2)
}

// CreateSessionGoogle mocks base method.
func (m *MockService) CreateSessionGoogle(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSessionGoogle", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSessionGoogle indicates an expected call of CreateSessionGoogle.
func (mr *MockServiceMockRecorder) CreateSessionGoogle(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSessionGoogle", reflect.TypeOf((*MockService)(nil).CreateSessionGoogle), arg0, arg1, arg2)
}

// GetCredentials mocks base method.
func (m *MockService) GetCredentials(arg0 context.Context) ([]model.CredentialsGoogle, []model.CredentialsDeepL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCredentials", arg0)
	ret0, _ := ret[0].([]model.CredentialsGoogle)
	ret1, _ := ret[1].([]model.CredentialsDeepL)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCredentials indicates an expected call of GetCredentials.
func (mr *MockServiceMockRecorder) GetCredentials(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCredentials", reflect.TypeOf((*MockService)(nil).GetCredentials), arg0)
}

// GetSessionGoogleURL mocks base method.
func (m *MockService) GetSessionGoogleURL(arg0 context.Context, arg1 uint, arg2 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSessionGoogleURL", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSessionGoogleURL indicates an expected call of GetSessionGoogleURL.
func (mr *MockServiceMockRecorder) GetSessionGoogleURL(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSessionGoogleURL", reflect.TypeOf((*MockService)(nil).GetSessionGoogleURL), arg0, arg1, arg2)
}

// GetSessionsGoogleByUser mocks base method.
func (m *MockService) GetSessionsGoogleByUser(arg0 context.Context, arg1 string) ([]model.SessionGoogle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSessionsGoogleByUser", arg0, arg1)
	ret0, _ := ret[0].([]model.SessionGoogle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSessionsGoogleByUser indicates an expected call of GetSessionsGoogleByUser.
func (mr *MockServiceMockRecorder) GetSessionsGoogleByUser(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSessionsGoogleByUser", reflect.TypeOf((*MockService)(nil).GetSessionsGoogleByUser), arg0, arg1)
}

// RemoveCredentialsDeepL mocks base method.
func (m *MockService) RemoveCredentialsDeepL(arg0 context.Context, arg1 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveCredentialsDeepL", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveCredentialsDeepL indicates an expected call of RemoveCredentialsDeepL.
func (mr *MockServiceMockRecorder) RemoveCredentialsDeepL(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveCredentialsDeepL", reflect.TypeOf((*MockService)(nil).RemoveCredentialsDeepL), arg0, arg1)
}

// RemoveCredentialsGoogle mocks base method.
func (m *MockService) RemoveCredentialsGoogle(arg0 context.Context, arg1 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveCredentialsGoogle", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveCredentialsGoogle indicates an expected call of RemoveCredentialsGoogle.
func (mr *MockServiceMockRecorder) RemoveCredentialsGoogle(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveCredentialsGoogle", reflect.TypeOf((*MockService)(nil).RemoveCredentialsGoogle), arg0, arg1)
}

// RemoveSessionGoogle mocks base method.
func (m *MockService) RemoveSessionGoogle(arg0 context.Context, arg1 string, arg2 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveSessionGoogle", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveSessionGoogle indicates an expected call of RemoveSessionGoogle.
func (mr *MockServiceMockRecorder) RemoveSessionGoogle(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveSessionGoogle", reflect.TypeOf((*MockService)(nil).RemoveSessionGoogle), arg0, arg1, arg2)
}
