package jsonpath

import (
	"testing"
)

func TestMixed_FilterUnionFilterSuccess(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0,1][?(@.b)]`,
		inputJSON:    `[[{"a":1},{"b":2}],[{"a":3},{"b":4}]]`,
		expectedJSON: `[{"b":2},{"b":4}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_compound_filter_union_filter_success")
}

func TestMixed_FilterUnionFilterNoMatch(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1][?(@.b)]`,
		inputJSON:   `[[{"a":1},{"x":2}],[{"a":3},{"x":4}]]`,
		expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_filter_union_filter_no_match_error")
}

func TestMixed_FilterUnionFilterEmptyRoot(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1][?(@.b)]`,
		inputJSON:   `[]`,
		expectedErr: createErrorMemberNotExist(`[0,1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_filter_union_filter_empty_root_error")
}

func TestMixed_FilterUnionFilterTypeMismatch(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1][?(@.b)]`,
		inputJSON:   `"x"`,
		expectedErr: createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_filter_union_filter_type_mismatch_error")
}

func TestMixed_FilterRecursivePropertySuccess(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.b)]..a`,
		inputJSON:    `[{"a":1},{"b":{"a":2}},{"c":3},{"b":[{"a":4}]}]`,
		expectedJSON: `[2,4]`,
	}
	runTestCase(t, testCase, "TestRetrieve_compound_filter_recursive_property_success")
}

func TestMixed_FilterRecursivePropertyNoMatch(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.b)]..a`,
		inputJSON:   `[{"a":1},{"b":{"x":2}},{"c":3},{"b":[{"x":4}]}]`,
		expectedErr: createErrorMemberNotExist(`a`),
	}
	runTestCase(t, testCase, "TestMixed_FilterRecursivePropertyNoMatch")
}

func TestMixed_FilterRecursivePropertyWrongType(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.b)]..a`,
		inputJSON:   `[{"a":1},{"b":"a"},{"c":3},{"b":"a"}]`,
		expectedErr: createErrorMemberNotExist(`a`),
	}
	runTestCase(t, testCase, "TestMixed_FilterRecursivePropertyWrongType")
}

func TestMixed_FilterRecursivePropertyFilterNoMatch(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.b)]..a`,
		inputJSON:   `[{"a":1},{"x":{"a":2}},{"c":3},{"x":[{"a":4}]}]`,
		expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
	}
	runTestCase(t, testCase, "TestMixed_FilterRecursivePropertyFilterNoMatch")
}

func TestMixed_FilterUnionPropertiesSuccess(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.b)]['a','c']`,
		inputJSON:    `[{"a":1},{"b":2},{"a":3,"b":4},{"c":5},{"a":6,"c":7},{"b":8,"c":9},{"a":10,"b":11,"c":12}]`,
		expectedJSON: `[3,9,10,12]`,
	}
	runTestCase(t, testCase, "TestMixed_FilterUnionPropertiesSuccess")
}

func TestMixed_FilterUnionPropertiesNoMatch(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.b)]['a','c']`,
		inputJSON:   `[{"a":1},{"b":2},{"x":3,"b":4},{"c":5},{"a":6,"c":7},{"b":8,"z":9},{"x":10,"b":11,"z":12}]`,
		expectedErr: createErrorMemberNotExist(`['a','c']`),
	}
	runTestCase(t, testCase, "TestMixed_FilterUnionPropertiesNoMatch")
}

func TestMixed_FilterUnionPropertiesFilterNoMatch(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.b)]['a','c']`,
		inputJSON:   `[{"a":1},{"x":2},{"a":3,"x":4},{"c":5},{"a":6,"c":7},{"x":8,"c":9},{"a":10,"x":11,"c":12}]`,
		expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
	}
	runTestCase(t, testCase, "TestMixed_FilterUnionPropertiesFilterNoMatch")
}

func TestMixed_FilterUnionPropertiesTypeMismatch(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.b)]['a','c']`,
		inputJSON:   `"x"`,
		expectedErr: createErrorTypeUnmatched(`[?(@.b)]`, `object/array`, `string`),
	}
	runTestCase(t, testCase, "TestMixed_FilterUnionPropertiesTypeMismatch")
}

