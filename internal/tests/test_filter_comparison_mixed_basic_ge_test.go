package tests

import (
	"testing"
)

func TestFilterComparisonGE_BasicOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a >= 1)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":2}]`,
			expectedJSON: `[{"a":1},{"a":2}]`,
		},
		{
			jsonpath:     `$[?(@.a >= 1.000001)]`,
			inputJSON:    `[{"a":0.5},{"a":1},{"a":1.000001},{"a":1.1}]`,
			expectedJSON: `[{"a":1.000001},{"a":1.1}]`,
		},
		{
			jsonpath:     `$[?(@.value >= 10)]`,
			inputJSON:    `[{"value":5},{"value":10},{"value":15}]`,
			expectedJSON: `[{"value":10},{"value":15}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonGE_BasicOperations", testCases)
}

func TestFilterComparisonGE_WithLiterals(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(1.000001 <= @.a)]`,
			inputJSON:    `[{"a":0.5},{"a":1},{"a":1.000001},{"a":1.1}]`,
			expectedJSON: `[{"a":1.000001},{"a":1.1}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonGE_WithLiterals", testCases)
}
