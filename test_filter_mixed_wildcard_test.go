package jsonpath

import (
	"fmt"
	"testing"
)

func TestFilter_WildcardValueGroupProhibitions(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@[*]==1)]`,
			inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[*]==1)]`},
		},
		{
			jsonpath:    `$[?(@[*].a==1)]`,
			inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[*].a==1)]`},
		},
		{
			jsonpath:    `$[?(@.a[*]==1)]`,
			inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[*]==1)]`},
		},

		{
			jsonpath:    `$[?(@.*==2)]`,
			inputJSON:   `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.*==2)]`},
		},
		{
			jsonpath:    `$[?(@.*[0]==2)]`,
			inputJSON:   `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.*[0]==2)]`},
		},
		{
			jsonpath:    `$[?(@.*.a==2)]`,
			inputJSON:   `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.*.a==2)]`},
		},
		{
			jsonpath:    `$[?(@.a.*==2)]`,
			inputJSON:   `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
			expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a.*==2)]`},
		},
		{
			jsonpath:    `$[?(@.*==[1,2])]`,
			inputJSON:   `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.*==[1,2])]`},
		},
		{
			jsonpath:    `$[?(@.*==['1','2'])]`,
			inputJSON:   `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.*==['1','2'])]`},
		},
	}

	for i, tc := range testCases {
		runTestCase(t, tc, fmt.Sprintf("TestFilter_WildcardValueGroupProhibitions[%d]", i))
	}
}
