package tests

import (
	"testing"
)

func TestFilterFunction_ChainedOperations(t *testing.T) {
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
	}

	runTestCases(t, "TestFilterFunction_ChainedOperations", testCases)
}

func TestFilterFunction_FilterOperations(t *testing.T) {
	testCases := []TestCase{
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
	}

	runTestCases(t, "TestFilterFunction_FilterOperations", testCases)
}

func TestFilterFunction_ErrorCases(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*.unknown()`,
		inputJSON:    `[123.456,256]`,
		expectedJSON: ``,
		expectedErr:  createErrorFunctionNotFound(`.unknown()`),
	}
	runSingleTestCase(t, "TestFilterFunction_ErrorCases_UnknownFunction", testCase)
}
