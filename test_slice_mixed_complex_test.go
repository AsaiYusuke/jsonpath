package jsonpath

import (
	"testing"
)

func TestRetrieve_sliceComplexNestedDeleted(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[0:2]..a`,
			inputJSON:    `[{"a":1},{"b":{"a":2}},{"a":3}]`,
			expectedJSON: `[1,2]`,
		},
		{
			jsonpath:     `$[0:2]..a`,
			inputJSON:    `[{"x":1},{"b":{"x":2}},{"a":3}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`a`),
		},
		{
			jsonpath:     `$[0:2]..a`,
			inputJSON:    `["a","b",{"a":3}]`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`..`, `object/array`, `string`),
		},
		{
			jsonpath:     `$[0:2]..a`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[0:2]`, `array`, `string`),
		},
		{
			jsonpath:     `$[0:2]['a','b']`,
			inputJSON:    `[{"a":1},{"b":2}]`,
			expectedJSON: `[1,2]`,
		},
		{
			jsonpath:     `$[0:2]['a','b']`,
			inputJSON:    `[{"x":1},{"x":2}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`['a','b']`),
		},
		{
			jsonpath:     `$[0:2]['a','b']`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[0:2]`),
		},
		{
			jsonpath:     `$[0:2]['a','b']`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[0:2]`, `array`, `string`),
		},
		{
			jsonpath:     `$[0:2].*`,
			inputJSON:    `[{"a":1,"c":2},{"d":3,"b":4}]`,
			expectedJSON: `[1,2,4,3]`,
		},
		{
			jsonpath:     `$[0:2].*`,
			inputJSON:    `[[],[]]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`.*`),
		},
		{
			jsonpath:     `$[0:2].*`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[0:2]`, `array`, `string`),
		},
		{
			jsonpath:     `$[0:2][0:2]`,
			inputJSON:    `[[1,2,3],[4,5,6]]`,
			expectedJSON: `[1,2,4,5]`,
		},
		{
			jsonpath:     `$[0:2][0:2]`,
			inputJSON:    `[[],[]]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[0:2]`),
		},
		{
			jsonpath:     `$[0:2][0:2]`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[0:2]`, `array`, `string`),
		},
		{
			jsonpath:     `$[0:2][*]`,
			inputJSON:    `[{"a":1,"c":3},{"d":4,"b":2},{"e":5}]`,
			expectedJSON: `[1,3,2,4]`,
		},
		{
			jsonpath:     `$[0:2][*]`,
			inputJSON:    `[{},{}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[*]`),
		},
		{
			jsonpath:     `$[0:2][*]`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[0:2]`, `array`, `string`),
		},
		{
			jsonpath:     `$[0:2][0,1]`,
			inputJSON:    `[[1,3,2],[4,6,5],[7]]`,
			expectedJSON: `[1,3,4,6]`,
		},
		{
			jsonpath:     `$[0:2][0,1]`,
			inputJSON:    `[[],[],[7]]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[0,1]`),
		},
		{
			jsonpath:     `$[0:2][0,1]`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[0:2]`),
		},
		{
			jsonpath:     `$[0:2][0,1]`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[0:2]`, `array`, `string`),
		},
		{
			jsonpath:     `$[0:2][?(@.b)]`,
			inputJSON:    `[[{"a":1},{"b":2}],[{"a":3},{"b":4}]]`,
			expectedJSON: `[{"b":2},{"b":4}]`,
		},
		{
			jsonpath:     `$[0:2][?(@.b)]`,
			inputJSON:    `[[{"a":1},{"x":2}],[{"a":3},{"x":4}]]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?(@.b)]`),
		},
	}

	for i, testCase := range testCases {
		runTestCase(t, testCase, "TestRetrieve_sliceComplexNestedDeleted_case_"+string(rune('A'+i)))
	}
}

