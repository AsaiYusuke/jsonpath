package jsonpath

import (
	"testing"
)

// Invalid literal value tests - restoring deleted tests from literal value error patterns

func TestRetrieve_invalid_literal_false_case_variations_error_1(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a==fAlse)]`,
		inputJSON:   `[{"a":false}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==fAlse)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_literal_false_case_variations_error_1")
}

func TestRetrieve_invalid_literal_false_case_variations_error_2(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a==faLse)]`,
		inputJSON:   `[{"a":false}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==faLse)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_literal_false_case_variations_error_2")
}

func TestRetrieve_invalid_literal_false_case_variations_error_3(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a==falsE)]`,
		inputJSON:   `[{"a":false}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==falsE)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_literal_false_case_variations_error_3")
}

func TestRetrieve_invalid_literal_false_case_variations_error_4(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a==FaLse)]`,
		inputJSON:   `[{"a":false}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==FaLse)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_literal_false_case_variations_error_4")
}

func TestRetrieve_invalid_literal_false_case_variations_error_5(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a==FalSe)]`,
		inputJSON:   `[{"a":false}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==FalSe)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_literal_false_case_variations_error_5")
}

func TestRetrieve_invalid_literal_false_case_variations_error_6(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a==FalsE)]`,
		inputJSON:   `[{"a":false}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==FalsE)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_literal_false_case_variations_error_6")
}

func TestRetrieve_invalid_literal_false_case_variations_error_7(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a==FAlSE)]`,
		inputJSON:   `[{"a":false}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==FAlSE)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_literal_false_case_variations_error_7")
}

func TestRetrieve_invalid_literal_false_case_variations_error_8(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a==FALsE)]`,
		inputJSON:   `[{"a":false}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==FALsE)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_literal_false_case_variations_error_8")
}

func TestRetrieve_invalid_literal_false_case_variations_error_9(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a==FALSe)]`,
		inputJSON:   `[{"a":false}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==FALSe)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_literal_false_case_variations_error_9")
}

func TestRetrieve_invalid_literal_true_case_variations_error_1(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a==tRue)]`,
		inputJSON:   `[{"a":true}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==tRue)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_literal_true_case_variations_error_1")
}

func TestRetrieve_invalid_literal_true_case_variations_error_2(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a==truE)]`,
		inputJSON:   `[{"a":true}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==truE)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_literal_true_case_variations_error_2")
}

func TestRetrieve_invalid_literal_true_case_variations_error_3(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a==TrUe)]`,
		inputJSON:   `[{"a":true}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==TrUe)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_literal_true_case_variations_error_3")
}

func TestRetrieve_invalid_literal_true_case_variations_error_4(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a==TruE)]`,
		inputJSON:   `[{"a":true}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==TruE)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_literal_true_case_variations_error_4")
}

func TestRetrieve_invalid_literal_true_case_variations_error_5(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a==TrUE)]`,
		inputJSON:   `[{"a":true}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==TrUE)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_literal_true_case_variations_error_5")
}

func TestRetrieve_invalid_literal_true_case_variations_error_6(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a==TRuE)]`,
		inputJSON:   `[{"a":true}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==TRuE)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_literal_true_case_variations_error_6")
}

func TestRetrieve_invalid_literal_true_case_variations_error_7(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a==TRUe)]`,
		inputJSON:   `[{"a":true}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==TRUe)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_literal_true_case_variations_error_7")
}

func TestRetrieve_invalid_literal_null_case_variations_error_1(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a==nUll)]`,
		inputJSON:   `[{"a":null}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==nUll)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_literal_null_case_variations_error_1")
}

func TestRetrieve_invalid_literal_null_case_variations_error_2(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a==nuLl)]`,
		inputJSON:   `[{"a":null}]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==nuLl)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_invalid_literal_null_case_variations_error_2")
}
