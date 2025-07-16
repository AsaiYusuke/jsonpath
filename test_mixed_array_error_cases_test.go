package jsonpath

import (
	"testing"
)

func TestMixed_ArrayIndexDotNotation(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0].[1]`,
		inputJSON:   `[["a","b"],["c"],["d"]]`,
		expectedErr: ErrorInvalidSyntax{position: 4, reason: `unrecognized input`, near: `.[1]`},
	}
	runTestCase(t, testCase, "TestMixed_ArrayIndexDotNotation")
}

func TestMixed_ArrayIndexSlice(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0].[1:3]`,
		inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
		expectedErr: ErrorInvalidSyntax{position: 4, reason: `unrecognized input`, near: `.[1:3]`},
	}
	runTestCase(t, testCase, "TestMixed_ArrayIndexSlice")
}
