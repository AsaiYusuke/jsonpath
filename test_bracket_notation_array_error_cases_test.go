package jsonpath

import (
	"fmt"
	"testing"
)

func TestBracketNotation_NumericErrorCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:  `$[10000000000000000000]`,
			inputJSON: `["first","second","third"]`,
			expectedErr: ErrorInvalidArgument{
				argument: `10000000000000000000`,
				err:      fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
			},
		},
		{
			jsonpath:  `$[0,10000000000000000000]`,
			inputJSON: `["first","second","third"]`,
			expectedErr: ErrorInvalidArgument{
				argument: `10000000000000000000`,
				err:      fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
			},
		},
		{
			jsonpath:  `$[10000000000000000000:1]`,
			inputJSON: `["first","second","third"]`,
			expectedErr: ErrorInvalidArgument{
				argument: `10000000000000000000`,
				err:      fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
			},
		},
		{
			jsonpath:  `$[1:10000000000000000000]`,
			inputJSON: `["first","second","third"]`,
			expectedErr: ErrorInvalidArgument{
				argument: `10000000000000000000`,
				err:      fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
			},
		},
		{
			jsonpath:  `$[0:3:10000000000000000000]`,
			inputJSON: `["first","second","third"]`,
			expectedErr: ErrorInvalidArgument{
				argument: `10000000000000000000`,
				err:      fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
			},
		},
		{
			jsonpath:  `$[?(@.a==1e1abc)]`,
			inputJSON: `{}`,
			expectedErr: ErrorInvalidArgument{
				argument: `1e1abc`,
				err:      fmt.Errorf(`strconv.ParseFloat: parsing "1e1abc": invalid syntax`),
			},
		},
	}

	for i, testCase := range testCases {
		runSingleTestCase(t, fmt.Sprintf("NumericError_%d", i), testCase)
	}
}
