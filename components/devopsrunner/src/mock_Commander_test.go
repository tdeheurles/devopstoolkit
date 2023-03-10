// Code generated by mockery v2.21.1. DO NOT EDIT.

package main

import mock "github.com/stretchr/testify/mock"

// MockCommander is an autogenerated mock type for the Commander type
type MockCommander struct {
	mock.Mock
}

type MockCommander_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCommander) EXPECT() *MockCommander_Expecter {
	return &MockCommander_Expecter{mock: &_m.Mock}
}

// Assert provides a mock function with given fields:
func (_m *MockCommander) Assert() {
	_m.Called()
}

// MockCommander_Assert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Assert'
type MockCommander_Assert_Call struct {
	*mock.Call
}

// Assert is a helper method to define mock.On call
func (_e *MockCommander_Expecter) Assert() *MockCommander_Assert_Call {
	return &MockCommander_Assert_Call{Call: _e.mock.On("Assert")}
}

func (_c *MockCommander_Assert_Call) Run(run func()) *MockCommander_Assert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCommander_Assert_Call) Return() *MockCommander_Assert_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockCommander_Assert_Call) RunAndReturn(run func()) *MockCommander_Assert_Call {
	_c.Call.Return(run)
	return _c
}

// Execute provides a mock function with given fields:
func (_m *MockCommander) Execute() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// MockCommander_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockCommander_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
func (_e *MockCommander_Expecter) Execute() *MockCommander_Execute_Call {
	return &MockCommander_Execute_Call{Call: _e.mock.On("Execute")}
}

func (_c *MockCommander_Execute_Call) Run(run func()) *MockCommander_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCommander_Execute_Call) Return(_a0 int) *MockCommander_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCommander_Execute_Call) RunAndReturn(run func() int) *MockCommander_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// Parse provides a mock function with given fields:
func (_m *MockCommander) Parse() {
	_m.Called()
}

// MockCommander_Parse_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Parse'
type MockCommander_Parse_Call struct {
	*mock.Call
}

// Parse is a helper method to define mock.On call
func (_e *MockCommander_Expecter) Parse() *MockCommander_Parse_Call {
	return &MockCommander_Parse_Call{Call: _e.mock.On("Parse")}
}

func (_c *MockCommander_Parse_Call) Run(run func()) *MockCommander_Parse_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCommander_Parse_Call) Return() *MockCommander_Parse_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockCommander_Parse_Call) RunAndReturn(run func()) *MockCommander_Parse_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockCommander interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockCommander creates a new instance of MockCommander. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockCommander(t mockConstructorTestingTNewMockCommander) *MockCommander {
	mock := &MockCommander{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
