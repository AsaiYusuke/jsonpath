package jsonpath

import (
	"fmt"
	"testing"
)

func TestBracketNotationWildcard_SimpleWildcard(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$[*]`,
			inputJSON:   `[]`,
			expectedErr: createErrorMemberNotExist(`[*]`),
		},
		{
			jsonpath:    `$[*]`,
			inputJSON:   `{}`,
			expectedErr: createErrorMemberNotExist(`[*]`),
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("SimpleWildcard_%d", i), test)
	}
}

func TestBracketNotationWildcard_ComplexWildcardCombinations(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[*]`,
			inputJSON:    `["a",123,true,{"b":"c"},[0,1],null]`,
			expectedJSON: `["a",123,true,{"b":"c"},[0,1],null]`,
		},
		{
			jsonpath:     `$[*]`,
			inputJSON:    `{"a":[1],"b":[2,3]}`,
			expectedJSON: `[[1],[2,3]]`,
		},
		{
			jsonpath:     `$['a',*]`,
			inputJSON:    `{"a":[1],"b":[2,3]}`,
			expectedJSON: `[[1],[1],[2,3]]`,
		},
		{
			jsonpath:     `$[*,'a']`,
			inputJSON:    `{"a":[1],"b":[2,3]}`,
			expectedJSON: `[[1],[2,3],[1]]`,
		},
		{
			jsonpath:     `$[*,*,*]`,
			inputJSON:    `{"a":[1],"b":[2,3]}`,
			expectedJSON: `[[1],[2,3],[1],[2,3],[1],[2,3]]`,
		},
		{
			jsonpath:     `$['a',*,*]`,
			inputJSON:    `{"a":[1],"b":[2,3]}`,
			expectedJSON: `[[1],[1],[2,3],[1],[2,3]]`,
		},
		{
			jsonpath:     `$[*,'a',*]`,
			inputJSON:    `{"a":[1],"b":[2,3]}`,
			expectedJSON: `[[1],[2,3],[1],[1],[2,3]]`,
		},
		{
			jsonpath:     `$['a','a',*]`,
			inputJSON:    `{"a":[1],"b":[2,3]}`,
			expectedJSON: `[[1],[1],[1],[2,3]]`,
		},
		{
			jsonpath:     `$[*,*,'a']`,
			inputJSON:    `{"a":[1],"b":[2,3]}`,
			expectedJSON: `[[1],[2,3],[1],[2,3],[1]]`,
		},
		{
			jsonpath:     `$['a',*,'a']`,
			inputJSON:    `{"a":[1],"b":[2,3]}`,
			expectedJSON: `[[1],[1],[2,3],[1]]`,
		},
		{
			jsonpath:     `$[*,'a','a']`,
			inputJSON:    `{"a":[1],"b":[2,3]}`,
			expectedJSON: `[[1],[2,3],[1],[1]]`,
		},
		{
			jsonpath:     `$['a','a','a']`,
			inputJSON:    `{"a":[1],"b":[2,3]}`,
			expectedJSON: `[[1],[1],[1]]`,
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("ComplexWildcardCombinations_%d", i), test)
	}
}

func TestBracketNotationWildcard_ArraySliceWithWildcard(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[0:2][*]`,
			inputJSON:    `[[1,2],[3,4],[5,6]]`,
			expectedJSON: `[1,2,3,4]`,
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("ArraySliceWithWildcard_%d", i), test)
	}
}

func TestBracketNotationWildcard_NestedWildcardAccess(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[*].a`,
			inputJSON:    `[{"a":1},{"b":2}]`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$[*].a`,
			inputJSON:    `[{"a":1},{"a":1}]`,
			expectedJSON: `[1,1]`,
		},
		{
			jsonpath:     `$[*].a`,
			inputJSON:    `[{"a":[1,[2]]},{"a":2}]`,
			expectedJSON: `[[1,[2]],2]`,
		},
		{
			jsonpath:     `$[*].a[*]`,
			inputJSON:    `[{"a":[1,[2]]},{"a":2}]`,
			expectedJSON: `[1,[2]]`,
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("NestedWildcardAccess_%d", i), test)
	}
}

func TestBracketNotationWildcard_NestedErrorCases(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$[*].a.b`,
			inputJSON:   `{"a":{"b":1}}`,
			expectedErr: createErrorMemberNotExist(`.a`),
		},
		{
			jsonpath:    `$[*].a.b`,
			inputJSON:   `[{"b":1}]`,
			expectedErr: createErrorMemberNotExist(`.a`),
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("NestedErrorCases_%d", i), test)
	}
}

func TestBracketNotationWildcard_DeepNestedErrorCases(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$[*].a.b.c`,
			inputJSON:   `{"a":{"b":1},"b":{"a":2}}`,
			expectedErr: createErrorTypeUnmatched(`.b`, `object`, `float64`),
		},
		{
			jsonpath:    `$[*].a.b.c`,
			inputJSON:   `[{"b":1},{"a":2}]`,
			expectedErr: createErrorTypeUnmatched(`.b`, `object`, `float64`),
		},
		{
			jsonpath:    `$[*].a.b.c`,
			inputJSON:   `{"a":{"a":1},"b":{"a":{"c":2}}}`,
			expectedErr: createErrorMemberNotExist(`.b`),
		},
		{
			jsonpath:    `$[*].a.b.c`,
			inputJSON:   `[{"a":1},{"a":{"c":2}}]`,
			expectedErr: createErrorMemberNotExist(`.b`),
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("DeepNestedErrorCases_%d", i), test)
	}
}
