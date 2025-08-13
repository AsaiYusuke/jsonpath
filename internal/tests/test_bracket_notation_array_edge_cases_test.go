package tests

import (
	"testing"
)

func TestBracketNotation_ArrayIndexEdgeCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[1]`,
			inputJSON:   `[]`,
			expectedErr: createErrorMemberNotExist(`[1]`),
		},
	}

	runTestCases(t, "TestBracketNotation_ArrayIndexEdgeCases", testCases)
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

	runTestCases(t, "TestBracketNotation_PropertyTypeMismatchErrors", testCases)
}

func TestBracketNotation_MixedDotPropertyAccess(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.a['b']`,
		inputJSON:    `{"b":2,"a":{"b":1}}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestBracketNotation_MixedDotPropertyAccess")
}
