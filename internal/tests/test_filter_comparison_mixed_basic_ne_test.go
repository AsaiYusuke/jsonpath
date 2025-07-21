package tests

import (
	"testing"
)

func TestFilterComparisonNE_BasicOperations(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.a != 2)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":2,"b":4},{"a":1.999999},{"a":2.000000000001},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
			expectedJSON: `[{"a":0},{"a":1},{"a":1.999999},{"a":2.000000000001},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
		},
		{
			jsonpath:     `$[?(2 != @.a)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":2,"b":4},{"a":1.999999},{"a":2.000000000001},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
			expectedJSON: `[{"a":0},{"a":1},{"a":1.999999},{"a":2.000000000001},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonNE_BasicOperations", tests)
}

func TestFilterComparisonNE_StringOperations(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.a!='ab')]`,
			inputJSON:    `[{"a":"abc"}]`,
			expectedJSON: `[{"a":"abc"}]`,
		},
		{
			jsonpath:    `$[?(@.a!='ab')]`,
			inputJSON:   `[{"a":"ab"}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.a!='ab')]`),
		},
		{
			jsonpath:     `$[?(@.a!=1)]`,
			inputJSON:    `[{"a":1},{"b":1}]`,
			expectedJSON: `[{"b":1}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonNE_StringOperations", tests)
}

func TestFilterComparisonNE_TypeMismatch(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.a!="string")]`,
			inputJSON:    `[{"a":"string"},{"a":123},{"a":true},{"a":null},{"a":[]},{"a":{}}]`,
			expectedJSON: `[{"a":123},{"a":true},{"a":null},{"a":[]},{"a":{}}]`,
		},
		{
			jsonpath:     `$[?(@.a!=123)]`,
			inputJSON:    `[{"a":"123"},{"a":123},{"a":true},{"a":null}]`,
			expectedJSON: `[{"a":"123"},{"a":true},{"a":null}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonNE_TypeMismatch", tests)
}
