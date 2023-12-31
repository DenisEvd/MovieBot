// Code generated by MockGen. DO NOT EDIT.
// Source: internal/storage/storage.go

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	events "MovieBot/internal/events"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// AddMovie mocks base method.
func (m *MockStorage) AddMovie(username string, movie events.Movie) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMovie", username, movie)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddMovie indicates an expected call of AddMovie.
func (mr *MockStorageMockRecorder) AddMovie(username, movie interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMovie", reflect.TypeOf((*MockStorage)(nil).AddMovie), username, movie)
}

// AddRequest mocks base method.
func (m *MockStorage) AddRequest(text string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRequest", text)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddRequest indicates an expected call of AddRequest.
func (mr *MockStorageMockRecorder) AddRequest(text interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRequest", reflect.TypeOf((*MockStorage)(nil).AddRequest), text)
}

// DeleteRequest mocks base method.
func (m *MockStorage) DeleteRequest(id int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRequest", id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteRequest indicates an expected call of DeleteRequest.
func (mr *MockStorageMockRecorder) DeleteRequest(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRequest", reflect.TypeOf((*MockStorage)(nil).DeleteRequest), id)
}

// GetAll mocks base method.
func (m *MockStorage) GetAll(username string) ([]events.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", username)
	ret0, _ := ret[0].([]events.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockStorageMockRecorder) GetAll(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockStorage)(nil).GetAll), username)
}

// GetNMovie mocks base method.
func (m *MockStorage) GetNMovie(username string, n int) (events.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNMovie", username, n)
	ret0, _ := ret[0].(events.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNMovie indicates an expected call of GetNMovie.
func (mr *MockStorageMockRecorder) GetNMovie(username, n interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNMovie", reflect.TypeOf((*MockStorage)(nil).GetNMovie), username, n)
}

// IsExistRecord mocks base method.
func (m *MockStorage) IsExistRecord(username string, movieID int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsExistRecord", username, movieID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsExistRecord indicates an expected call of IsExistRecord.
func (mr *MockStorageMockRecorder) IsExistRecord(username, movieID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsExistRecord", reflect.TypeOf((*MockStorage)(nil).IsExistRecord), username, movieID)
}

// IsWatched mocks base method.
func (m *MockStorage) IsWatched(username string, movieID int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsWatched", username, movieID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsWatched indicates an expected call of IsWatched.
func (mr *MockStorageMockRecorder) IsWatched(username, movieID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsWatched", reflect.TypeOf((*MockStorage)(nil).IsWatched), username, movieID)
}

// Remove mocks base method.
func (m *MockStorage) Remove(username string, movieID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", username, movieID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockStorageMockRecorder) Remove(username, movieID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockStorage)(nil).Remove), username, movieID)
}

// Watch mocks base method.
func (m *MockStorage) Watch(username string, movieID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Watch", username, movieID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Watch indicates an expected call of Watch.
func (mr *MockStorageMockRecorder) Watch(username, movieID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockStorage)(nil).Watch), username, movieID)
}
