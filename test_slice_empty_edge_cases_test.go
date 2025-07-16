package jsonpath

import (
	"testing"
)

func TestSliceEmptyArrayOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[0:2]..a`,
			inputJSON:   `[]`,
			expectedErr: createErrorMemberNotExist(`[0:2]`),
		},
		{
			jsonpath:    `$[0:2].*`,
			inputJSON:   `[]`,
			expectedErr: createErrorMemberNotExist(`[0:2]`),
		},
		{
			jsonpath:    `$[0:2][0:2]`,
			inputJSON:   `[]`,
			expectedErr: createErrorMemberNotExist(`[0:2]`),
		},
		{
			jsonpath:    `$[0:2][*]`,
			inputJSON:   `[]`,
			expectedErr: createErrorMemberNotExist(`[0:2]`),
		},
		{
			jsonpath:    `$[0:2][?(@.b)]`,
			inputJSON:   `[]`,
			expectedErr: createErrorMemberNotExist(`[0:2]`),
		},
		{
			jsonpath:    `$[*]..a`,
			inputJSON:   `[]`,
			expectedErr: createErrorMemberNotExist(`[*]`),
		},
	}

	for _, testCase := range testCases {
		runSingleTestCase(t, "TestSliceEmptyArrayOperations", testCase)
	}
}
