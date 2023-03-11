// Code generated by mockery v2.22.1. DO NOT EDIT.

package main

import mock "github.com/stretchr/testify/mock"

// MockDevopsRunnerer is an autogenerated mock type for the DevopsRunnerer type
type MockDevopsRunnerer struct {
	mock.Mock
}

type MockDevopsRunnerer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDevopsRunnerer) EXPECT() *MockDevopsRunnerer_Expecter {
	return &MockDevopsRunnerer_Expecter{mock: &_m.Mock}
}

// Run provides a mock function with given fields:
func (_m *MockDevopsRunnerer) Run() {
	_m.Called()
}

// MockDevopsRunnerer_Run_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Run'
type MockDevopsRunnerer_Run_Call struct {
	*mock.Call
}

// Run is a helper method to define mock.On call
func (_e *MockDevopsRunnerer_Expecter) Run() *MockDevopsRunnerer_Run_Call {
	return &MockDevopsRunnerer_Run_Call{Call: _e.mock.On("Run")}
}

func (_c *MockDevopsRunnerer_Run_Call) Run(run func()) *MockDevopsRunnerer_Run_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockDevopsRunnerer_Run_Call) Return() *MockDevopsRunnerer_Run_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockDevopsRunnerer_Run_Call) RunAndReturn(run func()) *MockDevopsRunnerer_Run_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockDevopsRunnerer interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockDevopsRunnerer creates a new instance of MockDevopsRunnerer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockDevopsRunnerer(t mockConstructorTestingTNewMockDevopsRunnerer) *MockDevopsRunnerer {
	mock := &MockDevopsRunnerer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
