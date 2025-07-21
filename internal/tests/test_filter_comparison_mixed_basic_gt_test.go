package tests

import (
	"testing"
)

func TestFilterComparisonGT_BasicOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a > 1)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":2}]`,
			expectedJSON: `[{"a":2}]`,
		},
		{
			jsonpath:     `$[?(@.a > 1.5)]`,
			inputJSON:    `[{"a":0.5},{"a":1},{"a":1.5},{"a":2}]`,
			expectedJSON: `[{"a":2}]`,
		},
		{
			jsonpath:     `$[?(@.value > 10)]`,
			inputJSON:    `[{"value":5},{"value":10},{"value":15}]`,
			expectedJSON: `[{"value":15}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonGT_BasicOperations", testCases)
}

func TestFilterComparisonGT_WithLiterals(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(1 < @.a)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":2}]`,
			expectedJSON: `[{"a":2}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonGT_WithLiterals", testCases)
}

func TestFilterComparisonGT_ErrorCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@.a > $.b)]`,
			inputJSON:   `[{"a":1},{"a":2}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.a > $.b)]`),
		},
	}

	runTestCases(t, "TestFilterComparisonGT_ErrorCases", testCases)
}
