package jsonpath

import (
	"fmt"
	"testing"
)

func TestBracketNotationMixedInvalidSyntax(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$['a\c']`,
			inputJSON:   `{"ac":1,"b":2}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['a\c']`},
		},
		{
			jsonpath:    `$['a'c']`,
			inputJSON:   `{"ac":1,"b":2}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['a'c']`},
		},

		{
			jsonpath:    `$["a\c"]`,
			inputJSON:   `{"ac":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `["a\c"]`},
		},
		{
			jsonpath:    `$["a"b"]`,
			inputJSON:   `{"ab":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `["a"b"]`},
		},
	}

	for i, tc := range tests {
		runTestCase(t, tc, fmt.Sprintf("TestBracketNotationMixedInvalidSyntax_case_%d", i))
	}
}
