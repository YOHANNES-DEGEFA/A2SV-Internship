package mocks

import (
	"task_manager/domain"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock type for the UserRepository type
type MockUserRepository struct {
	mock.Mock
}

// Save provides a mock function with given fields: user, passwordHash
func (_m *MockUserRepository) Save(user domain.User, passwordHash string) (domain.User, error) {
	ret := _m.Called(user, passwordHash)

	var r0 domain.User
	if rf, ok := ret.Get(0).(func(domain.User, string) domain.User); ok {
		r0 = rf(user, passwordHash)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(domain.User, string) error); ok {
		r1 = rf(user, passwordHash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByUsername provides a mock function with given fields: username
func (_m *MockUserRepository) FindByUsername(username string) (domain.User, string, error) {
	ret := _m.Called(username)

	var r0 domain.User
	if rf, ok := ret.Get(0).(func(string) domain.User); ok {
		r0 = rf(username)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(string) string); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(username)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// FindByID provides a mock function with given fields: id
func (_m *MockUserRepository) FindByID(id int) (domain.User, error) {
	ret := _m.Called(id)

	var r0 domain.User
	if rf, ok := ret.Get(0).(func(int) domain.User); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}


// UpdateRole provides a mock function with given fields: userID, newRole
func (_m *MockUserRepository) UpdateRole(userID int, newRole string) error {
	ret := _m.Called(userID, newRole)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, string) error); ok {
		r0 = rf(userID, newRole)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CountUsers provides a mock function with given fields:
func (_m *MockUserRepository) CountUsers() (int, error) {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockUserRepository creates a new instance of MockUserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUserRepository {
	mock := &MockUserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}