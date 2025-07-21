package tests

import (
	"fmt"
	"testing"
)

func TestBracketNotation_LargeIndexOutOfRange(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[1000000000000000000]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[1000000000000000000]`),
		},
	}

	runTestCases(t, "TestBracketNotation_LargeIndexOutOfRange", testCases)
}

func TestBracketNotation_NumericErrorCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:  `$[10000000000000000000]`,
			inputJSON: `["first","second","third"]`,
			expectedErr: createErrorInvalidArgument(
				`10000000000000000000`,
				fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
			),
		},
		{
			jsonpath:  `$[0,10000000000000000000]`,
			inputJSON: `["first","second","third"]`,
			expectedErr: createErrorInvalidArgument(
				`10000000000000000000`,
				fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
			),
		},
		{
			jsonpath:  `$[10000000000000000000:1]`,
			inputJSON: `["first","second","third"]`,
			expectedErr: createErrorInvalidArgument(
				`10000000000000000000`,
				fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
			),
		},
		{
			jsonpath:  `$[1:10000000000000000000]`,
			inputJSON: `["first","second","third"]`,
			expectedErr: createErrorInvalidArgument(
				`10000000000000000000`,
				fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
			),
		},
		{
			jsonpath:  `$[0:3:10000000000000000000]`,
			inputJSON: `["first","second","third"]`,
			expectedErr: createErrorInvalidArgument(
				`10000000000000000000`,
				fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
			),
		},
		{
			jsonpath:  `$[?(@.a==1e1abc)]`,
			inputJSON: `{}`,
			expectedErr: createErrorInvalidArgument(
				`1e1abc`,
				fmt.Errorf(`strconv.ParseFloat: parsing "1e1abc": invalid syntax`),
			),
		},
	}

	runTestCases(t, "TestBracketNotation_NumericErrorCases", testCases)
}
