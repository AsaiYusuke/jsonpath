package jsonpath

import (
	"testing"
)

func TestRetrieve_valueGroupCombinationWildcardIdentifier_recursive_descent_basic(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*..a`,
		inputJSON:    `{"a":{"a":1,"c":2},"b":{"d":{"e":3,"a":4}}}`,
		expectedJSON: `[1,4]`,
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_recursive_descent_basic")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_recursive_descent_not_found(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*..a`,
		inputJSON:   `{"x":{"x":1,"c":2},"b":{"d":{"e":3,"x":4}}}`,
		expectedErr: createErrorMemberNotExist(`a`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_recursive_descent_not_found")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_recursive_descent_type_unmatched(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*..a`,
		inputJSON:   `{"a":"a","b":"b"}`,
		expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_recursive_descent_type_unmatched")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_recursive_descent_empty_object(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*..a`,
		inputJSON:   `{}`,
		expectedErr: createErrorMemberNotExist(`.*`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_recursive_descent_empty_object")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_recursive_descent_string_root(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*..a`,
		inputJSON:   `"x"`,
		expectedErr: createErrorTypeUnmatched(`.*`, `object/array`, `string`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_recursive_descent_string_root")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_multiple_identifier_basic(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*['a','b']`,
		inputJSON:    `{"a":{"a":1},"c":{"c":3},"b":{"b":2}}`,
		expectedJSON: `[1,2]`,
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_multiple_identifier_basic")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_multiple_identifier_not_found(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*['a','b']`,
		inputJSON:   `{"a":{"x":1},"c":{"c":3},"b":{"x":2}}`,
		expectedErr: createErrorMemberNotExist(`['a','b']`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_multiple_identifier_not_found")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_multiple_identifier_empty_object(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*['a','b']`,
		inputJSON:   `{}`,
		expectedErr: createErrorMemberNotExist(`.*`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_multiple_identifier_empty_object")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_multiple_identifier_string_root(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*['a','b']`,
		inputJSON:   `"x"`,
		expectedErr: createErrorTypeUnmatched(`.*`, `object/array`, `string`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_multiple_identifier_string_root")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_wildcard_identifier_basic(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*.*`,
		inputJSON:    `{"a":{"a":1,"c":2},"b":{"d":3,"a":4}}`,
		expectedJSON: `[1,2,4,3]`,
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_wildcard_identifier_basic")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_wildcard_identifier_empty_objects(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*.*`,
		inputJSON:   `{"a":{},"b":{}}`,
		expectedErr: createErrorMemberNotExist(`.*`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_wildcard_identifier_empty_objects")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_wildcard_identifier_empty_object(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*.*`,
		inputJSON:   `{}`,
		expectedErr: createErrorMemberNotExist(`.*`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_wildcard_identifier_empty_object")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_wildcard_identifier_string_root(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*.*`,
		inputJSON:   `"x"`,
		expectedErr: createErrorTypeUnmatched(`.*`, `object/array`, `string`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_wildcard_identifier_string_root")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_slice_qualifier_basic(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*[0:2]`,
		inputJSON:    `{"a":[1,3,2],"b":[4,6,5]}`,
		expectedJSON: `[1,3,4,6]`,
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_slice_qualifier_basic")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_slice_qualifier_empty_arrays(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*[0:2]`,
		inputJSON:   `{"a":[],"b":[]}`,
		expectedErr: createErrorMemberNotExist(`[0:2]`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_slice_qualifier_empty_arrays")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_slice_qualifier_empty_object(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*[0:2]`,
		inputJSON:   `{}`,
		expectedErr: createErrorMemberNotExist(`.*`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_slice_qualifier_empty_object")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_slice_qualifier_string_root(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*[0:2]`,
		inputJSON:   `"x"`,
		expectedErr: createErrorTypeUnmatched(`.*`, `object/array`, `string`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_slice_qualifier_string_root")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_wildcard_qualifier_basic(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*[*]`,
		inputJSON:    `{"a":[1,3,2],"b":[4,6,5]}`,
		expectedJSON: `[1,3,2,4,6,5]`,
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_wildcard_qualifier_basic")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_wildcard_qualifier_empty_arrays(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*[*]`,
		inputJSON:   `{"a":[],"b":[]}`,
		expectedErr: createErrorMemberNotExist(`[*]`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_wildcard_qualifier_empty_arrays")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_wildcard_qualifier_empty_object(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*[*]`,
		inputJSON:   `{}`,
		expectedErr: createErrorMemberNotExist(`.*`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_wildcard_qualifier_empty_object")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_wildcard_qualifier_string_root(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*[*]`,
		inputJSON:   `"x"`,
		expectedErr: createErrorTypeUnmatched(`.*`, `object/array`, `string`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_wildcard_qualifier_string_root")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_union_qualifier_basic(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*[0,1]`,
		inputJSON:    `{"a":[1,3,2],"b":[4,6,5]}`,
		expectedJSON: `[1,3,4,6]`,
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_union_qualifier_basic")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_union_qualifier_empty_arrays(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*[0,1]`,
		inputJSON:   `{"a":[],"b":[]}`,
		expectedErr: createErrorMemberNotExist(`[0,1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_union_qualifier_empty_arrays")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_union_qualifier_empty_object(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*[0,1]`,
		inputJSON:   `{}`,
		expectedErr: createErrorMemberNotExist(`.*`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_union_qualifier_empty_object")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_union_qualifier_string_root(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*[0,1]`,
		inputJSON:   `"x"`,
		expectedErr: createErrorTypeUnmatched(`.*`, `object/array`, `string`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_union_qualifier_string_root")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_filter_qualifier_basic(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*[?(@.b)]`,
		inputJSON:    `{"a":[{"a":1},{"b":2}],"b":[{"a":3},{"b":4}]}`,
		expectedJSON: `[{"b":2},{"b":4}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_filter_qualifier_basic")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_filter_qualifier_not_found(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*[?(@.b)]`,
		inputJSON:   `{"a":[{"a":1},{"x":2}],"b":[{"a":3},{"x":4}]}`,
		expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_filter_qualifier_not_found")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_filter_qualifier_empty_object(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*[?(@.b)]`,
		inputJSON:   `{}`,
		expectedErr: createErrorMemberNotExist(`.*`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_filter_qualifier_empty_object")
}

func TestRetrieve_valueGroupCombinationWildcardIdentifier_filter_qualifier_string_root(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*[?(@.b)]`,
		inputJSON:   `"x"`,
		expectedErr: createErrorTypeUnmatched(`.*`, `object/array`, `string`),
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupCombinationWildcardIdentifier_filter_qualifier_string_root")
}
