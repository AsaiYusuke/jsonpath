package jsonpath

import (
	"fmt"
	"testing"
)

func TestAggregateFunction_SyntaxErrors(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$.func(a`,
			inputJSON:   `{}`,
			expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `(a`},
		},
		{
			jsonpath:    `$.func(a)`,
			inputJSON:   `{}`,
			expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `(a)`},
		},
	}

	for i, tc := range testCases {
		runTestCase(t, tc, fmt.Sprintf("TestAggregateFunction_SyntaxErrors_case_%d", i))
	}
}

func TestAggregateFunction_NotFoundErrors(t *testing.T) {
	tests := []struct {
		name     string
		testCase TestCase
	}{
		{
			name: "unknown function with invalid syntax",
			testCase: TestCase{
				jsonpath:    `$.func()(`,
				inputJSON:   `{}`,
				expectedErr: ErrorFunctionNotFound{function: `.func()`},
			},
		},
		{
			name: "unknown function call",
			testCase: TestCase{
				jsonpath:    `$.func(){}`,
				inputJSON:   `{}`,
				expectedErr: ErrorFunctionNotFound{function: `.func()`},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestCase(t, tt.testCase, tt.name)
		})
	}
}
