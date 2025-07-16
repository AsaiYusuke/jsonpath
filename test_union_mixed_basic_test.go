package jsonpath

import (
	"fmt"
	"testing"
)

func TestUnionBasic_SimpleUnion(t *testing.T) {
	tests := []TestCase{
		// Basic union tests
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

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("SimpleUnion_%d", i), test)
	}
}

func TestUnionBasic_SliceUnion(t *testing.T) {
	tests := []TestCase{
		// Slice with union tests (including the first deleted test case)
		{
			jsonpath:     `$[:2,0]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second","first"]`,
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("SliceUnion_%d", i), test)
	}
}
