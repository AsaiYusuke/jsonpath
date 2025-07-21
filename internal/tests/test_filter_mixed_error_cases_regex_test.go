package tests

import (
	"fmt"
	"testing"
)

func TestFilterRegexErrors(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@.a=~/abc)]`,
			inputJSON:   `[]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a=~/abc)]`),
		},

		{
			jsonpath:    `$[?(a=~/123/)]`,
			inputJSON:   `[{"a":"123"},{"a":123}]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(a=~/123/)]`),
		},

		{
			jsonpath:    `$[?(@.a=~/(?x)CASE/)]`,
			inputJSON:   `[{"a":"case"},{"a":"CASE"},{"a":"Case"},{"a":"abc"}]`,
			expectedErr: createErrorInvalidArgument(`(?x)CASE`, fmt.Errorf("error parsing regexp: invalid or unsupported Perl syntax: `(?x`")),
		},
	}

	runTestCases(t, "TestFilterRegexErrors", testCases)
}

func TestFilterRegexSyntaxErrors(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@.a>1 && ())]`,
			inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a>1 && ())]`),
		},
		{
			jsonpath:    `$[?(@.a=~///)]`,
			inputJSON:   `[]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a=~///)]`),
		},
		{
			jsonpath:    `$[?(@.a=~s/a/b/)]`,
			inputJSON:   `[]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a=~s/a/b/)]`),
		},
		{
			jsonpath:    `$[?(@.a=~@abc@)]`,
			inputJSON:   `[]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a=~@abc@)]`),
		},
	}

	runTestCases(t, "TestFilterRegexSyntaxErrors", testCases)
}

func TestFilterRegexMemberAccessErrors(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@.a.b=~/abc/)]`,
			inputJSON:   `[{"a":"abc"}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.a.b=~/abc/)]`),
		},
	}

	runTestCases(t, "TestFilterRegexMemberAccessErrors", testCases)
}
