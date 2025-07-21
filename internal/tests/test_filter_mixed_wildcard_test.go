package tests

import (
	"testing"
)

func TestFilter_WildcardValueGroupProhibitions(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@[*]==1)]`,
			inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
			expectedErr: createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@[*]==1)]`),
		},
		{
			jsonpath:    `$[?(@[*].a==1)]`,
			inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
			expectedErr: createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@[*].a==1)]`),
		},
		{
			jsonpath:    `$[?(@.a[*]==1)]`,
			inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
			expectedErr: createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@.a[*]==1)]`),
		},

		{
			jsonpath:    `$[?(@.*==2)]`,
			inputJSON:   `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
			expectedErr: createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@.*==2)]`),
		},
		{
			jsonpath:    `$[?(@.*[0]==2)]`,
			inputJSON:   `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
			expectedErr: createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@.*[0]==2)]`),
		},
		{
			jsonpath:    `$[?(@.*.a==2)]`,
			inputJSON:   `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
			expectedErr: createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@.*.a==2)]`),
		},
		{
			jsonpath:    `$[?(@.a.*==2)]`,
			inputJSON:   `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
			expectedErr: createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@.a.*==2)]`),
		},
		{
			jsonpath:    `$[?(@.*==[1,2])]`,
			inputJSON:   `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.*==[1,2])]`),
		},
		{
			jsonpath:    `$[?(@.*==['1','2'])]`,
			inputJSON:   `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.*==['1','2'])]`),
		},
	}

	runTestCases(t, "TestFilter_WildcardValueGroupProhibitions", testCases)
}
