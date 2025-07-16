package jsonpath

import (
	"testing"
)

func TestRetrieve_valueGroupCombination_Union_in_qualifier(t *testing.T) {
	testGroups := TestGroup{
		`Recursive-descent`: []TestCase{
			{
				jsonpath:     `$[0,1]..a`,
				inputJSON:    `[{"a":1},{"b":{"a":2}},{"a":3}]`,
				expectedJSON: `[1,2]`,
			},
			{
				jsonpath:    `$[0,1]..a`,
				inputJSON:   `[{"x":1},{"b":{"x":2}},{"a":3}]`,
				expectedErr: createErrorMemberNotExist(`a`),
			},
			{
				jsonpath:    `$[0,1]..a`,
				inputJSON:   `["a","b",{"a":3}]`,
				expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
			},
			{
				jsonpath:    `$[0,1]..a`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$[0,1]..a`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
			},
		},
		`Multiple-identifier`: []TestCase{
			{
				jsonpath:     `$[0,1]['a','b']`,
				inputJSON:    `[{"a":1},{"b":2}]`,
				expectedJSON: `[1,2]`,
			},
			{
				jsonpath:    `$[0,1]['a','b']`,
				inputJSON:   `[{"x":1},{"x":2}]`,
				expectedErr: createErrorMemberNotExist(`['a','b']`),
			},
			{
				jsonpath:    `$[0,1]['a','b']`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$[0,1]['a','b']`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
			},
		},
		`Wildcard-identifier`: []TestCase{
			{
				jsonpath:     `$[0,1].*`,
				inputJSON:    `[{"a":1,"c":2},{"d":3,"b":4}]`,
				expectedJSON: `[1,2,4,3]`,
			},
			{
				jsonpath:    `$[0,1].*`,
				inputJSON:   `[[],[]]`,
				expectedErr: createErrorMemberNotExist(`.*`),
			},
			{
				jsonpath:    `$[0,1].*`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$[0,1].*`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
			},
		},
		`Slice-qualifier`: []TestCase{
			{
				jsonpath:     `$[0,1][0:2]`,
				inputJSON:    `[[1,2,3],[4,5,6]]`,
				expectedJSON: `[1,2,4,5]`,
			},
			{
				jsonpath:    `$[0,1][0:2]`,
				inputJSON:   `[[],[]]`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
			{
				jsonpath:    `$[0,1][0:2]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$[0,1][0:2]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
			},
		},
		`Wildcard-qualifier`: []TestCase{
			{
				jsonpath:     `$[0,1][*]`,
				inputJSON:    `[{"a":1,"c":3},{"d":4,"b":2},{"e":5}]`,
				expectedJSON: `[1,3,2,4]`,
			},
			{
				jsonpath:    `$[0,1][*]`,
				inputJSON:   `[{},{}]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[0,1][*]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$[0,1][*]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
			},
		},
		`Union-in-qualifier`: []TestCase{
			{
				jsonpath:     `$[0,1][0,1]`,
				inputJSON:    `[[1,3,2],[4,6,5],[7]]`,
				expectedJSON: `[1,3,4,6]`,
			},
			{
				jsonpath:    `$[0,1][0,1]`,
				inputJSON:   `[[],[],[7]]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$[0,1][0,1]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$[0,1][0,1]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
			},
		},
		`Filter-qualifier`: []TestCase{
			{
				jsonpath:     `$[0,1][?(@.b)]`,
				inputJSON:    `[[{"a":1},{"b":2}],[{"a":3},{"b":4}]]`,
				expectedJSON: `[{"b":2},{"b":4}]`,
			},
			{
				jsonpath:    `$[0,1][?(@.b)]`,
				inputJSON:   `[[{"a":1},{"x":2}],[{"a":3},{"x":4}]]`,
				expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
			},
			{
				jsonpath:    `$[0,1][?(@.b)]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$[0,1][?(@.b)]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
			},
		},
	}

	runTestGroups(t, testGroups)
}
