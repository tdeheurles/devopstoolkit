// Code generated by mockery v2.21.1. DO NOT EDIT.

package main

import mock "github.com/stretchr/testify/mock"

// MockCommandFactorier is an autogenerated mock type for the CommandFactorier type
type MockCommandFactorier struct {
	mock.Mock
}

type MockCommandFactorier_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCommandFactorier) EXPECT() *MockCommandFactorier_Expecter {
	return &MockCommandFactorier_Expecter{mock: &_m.Mock}
}

// Parse provides a mock function with given fields:
func (_m *MockCommandFactorier) Parse() Commander {
	ret := _m.Called()

	var r0 Commander
	if rf, ok := ret.Get(0).(func() Commander); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Commander)
		}
	}

	return r0
}

// MockCommandFactorier_Parse_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Parse'
type MockCommandFactorier_Parse_Call struct {
	*mock.Call
}

// Parse is a helper method to define mock.On call
func (_e *MockCommandFactorier_Expecter) Parse() *MockCommandFactorier_Parse_Call {
	return &MockCommandFactorier_Parse_Call{Call: _e.mock.On("Parse")}
}

func (_c *MockCommandFactorier_Parse_Call) Run(run func()) *MockCommandFactorier_Parse_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCommandFactorier_Parse_Call) Return(_a0 Commander) *MockCommandFactorier_Parse_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCommandFactorier_Parse_Call) RunAndReturn(run func() Commander) *MockCommandFactorier_Parse_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockCommandFactorier interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockCommandFactorier creates a new instance of MockCommandFactorier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockCommandFactorier(t mockConstructorTestingTNewMockCommandFactorier) *MockCommandFactorier {
	mock := &MockCommandFactorier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
