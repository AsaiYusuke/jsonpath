package jsonpath

import (
	"fmt"
	"testing"
)

func TestFilterRegexErrors(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@.a=~/abc)]`,
			inputJSON:   `[]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=~/abc)]`},
		},

		{
			jsonpath:    `$[?(a=~/123/)]`,
			inputJSON:   `[{"a":"123"},{"a":123}]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(a=~/123/)]`},
		},

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

func TestFilterRegexSyntaxErrors(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@.a>1 && ())]`,
			inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a>1 && ())]`},
		},
		{
			jsonpath:    `$[?(@.a=~///)]`,
			inputJSON:   `[]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=~///)]`},
		},
		{
			jsonpath:    `$[?(@.a=~s/a/b/)]`,
			inputJSON:   `[]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=~s/a/b/)]`},
		},
		{
			jsonpath:    `$[?(@.a=~@abc@)]`,
			inputJSON:   `[]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=~@abc@)]`},
		},
	}

	for i, testCase := range testCases {
		runSingleTestCase(t, fmt.Sprintf("RegexSyntaxError_%d", i), testCase)
	}
}
