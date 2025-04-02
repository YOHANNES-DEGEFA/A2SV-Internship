package mocks

import (
	"task_manager/domain"
	"github.com/stretchr/testify/mock"
)

// MockTaskRepository is a mock type for the TaskRepository type
type MockTaskRepository struct {
	mock.Mock
}

// Save provides a mock function with given fields: task
func (_m *MockTaskRepository) Save(task domain.Task) (domain.Task, error) {
	ret := _m.Called(task)

	var r0 domain.Task
	if rf, ok := ret.Get(0).(func(domain.Task) domain.Task); ok {
		r0 = rf(task)
	} else {
		r0 = ret.Get(0).(domain.Task)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(domain.Task) error); ok {
		r1 = rf(task)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByID provides a mock function with given fields: id
func (_m *MockTaskRepository) FindByID(id int) (domain.Task, error) {
	ret := _m.Called(id)

	var r0 domain.Task
	if rf, ok := ret.Get(0).(func(int) domain.Task); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(domain.Task)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAll provides a mock function with given fields:
func (_m *MockTaskRepository) FindAll() ([]domain.Task, error) {
	ret := _m.Called()

	var r0 []domain.Task
	if rf, ok := ret.Get(0).(func() []domain.Task); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Task)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: task
func (_m *MockTaskRepository) Update(task domain.Task) (domain.Task, error) {
	ret := _m.Called(task)

	var r0 domain.Task
	if rf, ok := ret.Get(0).(func(domain.Task) domain.Task); ok {
		r0 = rf(task)
	} else {
		r0 = ret.Get(0).(domain.Task)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(domain.Task) error); ok {
		r1 = rf(task)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: id
func (_m *MockTaskRepository) Delete(id int) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockTaskRepository creates a new instance of MockTaskRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockTaskRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTaskRepository {
	mock := &MockTaskRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}