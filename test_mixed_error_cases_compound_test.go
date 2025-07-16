package jsonpath

import (
	"testing"
)

// Invalid syntax tests - restoring deleted tests from syntax error patterns

func TestRetrieve_invalid_syntax_dot_bracket_combination_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.["a"]`,
		inputJSON:   `{"a":1}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.["a"]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_syntax_dot_bracket_combination_error")
}

func TestRetrieve_invalid_syntax_union_dot_bracket_1_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1].[1]`,
		inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
		expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `.[1]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_syntax_union_dot_bracket_1_error")
}

func TestRetrieve_invalid_syntax_slice_dot_bracket_1_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:2].[1]`,
		inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
		expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `.[1]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_syntax_slice_dot_bracket_1_error")
}

func TestRetrieve_invalid_syntax_single_dot_bracket_union_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0].[1,2]`,
		inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
		expectedErr: ErrorInvalidSyntax{position: 4, reason: `unrecognized input`, near: `.[1,2]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_syntax_single_dot_bracket_union_error")
}

func TestRetrieve_invalid_syntax_union_dot_bracket_union_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1].[1,2]`,
		inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
		expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `.[1,2]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_syntax_union_dot_bracket_union_error")
}

func TestRetrieve_invalid_syntax_slice_dot_bracket_union_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:2].[1,2]`,
		inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
		expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `.[1,2]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_syntax_slice_dot_bracket_union_error")
}

func TestRetrieve_invalid_syntax_union_dot_bracket_slice_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1].[1:3]`,
		inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
		expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `.[1:3]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_syntax_union_dot_bracket_slice_error")
}

func TestRetrieve_invalid_syntax_slice_dot_bracket_slice_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:1].[1:3]`,
		inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
		expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `.[1:3]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_syntax_slice_dot_bracket_slice_error")
}

func TestRetrieve_invalid_syntax_multiple_filter_union_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a),?(@.b)]`,
		inputJSON:   `{}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a),?(@.b)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_syntax_multiple_filter_union_error")
}

func TestRetrieve_invalid_syntax_empty_bracket_object_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[]`,
		inputJSON:   `{"a":1}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_syntax_empty_bracket_object_error")
}

func TestRetrieve_invalid_syntax_empty_bracket_array_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_syntax_empty_bracket_array_error")
}

func TestRetrieve_invalid_syntax_unclosed_bracket_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_syntax_unclosed_bracket_error")
}

func TestRetrieve_invalid_syntax_incomplete_filter_quote_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a=='abc`,
		inputJSON:   `[]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=='abc`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_syntax_incomplete_filter_quote_error")
}

func TestRetrieve_invalid_syntax_incomplete_filter_double_quote_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a=="abc`,
		inputJSON:   `[]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=="abc`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_syntax_incomplete_filter_double_quote_error")
}

func TestRetrieve_invalid_syntax_malformed_filter_parentheses_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?((@.a>1 )]`,
		inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?((@.a>1 )]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_syntax_malformed_filter_parentheses_error")
}

func TestRetrieve_invalid_syntax_incomplete_filter_parentheses_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?((@.a>1`,
		inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?((@.a>1`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_syntax_incomplete_filter_parentheses_error")
}

func TestRetrieve_invalid_syntax_large_number_in_slice_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:10000000000000000000:a]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[0:10000000000000000000:a]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_syntax_large_number_in_slice_error")
}
