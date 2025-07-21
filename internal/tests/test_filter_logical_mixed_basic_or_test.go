package tests

import (
	"testing"
)

func TestFilterLogicalOr_BasicOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a || @.b)]`,
			inputJSON:    `[{"a":1},{"b":2},{"c":3},{"a":4,"b":5}]`,
			expectedJSON: `[{"a":1},{"b":2},{"a":4,"b":5}]`,
		},
		{
			jsonpath:     `$[?(@.a>2 || @.a<2)]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
			expectedJSON: `[{"a":1},{"a":3}]`,
		},
		{
			jsonpath:     `$[?(@.a<2 || @.a>2)]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
			expectedJSON: `[{"a":1},{"a":3}]`,
		},
	}

	runTestCases(t, "TestFilterLogicalOr_BasicOperations", testCases)
}

func TestFilterLogicalOr_WithLiterals(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?((1==2) || @.a>1)]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
			expectedJSON: `[{"a":2},{"a":3}]`,
		},
		{
			jsonpath:     `$[?((1==1) || @.a>1)]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
			expectedJSON: `[{"a":1},{"a":2},{"a":3}]`,
		},
		{
			jsonpath:     `$[?(@.a>1 || (1==2))]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
			expectedJSON: `[{"a":2},{"a":3}]`,
		},
		{
			jsonpath:     `$[?(@.a>1 || (1==1))]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
			expectedJSON: `[{"a":1},{"a":2},{"a":3}]`,
		},
	}

	runTestCases(t, "TestFilterLogicalOr_WithLiterals", testCases)
}

func TestFilterLogicalOr_PropertyExistence(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.x || @.b > 2)]`,
			inputJSON:    `[{"a":1,"b":3},{"x":1,"b":1},{"x":1,"b":3}]`,
			expectedJSON: `[{"a":1,"b":3},{"b":1,"x":1},{"b":3,"x":1}]`,
		},
		{
			jsonpath:     `$[?(@.b > 2 || @.x)]`,
			inputJSON:    `[{"a":1,"b":3},{"x":1,"b":1},{"x":1,"b":3}]`,
			expectedJSON: `[{"a":1,"b":3},{"b":1,"x":1},{"b":3,"x":1}]`,
		},
		{
			jsonpath:     `$[?(@.x || @.x)]`,
			inputJSON:    `[{"a":1,"b":3},{"x":1,"b":1}]`,
			expectedJSON: `[{"b":1,"x":1}]`,
		},
	}

	runTestCases(t, "TestFilterLogicalOr_PropertyExistence", testCases)
}

func TestFilterLogicalOr_AlwaysTrue(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.b > 2 || @.b < 2)]`,
			inputJSON:    `[{"a":1,"b":1},{"a":1,"b":3}]`,
			expectedJSON: `[{"a":1,"b":1},{"a":1,"b":3}]`,
		},
	}

	runTestCases(t, "TestFilterLogicalOr_AlwaysTrue", testCases)
}

func TestFilterLogicalOr_WithRootReference(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$.z[?($..x || @.b < 2)]`,
			inputJSON:    `{"x":1,"z":[{"a":1,"b":1},{"a":1,"b":3}]}`,
			expectedJSON: `[{"a":1,"b":1},{"a":1,"b":3}]`,
		},
		{
			jsonpath:     `$.z[?($..xx || @.b < 2)]`,
			inputJSON:    `{"x":1,"z":[{"a":1,"b":1},{"a":1,"b":3}]}`,
			expectedJSON: `[{"a":1,"b":1}]`,
		},
	}

	runTestCases(t, "TestFilterLogicalOr_WithRootReference", testCases)
}
