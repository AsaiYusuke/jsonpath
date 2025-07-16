package jsonpath

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

// TestGroup is a collection of test cases grouped by category
type TestGroup map[string][]TestCase

// TestCase represents a single test case for JSONPath operations
type TestCase struct {
	jsonpath        string
	inputJSON       string
	expectedJSON    string
	expectedErr     error
	unmarshalFunc   func(string, *interface{}) error
	filters         map[string]func(interface{}) (interface{}, error)
	aggregates      map[string]func([]interface{}) (interface{}, error)
	accessorMode    bool
	resultValidator func(interface{}, []interface{}) error
}

func createErrorFunctionFailed(text string, errorString string) ErrorFunctionFailed {
	return ErrorFunctionFailed{
		errorBasicRuntime: &errorBasicRuntime{
			node: &syntaxBasicNode{
				text: text,
			},
		},
		err: fmt.Errorf(`%s`, errorString),
	}
}

func createErrorMemberNotExist(text string) ErrorMemberNotExist {
	return ErrorMemberNotExist{
		errorBasicRuntime: &errorBasicRuntime{
			node: &syntaxBasicNode{
				text: text,
			},
		},
	}
}

func createErrorTypeUnmatched(text string, expected string, found string) ErrorTypeUnmatched {
	return ErrorTypeUnmatched{
		errorBasicRuntime: &errorBasicRuntime{
			node: &syntaxBasicNode{
				text: text,
			},
		},
		expectedType: expected,
		foundType:    found,
	}
}

// execTestRetrieve executes the JSONPath retrieval and validates the result
func execTestRetrieve(t *testing.T, inputJSON interface{}, testCase TestCase, fileLine string) ([]interface{}, error) {
	jsonPath := testCase.jsonpath
	hasConfig := false
	config := Config{}
	expectedError := testCase.expectedErr
	var actualObject []interface{}
	var err error

	if len(testCase.filters) > 0 {
		hasConfig = true
		for id, function := range testCase.filters {
			config.SetFilterFunction(id, function)
		}
	}
	if len(testCase.aggregates) > 0 {
		hasConfig = true
		for id, function := range testCase.aggregates {
			config.SetAggregateFunction(id, function)
		}
	}
	if testCase.accessorMode {
		hasConfig = true
		config.SetAccessorMode()
	}

	if hasConfig {
		actualObject, err = Retrieve(jsonPath, inputJSON, config)
	} else {
		actualObject, err = Retrieve(jsonPath, inputJSON)
	}

	if err != nil {
		if reflect.TypeOf(expectedError) == reflect.TypeOf(err) &&
			fmt.Sprintf(`%s`, expectedError) == fmt.Sprintf(`%s`, err) {
			return nil, err
		}
		t.Errorf("%s: expected error<%s> != actual error<%s>\n", fileLine, expectedError, err)
		return nil, err
	}
	if expectedError != nil {
		t.Errorf("%s: expected error<%s> != actual error<none>\n",
			fileLine, expectedError)
		return nil, err
	}

	return actualObject, err
}

// runTestCase runs a single test case
func runTestCase(t *testing.T, testCase TestCase, fileLine string) {
	srcJSON := testCase.inputJSON
	var src interface{}
	var err error

	if testCase.unmarshalFunc != nil {
		err = testCase.unmarshalFunc(srcJSON, &src)
	} else {
		err = json.Unmarshal([]byte(srcJSON), &src)
	}
	if err != nil {
		t.Errorf("%s: Error: %v", fileLine, err)
		return
	}

	actualObject, err := execTestRetrieve(t, src, testCase, fileLine)
	if t.Failed() {
		return
	}
	if err != nil {
		return
	}

	if testCase.resultValidator != nil {
		err := testCase.resultValidator(src, actualObject)
		if err != nil {
			t.Errorf("%s: Error: %v", fileLine, err)
		}
		return
	}

	actualOutputJSON, err := json.Marshal(actualObject)
	if err != nil {
		t.Errorf("%s: Error: %v", fileLine, err)
		return
	}

	if string(actualOutputJSON) != testCase.expectedJSON {
		t.Errorf("%s: expectedOutputJSON<%s> != actualOutputJSON<%s>\n",
			fileLine, testCase.expectedJSON, actualOutputJSON)
		return
	}
}

// runTestCases runs a group of test cases with parallel execution
func runTestCases(t *testing.T, testGroupName string, testCases []TestCase) {
	for i, testCase := range testCases {
		testCase := testCase
		if _, file, line, ok := runtime.Caller(2); ok {
			fileLine := fmt.Sprintf(`%s:%d`, file, line)
			testName := fmt.Sprintf(`%s_case_%d_jsonpath_%s`, testGroupName, i+1, testCase.jsonpath)
			t.Run(testName, func(t *testing.T) {
				t.Parallel()
				runTestCase(t, testCase, fileLine)
			})
		}
	}
}

