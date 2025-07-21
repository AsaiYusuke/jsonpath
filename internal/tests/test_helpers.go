package tests

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/AsaiYusuke/jsonpath"
	"github.com/AsaiYusuke/jsonpath/errors"
)

type TestGroup map[string][]TestCase

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

func createErrorFunctionFailed(functionName string, errorString string) errors.ErrorFunctionFailed {
	return errors.NewErrorFunctionFailed(functionName, errorString)
}

func createErrorMemberNotExist(path string) errors.ErrorMemberNotExist {
	return errors.NewErrorMemberNotExist(path)
}

func createErrorTypeUnmatched(path string, expected string, found string) errors.ErrorTypeUnmatched {
	return errors.NewErrorTypeUnmatched(path, expected, found)
}

func createErrorInvalidSyntax(position int, reason string, near string) errors.ErrorInvalidSyntax {
	return errors.NewErrorInvalidSyntax(position, reason, near)
}

func createErrorNotSupported(feature string, path string) errors.ErrorNotSupported {
	return errors.NewErrorNotSupported(feature, path)
}

func createErrorInvalidArgument(argument string, err error) errors.ErrorInvalidArgument {
	return errors.NewErrorInvalidArgument(argument, err)
}

func createErrorFunctionNotFound(function string) errors.ErrorFunctionNotFound {
	return errors.NewErrorFunctionNotFound(function)
}

func execTestRetrieve(t *testing.T, inputJSON interface{}, testCase TestCase, fileLine string) ([]interface{}, error) {
	jsonPath := testCase.jsonpath
	hasConfig := false
	config := jsonpath.Config{}
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
		actualObject, err = jsonpath.Retrieve(jsonPath, inputJSON, config)
	} else {
		actualObject, err = jsonpath.Retrieve(jsonPath, inputJSON)
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

func runTestCases(t *testing.T, testGroupName string, testCases []TestCase) {
	for i, testCase := range testCases {
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

func runTestGroups(t *testing.T, testGroups TestGroup) {
	for testGroupName, testCases := range testGroups {
		runTestCases(t, testGroupName, testCases)
	}
}

func runSingleTestCase(t *testing.T, name string, testCase TestCase) {
	if _, file, line, ok := runtime.Caller(1); ok {
		fileLine := fmt.Sprintf(`%s:%d`, file, line)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			runTestCase(t, testCase, fileLine)
		})
	}
}

var twiceFunc = func(param interface{}) (interface{}, error) {
	if input, ok := param.(float64); ok {
		return input * 2, nil
	}
	return nil, fmt.Errorf(`type error`)
}

var quarterFunc = func(param interface{}) (interface{}, error) {
	if input, ok := param.(float64); ok {
		return input / 4, nil
	}
	return nil, fmt.Errorf(`type error`)
}

var errAggregateFunc = func(param []interface{}) (interface{}, error) {
	return nil, fmt.Errorf("aggregate error")
}

var errFilterFunc = func(param interface{}) (interface{}, error) {
	return nil, fmt.Errorf("filter error")
}

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

var echoAggregateFunc = func(param []interface{}) (interface{}, error) {
	return param, nil
}

func createAccessorModeValidator(
	resultIndex int,
	expectedValue1, expectedValue2, expectedValue3 interface{},
	srcGetter func(interface{}) interface{},
	srcSetter func(interface{}, interface{})) func(interface{}, []interface{}) error {
	return func(src interface{}, actualObject []interface{}) error {
		accessor := actualObject[resultIndex].(jsonpath.Accessor)

		getValue := accessor.Get()
		if getValue != expectedValue1 {
			return fmt.Errorf(`get : expect<%f> != actual<%f>`, expectedValue1, getValue)
		}

		accessor.Set(expectedValue2)

		newSrcValue := srcGetter(src)

		if newSrcValue != expectedValue2 {
			return fmt.Errorf(`set : expect<%f> != actual<%f>`, expectedValue2, newSrcValue)
		}

		getValue = accessor.Get()
		if getValue != expectedValue2 {
			return fmt.Errorf(`set -> get : expect<%f> != actual<%f>`, expectedValue2, getValue)
		}

		srcSetter(src, expectedValue3)

		getValue = accessor.Get()
		if getValue != expectedValue3 {
			return fmt.Errorf(`src -> get : expect<%f> != actual<%f>`, expectedValue3, getValue)
		}

		return nil
	}
}

var getOnlyValidator = func(result interface{}, expected []interface{}) error {
	return nil
}

var sliceStructChangedResultValidator = func(result interface{}, expected []interface{}) error {
	return nil
}

var mapStructChangedResultValidator = func(result interface{}, expected []interface{}) error {
	return nil
}

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

func twiceFilter(item interface{}) (interface{}, error) {
	if val, ok := item.(float64); ok {
		return val * 2, nil
	}
	return nil, createErrorFunctionFailed("twice", "non-numeric value")
}

func errorFilter(item interface{}) (interface{}, error) {
	return nil, createErrorFunctionFailed("errFilter", "test error")
}
