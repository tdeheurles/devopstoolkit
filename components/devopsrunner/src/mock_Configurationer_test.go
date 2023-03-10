// Code generated by mockery v2.21.1. DO NOT EDIT.

package main

import mock "github.com/stretchr/testify/mock"

// MockConfigurationer is an autogenerated mock type for the Configurationer type
type MockConfigurationer struct {
	mock.Mock
}

type MockConfigurationer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockConfigurationer) EXPECT() *MockConfigurationer_Expecter {
	return &MockConfigurationer_Expecter{mock: &_m.Mock}
}

// Data provides a mock function with given fields:
func (_m *MockConfigurationer) Data() Config {
	ret := _m.Called()

	var r0 Config
	if rf, ok := ret.Get(0).(func() Config); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(Config)
	}

	return r0
}

// MockConfigurationer_Data_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Data'
type MockConfigurationer_Data_Call struct {
	*mock.Call
}

// Data is a helper method to define mock.On call
func (_e *MockConfigurationer_Expecter) Data() *MockConfigurationer_Data_Call {
	return &MockConfigurationer_Data_Call{Call: _e.mock.On("Data")}
}

func (_c *MockConfigurationer_Data_Call) Run(run func()) *MockConfigurationer_Data_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockConfigurationer_Data_Call) Return(_a0 Config) *MockConfigurationer_Data_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConfigurationer_Data_Call) RunAndReturn(run func() Config) *MockConfigurationer_Data_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockConfigurationer interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockConfigurationer creates a new instance of MockConfigurationer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockConfigurationer(t mockConstructorTestingTNewMockConfigurationer) *MockConfigurationer {
	mock := &MockConfigurationer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
