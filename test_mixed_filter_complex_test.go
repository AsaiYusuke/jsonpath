package jsonpath

import (
	"testing"
)

// Compound filter access tests - restoring deleted tests from filter combination patterns

func TestRetrieve_compound_filter_union_filter_success(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0,1][?(@.b)]`,
		inputJSON:    `[[{"a":1},{"b":2}],[{"a":3},{"b":4}]]`,
		expectedJSON: `[{"b":2},{"b":4}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_compound_filter_union_filter_success")
}

func TestRetrieve_compound_filter_union_filter_no_match_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1][?(@.b)]`,
		inputJSON:   `[[{"a":1},{"x":2}],[{"a":3},{"x":4}]]`,
		expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_filter_union_filter_no_match_error")
}

func TestRetrieve_compound_filter_union_filter_empty_root_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1][?(@.b)]`,
		inputJSON:   `[]`,
		expectedErr: createErrorMemberNotExist(`[0,1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_filter_union_filter_empty_root_error")
}

func TestRetrieve_compound_filter_union_filter_type_mismatch_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1][?(@.b)]`,
		inputJSON:   `"x"`,
		expectedErr: createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_filter_union_filter_type_mismatch_error")
}

func TestRetrieve_compound_filter_recursive_property_success(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.b)]..a`,
		inputJSON:    `[{"a":1},{"b":{"a":2}},{"c":3},{"b":[{"a":4}]}]`,
		expectedJSON: `[2,4]`,
	}
	runTestCase(t, testCase, "TestRetrieve_compound_filter_recursive_property_success")
}

func TestRetrieve_compound_filter_recursive_property_no_match_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.b)]..a`,
		inputJSON:   `[{"a":1},{"b":{"x":2}},{"c":3},{"b":[{"x":4}]}]`,
		expectedErr: createErrorMemberNotExist(`a`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_filter_recursive_property_no_match_error")
}

func TestRetrieve_compound_filter_recursive_property_wrong_type_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.b)]..a`,
		inputJSON:   `[{"a":1},{"b":"a"},{"c":3},{"b":"a"}]`,
		expectedErr: createErrorMemberNotExist(`a`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_filter_recursive_property_wrong_type_error")
}

func TestRetrieve_compound_filter_recursive_property_filter_no_match_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.b)]..a`,
		inputJSON:   `[{"a":1},{"x":{"a":2}},{"c":3},{"x":[{"a":4}]}]`,
		expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_filter_recursive_property_filter_no_match_error")
}

func TestRetrieve_compound_filter_union_properties_success(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.b)]['a','c']`,
		inputJSON:    `[{"a":1},{"b":2},{"a":3,"b":4},{"c":5},{"a":6,"c":7},{"b":8,"c":9},{"a":10,"b":11,"c":12}]`,
		expectedJSON: `[3,9,10,12]`,
	}
	runTestCase(t, testCase, "TestRetrieve_compound_filter_union_properties_success")
}

func TestRetrieve_compound_filter_union_properties_no_match_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.b)]['a','c']`,
		inputJSON:   `[{"a":1},{"b":2},{"x":3,"b":4},{"c":5},{"a":6,"c":7},{"b":8,"z":9},{"x":10,"b":11,"z":12}]`,
		expectedErr: createErrorMemberNotExist(`['a','c']`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_filter_union_properties_no_match_error")
}

func TestRetrieve_compound_filter_union_properties_filter_no_match_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.b)]['a','c']`,
		inputJSON:   `[{"a":1},{"x":2},{"a":3,"x":4},{"c":5},{"a":6,"c":7},{"x":8,"c":9},{"a":10,"x":11,"c":12}]`,
		expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_filter_union_properties_filter_no_match_error")
}

func TestRetrieve_compound_filter_union_properties_type_mismatch_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.b)]['a','c']`,
		inputJSON:   `"x"`,
		expectedErr: createErrorTypeUnmatched(`[?(@.b)]`, `object/array`, `string`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_filter_union_properties_type_mismatch_error")
}

