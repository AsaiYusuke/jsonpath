package jsonpath

import (
	"fmt"
	"testing"
)

func TestBracketNotation_WildcardSimple(t *testing.T) {
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
		runSingleTestCase(t, fmt.Sprintf("TestBracketNotation_WildcardSimple_%d", i), test)
	}
}

func TestBracketNotation_WildcardComplexCombinations(t *testing.T) {
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
		runSingleTestCase(t, fmt.Sprintf("TestBracketNotation_WildcardComplexCombinations_%d", i), test)
	}
}

func TestBracketNotation_WildcardArraySlice(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[0:2][*]`,
			inputJSON:    `[[1,2],[3,4],[5,6]]`,
			expectedJSON: `[1,2,3,4]`,
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("TestBracketNotation_WildcardArraySlice_%d", i), test)
	}
}

func TestBracketNotation_WildcardNestedAccess(t *testing.T) {
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
		runSingleTestCase(t, fmt.Sprintf("TestBracketNotation_WildcardNestedAccess_%d", i), test)
	}
}

func TestBracketNotation_WildcardNestedErrors(t *testing.T) {
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
		runSingleTestCase(t, fmt.Sprintf("TestBracketNotation_WildcardNestedErrors_%d", i), test)
	}
}

func TestBracketNotation_WildcardDeepNestedErrors(t *testing.T) {
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
		runSingleTestCase(t, fmt.Sprintf("TestBracketNotation_WildcardDeepNestedErrors_%d", i), test)
	}
}
