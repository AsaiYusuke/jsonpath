package tests

import (
	"testing"
)

func TestFilterComparisonLT_BasicOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a < 1)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":2}]`,
			expectedJSON: `[{"a":0}]`,
		},
		{
			jsonpath:     `$[?(@.a < 1.5)]`,
			inputJSON:    `[{"a":0.5},{"a":1},{"a":1.5},{"a":2}]`,
			expectedJSON: `[{"a":0.5},{"a":1}]`,
		},
		{
			jsonpath:     `$[?(@.value < 10)]`,
			inputJSON:    `[{"value":5},{"value":10},{"value":15}]`,
			expectedJSON: `[{"value":5}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonLT_BasicOperations", testCases)
}

func TestFilterComparisonLT_WithLiterals(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(1 > @.a)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":2}]`,
			expectedJSON: `[{"a":0}]`,
		},
		{
			jsonpath:     `$[?(1 < @.a)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":2}]`,
			expectedJSON: `[{"a":2}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonLT_WithLiterals", testCases)
}

func TestFilterComparisonLT_ErrorCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(1 < @.a)]`,
			inputJSON:   `[{"b":0},{"b":1},{"b":2}]`,
			expectedErr: createErrorMemberNotExist(`[?(1 < @.a)]`),
		},
	}

	runTestCases(t, "TestFilterComparisonLT_ErrorCases", testCases)
}

func TestFilterComparisonLT_WithRootReference(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@.a < $.b)]`,
			inputJSON:   `[{"a":1},{"a":2}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.a < $.b)]`),
		},
	}

	runTestCases(t, "TestFilterComparisonLT_WithRootReference", testCases)
}
