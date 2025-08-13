package tests

import (
	"testing"
)

func TestFilterComparisonLE_BasicOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a <= 1)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":2}]`,
			expectedJSON: `[{"a":0},{"a":1}]`,
		},
		{
			jsonpath:     `$[?(@.a <= 1.00001)]`,
			inputJSON:    `[{"a":0.5},{"a":1},{"a":1.00001},{"a":1.1}]`,
			expectedJSON: `[{"a":0.5},{"a":1},{"a":1.00001}]`,
		},
		{
			jsonpath:     `$[?(@.value <= 10)]`,
			inputJSON:    `[{"value":5},{"value":10},{"value":15}]`,
			expectedJSON: `[{"value":5},{"value":10}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonLE_BasicOperations", testCases)
}

func TestFilterComparisonLE_WithLiterals(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(1.00001 >= @.a)]`,
			inputJSON:    `[{"a":0.5},{"a":1},{"a":1.00001},{"a":1.1}]`,
			expectedJSON: `[{"a":0.5},{"a":1},{"a":1.00001}]`,
		},
		{
			jsonpath:     `$[?(1.000001 <= @.a)]`,
			inputJSON:    `[{"a":0.5},{"a":1},{"a":1.000001},{"a":1.1}]`,
			expectedJSON: `[{"a":1.000001},{"a":1.1}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonLE_WithLiterals", testCases)
}
