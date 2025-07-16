package jsonpath

import (
	"testing"
)

func TestRetrieve_dotNotationWildcardFilter_nested_array_property_filter(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*[?(@.a)]`,
		inputJSON:    `[[{"a":1},{"b":2}],[{"c":1},{"d":2}]]`,
		expectedJSON: `[{"a":1}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotationWildcardFilter_nested_array_property_filter")
}

func TestRetrieve_dotNotationWildcardFilter_array_property_filter(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*[?(@.a)]`,
		inputJSON:   `[{"a":1},{"b":2}]`,
		expectedErr: createErrorMemberNotExist(`[?(@.a)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotationWildcardFilter_array_property_filter")
}

func TestRetrieve_dotNotationWildcardFilter_object_array_property_filter(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*[?(@.a)]`,
		inputJSON:    `{"a":[{"a":1}],"b":[{"b":2}]}`,
		expectedJSON: `[{"a":1}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotationWildcardFilter_object_array_property_filter")
}

func TestRetrieve_dotNotationWildcardFilter_object_property_filter(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*[?(@.a)]`,
		inputJSON:   `{"a":{"a":1},"b":{"b":2}}`,
		expectedErr: createErrorMemberNotExist(`[?(@.a)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotationWildcardFilter_object_property_filter")
}

func TestRetrieve_dotNotationWildcardFilter_equals_one_in_array(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*[?(@==1)]`,
		inputJSON:    `[[1],{"b":2}]`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotationWildcardFilter_equals_one_in_array")
}

func TestRetrieve_dotNotationWildcardFilter_equals_one_object_values(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*[?(@==1)]`,
		inputJSON:    `[{"a":1},{"b":2}]`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotationWildcardFilter_equals_one_object_values")
}
