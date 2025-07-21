package tests

import (
	"testing"
)

func TestRecursiveBasic_ConditionalRecursive(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$.z[?($..x)]`,
			inputJSON:    `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
			expectedJSON: `[{"b":1},{"b":2},{"b":3}]`,
		},
		{
			jsonpath:     `$..[?(@.a==2)]`,
			inputJSON:    `{"a":2,"x":[{"a":2},{"b":{"a":2}},{"a":{"a":2}},[{"a":2}]]}`,
			expectedJSON: `[{"a":2},{"a":2},{"a":2},{"a":2}]`,
		},
	}

	runTestCases(t, "TestRecursiveBasic_ConditionalRecursive", tests)
}

func TestRecursiveBasic_FilterWithRecursive(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$..*[?(@.a>2)]`,
			inputJSON:   `[{"b":"1","a":1},{"c":"2","a":2},{"d":"3","a":3}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.a>2)]`),
		},
		{
			jsonpath:     `$..*[?(@.a>2)]`,
			inputJSON:    `{"z":[{"d":"1","a":1},{"c":"2","a":2},{"b":"3","a":3}],"y":{"b":"4","a":4}}`,
			expectedJSON: `[{"a":3,"b":"3"}]`,
		},
		{
			jsonpath:     `$..*[?(@.a>2)]`,
			inputJSON:    `{"x":{"z":[{"x":"1","a":1},{"z":"2","a":2},{"y":"3","a":3}],"y":{"b":"4","a":4}}}`,
			expectedJSON: `[{"a":4,"b":"4"},{"a":3,"y":"3"}]`,
		},
		{
			jsonpath:     `$..*[?(@.a>2)]`,
			inputJSON:    `[{"x":{"z":[{"b":"1","a":1},{"b":"2","a":2},{"b":"3","a":3},{"b":"6","a":6}],"y":{"b":"4","a":4}}},{"b":"5","a":5}]`,
			expectedJSON: `[{"a":4,"b":"4"},{"a":3,"b":"3"},{"a":6,"b":"6"}]`,
		},
	}

	runTestCases(t, "TestRecursiveBasic_FilterWithRecursive", tests)
}
