package jsonpath

import (
	"fmt"
	"testing"
)

func TestSliceBasic_SimpleSlice(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[0:1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first"]`,
		},
		{
			jsonpath:     `$[0:2]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second"]`,
		},
		{
			jsonpath:     `$[0:3]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second","third"]`,
		},
		{
			jsonpath:     `$[1:2]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["second"]`,
		},
		{
			jsonpath:     `$[1:3]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["second","third"]`,
		},
		{
			jsonpath:     `$[2:3]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["third"]`,
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("SimpleSlice_%d", i), test)
	}
}

func TestSliceBasic_EmptySlice(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$[0:0]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[0:0]`),
		},
		{
			jsonpath:    `$[1:1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[1:1]`),
		},
		{
			jsonpath:    `$[2:2]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[2:2]`),
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("EmptySlice_%d", i), test)
	}
}

func TestSliceBasic_NegativeIndices(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[-2:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["second"]`,
		},
		{
			jsonpath:     `$[-3:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second"]`,
		},
		{
			jsonpath:    `$[-1:-1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[-1:-1]`),
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("NegativeIndices_%d", i), test)
	}
}

func TestSliceBasic_ReverseSlice(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$[2:1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[2:1]`),
		},
		{
			jsonpath:    `$[2:0]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[2:0]`),
		},
		{
			jsonpath:    `$[3:2]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[3:2]`),
		},
		{
			jsonpath:    `$[3:3]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[3:3]`),
		},
		{
			jsonpath:    `$[3:4]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[3:4]`),
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("ReverseSlice_%d", i), test)
	}
}
