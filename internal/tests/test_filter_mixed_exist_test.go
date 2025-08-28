package tests

import (
	"testing"
)

func TestRetrieve_filterExist_current_root_exists(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@)]`,
		inputJSON:    `["a","b"]`,
		expectedJSON: `["a","b"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_current_root_exists")
}

func TestRetrieve_filterExist_current_root_not_exists(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(!@)]`,
		inputJSON:   `["a","b"]`,
		expectedErr: createErrorMemberNotExist(`[?(!@)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_current_root_not_exists")
}

func TestRetrieve_filterExist_child_property_exists(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a)]`,
		inputJSON:    `[{"b":2},{"a":1},{"a":"value"},{"a":""},{"a":true},{"a":false},{"a":null},{"a":{}},{"a":[]}]`,
		expectedJSON: `[{"a":1},{"a":"value"},{"a":""},{"a":true},{"a":false},{"a":null},{"a":{}},{"a":[]}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_child_property_exists")
}

func TestRetrieve_filterExist_child_property_not_exists(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@.a)]`,
		inputJSON:    `[{"b":2},{"a":1},{"a":"value"},{"a":""},{"a":true},{"a":false},{"a":null},{"a":{}},{"a":[]}]`,
		expectedJSON: `[{"b":2}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_child_property_not_exists")
}

func TestRetrieve_filterExist_child_missing_property_no_match(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.c)]`,
		inputJSON:   `[{"a":1},{"b":2}]`,
		expectedErr: createErrorMemberNotExist(`[?(@.c)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_child_missing_property_no_match")
}

func TestRetrieve_filterExist_child_missing_property_negated(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@.c)]`,
		inputJSON:    `[{"a":1},{"b":2}]`,
		expectedJSON: `[{"a":1},{"b":2}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_child_missing_property_negated")
}

func TestRetrieve_filterExist_index_exists(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@[1])]`,
		inputJSON:    `[[{"a":1}],[{"b":2},{"c":3}],[],{"d":4}]`,
		expectedJSON: `[[{"b":2},{"c":3}]]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_index_exists")
}

func TestRetrieve_filterExist_index_not_exists(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@[1])]`,
		inputJSON:    `[[{"a":1}],[{"b":2},{"c":3}],[],{"d":4}]`,
		expectedJSON: `[[{"a":1}],[],{"d":4}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_index_not_exists")
}

func TestRetrieve_filterExist_wildcard_exists(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.*)]`,
		inputJSON:    `[{"a":1},{"b":2}]`,
		expectedJSON: `[{"a":1},{"b":2}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_wildcard_exists")
}

func TestRetrieve_filterExist_wildcard_not_exists(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.*)]`,
		inputJSON:   `[1,2]`,
		expectedErr: createErrorMemberNotExist(`[?(@.*)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_wildcard_not_exists")
}

func TestRetrieve_filterExist_wildcard_qualifier_array(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@[*])]`,
		inputJSON:    `[[{"a":1}],[]]`,
		expectedJSON: `[[{"a":1}]]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_wildcard_qualifier_array")
}

func TestRetrieve_filterExist_wildcard_qualifier_object(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@[*])]`,
		inputJSON:    `[{"a":1},{}]`,
		expectedJSON: `[{"a":1}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_wildcard_qualifier_object")
}

func TestRetrieve_filterExist_wildcard_qualifier_primitive(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@[*])]`,
		inputJSON:   `[1,2]`,
		expectedErr: createErrorMemberNotExist(`[?(@[*])]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_wildcard_qualifier_primitive")
}

func TestRetrieve_filterExist_filter_chain(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@)][?(@)]`,
		inputJSON:    `[1,[21,[221,[222]]]]`,
		expectedJSON: `[21,[221,[222]]]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_filter_chain")
}

func TestRetrieve_filterExist_array_slice_check_two_elements(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@[1:3])]`,
		inputJSON:    `[[{"a":1}],[{"b":2},{"c":3}],[],{"d":4}]`,
		expectedJSON: `[[{"b":2},{"c":3}]]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_array_slice_check_two_elements")
}

func TestRetrieve_filterExist_array_slice_check_negated_two_elements(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@[1:3])]`,
		inputJSON:    `[[{"a":1}],[{"b":2},{"c":3}],[],{"d":4}]`,
		expectedJSON: `[[{"a":1}],[],{"d":4}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_array_slice_check_negated_two_elements")
}

