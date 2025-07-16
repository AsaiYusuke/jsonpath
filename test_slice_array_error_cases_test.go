package jsonpath

import (
	"testing"
)

// TestRetrieve_sliceFloatInvalid tests invalid slice syntax with float values
func TestRetrieve_sliceFloatInvalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0:2.0]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: ``,
		expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[0:2.0]`},
	}
	runTestCase(t, testCase, "TestRetrieve_sliceFloatInvalid")
}
