package tests

import (
	"testing"
)

func TestUnion_BasicOperations(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[0,3]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first"]`,
		},
		{
			jsonpath:    `$[3,3]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[3,3]`),
		},
	}

	runTestCases(t, "TestUnion_BasicOperations", tests)
}

func TestUnion_SliceOperations(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[:2,0]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second","first"]`,
		},
	}

	runTestCases(t, "TestUnion_SliceOperations", tests)
}
