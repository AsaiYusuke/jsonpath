package jsonpath

import (
	"testing"
)

func TestFilter_PropertyAccessFilterThenPropertyChain(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a)].b.c`,
		inputJSON:   `[{"a":1},{"b":2}]`,
		expectedErr: createErrorMemberNotExist(`.b`),
	}
	runTestCase(t, testCase, "TestFilter_PropertyAccessFilterThenPropertyChain")
}
