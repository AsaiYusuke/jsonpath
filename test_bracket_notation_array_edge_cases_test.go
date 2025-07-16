package jsonpath

import (
	"fmt"
	"testing"
)

func TestBracketNotation_ArrayIndexEdgeCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[1]`,
			inputJSON:   `[]`,
			expectedErr: createErrorMemberNotExist(`[1]`),
		},
		{
			jsonpath:    `$[1000000000000000000]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[1000000000000000000]`),
		},
	}

	for i, testCase := range testCases {
		runTestCase(t, testCase, fmt.Sprintf("TestBracketNotation_ArrayIndexEdgeCases[%d]", i))
	}
}

func TestBracketNotation_PropertyTypeMismatchErrors(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$['a']`,
			inputJSON:   `[]`,
			expectedErr: createErrorTypeUnmatched(`['a']`, `object`, `[]interface {}`),
		},
		{
			jsonpath:    `$['a']`,
			inputJSON:   `{}`,
			expectedErr: createErrorMemberNotExist(`['a']`),
		},
		{
			jsonpath:    `$['a']`,
			inputJSON:   `123`,
			expectedErr: createErrorTypeUnmatched(`['a']`, `object`, `float64`),
		},
		{
			jsonpath:    `$['a']`,
			inputJSON:   `true`,
			expectedErr: createErrorTypeUnmatched(`['a']`, `object`, `bool`),
		},
		{
			jsonpath:    `$['a']`,
			inputJSON:   `null`,
			expectedErr: createErrorTypeUnmatched(`['a']`, `object`, `null`),
		},
	}

	for i, testCase := range testCases {
		runTestCase(t, testCase, fmt.Sprintf("TestBracketNotation_PropertyTypeMismatchErrors[%d]", i))
	}
}

func TestBracketNotation_MixedDotPropertyAccess(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.a['b']`,
		inputJSON:    `{"b":2,"a":{"b":1}}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestBracketNotation_MixedDotPropertyAccess")
}
