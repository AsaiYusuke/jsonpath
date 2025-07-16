package jsonpath

import (
	"fmt"
	"testing"
)

func TestBracketNotationInvalidSyntax(t *testing.T) {
	tests := []TestCase{
		// Invalid escape sequences in single quotes
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

		// Invalid escape sequences in double quotes
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
		runTestCase(t, tc, fmt.Sprintf("TestBracketNotationInvalidSyntax_case_%d", i))
	}
}

func TestBracketNotationMemberNotExist(t *testing.T) {
	tests := []TestCase{
		// Only keeping unique test cases
	}

	for i, tc := range tests {
		runTestCase(t, tc, fmt.Sprintf("TestBracketNotationMemberNotExist_case_%d", i))
	}
}

func TestBracketNotationSpecialCharacters(t *testing.T) {
	tests := []TestCase{
		// This test should be unique and not duplicate any in other files
	}

	for i, tc := range tests {
		runTestCase(t, tc, fmt.Sprintf("TestBracketNotationSpecialCharacters_case_%d", i))
	}
}

func TestBracketNotationTypeErrors(t *testing.T) {
	tests := []TestCase{
		// Only keeping unique test cases that don't exist elsewhere
	}

	for i, tc := range tests {
		runTestCase(t, tc, fmt.Sprintf("TestBracketNotationTypeErrors_case_%d", i))
	}
}