func TestRetrieve_compound_filter_wildcard_properties_success(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.b)].*`,
		inputJSON:    `[{"a":1},{"b":2},{"a":3,"b":4},{"c":5},{"a":6,"c":7},{"b":8,"c":9},{"a":10,"b":11,"c":12}]`,
		expectedJSON: `[2,3,4,8,9,10,11,12]`,
	}
	runTestCase(t, testCase, "TestRetrieve_compound_filter_wildcard_properties_success")
}

func TestRetrieve_compound_filter_wildcard_properties_filter_no_match_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.b)].*`,
		inputJSON:   `[{"a":1},{"x":2},{"a":3,"x":4},{"c":5},{"a":6,"c":7},{"x":8,"c":9},{"a":10,"x":11,"c":12}]`,
		expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_filter_wildcard_properties_filter_no_match_error")
}

// Generic filter compound access tests

func TestRetrieve_compound_generic_filter_slice_success(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@)][0:2]`,
		inputJSON:    `[[1,2],[3,4],[5,6,7]]`,
		expectedJSON: `[1,2,3,4,5,6]`,
	}
	runTestCase(t, testCase, "TestRetrieve_compound_generic_filter_slice_success")
}

func TestRetrieve_compound_generic_filter_slice_empty_arrays_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@)][0:2]`,
		inputJSON:   `[[],[],[]]`,
		expectedErr: createErrorMemberNotExist(`[0:2]`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_generic_filter_slice_empty_arrays_error")
}

func TestRetrieve_compound_generic_filter_wildcard_success(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@)][*]`,
		inputJSON:    `[[1,2],[3,4],[5,6,7]]`,
		expectedJSON: `[1,2,3,4,5,6,7]`,
	}
	runTestCase(t, testCase, "TestRetrieve_compound_generic_filter_wildcard_success")
}

func TestRetrieve_compound_generic_filter_wildcard_empty_arrays_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@)][*]`,
		inputJSON:   `[[],[],[]]`,
		expectedErr: createErrorMemberNotExist(`[*]`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_generic_filter_wildcard_empty_arrays_error")
}

func TestRetrieve_compound_generic_filter_union_success(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@)][0,1]`,
		inputJSON:    `[[1,2],[3,4],[5,6,7]]`,
		expectedJSON: `[1,2,3,4,5,6]`,
	}
	runTestCase(t, testCase, "TestRetrieve_compound_generic_filter_union_success")
}

func TestRetrieve_compound_generic_filter_union_empty_arrays_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@)][0,1]`,
		inputJSON:   `[[],[],[]]`,
		expectedErr: createErrorMemberNotExist(`[0,1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_generic_filter_union_empty_arrays_error")
}

func TestRetrieve_compound_generic_filter_union_empty_root_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@)][0,1]`,
		inputJSON:   `[]`,
		expectedErr: createErrorMemberNotExist(`[?(@)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_generic_filter_union_empty_root_error")
}

func TestRetrieve_compound_generic_filter_union_type_mismatch_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@)][0,1]`,
		inputJSON:   `"x"`,
		expectedErr: createErrorTypeUnmatched(`[?(@)]`, `object/array`, `string`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_generic_filter_union_type_mismatch_error")
}

func TestRetrieve_compound_nested_filter_success(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a)][?(@.b)]`,
		inputJSON:    `[{"a":{"b":2}},{"b":{"a":1}},{"a":{"a":3}}]`,
		expectedJSON: `[{"b":2}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_compound_nested_filter_success")
}

func TestRetrieve_compound_nested_filter_inner_no_match_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a)][?(@.b)]`,
		inputJSON:   `[{"a":{"x":2}},{"b":{"a":1}},{"a":{"a":3}}]`,
		expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_nested_filter_inner_no_match_error")
}

func TestRetrieve_compound_nested_filter_outer_no_match_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a)][?(@.b)]`,
		inputJSON:   `[{"x":{"b":2}},{"b":{"a":1}},{"x":{"a":3}}]`,
		expectedErr: createErrorMemberNotExist(`[?(@.a)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_nested_filter_outer_no_match_error")
}
