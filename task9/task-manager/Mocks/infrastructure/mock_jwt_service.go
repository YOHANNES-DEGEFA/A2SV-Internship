package mocks

import (
	"task_manager/infrastructure" // Need JWTClaims definition
	"github.com/stretchr/testify/mock"
)

// MockJWTService is a mock type for the JWTService type
type MockJWTService struct {
	mock.Mock
}

// GenerateToken provides a mock function with given fields: userID, username, role
func (_m *MockJWTService) GenerateToken(userID int, username string, role string) (string, error) {
	ret := _m.Called(userID, username, role)

	var r0 string
	if rf, ok := ret.Get(0).(func(int, string, string) string); ok {
		r0 = rf(userID, username, role)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, string, string) error); ok {
		r1 = rf(userID, username, role)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidateToken provides a mock function with given fields: tokenString
func (_m *MockJWTService) ValidateToken(tokenString string) (*infrastructure.JWTClaims, error) {
	ret := _m.Called(tokenString)

	var r0 *infrastructure.JWTClaims
	if rf, ok := ret.Get(0).(func(string) *infrastructure.JWTClaims); ok {
		r0 = rf(tokenString)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*infrastructure.JWTClaims)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(tokenString)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockJWTService creates a new instance of MockJWTService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockJWTService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockJWTService {
	mock := &MockJWTService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}