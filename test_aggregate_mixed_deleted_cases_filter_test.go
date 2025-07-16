package jsonpath

import (
	"fmt"
	"testing"
)

// TestFilterFunctionDeletedCases tests deleted filter function cases
func TestFilterFunctionDeletedCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$.*.twice()`,
			inputJSON:    `[123.456,256]`,
			expectedJSON: `[246.912,512]`,
			filters: map[string]func(interface{}) (interface{}, error){
				`twice`: twiceFunc,
			},
		},
		{
			jsonpath:     `$.*.twice().twice()`,
			inputJSON:    `[123.456,256]`,
			expectedJSON: `[493.824,1024]`,
			filters: map[string]func(interface{}) (interface{}, error){
				`twice`: twiceFunc,
			},
		},
		{
			jsonpath:     `$.*.twice().quarter()`,
			inputJSON:    `[123.456,256]`,
			expectedJSON: `[61.728,128]`,
			filters: map[string]func(interface{}) (interface{}, error){
				`twice`:   twiceFunc,
				`quarter`: quarterFunc,
			},
		},
		{
			jsonpath:     `$.*.quarter().twice()`,
			inputJSON:    `[123.456,256]`,
			expectedJSON: `[61.728,128]`,
			filters: map[string]func(interface{}) (interface{}, error){
				`twice`:   twiceFunc,
				`quarter`: quarterFunc,
			},
		},
		{
			jsonpath:     `$[?(@.twice())]`,
			inputJSON:    `[123.456,256]`,
			expectedJSON: `[123.456,256]`,
			filters: map[string]func(interface{}) (interface{}, error){
				`twice`: twiceFunc,
			},
		},
		{
			jsonpath:     `$[?(@.twice() == 512)]`,
			inputJSON:    `[123.456,256]`,
			expectedJSON: `[256]`,
			filters: map[string]func(interface{}) (interface{}, error){
				`twice`: twiceFunc,
			},
		},
		{
			jsonpath:     `$[?(512 != @.twice())]`,
			inputJSON:    `[123.456,256]`,
			expectedJSON: `[123.456]`,
			filters: map[string]func(interface{}) (interface{}, error){
				`twice`: twiceFunc,
			},
		},
		{
			jsonpath:     `$[?(@.twice() == $[0].twice())]`,
			inputJSON:    `[123.456,256]`,
			expectedJSON: `[123.456]`,
			filters: map[string]func(interface{}) (interface{}, error){
				`twice`: twiceFunc,
			},
		},
		{
			jsonpath:     `$.*.unknown()`,
			inputJSON:    `[123.456,256]`,
			expectedJSON: ``,
			expectedErr:  ErrorFunctionNotFound{function: `.unknown()`},
		},
	}

	for i, testCase := range testCases {
		runSingleTestCase(t, fmt.Sprintf("FilterFunctionDeleted_%d", i), testCase)
	}
}

// TestAggregateFunctionDeletedCases tests deleted aggregate function cases
func TestAggregateFunctionDeletedCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$.*.max().max()`,
			inputJSON:    `[122.345,123.45,123.456]`,
			expectedJSON: `[123.456]`,
			aggregates: map[string]func([]interface{}) (interface{}, error){
				`max`: maxFunc,
			},
		},
		{
			jsonpath:     `$.*.max().min()`,
			inputJSON:    `[122.345,123.45,123.456]`,
			expectedJSON: `[123.456]`,
			aggregates: map[string]func([]interface{}) (interface{}, error){
				`max`: maxFunc,
				`min`: minFunc,
			},
		},
		{
			jsonpath:     `$.*.min().max()`,
			inputJSON:    `[122.345,123.45,123.456]`,
			expectedJSON: `[122.345]`,
			aggregates: map[string]func([]interface{}) (interface{}, error){
				`max`: maxFunc,
				`min`: minFunc,
			},
		},
	}

	for i, testCase := range testCases {
		runSingleTestCase(t, fmt.Sprintf("AggregateFunctionDeleted_%d", i), testCase)
	}
}
