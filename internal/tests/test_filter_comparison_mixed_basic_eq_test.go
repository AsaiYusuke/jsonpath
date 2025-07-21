package tests

import (
	"testing"
)

func TestFilterComparisonEQ_DecimalNumbers(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.a == 2.1)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":2.0,"b":4},{"a":2.1,"b":5},{"a":2.2,"b":6},{"a":"2.1"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
			expectedJSON: `[{"a":2.1,"b":5}]`,
		},
		{
			jsonpath:     `$[?(2.1 == @.a)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":2.0,"b":4},{"a":2.1,"b":5},{"a":2.2,"b":6},{"a":"2.1"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
			expectedJSON: `[{"a":2.1,"b":5}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonEQ_DecimalNumbers", tests)
}

func TestFilterComparisonEQ_StringEquality(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.a=='ab')]`,
			inputJSON:    `[{"a":"ab"}]`,
			expectedJSON: `[{"a":"ab"}]`,
		},
		{
			jsonpath:     `$[?(@.a==1)]`,
			inputJSON:    `[{"a":1},{"b":1}]`,
			expectedJSON: `[{"a":1}]`,
		},
		{
			jsonpath:    `$[?(@.a=='ab')]`,
			inputJSON:   `[{"a":"abc"}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.a=='ab')]`),
		},
	}

	runTestCases(t, "TestFilterComparisonEQ_StringEquality", tests)
}

func TestFilterComparisonEQ_StringLiterals(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.a=="value")]`,
			inputJSON:    `[{"a":"value"},{"a":0},{"a":1},{"a":-1},{"a":"val"},{"a":true},{"a":{}},{"a":[]},{"a":["b"]},{"a":{"a":"value"}},{"b":"value"}]`,
			expectedJSON: `[{"a":"value"}]`,
		},
		{
			jsonpath:     `$[?(@.a=='value')]`,
			inputJSON:    `[{"a":"value"},{"a":0},{"a":1},{"a":-1},{"a":"val"},{"a":{}},{"a":[]},{"a":["b"]},{"a":{"a":"value"}},{"b":"value"}]`,
			expectedJSON: `[{"a":"value"}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonEQ_StringLiterals", tests)
}

func TestFilterComparisonEQ_BooleanValues(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.a==true)]`,
			inputJSON:    `[{"a":true},{"a":false},{"a":1},{"a":0},{"a":"true"}]`,
			expectedJSON: `[{"a":true}]`,
		},
		{
			jsonpath:     `$[?(@.a==false)]`,
			inputJSON:    `[{"a":true},{"a":false},{"a":1},{"a":0},{"a":"false"}]`,
			expectedJSON: `[{"a":false}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonEQ_BooleanValues", tests)
}

func TestFilterComparisonEQ_NullValues(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.a==null)]`,
			inputJSON:    `[{"a":null},{"a":0},{"a":""},{"a":"null"},{"b":"value"}]`,
			expectedJSON: `[{"a":null}]`,
		},
	}

	runTestCases(t, "TestFilterComparisonEQ_NullValues", tests)
}
