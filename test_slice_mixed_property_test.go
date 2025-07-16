package jsonpath

import "testing"

func TestRetrieve_arraySliceProperty_missing_property1(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[1:2].a.b`,
		inputJSON:   `[0]`,
		expectedErr: createErrorMemberNotExist(`[1:2]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySliceProperty_missing_property1")
}

func TestRetrieve_arraySliceProperty_missing_property2(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:2].a.b`,
		inputJSON:   `[{"b":1}]`,
		expectedErr: createErrorMemberNotExist(`.a`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySliceProperty_missing_property2")
}

func TestRetrieve_arraySliceProperty_missing_property3(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:2].a.b`,
		inputJSON:   `[{"b":1},{"c":2}]`,
		expectedErr: createErrorMemberNotExist(`.a`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySliceProperty_missing_property3")
}

func TestRetrieve_arraySliceProperty_deeply_nested_missing(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:2].a.b.c`,
		inputJSON:   `[{"a":1},{"b":2}]`,
		expectedErr: createErrorTypeUnmatched(`.b`, `object`, `float64`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySliceProperty_deeply_nested_missing")
}
