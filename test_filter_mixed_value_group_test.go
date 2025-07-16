package jsonpath

import (
	"testing"
)

func TestRetrieve_filterValueGroup_array_wildcard_comparison_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.x[?(@[*]>=$.y[*])]`,
		inputJSON:   `{"x":[[1,2],[3,4],[5,6]],"y":[3,4,5]}`,
		expectedErr: ErrorInvalidSyntax{position: 6, reason: `JSONPath that returns a value group is prohibited`, near: `@[*]>=$.y[*])]`},
	}
	runTestCase(t, testCase, "TestRetrieve_filterValueGroup_array_wildcard_comparison_error")
}
