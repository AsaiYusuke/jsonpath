package tests

import (
	"testing"
)

func TestSlice_FloatInvalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0:2.0]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: ``,
		expectedErr:  createErrorInvalidSyntax(1, `unrecognized input`, `[0:2.0]`),
	}
	runTestCase(t, testCase, "TestSlice_FloatInvalid")
}
