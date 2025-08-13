package tests

import (
	"testing"
)

func TestMixed_ArrayIndexDotNotation(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0].[1]`,
		inputJSON:   `[["a","b"],["c"],["d"]]`,
		expectedErr: createErrorInvalidSyntax(4, `unrecognized input`, `.[1]`),
	}
	runTestCase(t, testCase, "TestMixed_ArrayIndexDotNotation")
}

func TestMixed_ArrayIndexSlice(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0].[1:3]`,
		inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
		expectedErr: createErrorInvalidSyntax(4, `unrecognized input`, `.[1:3]`),
	}
	runTestCase(t, testCase, "TestMixed_ArrayIndexSlice")
}
