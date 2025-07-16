package jsonpath

import (
	"fmt"
	"testing"
)

func TestAdvancedFilterOperations(t *testing.T) {
	testCases := []TestCase{
		// Root node comparison with member access
		{
			jsonpath:    `$[?(@.a == $.b)]`,
			inputJSON:   `[{"a":1},{"a":2}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.a == $.b)]`),
		},
		{
			jsonpath:    `$[?($.b == @.a)]`,
			inputJSON:   `[{"a":1},{"a":2}]`,
			expectedErr: createErrorMemberNotExist(`[?($.b == @.a)]`),
		},

		// Array slice operations in filters
		{
			jsonpath:    `$[?(@[0:1]==1)]`,
			inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0:1]==1)]`},
		},
		{
			jsonpath:    `$[?(@[0:2]==1)]`,
			inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0:2]==1)]`},
		},
		{
			jsonpath:    `$[?(@[0:2].a==1)]`,
			inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0:2].a==1)]`},
		},
		{
			jsonpath:    `$[?(@.a[0:2]==1)]`,
			inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[0:2]==1)]`},
		},

		// Multi-index operations in filters
		{
			jsonpath:    `$[?(@[0,1]==1)]`,
			inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0,1]==1)]`},
		},
		{
			jsonpath:    `$[?(@[0,1].a==1)]`,
			inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0,1].a==1)]`},
		},
		{
			jsonpath:    `$[?(@.a[0,1]==1)]`,
			inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[0,1]==1)]`},
		},

		// Recursive descent operations in filters
		{
			jsonpath:    `$[?(@..a==123)]`,
			inputJSON:   `[{"a":"123"},{"a":123}]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@..a==123)]`},
		},
		{
			jsonpath:    `$[?(@..a.b==123)]`,
			inputJSON:   `[{"a":"123"},{"a":123}]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@..a.b==123)]`},
		},
		{
			jsonpath:    `$[?(@.a..b==123)]`,
			inputJSON:   `[{"a":"123"},{"a":123}]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a..b==123)]`},
		},
		{
			jsonpath:    `$[?(@..a..b==123)]`,
			inputJSON:   `[{"a":"123"},{"a":123}]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@..a..b==123)]`},
		},

		// Quoted key operations in filters
		{
			jsonpath:    `$[?(@['a','b']==123)]`,
			inputJSON:   `[{"a":"123"},{"a":123}]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@['a','b']==123)]`},
		},
		{
			jsonpath:    `$[?(@['a','b','c']==123)]`,
			inputJSON:   `[{"a":"123"},{"a":123}]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@['a','b','c']==123)]`},
		},
		{
			jsonpath:    `$[?(@['a','b']['a']==123)]`,
			inputJSON:   `[{"a":"123"},{"a":123}]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@['a','b']['a']==123)]`},
		},
	}

	for i, tc := range testCases {
		runTestCase(t, tc, fmt.Sprintf("TestAdvancedFilterOperations_case_%d", i))
	}
}
