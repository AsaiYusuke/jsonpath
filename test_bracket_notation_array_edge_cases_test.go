package jsonpath

import "testing"

func TestRetrieve_arrayAccess_on_empty_array(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[1]`,
		inputJSON:   `[]`,
		expectedErr: createErrorMemberNotExist(`[1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arrayAccess_on_empty_array")
}

func TestRetrieve_arrayAccess_large_index_overflow(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[1000000000000000000]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[1000000000000000000]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arrayAccess_large_index_overflow")
}

func TestRetrieve_bracketNotation_property_on_array(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$['a']`,
		inputJSON:   `[]`,
		expectedErr: createErrorTypeUnmatched(`['a']`, `object`, `[]interface {}`),
	}
	runTestCase(t, testCase, "TestRetrieve_bracketNotation_property_on_array")
}

func TestRetrieve_bracketNotation_property_on_object(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$['a']`,
		inputJSON:   `{}`,
		expectedErr: createErrorMemberNotExist(`['a']`),
	}
	runTestCase(t, testCase, "TestRetrieve_bracketNotation_property_on_object")
}

func TestRetrieve_bracketNotation_property_on_number(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$['a']`,
		inputJSON:   `123`,
		expectedErr: createErrorTypeUnmatched(`['a']`, `object`, `float64`),
	}
	runTestCase(t, testCase, "TestRetrieve_bracketNotation_property_on_number")
}

func TestRetrieve_bracketNotation_property_on_boolean(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$['a']`,
		inputJSON:   `true`,
		expectedErr: createErrorTypeUnmatched(`['a']`, `object`, `bool`),
	}
	runTestCase(t, testCase, "TestRetrieve_bracketNotation_property_on_boolean")
}

func TestRetrieve_bracketNotation_property_on_null(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$['a']`,
		inputJSON:   `null`,
		expectedErr: createErrorTypeUnmatched(`['a']`, `object`, `null`),
	}
	runTestCase(t, testCase, "TestRetrieve_bracketNotation_property_on_null")
}

func TestRetrieve_bracketNotation_mixed_dot_property_access(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.a['b']`,
		inputJSON:    `{"b":2,"a":{"b":1}}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestRetrieve_bracketNotation_mixed_dot_property_access")
}
