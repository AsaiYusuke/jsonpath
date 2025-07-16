package jsonpath

import (
	"testing"
)

// Compound array access tests - restoring deleted tests from $[0,1][0:2] pattern group

func TestRetrieve_compound_array_union_slice_success(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0,1][0:2]`,
		inputJSON:    `[[1,2,3],[4,5,6]]`,
		expectedJSON: `[1,2,4,5]`,
	}
	runTestCase(t, testCase, "TestRetrieve_compound_array_union_slice_success")
}

func TestRetrieve_compound_array_union_slice_empty_arrays_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1][0:2]`,
		inputJSON:   `[[],[]]`,
		expectedErr: createErrorMemberNotExist(`[0:2]`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_array_union_slice_empty_arrays_error")
}

func TestRetrieve_compound_array_union_slice_empty_root_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1][0:2]`,
		inputJSON:   `[]`,
		expectedErr: createErrorMemberNotExist(`[0,1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_array_union_slice_empty_root_error")
}

func TestRetrieve_compound_array_union_slice_type_mismatch_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1][0:2]`,
		inputJSON:   `"x"`,
		expectedErr: createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_array_union_slice_type_mismatch_error")
}

func TestRetrieve_compound_array_union_wildcard_success(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0,1][*]`,
		inputJSON:    `[{"a":1,"c":3},{"d":4,"b":2},{"e":5}]`,
		expectedJSON: `[1,3,2,4]`,
	}
	runTestCase(t, testCase, "TestRetrieve_compound_array_union_wildcard_success")
}

func TestRetrieve_compound_array_union_wildcard_empty_objects_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1][*]`,
		inputJSON:   `[{},{}]`,
		expectedErr: createErrorMemberNotExist(`[*]`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_array_union_wildcard_empty_objects_error")
}

func TestRetrieve_compound_array_union_wildcard_empty_root_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1][*]`,
		inputJSON:   `[]`,
		expectedErr: createErrorMemberNotExist(`[0,1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_array_union_wildcard_empty_root_error")
}

func TestRetrieve_compound_array_union_wildcard_type_mismatch_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1][*]`,
		inputJSON:   `"x"`,
		expectedErr: createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_array_union_wildcard_type_mismatch_error")
}

func TestRetrieve_compound_array_union_union_success(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0,1][0,1]`,
		inputJSON:    `[[1,3,2],[4,6,5],[7]]`,
		expectedJSON: `[1,3,4,6]`,
	}
	runTestCase(t, testCase, "TestRetrieve_compound_array_union_union_success")
}

func TestRetrieve_compound_array_union_union_empty_arrays_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1][0,1]`,
		inputJSON:   `[[],[],[7]]`,
		expectedErr: createErrorMemberNotExist(`[0,1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_array_union_union_empty_arrays_error")
}

func TestRetrieve_compound_array_union_union_empty_root_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1][0,1]`,
		inputJSON:   `[]`,
		expectedErr: createErrorMemberNotExist(`[0,1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_array_union_union_empty_root_error")
}

func TestRetrieve_compound_array_union_union_type_mismatch_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1][0,1]`,
		inputJSON:   `"x"`,
		expectedErr: createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
	}
	runTestCase(t, testCase, "TestRetrieve_compound_array_union_union_type_mismatch_error")
}
