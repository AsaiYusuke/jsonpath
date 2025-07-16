package jsonpath

import (
	"testing"
)

func TestFilterLogicalAND_BasicCombinations(t *testing.T) {
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
		{
			jsonpath:    `$[?(@.x && @.b > 2)]`,
			inputJSON:   `[{"a":"a"},{"b":2},{"b":3}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.x && @.b > 2)]`),
		},
		{
			jsonpath:    `$[?(@.b > 2 && @.x)]`,
			inputJSON:   `[{"a":"a"},{"b":2},{"b":3}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.b > 2 && @.x)]`),
		},
		{
			jsonpath:    `$[?(@.x && @.x)]`,
			inputJSON:   `[{"a":"a"},{"b":2},{"b":3}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.x && @.x)]`),
		},
		{
			jsonpath:    `$[?(@.b > 2 && @.b < 2)]`,
			inputJSON:   `[{"b":1},{"b":2},{"b":3}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.b > 2 && @.b < 2)]`),
		},
		{
			jsonpath:     `$.z[?($..x && @.b < 2)]`,
			inputJSON:    `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
			expectedJSON: `[{"b":1}]`,
		},
		{
			jsonpath:    `$.z[?($..xx && @.b < 2)]`,
			inputJSON:   `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
			expectedErr: createErrorMemberNotExist(`[?($..xx && @.b < 2)]`),
		},
	}

	for _, testCase := range testCases {
		runSingleTestCase(t, "TestFilterLogicalAND_BasicCombinations", testCase)
	}
}
