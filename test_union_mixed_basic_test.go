package jsonpath

import (
	"fmt"
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

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("TestUnion_BasicOperations_%d", i), test)
	}
}

func TestUnion_SliceOperations(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[:2,0]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second","first"]`,
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("TestUnion_SliceOperations_%d", i), test)
	}
}
