package jsonpath

import (
	"testing"
)

func TestFilterLogicalOR_BasicCombinations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a || @.b)]`,
			inputJSON:    `[{"a":1},{"b":2},{"c":3}]`,
			expectedJSON: `[{"a":1},{"b":2}]`,
		},
		{
			jsonpath:     `$[?(@.a>2 || @.a<2)]`,
			inputJSON:    `[{"a":1},{"a":1.9},{"a":2},{"a":2.1},{"a":3}]`,
			expectedJSON: `[{"a":1},{"a":1.9},{"a":2.1},{"a":3}]`,
		},
		{
			jsonpath:     `$[?(@.a<2 || @.a>2)]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
			expectedJSON: `[{"a":1},{"a":3}]`,
		},
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
		{
			jsonpath:     `$[?(@.x || @.b > 2)]`,
			inputJSON:    `[{"a":"a"},{"b":2},{"b":3}]`,
			expectedJSON: `[{"b":3}]`,
		},
		{
			jsonpath:     `$[?(@.b > 2 || @.x)]`,
			inputJSON:    `[{"a":"a"},{"b":2},{"b":3}]`,
			expectedJSON: `[{"b":3}]`,
		},
		{
			jsonpath:    `$[?(@.x || @.x)]`,
			inputJSON:   `[{"a":"a"},{"b":2},{"b":3}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.x || @.x)]`),
		},
		{
			jsonpath:     `$[?(@.b > 2 || @.b < 2)]`,
			inputJSON:    `[{"b":1},{"b":2},{"b":3}]`,
			expectedJSON: `[{"b":1},{"b":3}]`,
		},
		{
			jsonpath:     `$.z[?($..x || @.b < 2)]`,
			inputJSON:    `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
			expectedJSON: `[{"b":1},{"b":2},{"b":3}]`,
		},
		{
			jsonpath:     `$.z[?($..xx || @.b < 2)]`,
			inputJSON:    `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
			expectedJSON: `[{"b":1}]`,
		},
	}

	for _, testCase := range testCases {
		runSingleTestCase(t, "TestFilterLogicalOR_BasicCombinations", testCase)
	}
}