func TestRetrieve_filterExist_array_slice_check_three_elements(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@[1:3])]`,
		inputJSON:    `[[{"a":1}],[{"b":2},{"c":3},{"e":5}],[],{"d":4}]`,
		expectedJSON: `[[{"b":2},{"c":3},{"e":5}]]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_array_slice_check_three_elements")
}

func TestRetrieve_filterExist_array_slice_check_negated_three_elements(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@[1:3])]`,
		inputJSON:    `[[{"a":1}],[{"b":2},{"c":3},{"e":5}],[],{"d":4}]`,
		expectedJSON: `[[{"a":1}],[],{"d":4}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_array_slice_check_negated_three_elements")
}

func TestRetrieve_filterExist_array_slice_start_range(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@[0:1])]`,
		inputJSON:    `[[{"a":1}],[]]`,
		expectedJSON: `[[{"a":1}]]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_array_slice_start_range")
}

func TestRetrieve_filterExist_root_reference_object(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?($)]`,
		inputJSON:    `{"a":1,"b":2}`,
		expectedJSON: `[1,2]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_root_reference_object")
}

func TestRetrieve_filterExist_current_node_object(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@)]`,
		inputJSON:    `{"a":1,"b":2}`,
		expectedJSON: `[1,2]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_current_node_object")
}

func TestRetrieve_filterExist_nested_property_a1(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a1)]`,
		inputJSON:    `{"a":{"a1":1},"b":{"b1":2}}`,
		expectedJSON: `[{"a1":1}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_nested_property_a1")
}

func TestRetrieve_filterExist_nested_property_a1_negated(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@.a1)]`,
		inputJSON:    `{"a":{"a1":1},"b":{"b1":2}}`,
		expectedJSON: `[{"b1":2}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_nested_property_a1_negated")
}

func TestRetrieve_filterExist_recursive_descent_property(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@..a)]`,
		inputJSON:    `[{"a":1},{"b":2},{"c":{"a":3}},{"a":{"a":4}}]`,
		expectedJSON: `[{"a":1},{"c":{"a":3}},{"a":{"a":4}}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_recursive_descent_property")
}

func TestRetrieve_filterExist_recursive_descent_property_negated(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@..a)]`,
		inputJSON:    `[{"a":1},{"b":2},{"c":{"a":3}},{"a":{"a":4}}]`,
		expectedJSON: `[{"b":2}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_recursive_descent_property_negated")
}

func TestRetrieve_filterExist_array_index_check(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@[1])]`,
		inputJSON:    `{"a":["a1"],"b":["b1","b2"],"c":[],"d":4}`,
		expectedJSON: `[["b1","b2"]]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_array_index_check")
}

func TestRetrieve_filterExist_array_index_check_negated(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@[1])]`,
		inputJSON:    `{"a":["a1"],"b":["b1","b2"],"c":[],"d":4}`,
		expectedJSON: `[["a1"],[],4]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_array_index_check_negated")
}

func TestRetrieve_filterExist_array_slice_object_property_multiple(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@[1:3])]`,
		inputJSON:    `{"a":[],"b":[2],"c":[3,4,5,6],"d":4}`,
		expectedJSON: `[[3,4,5,6]]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_array_slice_object_property_multiple")
}

func TestRetrieve_filterExist_array_slice_object_property_multiple_negated(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@[1:3])]`,
		inputJSON:    `{"a":[],"b":[2],"c":[3,4,5,6],"d":4}`,
		expectedJSON: `[[],[2],4]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_array_slice_object_property_multiple_negated")
}

func TestRetrieve_filterExist_array_slice_object_property_exact(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@[1:3])]`,
		inputJSON:    `{"a":[],"b":[2],"c":[3,4],"d":4}`,
		expectedJSON: `[[3,4]]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_array_slice_object_property_exact")
}

func TestRetrieve_filterExist_array_slice_object_property_exact_negated(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@[1:3])]`,
		inputJSON:    `{"a":[],"b":[2],"c":[3,4],"d":4}`,
		expectedJSON: `[[],[2],4]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_array_slice_object_property_exact_negated")
}

func TestRetrieve_filterExist_root_reference_array(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?($)]`,
		inputJSON:    `[{"a":1},{"b":2}]`,
		expectedJSON: `[{"a":1},{"b":2}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_root_reference_array")
}

