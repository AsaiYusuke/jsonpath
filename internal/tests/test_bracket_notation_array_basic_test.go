package tests

import (
	"testing"
)

func TestRetrieve_arrayIndex_basic_positive_index_0(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_basic_positive_index_0")
}

func TestRetrieve_arrayIndex_basic_positive_index_1(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["second"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_basic_positive_index_1")
}

func TestRetrieve_arrayIndex_basic_positive_index_2(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[2]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_basic_positive_index_2")
}

func TestRetrieve_arrayIndex_basic_positive_index_out_of_bounds(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[3]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[3]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_basic_positive_index_out_of_bounds")
}

func TestRetrieve_arrayIndex_basic_negative_index_minus_1(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_basic_negative_index_minus_1")
}

func TestRetrieve_arrayIndex_basic_negative_index_minus_2(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[-2]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["second"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_basic_negative_index_minus_2")
}

func TestRetrieve_arrayIndex_basic_negative_index_minus_3(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[-3]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_basic_negative_index_minus_3")
}

func TestRetrieve_arrayIndex_basic_negative_index_out_of_bounds(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[-4]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[-4]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_basic_negative_index_out_of_bounds")
}

func TestRetrieve_arrayIndex_syntax_positive_plus_sign(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[+1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["second"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_syntax_positive_plus_sign")
}

func TestRetrieve_arrayIndex_syntax_leading_zero(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[01]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["second"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_syntax_leading_zero")
}

func TestRetrieve_arrayIndex_syntax_decimal_invalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[1.0]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[1.0]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_syntax_decimal_invalid")
}

func TestRetrieve_arrayIndex_type_unmatched_object(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0]`,
		inputJSON:   `{}`,
		expectedErr: createErrorTypeUnmatched(`[0]`, `array`, `map[string]interface {}`),
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_type_unmatched_object")
}

func TestRetrieve_arrayIndex_empty_array(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0]`,
		inputJSON:   `[]`,
		expectedErr: createErrorMemberNotExist(`[0]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_empty_array")
}

func TestRetrieve_arrayIndex_negative_empty_array_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[-1]`,
		inputJSON:   `[]`,
		expectedErr: createErrorMemberNotExist(`[-1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_negative_empty_array_error")
}

func TestRetrieve_arrayIndex_negative_very_large_number_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[-1000000000000000000]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[-1000000000000000000]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_negative_very_large_number_error")
}

func TestRetrieve_arrayIndex_nested_array_access(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0][1]`,
		inputJSON:    `[["a","b"],["c"],["d"]]`,
		expectedJSON: `["b"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_nested_array_access")
}

func TestRetrieve_arrayIndex_triple_nested_array_access(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0][1][2]`,
		inputJSON:    `[["a",["b","c","d"]],["e"],["f"]]`,
		expectedJSON: `["d"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_triple_nested_array_access")
}

func TestRetrieve_arrayIndex_on_object(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0]`,
		inputJSON:   `{"a":1,"b":2}`,
		expectedErr: createErrorTypeUnmatched(`[0]`, `array`, `map[string]interface {}`),
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_on_object")
}

func TestRetrieve_arrayIndex_on_string(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0]`,
		inputJSON:   `"abc"`,
		expectedErr: createErrorTypeUnmatched(`[0]`, `array`, `string`),
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_on_string")
}

func TestRetrieve_arrayIndex_on_number(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0]`,
		inputJSON:   `123`,
		expectedErr: createErrorTypeUnmatched(`[0]`, `array`, `float64`),
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_on_number")
}

func TestRetrieve_arrayIndex_on_boolean(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0]`,
		inputJSON:   `true`,
		expectedErr: createErrorTypeUnmatched(`[0]`, `array`, `bool`),
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_on_boolean")
}

func TestRetrieve_arrayIndex_on_null(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0]`,
		inputJSON:   `null`,
		expectedErr: createErrorTypeUnmatched(`[0]`, `array`, `null`),
	}
	runTestCase(t, testCase, "TestRetrieve_arrayIndex_on_null")
}
