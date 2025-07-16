package jsonpath

import (
	"testing"
)

func TestRetrieve_mixedNotation_array_index_dot_notation(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0].[1]`,
		inputJSON:   `[["a","b"],["c"],["d"]]`,
		expectedErr: ErrorInvalidSyntax{position: 4, reason: `unrecognized input`, near: `.[1]`},
	}
	runTestCase(t, testCase, "TestRetrieve_mixedNotation_array_index_dot_notation")
}

func TestRetrieve_mixedNotation_array_index_slice(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0].[1:3]`,
		inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
		expectedErr: ErrorInvalidSyntax{position: 4, reason: `unrecognized input`, near: `.[1:3]`},
	}
	runTestCase(t, testCase, "TestRetrieve_mixedNotation_array_index_slice")
}
