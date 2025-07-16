package jsonpath

import (
	"testing"
)

func TestSlice_FloatInvalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0:2.0]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: ``,
		expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[0:2.0]`},
	}
	runTestCase(t, testCase, "TestSlice_FloatInvalid")
}
