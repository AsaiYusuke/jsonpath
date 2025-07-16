package jsonpath

import (
	"fmt"
	"testing"
)

// TestInvalidSyntax_numberingDeletedCases tests deleted invalid syntax cases for large numbers and numeric errors
func TestInvalidSyntax_numberingDeletedCases(t *testing.T) {
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
		runSingleTestCase(t, fmt.Sprintf("NumberingDeleted_%d", i), testCase)
	}
}

// TestInvalidSyntax_booleanValueDeletedCases tests deleted invalid syntax cases for boolean value variations
func TestInvalidSyntax_booleanValueDeletedCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a==falSe)]`,
			inputJSON:    `[{"a":false}]`,
			expectedJSON: `[{"a":false}]`,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==falSe)]`},
		},
		{
			jsonpath:     `$[?(@.a==FaLSE)]`,
			inputJSON:    `[{"a":false}]`,
			expectedJSON: `[{"a":false}]`,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==FaLSE)]`},
		},
		{
			jsonpath:     `$[?(@.a==trUe)]`,
			inputJSON:    `[{"a":true}]`,
			expectedJSON: `[{"a":true}]`,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==trUe)]`},
		},
		{
			jsonpath:     `$[?(@.a==NuLl)]`,
			inputJSON:    `[{"a":null}]`,
			expectedJSON: `[{"a":null}]`,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==NuLl)]`},
		},
	}

	for i, testCase := range testCases {
		runSingleTestCase(t, fmt.Sprintf("BooleanValueDeleted_%d", i), testCase)
	}
}

// TestInvalidSyntax_regexDeletedCases tests deleted invalid syntax cases for regex patterns
func TestInvalidSyntax_regexDeletedCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a>1 && ())]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a>1 && ())]`},
		},
		{
			jsonpath:     `$[?(@.a=~///)]`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=~///)]`},
		},
		{
			jsonpath:     `$[?(@.a=~s/a/b/)]`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=~s/a/b/)]`},
		},
		{
			jsonpath:     `$[?(@.a=~@abc@)]`,
			inputJSON:    `[]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=~@abc@)]`},
		},
	}

	for i, testCase := range testCases {
		runSingleTestCase(t, fmt.Sprintf("RegexDeleted_%d", i), testCase)
	}
}

// TestInvalidSyntax_functionDeletedCases tests deleted invalid syntax cases for function calls
func TestInvalidSyntax_functionDeletedCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$.func(`,
			inputJSON:    `{}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `(`},
		},
		{
			jsonpath:     `$.func()(`,
			inputJSON:    `{}`,
			expectedJSON: ``,
			expectedErr:  ErrorFunctionNotFound{function: `.func()`},
		},
		{
			jsonpath:     `$.func(){}`,
			inputJSON:    `{}`,
			expectedJSON: ``,
			expectedErr:  ErrorFunctionNotFound{function: `.func()`},
		},
	}

	for i, testCase := range testCases {
		runSingleTestCase(t, fmt.Sprintf("FunctionDeleted_%d", i), testCase)
	}
}
