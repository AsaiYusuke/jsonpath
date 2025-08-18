package tests

import (
	"testing"
)

func TestFilterComparisonGE_BasicNumberOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a >= 1)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":2},{"b":0},{"b":1},{"b":2}]`,
			expectedJSON: `[{"a":1},{"a":2}]`,
		},
		{
			jsonpath:     `$[?(@.a >= 1.000001)]`,
			inputJSON:    `[{"a":0.5},{"a":1},{"a":1.000001},{"a":1.1}]`,
			expectedJSON: `[{"a":1.000001},{"a":1.1}]`,
		},
		{
			jsonpath:     `$[?(1.000001 <= @.a)]`,
			inputJSON:    `[{"a":0.5},{"a":1},{"a":1.000001},{"a":1.1}]`,
			expectedJSON: `[{"a":1.000001},{"a":1.1}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonGE_BasicNumberOperations", testCases)
}

func TestFilterComparisonGE_BasicStringOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a >= "ABD")]`,
			inputJSON:    `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"},{"b":"abc"},{"b":"abd"},{"b":"abe"}]`,
			expectedJSON: `[{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"}]`,
		},
		{
			jsonpath:     `$[?(@.a >= "abd")]`,
			inputJSON:    `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"}]`,
			expectedJSON: `[{"a":"abd"},{"a":"abe"}]`,
		},
		{
			jsonpath:     `$[?("abd" <= @.a)]`,
			inputJSON:    `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"}]`,
			expectedJSON: `[{"a":"abd"},{"a":"abe"}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonGE_BasicStringOperations", testCases)
}

func TestFilterComparisonGE_BasicJSONPathNumberOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a >= $[1].a)]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3},{"b":1},{"b":2},{"b":3}]`,
			expectedJSON: `[{"a":2},{"a":3}]`,
		},
		{
			jsonpath:     `$[?($[1].a <= @.a)]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3},{"b":1},{"b":2},{"b":3}]`,
			expectedJSON: `[{"a":2},{"a":3}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonGE_BasicJSONPathNumberOperations", testCases)
}

func TestFilterComparisonGE_BasicJSONPathStringOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a >= $[1].a)]`,
			inputJSON:    `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"},{"b":"abc"},{"b":"abd"},{"b":"abe"}]`,
			expectedJSON: `[{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"}]`,
		},
		{
			jsonpath:     `$[?($[1].a <= @.a)]`,
			inputJSON:    `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"}]`,
			expectedJSON: `[{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonGE_BasicJSONPathNumberOperations", testCases)
}

func TestFilterComparisonGE_BasicLiteralVariations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a >= 1)]`,
			inputJSON:    `[{"a":"ABE"},{"a":null},{"a":123},{"a":true},{"a":[]},{"a":{}}]`,
			expectedJSON: `[{"a":123}]`,
		},
		{
			jsonpath:     `$[?(@.a >= "ABD")]`,
			inputJSON:    `[{"a":"ABE"},{"a":null},{"a":123},{"a":true},{"a":[]},{"a":{}}]`,
			expectedJSON: `[{"a":"ABE"}]`,
		},
		{
			jsonpath:     `$[?(@.a >= $[0].a)]`,
			inputJSON:    `[{"a":"ABE"},{"a":null},{"a":123},{"a":true},{"a":[]},{"a":{}}]`,
			expectedJSON: `[{"a":"ABE"}]`,
		},
		{
			jsonpath:     `$[?(@.a >= $[2].a)]`,
			inputJSON:    `[{"a":"ABE"},{"a":null},{"a":123},{"a":true},{"a":[]},{"a":{}}]`,
			expectedJSON: `[{"a":123}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonGE_BasicLiteralVariations", testCases)
}