func TestRetrieve_sliceWildcardDeleted(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[0:2][?(@.b)]`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[0:2]`, `array`, `string`),
		},
		{
			jsonpath:     `$[*]..a`,
			inputJSON:    `[{"a":1},{"b":{"a":2}},{"c":3}]`,
			expectedJSON: `[1,2]`,
		},
		{
			jsonpath:     `$[*]..a`,
			inputJSON:    `[{"x":1},{"b":{"x":2}},{"c":3}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`a`),
		},
		{
			jsonpath:     `$[*]..a`,
			inputJSON:    `["a","b","c"]`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`..`, `object/array`, `string`),
		},
		{
			jsonpath:     `$[*]..a`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[*]`, `object/array`, `string`),
		},
		{
			jsonpath:     `$[*]['a','b']`,
			inputJSON:    `[{"c":4},{"b":2,"a":1},{"a":3}]`,
			expectedJSON: `[1,2,3]`,
		},
		{
			jsonpath:     `$[*]['a','b']`,
			inputJSON:    `[{"c":4},{"x":2},{"x":1}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`['a','b']`),
		},
		{
			jsonpath:     `$[*]['a','b']`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[*]`),
		},
		{
			jsonpath:     `$[*]['a','b']`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[*]`, `object/array`, `string`),
		},
		{
			jsonpath:     `$[*].*`,
			inputJSON:    `[{"c":4},{"b":2,"a":1},{"a":3}]`,
			expectedJSON: `[4,1,2,3]`,
		},
		{
			jsonpath:     `$[*].*`,
			inputJSON:    `[{},{},{}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`.*`),
		},
		{
			jsonpath:     `$[*].*`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[*]`, `object/array`, `string`),
		},
		{
			jsonpath:     `$[*][0:2]`,
			inputJSON:    `[[1,2,3],[4,5],[6]]`,
			expectedJSON: `[1,2,4,5,6]`,
		},
		{
			jsonpath:     `$[*][0:2]`,
			inputJSON:    `[[],[],[]]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[0:2]`),
		},
		{
			jsonpath:     `$[*][0:2]`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[*]`, `object/array`, `string`),
		},
		{
			jsonpath:     `$[*][*]`,
			inputJSON:    `[[1,2,3],[4,5],[6]]`,
			expectedJSON: `[1,2,3,4,5,6]`,
		},
		{
			jsonpath:     `$[*][*]`,
			inputJSON:    `[[],[],[]]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[*]`),
		},
		{
			jsonpath:     `$[*][*]`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[*]`, `object/array`, `string`),
		},
		{
			jsonpath:     `$[*][0,1]`,
			inputJSON:    `[[1,3,2],[4,6,5],[7]]`,
			expectedJSON: `[1,3,4,6,7]`,
		},
		{
			jsonpath:     `$[*][0,1]`,
			inputJSON:    `[[],[],[]]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[0,1]`),
		},
		{
			jsonpath:     `$[*][0,1]`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[*]`),
		},
		{
			jsonpath:     `$[*][0,1]`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[*]`, `object/array`, `string`),
		},
		{
			jsonpath:     `$[*][?(@.b)]`,
			inputJSON:    `[[{"a":1},{"b":2}],[{"a":3},{"b":4}]]`,
			expectedJSON: `[{"b":2},{"b":4}]`,
		},
		{
			jsonpath:     `$[*][?(@.b)]`,
			inputJSON:    `[[{"a":1},{"x":2}],[{"a":3},{"x":4}]]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[?(@.b)]`),
		},
		{
			jsonpath:     `$[*][?(@.b)]`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[*]`, `object/array`, `string`),
		},
	}

	for i, testCase := range testCases {
		runTestCase(t, testCase, "TestRetrieve_sliceWildcardDeleted_case_"+string(rune('A'+i)))
	}
}

func TestRetrieve_sliceUnionDeleted(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[0,1]..a`,
			inputJSON:    `[{"a":1},{"b":{"a":2}},{"a":3}]`,
			expectedJSON: `[1,2]`,
		},
		{
			jsonpath:     `$[0,1]..a`,
			inputJSON:    `[{"x":1},{"b":{"x":2}},{"a":3}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`a`),
		},
		{
			jsonpath:     `$[0,1]..a`,
			inputJSON:    `["a","b",{"a":3}]`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`..`, `object/array`, `string`),
		},
	}

	for i, testCase := range testCases {
		runTestCase(t, testCase, "TestRetrieve_sliceUnionDeleted_case_"+string(rune('A'+i)))
	}
}

func TestRetrieve_sliceUnionAdditional(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[0,1]..a`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[0,1]`),
		},
		{
			jsonpath:     `$[0,1]..a`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
		},
		{
			jsonpath:     `$[0,1]['a','b']`,
			inputJSON:    `[{"a":1},{"b":2}]`,
			expectedJSON: `[1,2]`,
		},
		{
			jsonpath:     `$[0,1]['a','b']`,
			inputJSON:    `[{"x":1},{"x":2}]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`['a','b']`),
		},
		{
			jsonpath:     `$[0,1]['a','b']`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[0,1]`),
		},
		{
			jsonpath:     `$[0,1]['a','b']`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
		},
	}

	for i, testCase := range testCases {
		runTestCase(t, testCase, "TestRetrieve_sliceUnionAdditional_case_"+string(rune('A'+i)))
	}
}

func TestRetrieve_sliceUnionWildcard(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[0,1].*`,
			inputJSON:    `[{"a":1,"c":2},{"d":3,"b":4}]`,
			expectedJSON: `[1,2,4,3]`,
		},
		{
			jsonpath:     `$[0,1].*`,
			inputJSON:    `[[],[]]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`.*`),
		},
		{
			jsonpath:     `$[0,1].*`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  createErrorMemberNotExist(`[0,1]`),
		},
		{
			jsonpath:     `$[0,1].*`,
			inputJSON:    `"x"`,
			expectedJSON: ``,
			expectedErr:  createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
		},
	}

	for i, testCase := range testCases {
		runTestCase(t, testCase, "TestRetrieve_sliceUnionWildcard_case_"+string(rune('A'+i)))
	}
}
