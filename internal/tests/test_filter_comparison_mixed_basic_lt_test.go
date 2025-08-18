package tests

import (
	"testing"
)

func TestFilterComparisonLT_BasicNumberOperations(t *testing.T) {
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
			jsonpath:     `$[?(1 > @.a)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":2}]`,
			expectedJSON: `[{"a":0}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonLT_BasicNumberOperations", testCases)
}

func TestFilterComparisonLT_BasicStringOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a < "ABD")]`,
			inputJSON:    `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"},{"b":"abc"},{"b":"abd"},{"b":"abe"}]`,
			expectedJSON: `[{"a":"ABC"}]`,
		},
		{
			jsonpath:     `$[?(@.a < "abd")]`,
			inputJSON:    `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"}]`,
			expectedJSON: `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"}]`,
		},
		{
			jsonpath:     `$[?("abd" > @.a)]`,
			inputJSON:    `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"}]`,
			expectedJSON: `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonLT_BasicStringOperations", testCases)
}

func TestFilterComparisonLT_BasicJSONPathNumberOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a < $[1].a)]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3},{"b":1},{"b":2},{"b":3}]`,
			expectedJSON: `[{"a":1}]`,
		},
		{
			jsonpath:     `$[?($[1].a > @.a)]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3},{"b":1},{"b":2},{"b":3}]`,
			expectedJSON: `[{"a":1}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonLT_BasicJSONPathNumberOperations", testCases)
}

func TestFilterComparisonLT_BasicJSONPathStringOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a < $[1].a)]`,
			inputJSON:    `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"},{"b":"abc"},{"b":"abd"},{"b":"abe"}]`,
			expectedJSON: `[{"a":"ABC"}]`,
		},
		{
			jsonpath:     `$[?($[1].a > @.a)]`,
			inputJSON:    `[{"a":"ABC"},{"a":"ABD"},{"a":"ABE"},{"a":"abc"},{"a":"abd"},{"a":"abe"}]`,
			expectedJSON: `[{"a":"ABC"}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonLT_BasicJSONPathStringOperations", testCases)
}
func TestFilterComparisonLT_BasicLiteralVariations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a < 124)]`,
			inputJSON:    `[{"a":"ABE"},{"a":null},{"a":123},{"a":true},{"a":[]},{"a":{}}]`,
			expectedJSON: `[{"a":123}]`,
		},
		{
			jsonpath:     `$[?(@.a < "ABF")]`,
			inputJSON:    `[{"a":"ABE"},{"a":null},{"a":123},{"a":true},{"a":[]},{"a":{}}]`,
			expectedJSON: `[{"a":"ABE"}]`,
		},
		{
			jsonpath:     `$[?(@.a < $[0].a)]`,
			inputJSON:    `[{"a":"ABE"},{"a":"ABD"},{"a":null},{"a":123},{"a":true},{"a":[]},{"a":{}}]`,
			expectedJSON: `[{"a":"ABD"}]`,
		},
		{
			jsonpath:     `$[?(@.a < $[2].a)]`,
			inputJSON:    `[{"a":"ABE"},{"a":null},{"a":123},{"a":122},{"a":true},{"a":[]},{"a":{}}]`,
			expectedJSON: `[{"a":122}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonLT_BasicLiteralVariations", testCases)
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
