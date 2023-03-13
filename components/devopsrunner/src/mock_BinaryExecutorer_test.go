// Code generated by mockery v2.21.1. DO NOT EDIT.

package main

import mock "github.com/stretchr/testify/mock"

// MockBinaryExecutorer is an autogenerated mock type for the BinaryExecutorer type
type MockBinaryExecutorer struct {
	mock.Mock
}

type MockBinaryExecutorer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockBinaryExecutorer) EXPECT() *MockBinaryExecutorer_Expecter {
	return &MockBinaryExecutorer_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: version
func (_m *MockBinaryExecutorer) Execute(version string) (bool, int) {
	ret := _m.Called(version)

	var r0 bool
	var r1 int
	if rf, ok := ret.Get(0).(func(string) (bool, int)); ok {
		return rf(version)
	}
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(version)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string) int); ok {
		r1 = rf(version)
	} else {
		r1 = ret.Get(1).(int)
	}

	return r0, r1
}

// MockBinaryExecutorer_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockBinaryExecutorer_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - version string
func (_e *MockBinaryExecutorer_Expecter) Execute(version interface{}) *MockBinaryExecutorer_Execute_Call {
	return &MockBinaryExecutorer_Execute_Call{Call: _e.mock.On("Execute", version)}
}

func (_c *MockBinaryExecutorer_Execute_Call) Run(run func(version string)) *MockBinaryExecutorer_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockBinaryExecutorer_Execute_Call) Return(stop bool, exitCode int) *MockBinaryExecutorer_Execute_Call {
	_c.Call.Return(stop, exitCode)
	return _c
}

func (_c *MockBinaryExecutorer_Execute_Call) RunAndReturn(run func(string) (bool, int)) *MockBinaryExecutorer_Execute_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockBinaryExecutorer interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockBinaryExecutorer creates a new instance of MockBinaryExecutorer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockBinaryExecutorer(t mockConstructorTestingTNewMockBinaryExecutorer) *MockBinaryExecutorer {
	mock := &MockBinaryExecutorer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
