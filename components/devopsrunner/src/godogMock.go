package main

// import (
// 	"fmt"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// type GodogMock struct {
// 	mock.Mock
// }

// type GodogCall struct {
// 	mock.Call
// }

// func (m *GodogMock) AssertExpectations() bool {
// 	var failedExpectations int

// 	// iterate through each expectation
// 	expectedCalls := m.ExpectedCalls()
// 	for _, expectedCall := range expectedCalls {
// 		satisfied, reason := m.CheckExpectation(expectedCall)
// 		if !satisfied {
// 			failedExpectations++
// 		}
// 		t.Logf(reason)
// 	}

// 	if failedExpectations != 0 {
// 		t.Errorf("FAIL: %d out of %d expectation(s) were met.\n\tThe code you are testing needs to make %d more call(s).\n\tat: %s", len(expectedCalls)-failedExpectations, len(expectedCalls), failedExpectations, assert.CallerInfo())
// 	}

// 	return failedExpectations == 0
// }

// func (m *GodogMock) ExpectedCalls() []*GodogCall {
// 	return append([]*GodogCall{}, m.ExpectedCalls...)
// }

// func (m *GodogMock) CheckExpectation(call *GodogCall) (bool, string) {
// 	if !call.optional && !m.methodWasCalled(call.Method, call.Arguments) && call.totalCalls == 0 {
// 		return false, fmt.Sprintf("FAIL:\t%s(%s)\n\t\tat: %s", call.Method, call.Arguments.String(), call.callerInfo)
// 	}
// 	if call.Repeatability > 0 {
// 		return false, fmt.Sprintf("FAIL:\t%s(%s)\n\t\tat: %s", call.Method, call.Arguments.String(), call.callerInfo)
// 	}
// 	return true, fmt.Sprintf("PASS:\t%s(%s)", call.Method, call.Arguments.String())
// }
