package jsonpath

import "testing"

func TestRetrieve_filterLogicalCombination_NOT_basic(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@.a)]`,
		inputJSON:    `[{"a":1},{"b":2},{"a":3,"b":4}]`,
		expectedJSON: `[{"b":2}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_NOT_basic")
}

func TestRetrieve_filterLogicalCombination_NOT_missing_field(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@.c)]`,
		inputJSON:    `[{"a":1},{"b":2},{"a":3,"b":4}]`,
		expectedJSON: `[{"a":1},{"b":2},{"a":3,"b":4}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_NOT_missing_field")
}

func TestRetrieve_filterLogicalCombination_NOT_recursive_descent_exists_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.z[?(!$..x)]`,
		inputJSON:   `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
		expectedErr: createErrorMemberNotExist(`[?(!$..x)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_NOT_recursive_descent_exists_error")
}

func TestRetrieve_filterLogicalCombination_NOT_recursive_descent_missing(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.z[?(!$..xx)]`,
		inputJSON:    `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
		expectedJSON: `[{"b":1},{"b":2},{"b":3}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_NOT_recursive_descent_missing")
}

func TestRetrieve_filterLogicalCombination_NOT_root_reference_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(!$)]`,
		inputJSON:   `{"a":1,"b":2}`,
		expectedErr: createErrorMemberNotExist(`[?(!$)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_NOT_root_reference_error")
}

func TestRetrieve_filterLogicalCombination_NOT_current_node_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(!@)]`,
		inputJSON:   `{"a":1}`,
		expectedErr: createErrorMemberNotExist(`[?(!@)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_NOT_current_node_error")
}
