package tests

import (
	"testing"
)

func TestBracketNotationMixedInvalidSyntax(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$['a\c']`,
			inputJSON:   `{"ac":1,"b":2}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `['a\c']`),
		},
		{
			jsonpath:    `$['a'c']`,
			inputJSON:   `{"ac":1,"b":2}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `['a'c']`),
		},

		{
			jsonpath:    `$["a\c"]`,
			inputJSON:   `{"ac":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `["a\c"]`),
		},
		{
			jsonpath:    `$["a"b"]`,
			inputJSON:   `{"ab":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `["a"b"]`),
		},
	}

	runTestCases(t, "TestBracketNotationMixedInvalidSyntax", tests)
}
