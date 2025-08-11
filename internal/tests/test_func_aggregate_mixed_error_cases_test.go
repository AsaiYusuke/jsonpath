package tests

import (
	"testing"
)

func TestAggregateFunction_SyntaxErrors(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$.func(a`,
			inputJSON:   `{}`,
			expectedErr: createErrorInvalidSyntax(6, `unrecognized input`, `(a`),
		},
		{
			jsonpath:    `$.func(a)`,
			inputJSON:   `{}`,
			expectedErr: createErrorInvalidSyntax(6, `unrecognized input`, `(a)`),
		},
		{
			jsonpath:    `$.func@()`,
			inputJSON:   `{}`,
			expectedErr: createErrorInvalidSyntax(6, `unrecognized input`, `@()`),
		},
	}

	runTestCases(t, "TestAggregateFunction_SyntaxErrors", testCases)
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
				expectedErr: createErrorFunctionNotFound(`.func()`),
			},
		},
		{
			name: "unknown function call",
			testCase: TestCase{
				jsonpath:    `$.func(){}`,
				inputJSON:   `{}`,
				expectedErr: createErrorFunctionNotFound(`.func()`),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestCase(t, tt.testCase, tt.name)
		})
	}
}
