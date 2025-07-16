package jsonpath

import (
	"testing"
)

func TestRetrieve_filterSubFilter_allowed_basic(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a[?(@.b>1)])]`,
		inputJSON:    `[{"a":[{"b":1},{"b":2}]},{"a":[{"b":1}]}]`,
		expectedJSON: `[{"a":[{"b":1},{"b":2}]}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterSubFilter_allowed_basic")
}

func TestRetrieve_filterSubFilter_prohibited_comparison_1(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a[?(@.b)] > 1)]`,
		inputJSON:   `[{"a":[{"b":1},{"b":2}]},{"a":[{"b":1}]}]`,
		expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[?(@.b)] > 1)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_filterSubFilter_prohibited_comparison_1")
}

func TestRetrieve_filterSubFilter_prohibited_comparison_2(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a[?(@.b)] > 1)]`,
		inputJSON:   `[{"a":[{"b":2}]},{"a":[{"b":1}]}]`,
		expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[?(@.b)] > 1)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_filterSubFilter_prohibited_comparison_2")
}

func TestRetrieve_filterSubFilter_prohibited_comparison_no_match(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a[?(@.b)] > 1)]`,
		inputJSON:   `[{"a":[{"c":2}]},{"a":[{"d":1}]}]`,
		expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[?(@.b)] > 1)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_filterSubFilter_prohibited_comparison_no_match")
}
