package tests

import (
	"testing"
)

func TestFilterLogicalNot_BasicOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(!@.a)]`,
			inputJSON:    `[{"a":1,"b":2},{"b":3},{"a":4}]`,
			expectedJSON: `[{"b":3}]`,
		},
		{
			jsonpath:     `$[?(!@.c)]`,
			inputJSON:    `[{"a":1,"b":2},{"b":3},{"c":4}]`,
			expectedJSON: `[{"a":1,"b":2},{"b":3}]`,
		},
	}

	runTestCases(t, "TestFilterLogicalNot_BasicOperations", testCases)
}

func TestFilterLogicalNot_WithRootReference(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$.z[?(!$..x)]`,
			inputJSON:    `{"z":[{"a":1,"b":1},{"a":1,"b":3}]}`,
			expectedJSON: `[{"a":1,"b":1},{"a":1,"b":3}]`,
		},
		{
			jsonpath:     `$.z[?(!$..xx)]`,
			inputJSON:    `{"x":1,"z":[{"a":1,"b":1},{"a":1,"b":3}]}`,
			expectedJSON: `[{"a":1,"b":1},{"a":1,"b":3}]`,
		},
	}

	runTestCases(t, "TestFilterLogicalNot_WithRootReference", testCases)
}

func TestFilterLogicalNot_InvalidContexts(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(!$)]`,
			inputJSON:   `[1,2,3]`,
			expectedErr: createErrorMemberNotExist(`[?(!$)]`),
		},
		{
			jsonpath:    `$[?(!@)]`,
			inputJSON:   `[1,2,3]`,
			expectedErr: createErrorMemberNotExist(`[?(!@)]`),
		},
	}

	runTestCases(t, "TestFilterLogicalNot_InvalidContexts", testCases)
}

func TestFilterLogicalNot_ComplexExpressions(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(!@.a && @.b || @.c)]`,
			inputJSON:    `[{"a":1,"b":2},{"b":3},{"c":4},{"a":5,"c":6}]`,
			expectedJSON: `[{"b":3},{"c":4},{"a":5,"c":6}]`,
		},
		{
			jsonpath:     `$[?(@.a && !@.b || @.c)]`,
			inputJSON:    `[{"a":1,"b":2},{"a":3},{"c":4},{"a":5,"c":6}]`,
			expectedJSON: `[{"a":3},{"c":4},{"a":5,"c":6}]`,
		},
		{
			jsonpath:     `$[?(!@.a && !@.b || @.c)]`,
			inputJSON:    `[{"a":1,"b":2},{"a":3},{"c":4},{"d":5}]`,
			expectedJSON: `[{"c":4},{"d":5}]`,
		},
		{
			jsonpath:     `$[?(@.a && @.b || !@.c)]`,
			inputJSON:    `[{"a":1,"b":2},{"a":3},{"c":4},{"a":5,"b":6}]`,
			expectedJSON: `[{"a":1,"b":2},{"a":3},{"a":5,"b":6}]`,
		},
		{
			jsonpath:     `$[?(!@.a && @.b || !@.c)]`,
			inputJSON:    `[{"a":1,"b":2},{"b":3},{"c":4},{"d":5}]`,
			expectedJSON: `[{"a":1,"b":2},{"b":3},{"d":5}]`,
		},
		{
			jsonpath:     `$[?(@.a && !@.b || !@.c)]`,
			inputJSON:    `[{"a":1,"b":2},{"a":3},{"c":4},{"a":5,"d":6}]`,
			expectedJSON: `[{"a":1,"b":2},{"a":3},{"a":5,"d":6}]`,
		},
		{
			jsonpath:     `$[?(!@.a && !@.b || !@.c)]`,
			inputJSON:    `[{"a":1,"b":2},{"a":3},{"c":4},{"d":5}]`,
			expectedJSON: `[{"a":1,"b":2},{"a":3},{"c":4},{"d":5}]`,
		},
	}

	runTestCases(t, "TestFilterLogicalNot_ComplexExpressions", testCases)
}
