package tests

import (
	"testing"
)

func TestMixed_ValueGroupCombinationWildcardQualifier(t *testing.T) {
	testGroups := TestGroup{
		`Recursive-descent`: []TestCase{
			{
				jsonpath:     `$[*]..a`,
				inputJSON:    `[{"a":1},{"b":{"a":2}},{"c":3}]`,
				expectedJSON: `[1,2]`,
			},
			{
				jsonpath:    `$[*]..a`,
				inputJSON:   `[{"x":1},{"b":{"x":2}},{"c":3}]`,
				expectedErr: createErrorMemberNotExist(`a`),
			},
			{
				jsonpath:    `$[*]..a`,
				inputJSON:   `["a","b","c"]`,
				expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
			},
			{
				jsonpath:    `$[*]..a`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[*]..a`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[*]`, `object/array`, `string`),
			},
		},
		`Multiple-identifier`: []TestCase{
			{
				jsonpath:     `$[*]['a','b']`,
				inputJSON:    `[{"c":4},{"b":2,"a":1},{"a":3}]`,
				expectedJSON: `[1,2,3]`,
			},
			{
				jsonpath:    `$[*]['a','b']`,
				inputJSON:   `[{"c":4},{"x":2},{"x":1}]`,
				expectedErr: createErrorMemberNotExist(`['a','b']`),
			},
			{
				jsonpath:    `$[*]['a','b']`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[*]['a','b']`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[*]`, `object/array`, `string`),
			},
		},
		`Wildcard-identifier`: []TestCase{
			{
				jsonpath:     `$[*].*`,
				inputJSON:    `[{"c":4},{"b":2,"a":1},{"a":3}]`,
				expectedJSON: `[4,1,2,3]`,
			},
			{
				jsonpath:    `$[*].*`,
				inputJSON:   `[{},{},{}]`,
				expectedErr: createErrorMemberNotExist(`.*`),
			},
			{
				jsonpath:    `$[*].*`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[*].*`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[*]`, `object/array`, `string`),
			},
		},
		`Slice-qualifier`: []TestCase{
			{
				jsonpath:     `$[*][0:2]`,
				inputJSON:    `[[1,2,3],[4,5],[6]]`,
				expectedJSON: `[1,2,4,5,6]`,
			},
			{
				jsonpath:    `$[*][0:2]`,
				inputJSON:   `[[],[],[]]`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
			{
				jsonpath:    `$[*][0:2]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[*][0:2]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[*]`, `object/array`, `string`),
			},
		},
		`Wildcard-qualifier`: []TestCase{
			{
				jsonpath:     `$[*][*]`,
				inputJSON:    `[[1,2,3],[4,5],[6]]`,
				expectedJSON: `[1,2,3,4,5,6]`,
			},
			{
				jsonpath:    `$[*][*]`,
				inputJSON:   `[[],[],[]]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[*][*]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[*][*]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[*]`, `object/array`, `string`),
			},
		},
		`Union-in-qualifier`: []TestCase{
			{
				jsonpath:     `$[*][0,1]`,
				inputJSON:    `[[1,3,2],[4,6,5],[7]]`,
				expectedJSON: `[1,3,4,6,7]`,
			},
			{
				jsonpath:    `$[*][0,1]`,
				inputJSON:   `[[],[],[]]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$[*][0,1]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[*][0,1]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[*]`, `object/array`, `string`),
			},
		},
		`Filter-qualifier`: []TestCase{
			{
				jsonpath:     `$[*][?(@.b)]`,
				inputJSON:    `[[{"a":1},{"b":2}],[{"a":3},{"b":4}]]`,
				expectedJSON: `[{"b":2},{"b":4}]`,
			},
			{
				jsonpath:    `$[*][?(@.b)]`,
				inputJSON:   `[[{"a":1},{"x":2}],[{"a":3},{"x":4}]]`,
				expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
			},
			{
				jsonpath:    `$[*][?(@.b)]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[*][?(@.b)]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[*]`, `object/array`, `string`),
			},
		},
	}

	runTestGroups(t, testGroups)
}
