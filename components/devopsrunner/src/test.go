package main

// import (
// 	"context"
// 	"fmt"

// 	"github.com/stretchr/testify/mock"
// )

// func AssertNumberOfCalls(calls []mock.Call, methodName string, expectedCalls int, ctx context.Context) (context.Context, error) {
// 	var actualCalls int
// 	for _, call := range calls {
// 		if call.Method == methodName {
// 			actualCalls++
// 		}
// 	}

// 	if expectedCalls != actualCalls {
// 		return ctx, fmt.Errorf("expected number of calls (%d) does not match the actual number of calls (%d)", expectedCalls, actualCalls)
// 	}
// 	return ctx, nil
// }

// func AssertExpectations(expectedCalls []mock.Call) (context.Context, error) {
// 	var failedExpectations int

// 	// iterate through each expectation
// 	for _, expectedCall := range expectedCalls {
// 		satisfied, reason := m.checkExpectation(expectedCall)
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

// func checkExpectation(call *mock.Call) (bool, string) {
// 	if !call.optional && !m.methodWasCalled(call.Method, call.Arguments) && call.totalCalls == 0 {
// 		return false, fmt.Sprintf("FAIL:\t%s(%s)\n\t\tat: %s", call.Method, call.Arguments.String(), call.callerInfo)
// 	}
// 	if call.Repeatability > 0 {
// 		return false, fmt.Sprintf("FAIL:\t%s(%s)\n\t\tat: %s", call.Method, call.Arguments.String(), call.callerInfo)
// 	}
// 	return true, fmt.Sprintf("PASS:\t%s(%s)", call.Method, call.Arguments.String())
// }
