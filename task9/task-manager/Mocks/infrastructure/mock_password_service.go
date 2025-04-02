package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockPasswordService is a mock type for the PasswordService type
type MockPasswordService struct {
	mock.Mock
}

// Hash provides a mock function with given fields: password
func (_m *MockPasswordService) Hash(password string) (string, error) {
	ret := _m.Called(password)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(password)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Compare provides a mock function with given fields: hashedPassword, password
func (_m *MockPasswordService) Compare(hashedPassword string, password string) error {
	ret := _m.Called(hashedPassword, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(hashedPassword, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockPasswordService creates a new instance of MockPasswordService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockPasswordService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPasswordService {
	mock := &MockPasswordService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}