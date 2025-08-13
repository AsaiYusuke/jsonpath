package tests

import (
	"testing"
)

func TestInvalidSyntax_BracketNotationErrors(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$[0,]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[0,]`),
		},
		{
			jsonpath:    `$[0,a]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[0,a]`),
		},
		{
			jsonpath:    `$[a:1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[a:1]`),
		},
		{
			jsonpath:    `$[0:a]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[0:a]`),
		},
	}

	runTestCases(t, "TestInvalidSyntax_BracketNotationErrors", tests)
}

func TestInvalidSyntax_FunctionErrors(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$.func(`,
			inputJSON:   `{}`,
			expectedErr: createErrorInvalidSyntax(6, `unrecognized input`, `(`),
		},
	}

	runTestCases(t, "TestInvalidSyntax_FunctionErrors", testCases)
}
