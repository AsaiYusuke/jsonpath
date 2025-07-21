package tests

import "testing"

func TestRetrieve_filterLogicalInvalid_missing_right_operand(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a>1 && )]`,
		inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
		expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a>1 && )]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalInvalid_missing_right_operand")
}

func TestRetrieve_filterLogicalInvalid_missing_left_operand(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?( && @.a>1 )]`,
		inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
		expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?( && @.a>1 )]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalInvalid_missing_left_operand")
}

func TestRetrieve_filterLogical_and_with_false_literal(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a>1 && false)]`,
		inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
		expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a>1 && false)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogical_and_with_false_literal")
}

func TestRetrieve_filterLogical_and_with_true_literal(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a>1 && true)]`,
		inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
		expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a>1 && true)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogical_and_with_true_literal")
}

func TestRetrieve_filterLogicalInvalid_or_missing_right_operand(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a>1 || )]`,
		inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
		expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a>1 || )]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalInvalid_or_missing_right_operand")
}

func TestRetrieve_filterLogicalInvalid_or_missing_left_operand(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?( || @.a>1 )]`,
		inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
		expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?( || @.a>1 )]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalInvalid_or_missing_left_operand")
}

func TestRetrieve_filterLogical_or_with_false_literal(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a>1 || false)]`,
		inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
		expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a>1 || false)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogical_or_with_false_literal")
}

func TestRetrieve_filterLogical_or_with_true_literal(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a>1 || true)]`,
		inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
		expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a>1 || true)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogical_or_with_true_literal")
}
