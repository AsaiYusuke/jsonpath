package jsonpath

import (
	"fmt"
	"testing"
)

// TestInvalidSyntax_wildcardAccessorDeletedCases tests deleted invalid syntax cases for wildcard accessor
func TestInvalidSyntax_wildcardAccessorDeletedCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[*].*`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[*]`),
		},
		{
			jsonpath:     `$[*][0:2]`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[*]`),
		},
		{
			jsonpath:     `$[*][*]`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[*]`),
		},
		{
			jsonpath:     `$[*][?(@.b)]`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[*]`),
		},
	}

	for i, testCase := range testCases {
		runSingleTestCase(t, fmt.Sprintf("WildcardAccessorDeleted_%d", i), testCase)
	}
}

// TestInvalidSyntax_filterChainDeletedCases tests deleted invalid syntax cases for filter chains
func TestInvalidSyntax_filterChainDeletedCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.b)]..a`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[?(@.b)]`, `object/array`, `string`),
		},
		{
			jsonpath:     `$[?(@.b)].*`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[?(@.b)]`, `object/array`, `string`),
		},
		{
			jsonpath:     `$[?(@)][0:2]`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?(@)]`),
		},
		{
			jsonpath:     `$[?(@)][0:2]`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[?(@)]`, `object/array`, `string`),
		},
		{
			jsonpath:     `$[?(@)][*]`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?(@)]`),
		},
		{
			jsonpath:     `$[?(@)][*]`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[?(@)]`, `object/array`, `string`),
		},
		{
			jsonpath:     `$[?(@.a)][?(@.b)]`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[?(@.a)]`, `object/array`, `string`),
		},
	}

	for i, testCase := range testCases {
		runSingleTestCase(t, fmt.Sprintf("FilterChainDeleted_%d", i), testCase)
	}
}

// TestInvalidSyntax_basicSyntaxDeletedCases tests deleted invalid syntax cases for basic syntax errors
func TestInvalidSyntax_basicSyntaxDeletedCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$$`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `$`},
		},
		{
			jsonpath:     `a.`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.`},
		},
		{
			jsonpath:     `b.`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.`},
		},
		{
			jsonpath:     `$a`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `a`},
		},
		{
			jsonpath:     `.`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 0, reason: `unrecognized input`, near: `.`},
		},
		{
			jsonpath:     `$.`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.`},
		},
		{
			jsonpath:     `..`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 0, reason: `unrecognized input`, near: `..`},
		},
		{
			jsonpath:     `$..`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `..`},
		},
		{
			jsonpath:     `$.a.`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: `.`},
		},
		{
			jsonpath:     `$.a..`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: `..`},
		},
		{
			jsonpath:     `$..a.`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `unrecognized input`, near: `.`},
		},
		{
			jsonpath:     `$..a..`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `unrecognized input`, near: `..`},
		},
		{
			jsonpath:     `$...a`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `...a`},
		},
	}

	for i, testCase := range testCases {
		runSingleTestCase(t, fmt.Sprintf("BasicSyntaxDeleted_%d", i), testCase)
	}
}

// TestInvalidSyntax_bracketNotationDeletedCases tests deleted invalid syntax cases for bracket notation
func TestInvalidSyntax_bracketNotationDeletedCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$['a]`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['a]`},
		},
		{
			jsonpath:     `$["a]`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `["a]`},
		},
		{
			jsonpath:     `$[a']`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[a']`},
		},
		{
			jsonpath:     `$[a"]`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[a"]`},
		},
		{
			jsonpath:     `$[a]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[a]`},
		},
		{
			jsonpath:     `$.[a]`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.[a]`},
		},
		{
			jsonpath:     `$['a'.'b']`,
			inputJSON:    `["a"]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['a'.'b']`},
		},
		{
			jsonpath:     `$[a.b]`,
			inputJSON:    `[{"a":{"b":1}}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[a.b]`},
		},
		{
			jsonpath:     `$['a'b']`,
			inputJSON:    `["a"]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['a'b']`},
		},
		{
			jsonpath:     `$.a[]`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: `[]`},
		},
		{
			jsonpath:     `$.a.b[]`,
			inputJSON:    `{"a":1}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 5, reason: `unrecognized input`, near: `[]`},
		},
	}

	for i, testCase := range testCases {
		runSingleTestCase(t, fmt.Sprintf("BracketNotationDeleted_%d", i), testCase)
	}
}

// TestInvalidSyntax_filterInvalidDeletedCases tests deleted invalid syntax cases for filter expressions
func TestInvalidSyntax_filterInvalidDeletedCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?()]`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?()]`},
		},
		{
			jsonpath:     `$[()]`,
			inputJSON:    `{}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[()]`},
		},
		{
			jsonpath:     `$()`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `()`},
		},
		{
			jsonpath:     `$(a)`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `(a)`},
		},
		{
			jsonpath:     `$[(`,
			inputJSON:    `{}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[(`},
		},
		{
			jsonpath:     `$[(]`,
			inputJSON:    `{}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[(]`},
		},
		{
			jsonpath:     `$[0`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[0`},
		},
		{
			jsonpath:     `$[?@a]`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?@a]`},
		},
		{
			jsonpath:     `$[0,10000000000000000000,]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[0,10000000000000000000,]`},
		},
		{
			jsonpath:     `$[?(<@.a)]`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(<@.a)]`},
		},
	}

	for i, testCase := range testCases {
		runSingleTestCase(t, fmt.Sprintf("FilterInvalidDeleted_%d", i), testCase)
	}
}
