package tests

import (
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

	runTestCases(t, "TestSliceBasic_SimpleSlice", tests)
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

	runTestCases(t, "TestSliceBasic_EmptySlice", tests)
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

	runTestCases(t, "TestSliceBasic_NegativeIndices", tests)
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

	runTestCases(t, "TestSliceBasic_ReverseSlice", tests)
}
