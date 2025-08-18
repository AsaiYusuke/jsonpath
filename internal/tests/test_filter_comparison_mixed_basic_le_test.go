package tests

import (
	"testing"
)

func TestFilterComparisonLE_BasicNumberOperations(t *testing.T) {
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
			jsonpath:     `$[?(1.00001 >= @.a)]`,
			inputJSON:    `[{"a":0.5},{"a":1},{"a":1.00001},{"a":1.1}]`,
			expectedJSON: `[{"a":0.5},{"a":1},{"a":1.00001}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonLE_BasicNumberOperations", testCases)
}

func TestFilterComparisonLE_BasicStringOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a <= "ABD")]`,
			inputJSON:    `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"},{"b":"abc"},{"b":"abd"},{"b":"abe"}]`,
			expectedJSON: `[{"a":"ABC"},{"a":"ABD"}]`,
		},
		{
			jsonpath:     `$[?(@.a <= "abd")]`,
			inputJSON:    `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"}]`,
			expectedJSON: `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"}]`,
		},
		{
			jsonpath:     `$[?("abd" >= @.a)]`,
			inputJSON:    `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"}]`,
			expectedJSON: `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonLE_BasicStringOperations", testCases)
}

func TestFilterComparisonLE_BasicJSONPathNumberOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a <= $[1].a)]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3},{"b":1},{"b":2},{"b":3}]`,
			expectedJSON: `[{"a":1},{"a":2}]`,
		},
		{
			jsonpath:     `$[?($[1].a >= @.a)]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3},{"b":1},{"b":2},{"b":3}]`,
			expectedJSON: `[{"a":1},{"a":2}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonLE_BasicJSONPathNumberOperations", testCases)
}

func TestFilterComparisonLE_BasicJSONPathStringOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a <= $[1].a)]`,
			inputJSON:    `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"},{"b":"abc"},{"b":"abd"},{"b":"abe"}]`,
			expectedJSON: `[{"a":"ABC"},{"a":"ABD"}]`,
		},
		{
			jsonpath:     `$[?($[1].a >= @.a)]`,
			inputJSON:    `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"}]`,
			expectedJSON: `[{"a":"ABC"},{"a":"ABD"}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonLE_BasicJSONPathStringOperations", testCases)
}

func TestFilterComparisonLE_BasicLiteralVariations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a <= 123)]`,
			inputJSON:    `[{"a":"ABE"},{"a":null},{"a":123},{"a":true},{"a":[]},{"a":{}}]`,
			expectedJSON: `[{"a":123}]`,
		},
		{
			jsonpath:     `$[?(@.a <= "ABE")]`,
			inputJSON:    `[{"a":"ABE"},{"a":null},{"a":123},{"a":true},{"a":[]},{"a":{}}]`,
			expectedJSON: `[{"a":"ABE"}]`,
		},
		{
			jsonpath:     `$[?(@.a <= $[0].a)]`,
			inputJSON:    `[{"a":"ABE"},{"a":null},{"a":123},{"a":true},{"a":[]},{"a":{}}]`,
			expectedJSON: `[{"a":"ABE"}]`,
		},
		{
			jsonpath:     `$[?(@.a <= $[2].a)]`,
			inputJSON:    `[{"a":"ABE"},{"a":null},{"a":123},{"a":true},{"a":[]},{"a":{}}]`,
			expectedJSON: `[{"a":123}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonLE_BasicLiteralVariations", testCases)
}
