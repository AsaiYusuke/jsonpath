package tests

import (
	"testing"
)

func TestFilterComparisonGT_BasicNumberOperations(t *testing.T) {
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
			jsonpath:     `$[?(1 < @.a)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":2}]`,
			expectedJSON: `[{"a":2}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonGT_BasicNumberOperations", testCases)
}

func TestFilterComparisonGT_BasicStringOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a > "ABD")]`,
			inputJSON:    `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"},{"b":"abc"},{"b":"abd"},{"b":"abe"}]`,
			expectedJSON: `[{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"}]`,
		},
		{
			jsonpath:     `$[?(@.a > "abd")]`,
			inputJSON:    `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"}]`,
			expectedJSON: `[{"a":"abe"}]`,
		},
		{
			jsonpath:     `$[?("abd" < @.a)]`,
			inputJSON:    `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"}]`,
			expectedJSON: `[{"a":"abe"}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonGT_BasicStringOperations", testCases)
}

func TestFilterComparisonGT_BasicJSONPathNumberOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a > $[1].a)]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3},{"b":1},{"b":2},{"b":3}]`,
			expectedJSON: `[{"a":3}]`,
		},
		{
			jsonpath:     `$[?($[1].a < @.a)]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3},{"b":1},{"b":2},{"b":3}]`,
			expectedJSON: `[{"a":3}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonGT_BasicJSONPathNumberOperations", testCases)
}

func TestFilterComparisonGT_BasicJSONPathStringOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a > $[1].a)]`,
			inputJSON:    `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"},{"b":"abc"},{"b":"abd"},{"b":"abe"}]`,
			expectedJSON: `[{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"}]`,
		},
		{
			jsonpath:     `$[?($[1].a < @.a)]`,
			inputJSON:    `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"}]`,
			expectedJSON: `[{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonGT_BasicJSONPathStringOperations", testCases)
}

func TestFilterComparisonGT_BasicLiteralVariations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a > 1)]`,
			inputJSON:    `[{"a":"ABE"},{"a":null},{"a":123},{"a":true},{"a":[]},{"a":{}}]`,
			expectedJSON: `[{"a":123}]`,
		},
		{
			jsonpath:     `$[?(@.a > "ABD")]`,
			inputJSON:    `[{"a":"ABE"},{"a":null},{"a":123},{"a":true},{"a":[]},{"a":{}}]`,
			expectedJSON: `[{"a":"ABE"}]`,
		},
		{
			jsonpath:     `$[?(@.a > $[0].a)]`,
			inputJSON:    `[{"a":"ABD"},{"a":"ABE"},{"a":null},{"a":123},{"a":true},{"a":[]},{"a":{}}]`,
			expectedJSON: `[{"a":"ABE"}]`,
		},
		{
			jsonpath:     `$[?(@.a > $[2].a)]`,
			inputJSON:    `[{"a":"ABE"},{"a":null},{"a":123},{"a":124},{"a":true},{"a":[]},{"a":{}}]`,
			expectedJSON: `[{"a":124}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonGT_BasicLiteralVariations", testCases)
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