func TestMixed_FilterWildcardPropertiesSuccess(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.b)].*`,
		inputJSON:    `[{"a":1},{"b":2},{"a":3,"b":4},{"c":5},{"a":6,"c":7},{"b":8,"c":9},{"a":10,"b":11,"c":12}]`,
		expectedJSON: `[2,3,4,8,9,10,11,12]`,
	}
	runTestCase(t, testCase, "TestMixed_FilterWildcardPropertiesSuccess")
}

func TestMixed_FilterWildcardPropertiesFilterNoMatch(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.b)].*`,
		inputJSON:   `[{"a":1},{"x":2},{"a":3,"x":4},{"c":5},{"a":6,"c":7},{"x":8,"c":9},{"a":10,"x":11,"c":12}]`,
		expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
	}
	runTestCase(t, testCase, "TestMixed_FilterWildcardPropertiesFilterNoMatch")
}

func TestMixed_GenericFilterSliceSuccess(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@)][0:2]`,
		inputJSON:    `[[1,2],[3,4],[5,6,7]]`,
		expectedJSON: `[1,2,3,4,5,6]`,
	}
	runTestCase(t, testCase, "TestMixed_GenericFilterSliceSuccess")
}

func TestMixed_GenericFilterSliceEmptyArrays(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@)][0:2]`,
		inputJSON:   `[[],[],[]]`,
		expectedErr: createErrorMemberNotExist(`[0:2]`),
	}
	runTestCase(t, testCase, "TestMixed_GenericFilterSliceEmptyArrays")
}

func TestMixed_GenericFilterWildcardSuccess(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@)][*]`,
		inputJSON:    `[[1,2],[3,4],[5,6,7]]`,
		expectedJSON: `[1,2,3,4,5,6,7]`,
	}
	runTestCase(t, testCase, "TestMixed_GenericFilterWildcardSuccess")
}

func TestMixed_GenericFilterWildcardEmptyArrays(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@)][*]`,
		inputJSON:   `[[],[],[]]`,
		expectedErr: createErrorMemberNotExist(`[*]`),
	}
	runTestCase(t, testCase, "TestMixed_GenericFilterWildcardEmptyArrays")
}

func TestMixed_GenericFilterUnionSuccess(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@)][0,1]`,
		inputJSON:    `[[1,2],[3,4],[5,6,7]]`,
		expectedJSON: `[1,2,3,4,5,6]`,
	}
	runTestCase(t, testCase, "TestMixed_GenericFilterUnionSuccess")
}

func TestMixed_GenericFilterUnionEmptyArrays(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@)][0,1]`,
		inputJSON:   `[[],[],[]]`,
		expectedErr: createErrorMemberNotExist(`[0,1]`),
	}
	runTestCase(t, testCase, "TestMixed_GenericFilterUnionEmptyArrays")
}

func TestMixed_GenericFilterUnionEmptyRoot(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@)][0,1]`,
		inputJSON:   `[]`,
		expectedErr: createErrorMemberNotExist(`[?(@)]`),
	}
	runTestCase(t, testCase, "TestMixed_GenericFilterUnionEmptyRoot")
}

func TestMixed_GenericFilterUnionTypeMismatch(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@)][0,1]`,
		inputJSON:   `"x"`,
		expectedErr: createErrorTypeUnmatched(`[?(@)]`, `object/array`, `string`),
	}
	runTestCase(t, testCase, "TestMixed_GenericFilterUnionTypeMismatch")
}

func TestMixed_NestedFilterSuccess(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a)][?(@.b)]`,
		inputJSON:    `[{"a":{"b":2}},{"b":{"a":1}},{"a":{"a":3}}]`,
		expectedJSON: `[{"b":2}]`,
	}
	runTestCase(t, testCase, "TestMixed_NestedFilterSuccess")
}

func TestMixed_NestedFilterInnerNoMatch(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a)][?(@.b)]`,
		inputJSON:   `[{"a":{"x":2}},{"b":{"a":1}},{"a":{"a":3}}]`,
		expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
	}
	runTestCase(t, testCase, "TestMixed_NestedFilterInnerNoMatch")
}

func TestMixed_NestedFilterOuterNoMatch(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a)][?(@.b)]`,
		inputJSON:   `[{"x":{"b":2}},{"b":{"a":1}},{"x":{"a":3}}]`,
		expectedErr: createErrorMemberNotExist(`[?(@.a)]`),
	}
	runTestCase(t, testCase, "TestMixed_NestedFilterOuterNoMatch")
}
