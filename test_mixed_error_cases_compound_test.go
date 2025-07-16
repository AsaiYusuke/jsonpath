package jsonpath

import (
	"testing"
)

func TestError_InvalidSyntaxDotBracketCombination(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.["a"]`,
		inputJSON:   `{"a":1}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.["a"]`},
	}
	runTestCase(t, testCase, "TestError_InvalidSyntaxDotBracketCombination")
}

func TestError_InvalidSyntaxUnionDotBracket1(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1].[1]`,
		inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
		expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `.[1]`},
	}
	runTestCase(t, testCase, "TestError_InvalidSyntaxUnionDotBracket1")
}

func TestError_InvalidSyntaxSliceDotBracket1(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:2].[1]`,
		inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
		expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `.[1]`},
	}
	runTestCase(t, testCase, "TestError_InvalidSyntaxSliceDotBracket1")
}

func TestError_InvalidSyntaxSingleDotBracketUnion(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0].[1,2]`,
		inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
		expectedErr: ErrorInvalidSyntax{position: 4, reason: `unrecognized input`, near: `.[1,2]`},
	}
	runTestCase(t, testCase, "TestError_InvalidSyntaxSingleDotBracketUnion")
}

func TestError_InvalidSyntaxUnionDotBracketUnion(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1].[1,2]`,
		inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
		expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `.[1,2]`},
	}
	runTestCase(t, testCase, "TestError_InvalidSyntaxUnionDotBracketUnion")
}

func TestError_InvalidSyntaxSliceDotBracketUnion(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:2].[1,2]`,
		inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
		expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `.[1,2]`},
	}
	runTestCase(t, testCase, "TestError_InvalidSyntaxSliceDotBracketUnion")
}

func TestError_InvalidSyntaxunion_dot_bracket_slice_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1].[1:3]`,
		inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
		expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `.[1:3]`},
	}
	runTestCase(t, testCase, "TestError_InvalidSyntaxunion_dot_bracket_slice_error")
}

func TestError_InvalidSyntaxslice_dot_bracket_slice_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:1].[1:3]`,
		inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
		expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `.[1:3]`},
	}
	runTestCase(t, testCase, "TestError_InvalidSyntaxslice_dot_bracket_slice_error")
}

func TestError_InvalidSyntaxmultiple_filter_union_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a),?(@.b)]`,
		inputJSON:   `{}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a),?(@.b)]`},
	}
	runTestCase(t, testCase, "TestError_InvalidSyntaxmultiple_filter_union_error")
}

func TestError_InvalidSyntaxempty_bracket_object_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[]`,
		inputJSON:   `{"a":1}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[]`},
	}
	runTestCase(t, testCase, "TestError_InvalidSyntaxempty_bracket_object_error")
}

func TestError_InvalidSyntaxempty_bracket_array_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[]`},
	}
	runTestCase(t, testCase, "TestError_InvalidSyntaxempty_bracket_array_error")
}

func TestError_InvalidSyntaxunclosed_bracket_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[`},
	}
	runTestCase(t, testCase, "TestError_InvalidSyntaxunclosed_bracket_error")
}

func TestError_InvalidSyntaxincomplete_filter_quote_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a=='abc`,
		inputJSON:   `[]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=='abc`},
	}
	runTestCase(t, testCase, "TestError_InvalidSyntaxincomplete_filter_quote_error")
}

func TestError_InvalidSyntaxincomplete_filter_double_quote_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a=="abc`,
		inputJSON:   `[]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=="abc`},
	}
	runTestCase(t, testCase, "TestError_InvalidSyntaxincomplete_filter_double_quote_error")
}

func TestError_InvalidSyntaxmalformed_filter_parentheses_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?((@.a>1 )]`,
		inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?((@.a>1 )]`},
	}
	runTestCase(t, testCase, "TestError_InvalidSyntaxmalformed_filter_parentheses_error")
}

func TestError_InvalidSyntaxincomplete_filter_parentheses_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?((@.a>1`,
		inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?((@.a>1`},
	}
	runTestCase(t, testCase, "TestError_InvalidSyntaxincomplete_filter_parentheses_error")
}

func TestError_InvalidSyntaxlarge_number_in_slice_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:10000000000000000000:a]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[0:10000000000000000000:a]`},
	}
	runTestCase(t, testCase, "TestError_InvalidSyntaxlarge_number_in_slice_error")
}
