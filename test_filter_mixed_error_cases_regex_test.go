package jsonpath

import (
	"fmt"
	"testing"
)

func TestFilterRegexErrors(t *testing.T) {
	testCases := []TestCase{
		// Missing regex field error
		{
			jsonpath:    `$[?(@.a=~/abc)]`,
			inputJSON:   `[]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=~/abc)]`},
		},

		// Invalid regex syntax - missing @ prefix
		{
			jsonpath:    `$[?(a=~/123/)]`,
			inputJSON:   `[{"a":"123"},{"a":123}]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(a=~/123/)]`},
		},

		// Invalid regex pattern - unsupported Perl syntax
		{
			jsonpath:    `$[?(@.a=~/(?x)CASE/)]`,
			inputJSON:   `[{"a":"case"},{"a":"CASE"},{"a":"Case"},{"a":"abc"}]`,
			expectedErr: ErrorInvalidArgument{argument: `(?x)CASE`, err: fmt.Errorf("error parsing regexp: invalid or unsupported Perl syntax: `(?x`")},
		},
	}

	for i, tc := range testCases {
		runTestCase(t, tc, fmt.Sprintf("TestFilterRegexErrors_case_%d", i))
	}
}
