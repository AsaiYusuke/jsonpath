package jsonpath

import (
	"fmt"
	"testing"
)

func TestFunctionSyntaxErrors(t *testing.T) {
	testCases := []TestCase{
		// Invalid function argument syntax
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
		runTestCase(t, tc, fmt.Sprintf("TestFunctionSyntaxErrors_case_%d", i))
	}
}