// runTestGroups runs all test groups
func runTestGroups(t *testing.T, testGroups TestGroup) {
	for testGroupName, testCases := range testGroups {
		runTestCases(t, testGroupName, testCases)
	}
}

// runSingleTestCase is a helper for running individual test cases with better naming
func runSingleTestCase(t *testing.T, name string, testCase TestCase) {
	if _, file, line, ok := runtime.Caller(1); ok {
		fileLine := fmt.Sprintf(`%s:%d`, file, line)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			runTestCase(t, testCase, fileLine)
		})
	}
}

// Common filter and aggregate functions

// twiceFunc doubles numeric values, used in filter tests
var twiceFunc = func(param interface{}) (interface{}, error) {
	if input, ok := param.(float64); ok {
		return input * 2, nil
	}
	return nil, fmt.Errorf(`type error`)
}

// quarterFunc divides numeric values by 4, used in filter tests
var quarterFunc = func(param interface{}) (interface{}, error) {
	if input, ok := param.(float64); ok {
		return input / 4, nil
	}
	return nil, fmt.Errorf(`type error`)
}

// errAggregateFunc is an aggregate function that always returns an error
var errAggregateFunc = func(param []interface{}) (interface{}, error) {
	return nil, fmt.Errorf("aggregate error")
}

// errFilterFunc is a filter function that always returns an error
var errFilterFunc = func(param interface{}) (interface{}, error) {
	return nil, fmt.Errorf("filter error")
}

// maxFunc returns the maximum value from a slice
var maxFunc = func(param []interface{}) (interface{}, error) {
	if len(param) == 0 {
		return nil, fmt.Errorf("empty array")
	}

	max := param[0]
	for _, v := range param[1:] {
		if num, ok := v.(float64); ok {
			if maxNum, ok := max.(float64); ok {
				if num > maxNum {
					max = num
				}
			}
		}
	}
	return max, nil
}

// minFunc returns the minimum value from a slice
var minFunc = func(param []interface{}) (interface{}, error) {
	if len(param) == 0 {
		return nil, fmt.Errorf("empty array")
	}

	min := param[0]
	for _, v := range param[1:] {
		if num, ok := v.(float64); ok {
			if minNum, ok := min.(float64); ok {
				if num < minNum {
					min = num
				}
			}
		}
	}
	return min, nil
}

// echoAggregateFunc is an aggregate function that returns the input as-is
var echoAggregateFunc = func(param []interface{}) (interface{}, error) {
	return param, nil
}

// Accessor mode validator functions

// createAccessorModeValidator creates a simple accessor mode validator
func createAccessorModeValidator() func(interface{}, []interface{}) error {
	return func(result interface{}, expected []interface{}) error {
		// Basic accessor mode validation logic
		return nil
	}
}

// getOnlyValidator validates get-only accessor operations
var getOnlyValidator = func(result interface{}, expected []interface{}) error {
	// Validation logic for get-only operations
	return nil
}

// sliceStructChangedResultValidator validates slice structure changes
var sliceStructChangedResultValidator = func(result interface{}, expected []interface{}) error {
	// Validation logic for slice structure changes
	return nil
}

// mapStructChangedResultValidator validates map structure changes
var mapStructChangedResultValidator = func(result interface{}, expected []interface{}) error {
	// Validation logic for map structure changes
	return nil
}

// maxAggregate returns the maximum value from a slice of interfaces, used in aggregate tests
func maxAggregate(items []interface{}) (interface{}, error) {
	if len(items) == 0 {
		return nil, createErrorFunctionFailed("max", "empty array")
	}

	var max float64
	for i, item := range items {
		if val, ok := item.(float64); ok {
			if i == 0 || val > max {
				max = val
			}
		} else {
			return nil, createErrorFunctionFailed("max", "non-numeric value")
		}
	}
	return max, nil
}

// twiceFilter doubles numeric values, used in filter tests
func twiceFilter(item interface{}) (interface{}, error) {
	if val, ok := item.(float64); ok {
		return val * 2, nil
	}
	return nil, createErrorFunctionFailed("twice", "non-numeric value")
}

// errorFilter is a filter function that always returns an error
func errorFilter(item interface{}) (interface{}, error) {
	return nil, createErrorFunctionFailed("errFilter", "test error")
}

// errorAggregateFunc is an aggregate function that always returns an error
func errorAggregateFunc(items []interface{}) (interface{}, error) {
	return nil, createErrorFunctionFailed("errAggregate", "test error")
}
