package jsonpath

import (
	"fmt"
	"testing"
)

func TestFilterArrayIndex(t *testing.T) {
	testCases := []TestCase{
		// Array negative index filtering
		{
			jsonpath:     `$[?(@[-1]==2)]`,
			inputJSON:    `[[0,1],[0,2],[2],["2"],["a","b"],["b"]]`,
			expectedJSON: `[[0,2],[2]]`,
		},

		// Nested property access in filters
		{
			jsonpath:     `$[?(@.a.b == 1)]`,
			inputJSON:    `[{"a":1},{"a":{"b":1}},{"a":{"a":1}}]`,
			expectedJSON: `[{"a":{"b":1}}]`,
		},

		// Root array element reference in filters - these work with proper JSON structure
		{
			jsonpath:     `$[?(@.a == $[2].b)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":2,"b":1},{"b":1}]`,
			expectedJSON: `[{"a":1}]`,
		},
		{
			jsonpath:     `$[?($[2].b == @.a)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":2,"b":1},{"b":1}]`,
			expectedJSON: `[{"a":1}]`,
		},
		{
			jsonpath:    `$[?(@.b == $[0].a)]`,
			inputJSON:   `[{"a":1},{"a":2}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.b == $[0].a)]`),
		},
		{
			jsonpath:    `$[?($[0].a == @.b)]`,
			inputJSON:   `[{"a":1},{"a":2}]`,
			expectedErr: createErrorMemberNotExist(`[?($[0].a == @.b)]`),
		},

		// Filter with property access on result
		{
			jsonpath:     `$[?(@.a == 2)].b`,
			inputJSON:    `[{"a":0},{"a":1},{"a":2,"b":4}]`,
			expectedJSON: `[4]`,
		},

		// Literal comparisons
		{
			jsonpath:     `$[?(10==10)]`,
			inputJSON:    `[{"a":10},{"a":20},{"a":30},{"a+10":20}]`,
			expectedJSON: `[{"a":10},{"a":20},{"a":30},{"a+10":20}]`,
		},
		{
			jsonpath:    `$[?(@.a==@.a)]`,
			inputJSON:   `[{"a":10},{"a":20},{"a":30},{"a+10":20}]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `comparison between two current nodes is prohibited`, near: `@.a==@.a)]`},
		},
	}

	for i, tc := range testCases {
		runTestCase(t, tc, fmt.Sprintf("TestFilterArrayIndex_case_%d", i))
	}
}
