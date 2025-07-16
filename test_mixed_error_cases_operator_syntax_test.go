package jsonpath

import (
	"testing"
)

func TestRetrieve_invalid_operator_double_not_equal_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a!!=1)]`,
		inputJSON:   `[]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a!!=1)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_operator_double_not_equal_error")
}

func TestRetrieve_invalid_operator_missing_value_equal_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a==)]`,
		inputJSON:   `[]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_operator_missing_value_equal_error")
}

func TestRetrieve_invalid_operator_missing_value_not_equal_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a!=)]`,
		inputJSON:   `[]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a!=)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_operator_missing_value_not_equal_error")
}

func TestRetrieve_invalid_operator_missing_value_less_equal_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a<=)]`,
		inputJSON:   `[]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a<=)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_operator_missing_value_less_equal_error")
}

func TestRetrieve_invalid_operator_missing_value_less_than_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a<)]`,
		inputJSON:   `[]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a<)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_operator_missing_value_less_than_error")
}

func TestRetrieve_invalid_operator_missing_value_greater_equal_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a>=)]`,
		inputJSON:   `[]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a>=)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_operator_missing_value_greater_equal_error")
}

func TestRetrieve_invalid_operator_missing_value_greater_than_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a>)]`,
		inputJSON:   `[]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a>)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_operator_missing_value_greater_than_error")
}

func TestRetrieve_invalid_operator_leading_equal_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(==@.a)]`,
		inputJSON:   `[]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(==@.a)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_operator_leading_equal_error")
}

func TestRetrieve_invalid_operator_leading_not_equal_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(!=@.a)]`,
		inputJSON:   `[]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(!=@.a)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_operator_leading_not_equal_error")
}

func TestRetrieve_invalid_operator_leading_less_equal_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(<=@.a)]`,
		inputJSON:   `[]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(<=@.a)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_operator_leading_less_equal_error")
}

func TestRetrieve_invalid_operator_leading_greater_equal_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(>=@.a)]`,
		inputJSON:   `[]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(>=@.a)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_operator_leading_greater_equal_error")
}

func TestRetrieve_invalid_operator_leading_greater_than_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(>@.a)]`,
		inputJSON:   `[]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(>@.a)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_operator_leading_greater_than_error")
}

func TestRetrieve_invalid_operator_triple_equal_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a===1)]`,
		inputJSON:   `[]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a===1)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_operator_triple_equal_error")
}

func TestRetrieve_invalid_operator_single_equal_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a=2)]`,
		inputJSON:   `[{"a":2}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=2)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_operator_single_equal_error")
}

func TestRetrieve_invalid_operator_not_equal_variation_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a<>2)]`,
		inputJSON:   `[{"a":2}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a<>2)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_operator_not_equal_variation_error")
}

func TestRetrieve_invalid_operator_equal_less_than_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a=<2)]`,
		inputJSON:   `[{"a":2}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=<2)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_operator_equal_less_than_error")
}

func TestRetrieve_invalid_literal_false_filter_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(false)]`,
		inputJSON:   `[0,1,false,true,null,{},[]]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(false)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_literal_false_filter_error")
}

func TestRetrieve_invalid_literal_true_filter_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(true)]`,
		inputJSON:   `[0,1,false,true,null,{},[]]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(true)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_literal_true_filter_error")
}

func TestRetrieve_invalid_literal_null_filter_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(null)]`,
		inputJSON:   `[0,1,false,true,null,{},[]]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(null)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_literal_null_filter_error")
}

func TestRetrieve_invalid_array_comparison_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a==["b"])]`,
		inputJSON:   `[{"a":["b"]}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==["b"])]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_array_comparison_error")
}

func TestRetrieve_invalid_slice_comparison_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@[0:1]==[1])]`,
		inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@[0:1]==[1])]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_slice_comparison_error")
}

func TestRetrieve_invalid_object_comparison_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@=={"k":"v"})]`,
		inputJSON:   `{}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@=={"k":"v"})]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_object_comparison_error")
}
