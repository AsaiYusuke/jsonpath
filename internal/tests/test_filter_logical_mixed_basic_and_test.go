package tests

import (
	"testing"
)

func TestFilterLogicalAnd_BasicOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a && @.b)]`,
			inputJSON:    `[{"a":1},{"b":2},{"a":3,"b":4}]`,
			expectedJSON: `[{"a":3,"b":4}]`,
		},
		{
			jsonpath:     `$[?(@.a>1 && @.a<3)]`,
			inputJSON:    `[{"a":1},{"a":1.1},{"a":2.9},{"a":3}]`,
			expectedJSON: `[{"a":1.1},{"a":2.9}]`,
		},
		{
			jsonpath:     `$[?(@.a<3 && @.a>1)]`,
			inputJSON:    `[{"a":1},{"a":1.1},{"a":2.9},{"a":3}]`,
			expectedJSON: `[{"a":1.1},{"a":2.9}]`,
		},
	}

	runTestCases(t, "TestFilterLogicalAnd_BasicOperations", testCases)
}

func TestFilterLogicalAnd_WithLiterals(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?((1==2) && @.a>1)]`,
			inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
			expectedErr: createErrorMemberNotExist(`[?((1==2) && @.a>1)]`),
		},
		{
			jsonpath:     `$[?((1==1) && @.a>1)]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
			expectedJSON: `[{"a":2},{"a":3}]`,
		},
		{
			jsonpath:    `$[?(@.a>1 && (1==2))]`,
			inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.a>1 && (1==2))]`),
		},
		{
			jsonpath:     `$[?(@.a>1 && (1==1))]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
			expectedJSON: `[{"a":2},{"a":3}]`,
		},
	}

	runTestCases(t, "TestFilterLogicalAnd_WithLiterals", testCases)
}

func TestFilterLogicalAnd_PropertyExistence(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.x && @.b > 2)]`,
			inputJSON:    `[{"a":1,"b":3},{"x":1,"b":1},{"x":1,"b":3}]`,
			expectedJSON: `[{"b":3,"x":1}]`,
		},
		{
			jsonpath:     `$[?(@.b > 2 && @.x)]`,
			inputJSON:    `[{"a":1,"b":3},{"x":1,"b":1},{"x":1,"b":3}]`,
			expectedJSON: `[{"b":3,"x":1}]`,
		},
		{
			jsonpath:     `$[?(@.x && @.x)]`,
			inputJSON:    `[{"a":1,"b":3},{"x":1,"b":1}]`,
			expectedJSON: `[{"b":1,"x":1}]`,
		},
	}

	runTestCases(t, "TestFilterLogicalAnd_PropertyExistence", testCases)
}

func TestFilterLogicalAnd_ImpossibleConditions(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@.b > 2 && @.b < 2)]`,
			inputJSON:   `[{"a":1,"b":3},{"a":1,"b":1}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.b > 2 && @.b < 2)]`),
		},
	}

	runTestCases(t, "TestFilterLogicalAnd_ImpossibleConditions", testCases)
}

func TestFilterLogicalAnd_WithRootReference(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$.z[?($..x && @.b < 2)]`,
			inputJSON:    `{"x":1,"z":[{"a":1,"b":1},{"a":1,"b":3}]}`,
			expectedJSON: `[{"a":1,"b":1}]`,
		},
		{
			jsonpath:    `$.z[?($..xx && @.b < 2)]`,
			inputJSON:   `{"x":1,"z":[{"a":1,"b":1},{"a":1,"b":3}]}`,
			expectedErr: createErrorMemberNotExist(`[?($..xx && @.b < 2)]`),
		},
	}

	runTestCases(t, "TestFilterLogicalAnd_WithRootReference", testCases)
}
