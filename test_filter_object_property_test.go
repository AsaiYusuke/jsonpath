package jsonpath

import (
	"testing"
)

func TestRetrieve_filterPropertyAccess_filter_then_property_chain(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a)].b.c`,
		inputJSON:   `[{"a":1},{"b":2}]`,
		expectedErr: createErrorMemberNotExist(`.b`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterPropertyAccess_filter_then_property_chain")
}