func TestRetrieve_filterExist_root_property_reference(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?($[0].a)]`,
		inputJSON:    `[{"a":1},{"b":2}]`,
		expectedJSON: `[{"a":1},{"b":2}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_root_property_reference")
}

func TestRetrieve_filterExist_root_property_reference_negated(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(!$[0].a)]`,
		inputJSON:   `[{"a":1},{"b":2}]`,
		expectedErr: createErrorMemberNotExist(`[?(!$[0].a)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_root_property_reference_negated")
}

func TestRetrieve_filterExist_union_property_access(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@['a','b'])]`,
		inputJSON:    `[{"a":1},{"b":2}]`,
		expectedJSON: `[{"a":1},{"b":2}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_union_property_access")
}

func TestRetrieve_filterExist_union_property_partial_match(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@['a','c'])]`,
		inputJSON:    `[{"a":1},{"b":2}]`,
		expectedJSON: `[{"a":1}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_union_property_partial_match")
}

func TestRetrieve_filterExist_union_property_no_match(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@['c','d'])]`,
		inputJSON:   `[{"a":1},{"b":2}]`,
		expectedErr: createErrorMemberNotExist(`[?(@['c','d'])]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_union_property_no_match")
}

func TestRetrieve_filterExist_union_array_indices_with_values(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@[0,1])]`,
		inputJSON:    `[[{"a":1}],[0,1]]`,
		expectedJSON: `[[{"a":1}],[0,1]]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_union_array_indices_with_values")
}

func TestRetrieve_filterExist_union_array_indices_partial_empty(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@[0,1])]`,
		inputJSON:    `[[{"a":1}],[]]`,
		expectedJSON: `[[{"a":1}]]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_union_array_indices_partial_empty")
}

func TestRetrieve_filterExist_nested_filter_condition(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a[?(@.b > 1)])]`,
		inputJSON:    `[{"a":[{"b":1},{"c":3}]},{"a":[{"b":2},{"c":5}]},{"b":4}]`,
		expectedJSON: `[{"a":[{"b":2},{"c":5}]}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_nested_filter_condition")
}

func TestRetrieve_filterExist_single_parentheses(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?((@.a>1))]`,
		inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
		expectedJSON: `[{"a":2},{"a":3}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_single_parentheses")
}

func TestRetrieve_filterExist_double_parentheses(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(((@.a>1)))]`,
		inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
		expectedJSON: `[{"a":2},{"a":3}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_double_parentheses")
}

func TestRetrieve_filterExist_chained_filters_triple(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@)][?(@)][?(@)]`,
		inputJSON:    `[1,[21,[221,[222]]]]`,
		expectedJSON: `[221,[222]]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_chained_filters_triple")
}

func TestRetrieve_filterExist_chained_filters_quadruple(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@)][?(@)][?(@)][?(@)]`,
		inputJSON:    `[1,[21,[221,[222]]]]`,
		expectedJSON: `[222]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_chained_filters_quadruple")
}

func TestRetrieve_filterExist_nested_filter_basic(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a[?(@.b)])]`,
		inputJSON:    `[{"a":[{"b":2},{"c":3}]},{"b":4}]`,
		expectedJSON: `[{"a":[{"b":2},{"c":3}]}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterExist_nested_filter_basic")
}

func TestRetrieve_filterExistsFromDeleted(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@)]`,
			inputJSON:    `{"a":1,"b":null}`,
			expectedJSON: `[1,null]`,
		},
		{
			jsonpath:     `$[?(@.a)]`,
			inputJSON:    `{"a":{"a":1},"b":{"b":2}}`,
			expectedJSON: `[{"a":1}]`,
		},
	}

	runTestCases(t, "TestRetrieve_filterExistsFromDeleted", testCases)
}

func TestRetrieve_filterEmptyInputFastPath(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@)]`,
			inputJSON:   `{}`,
			expectedErr: createErrorMemberNotExist(`[?(@)]`),
		},
		{
			jsonpath:    `$[?(@)]`,
			inputJSON:   `[]`,
			expectedErr: createErrorMemberNotExist(`[?(@)]`),
		},
	}

	runTestCases(t, "TestFilterQualifier_EmptyMapAndList_ReturnsMemberNotExist", testCases)
}
