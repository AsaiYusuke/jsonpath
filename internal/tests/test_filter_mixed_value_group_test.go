package tests

import (
	"testing"
)

func TestRetrieve_filterValueGroup_array_wildcard_comparison_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.x[?(@[*]>=$.y[*])]`,
		inputJSON:   `{"x":[[1,2],[3,4],[5,6]],"y":[3,4,5]}`,
		expectedErr: createErrorInvalidSyntax(6, `JSONPath that returns a value group is prohibited`, `@[*]>=$.y[*])]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterValueGroup_array_wildcard_comparison_error")
}
