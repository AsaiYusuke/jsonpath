package jsonpath

import (
	"testing"
)

// TestRetrieve_filterComplexNestedDeleted tests deleted complex nested filter cases
func TestRetrieve_filterComplexNestedDeleted(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@[1][0]>1)]`,
			inputJSON:    `[1,[21,[221,[222]]]]`,
			expectedJSON: `[[21,[221,[222]]]]`,
		},
		{
			jsonpath:     `$[?(@[1][0]>1)][?(@[1][0]>1)]`,
			inputJSON:    `[1,[21,[221,[222]]]]`,
			expectedJSON: `[[221,[222]]]`,
		},
		{
			jsonpath:     `$[?(@[1][0]>1)][?(@[1][0]>1)][?(@[0]>1)]`,
			inputJSON:    `[1,[21,[221,[222]]]]`,
			expectedJSON: `[[222]]`,
		},
		{
			jsonpath:     `$[?(@[1][0]>1)][?(@[1][0]>1)][?(@[1]>1)]`,
			inputJSON:    `[1,[21,[221,[222]]]]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?(@[1]>1)]`),
		},
	}

	for i, testCase := range testCases {
		runTestCase(t, testCase, "TestRetrieve_filterComplexNestedDeleted_case_"+string(rune('A'+i)))
	}
}

// TestRetrieve_filterRootComparisonErrors tests deleted filter cases with root comparison errors
func TestRetrieve_filterRootComparisonErrors(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a == $.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?(@.a == $.b)]`),
		},
		{
			jsonpath:     `$[?(@.a != $.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: `[{"a":0},{"a":1}]`,
		},
		{
			jsonpath:     `$[?(@.a <= $.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?(@.a <= $.b)]`),
		},
		{
			jsonpath:     `$[?($.b == @.a)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?($.b == @.a)]`),
		},
	}

	for i, testCase := range testCases {
		runTestCase(t, testCase, "TestRetrieve_filterRootComparisonErrors_case_"+string(rune('A'+i)))
	}
}

// TestRetrieve_filterRootComparisonMore tests additional deleted filter cases with root comparison
func TestRetrieve_filterRootComparisonMore(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?($.b != @.a)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: `[{"a":0},{"a":1}]`,
		},
		{
			jsonpath:     `$[?(@.b != $[0].a)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: `[{"a":0},{"a":1}]`,
		},
		{
			jsonpath:     `$[?($[0].a != @.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: `[{"a":0},{"a":1}]`,
		},
		{
			jsonpath:     `$[?(@.b == $.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: `[{"a":0},{"a":1}]`,
		},
		{
			jsonpath:     `$[?($.b == @.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: `[{"a":0},{"a":1}]`,
		},
		{
			jsonpath:     `$[?($.b < @.a)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?($.b < @.a)]`),
		},
		{
			jsonpath:     `$[?(@.b < $[0].a)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?(@.b < $[0].a)]`),
		},
		{
			jsonpath:     `$[?($[0].a < @.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?($[0].a < @.b)]`),
		},
		{
			jsonpath:     `$[?(@.b != $.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?(@.b != $.b)]`),
		},
		{
			jsonpath:     `$[?($.b != @.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?($.b != @.b)]`),
		},
		{
			jsonpath:     `$[?(@.b < $.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?(@.b < $.b)]`),
		},
		{
			jsonpath:     `$[?(@.b <= $.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?(@.b <= $.b)]`),
		},
		{
			jsonpath:     `$[?(@.b > $.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?(@.b > $.b)]`),
		},
		{
			jsonpath:     `$[?(@.b >= $.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?(@.b >= $.b)]`),
		},
		{
			jsonpath:     `$[?($.b < @.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?($.b < @.b)]`),
		},
		{
			jsonpath:     `$[?($.b <= @.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?($.b <= @.b)]`),
		},
		{
			jsonpath:     `$[?($.b > @.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?($.b > @.b)]`),
		},
		{
			jsonpath:     `$[?($.b >= @.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?($.b >= @.b)]`),
		},
		{
			jsonpath:     `$[?(@.b == $[0].a)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?(@.b == $[0].a)]`),
		},
		{
			jsonpath:     `$[?(@.b <= $[0].a)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?(@.b <= $[0].a)]`),
		},
		{
			jsonpath:     `$[?(@.b > $[0].a)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?(@.b > $[0].a)]`),
		},
		{
			jsonpath:     `$[?($[0].a == @.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?($[0].a == @.b)]`),
		},
		{
			jsonpath:     `$[?($[0].a <= @.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?($[0].a <= @.b)]`),
		},
		{
			jsonpath:     `$[?($[0].a > @.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?($[0].a > @.b)]`),
		},
		{
			jsonpath:     `$[?($[0].a >= @.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?($[0].a >= @.b)]`),
		},
	}

	for i, testCase := range testCases {
		runTestCase(t, testCase, "TestRetrieve_filterRootComparisonMore_case_"+string(rune('A'+i)))
	}
}

// TestRetrieve_filterRootComparisonAdditional tests remaining deleted filter cases with root comparison
func TestRetrieve_filterRootComparisonAdditional(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?($.b <= @.a)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?($.b <= @.a)]`),
		},
		{
			jsonpath:     `$[?($.b > @.a)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?($.b > @.a)]`),
		},
		{
			jsonpath:     `$[?($.b >= @.a)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?($.b >= @.a)]`),
		},
	}

	for i, testCase := range testCases {
		runTestCase(t, testCase, "TestRetrieve_filterRootComparisonAdditional_case_"+string(rune('A'+i)))
	}
}
