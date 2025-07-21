package tests

import (
	"testing"
)

func TestInvalidSyntax_basicSyntaxDeletedCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$$`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `$`),
		},
		{
			jsonpath:     `a.`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `.`),
		},
		{
			jsonpath:     `b.`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `.`),
		},
		{
			jsonpath:     `$a`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `a`),
		},
		{
			jsonpath:     `.`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(0, `unrecognized input`, `.`),
		},
		{
			jsonpath:     `$.`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `.`),
		},
		{
			jsonpath:     `..`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(0, `unrecognized input`, `..`),
		},
		{
			jsonpath:     `$..`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `..`),
		},
		{
			jsonpath:     `$.a.`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(3, `unrecognized input`, `.`),
		},
		{
			jsonpath:     `$.a..`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(3, `unrecognized input`, `..`),
		},
		{
			jsonpath:     `$..a.`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `unrecognized input`, `.`),
		},
		{
			jsonpath:     `$..a..`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `unrecognized input`, `..`),
		},
		{
			jsonpath:     `$...a`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `...a`),
		},
	}

	runTestCases(t, "TestInvalidSyntax_basicSyntaxDeletedCases", testCases)
}

func TestInvalidSyntax_bracketNotationDeletedCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$['a]`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `['a]`),
		},
		{
			jsonpath:     `$["a]`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `["a]`),
		},
		{
			jsonpath:     `$[a']`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `[a']`),
		},
		{
			jsonpath:     `$[a"]`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `[a"]`),
		},
		{
			jsonpath:     `$[a]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `[a]`),
		},
		{
			jsonpath:     `$.[a]`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `.[a]`),
		},
		{
			jsonpath:     `$['a'.'b']`,
			inputJSON:    `["a"]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `['a'.'b']`),
		},
		{
			jsonpath:     `$[a.b]`,
			inputJSON:    `[{"a":{"b":1}}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `[a.b]`),
		},
		{
			jsonpath:     `$['a'b']`,
			inputJSON:    `["a"]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `['a'b']`),
		},
		{
			jsonpath:     `$.a[]`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(3, `unrecognized input`, `[]`),
		},
		{
			jsonpath:     `$.a.b[]`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(5, `unrecognized input`, `[]`),
		},
	}

	runTestCases(t, "TestInvalidSyntax_bracketNotationDeletedCases", testCases)
}

func TestInvalidSyntax_filterInvalidDeletedCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?()]`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `[?()]`),
		},
		{
			jsonpath:     `$[()]`,
			inputJSON:    `{}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `[()]`),
		},
		{
			jsonpath:     `$()`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `()`),
		},
		{
			jsonpath:     `$(a)`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `(a)`),
		},
		{
			jsonpath:     `$[(`,
			inputJSON:    `{}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `[(`),
		},
		{
			jsonpath:     `$[(]`,
			inputJSON:    `{}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `[(]`),
		},
		{
			jsonpath:     `$[0`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `[0`),
		},
		{
			jsonpath:     `$[?@a]`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `[?@a]`),
		},
		{
			jsonpath:     `$[0,10000000000000000000,]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `[0,10000000000000000000,]`),
		},
		{
			jsonpath:     `$[?(<@.a)]`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `[?(<@.a)]`),
		},
	}

	runTestCases(t, "TestInvalidSyntax_filterInvalidDeletedCases", testCases)
}
