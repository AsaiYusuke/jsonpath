package jsonpath

import "testing"

func TestRetrieve_filterProperty_with_missing_property1(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a)].b`,
		inputJSON:   `[{"b":1}]`,
		expectedErr: createErrorMemberNotExist(`[?(@.a)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterProperty_with_missing_property1")
}

func TestRetrieve_filterProperty_with_missing_nested_property1(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a)].b.c`,
		inputJSON:   `[{"a":1}]`,
		expectedErr: createErrorMemberNotExist(`.b`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterProperty_with_missing_nested_property1")
}

func TestRetrieve_filterProperty_with_missing_nested_property2(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a)].b.c`,
		inputJSON:   `[{"a":1},{"a":1}]`,
		expectedErr: createErrorMemberNotExist(`.b`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterProperty_with_missing_nested_property2")
}

func TestRetrieve_filterProperty_with_complex_missing_path1(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a)].a.b.c`,
		inputJSON:   `[{"a":1},{"a":{"c":2}}]`,
		expectedErr: createErrorMemberNotExist(`.b`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterProperty_with_complex_missing_path1")
}

func TestRetrieve_filterProperty_with_complex_missing_path2(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a)].a.b.c`,
		inputJSON:   `[{"a":1},{"a":{"c":2}},{"b":3}]`,
		expectedErr: createErrorMemberNotExist(`.b`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterProperty_with_complex_missing_path2")
}

func TestRetrieve_filterProperty_object_root_missing_property1(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a)].b`,
		inputJSON:   `{"a":{"b":1}}`,
		expectedErr: createErrorMemberNotExist(`[?(@.a)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterProperty_object_root_missing_property1")
}

func TestRetrieve_filterProperty_object_root_missing_nested1(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a)].b.c`,
		inputJSON:   `{"a":{"a":1}}`,
		expectedErr: createErrorMemberNotExist(`.b`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterProperty_object_root_missing_nested1")
}

func TestRetrieve_filterProperty_object_root_missing_nested2(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a)].b.c`,
		inputJSON:   `{"a":{"a":1},"b":{"b":2}}`,
		expectedErr: createErrorMemberNotExist(`.b`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterProperty_object_root_missing_nested2")
}

func TestRetrieve_filterProperty_object_root_missing_nested3(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a)].b.c`,
		inputJSON:   `{"a":{"a":1},"b":{"a":1}}`,
		expectedErr: createErrorMemberNotExist(`.b`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterProperty_object_root_missing_nested3")
}

func TestRetrieve_filterPropertyChain_complex_missing_b_1(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a)].a.b.c`,
		inputJSON:   `{"a":{"a":1},"b":{"a":{"c":2}}}`,
		expectedErr: createErrorMemberNotExist(`.b`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterPropertyChain_complex_missing_b_1")
}

func TestRetrieve_filterPropertyChain_complex_missing_b_2(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a)].a.b.c`,
		inputJSON:   `{"a":{"a":1},"b":{"a":{"c":2}},"c":{"b":3}}`,
		expectedErr: createErrorMemberNotExist(`.b`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterPropertyChain_complex_missing_b_2")
}
