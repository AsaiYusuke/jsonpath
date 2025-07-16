package jsonpath

import (
	"fmt"
	"testing"
)

func TestInvalidSyntax_BracketNotationErrors(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$[0,]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[0,]`},
		},
		{
			jsonpath:    `$[0,a]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[0,a]`},
		},
		{
			jsonpath:    `$[a:1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[a:1]`},
		},
		{
			jsonpath:    `$[0:a]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[0:a]`},
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("BracketNotationErrors_%d", i), test)
	}
}

func TestInvalidSyntax_FunctionErrors(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$.func(`,
			inputJSON:   `{}`,
			expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `(`},
		},
	}

	for i, testCase := range testCases {
		runSingleTestCase(t, fmt.Sprintf("FunctionError_%d", i), testCase)
	}
}
