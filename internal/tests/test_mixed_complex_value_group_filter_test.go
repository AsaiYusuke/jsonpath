package tests

import (
	"testing"
)

func TestRetrieve_valueGroupCombination_Filter_qualifier(t *testing.T) {
	testGroups := TestGroup{
		`Recursive-descent`: []TestCase{
			{
				jsonpath:     `$[?(@.b)]..a`,
				inputJSON:    `[{"a":1},{"b":{"a":2}},{"c":3},{"b":[{"a":4}]}]`,
				expectedJSON: `[2,4]`,
			},
			{
				jsonpath:    `$[?(@.b)]..a`,
				inputJSON:   `[{"a":1},{"b":{"x":2}},{"c":3},{"b":[{"x":4}]}]`,
				expectedErr: createErrorMemberNotExist(`a`),
			},
			{
				jsonpath:    `$[?(@.b)]..a`,
				inputJSON:   `[{"a":1},{"b":"a"},{"c":3},{"b":"a"}]`,
				expectedErr: createErrorMemberNotExist(`a`),
			},
			{
				jsonpath:    `$[?(@.b)]..a`,
				inputJSON:   `[{"a":1},{"x":{"a":2}},{"c":3},{"x":[{"a":4}]}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
			},
			{
				jsonpath:    `$[?(@.b)]..a`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[?(@.b)]`, `object/array`, `string`),
			},
		},
		`Multiple-identifier`: []TestCase{
			{
				jsonpath:     `$[?(@.b)]['a','c']`,
				inputJSON:    `[{"a":1},{"b":2},{"a":3,"b":4},{"c":5},{"a":6,"c":7},{"b":8,"c":9},{"a":10,"b":11,"c":12}]`,
				expectedJSON: `[3,9,10,12]`,
			},
			{
				jsonpath:    `$[?(@.b)]['a','c']`,
				inputJSON:   `[{"a":1},{"b":2},{"x":3,"b":4},{"c":5},{"a":6,"c":7},{"b":8,"z":9},{"x":10,"b":11,"z":12}]`,
				expectedErr: createErrorMemberNotExist(`['a','c']`),
			},
			{
				jsonpath:    `$[?(@.b)]['a','c']`,
				inputJSON:   `[{"a":1},{"x":2},{"a":3,"x":4},{"c":5},{"a":6,"c":7},{"x":8,"c":9},{"a":10,"x":11,"c":12}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
			},
			{
				jsonpath:    `$[?(@.b)]['a','c']`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[?(@.b)]`, `object/array`, `string`),
			},
		},
		`Wildcard-identifier`: []TestCase{
			{
				jsonpath:     `$[?(@.b)].*`,
				inputJSON:    `[{"a":1},{"b":2},{"a":3,"b":4},{"c":5},{"a":6,"c":7},{"b":8,"c":9},{"a":10,"b":11,"c":12}]`,
				expectedJSON: `[2,3,4,8,9,10,11,12]`,
			},
			{
				jsonpath:    `$[?(@.b)].*`,
				inputJSON:   `[{"a":1},{"x":2},{"a":3,"x":4},{"c":5},{"a":6,"c":7},{"x":8,"c":9},{"a":10,"x":11,"c":12}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
			},
			{
				jsonpath:    `$[?(@.b)].*`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[?(@.b)]`, `object/array`, `string`),
			},
		},
		`Slice-qualifier`: []TestCase{
			{
				jsonpath:     `$[?(@)][0:2]`,
				inputJSON:    `[[1,2],[3,4],[5,6,7]]`,
				expectedJSON: `[1,2,3,4,5,6]`,
			},
			{
				jsonpath:    `$[?(@)][0:2]`,
				inputJSON:   `[[],[],[]]`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
			{
				jsonpath:    `$[?(@)][0:2]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[?(@)]`),
			},
			{
				jsonpath:    `$[?(@)][0:2]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[?(@)]`, `object/array`, `string`),
			},
		},
		`Wildcard-qualifier`: []TestCase{
			{
				jsonpath:     `$[?(@)][*]`,
				inputJSON:    `[[1,2],[3,4],[5,6,7]]`,
				expectedJSON: `[1,2,3,4,5,6,7]`,
			},
			{
				jsonpath:    `$[?(@)][*]`,
				inputJSON:   `[[],[],[]]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[?(@)][*]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[?(@)]`),
			},
			{
				jsonpath:    `$[?(@)][*]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[?(@)]`, `object/array`, `string`),
			},
		},
		`Union-in-qualifier`: []TestCase{
			{
				jsonpath:     `$[?(@)][0,1]`,
				inputJSON:    `[[1,2],[3,4],[5,6,7]]`,
				expectedJSON: `[1,2,3,4,5,6]`,
			},
			{
				jsonpath:    `$[?(@)][0,1]`,
				inputJSON:   `[[],[],[]]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$[?(@)][0,1]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[?(@)]`),
			},
			{
				jsonpath:    `$[?(@)][0,1]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[?(@)]`, `object/array`, `string`),
			},
		},
		`Filter-qualifier`: []TestCase{
			{
				jsonpath:     `$[?(@.a)][?(@.b)]`,
				inputJSON:    `[{"a":{"b":2}},{"b":{"a":1}},{"a":{"a":3}}]`,
				expectedJSON: `[{"b":2}]`,
			},
			{
				jsonpath:    `$[?(@.a)][?(@.b)]`,
				inputJSON:   `[{"a":{"x":2}},{"b":{"a":1}},{"a":{"a":3}}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
			},
			{
				jsonpath:    `$[?(@.a)][?(@.b)]`,
				inputJSON:   `[{"x":{"b":2}},{"b":{"a":1}},{"x":{"a":3}}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.a)]`),
			},
			{
				jsonpath:    `$[?(@.a)]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[?(@.a)]`, `object/array`, `string`),
			},
			{
				jsonpath:    `$[?(@.a)]`,
				inputJSON:   `123`,
				expectedErr: createErrorTypeUnmatched(`[?(@.a)]`, `object/array`, `float64`),
			},
			{
				jsonpath:    `$[?(@.a)]`,
				inputJSON:   `true`,
				expectedErr: createErrorTypeUnmatched(`[?(@.a)]`, `object/array`, `bool`),
			},
			{
				jsonpath:    `$[?(@.a)]`,
				inputJSON:   `null`,
				expectedErr: createErrorTypeUnmatched(`[?(@.a)]`, `object/array`, `null`),
			},
		},
	}

	runTestGroups(t, testGroups)
}
