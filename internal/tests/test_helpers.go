package tests

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/AsaiYusuke/jsonpath/v2"
	"github.com/AsaiYusuke/jsonpath/v2/config"
	"github.com/AsaiYusuke/jsonpath/v2/errors"
)

type TestGroup map[string][]TestCase

type TestCase struct {
	jsonpath        string
	inputJSON       string
	expectedJSON    string
	expectedErr     error
	unmarshalFunc   func(string, *any) error
	filters         map[string]func(any) (any, error)
	aggregates      map[string]func([]any) (any, error)
	accessorMode    bool
	resultValidator func(any, []any) error
}

func createErrorFunctionFailed(functionName string, errorString string) errors.ErrorFunctionFailed {
	return errors.NewErrorFunctionFailed(functionName, len(functionName), fmt.Errorf("%s", errorString))
}

func createErrorMemberNotExist(path string) errors.ErrorMemberNotExist {
	errBasicError := errors.NewErrorBasicRuntime(path, len(path))
	return errors.NewErrorMemberNotExist(&errBasicError)
}

func createErrorTypeUnmatched(path string, expected string, found string) errors.ErrorTypeUnmatched {
	errBasicError := errors.NewErrorBasicRuntime(path, len(path))
	return errors.NewErrorTypeUnmatched(&errBasicError, expected, found)
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

func execTestRetrieve(t *testing.T, inputJSON any, testCase TestCase, fileLine string) ([]any, error) {
	jsonPath := testCase.jsonpath
	hasConfig := false
	config := config.Config{}
	expectedError := testCase.expectedErr
	var actualObject []any
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
	var src any
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

// runTestCasesSerial runs the given test cases sequentially without t.Parallel.
func runTestCasesSerial(t *testing.T, testGroupName string, testCases []TestCase, init func()) {
	for i, testCase := range testCases {
		if _, file, line, ok := runtime.Caller(2); ok {
			fileLine := fmt.Sprintf(`%s:%d`, file, line)
			testName := fmt.Sprintf(`%s_case_%d_jsonpath_%s`, testGroupName, i+1, testCase.jsonpath)
			t.Run(testName, func(t *testing.T) {
				if init != nil {
					init()
				}
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

var twiceFunc = func(param any) (any, error) {
	if input, ok := param.(float64); ok {
		return input * 2, nil
	}
	return nil, fmt.Errorf(`type error`)
}

var quarterFunc = func(param any) (any, error) {
	if input, ok := param.(float64); ok {
		return input / 4, nil
	}
	return nil, fmt.Errorf(`type error`)
}

var errAggregateFunc = func(param []any) (any, error) {
	return nil, fmt.Errorf("aggregate error")
}

var errFilterFunc = func(param any) (any, error) {
	return nil, fmt.Errorf("filter error")
}

var maxFunc = func(param []any) (any, error) {
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

var minFunc = func(param []any) (any, error) {
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

var echoAggregateFunc = func(param []any) (any, error) {
	return param, nil
}

func createAccessorModeValidator(
	resultIndex int,
	expectedValue1, expectedValue2, expectedValue3 any,
	srcGetter func(any) any,
	srcSetter func(any, any)) func(any, []any) error {
	return func(src any, actualObject []any) error {
		accessor := actualObject[resultIndex].(config.Accessor)

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

func createGetOnlyValidator(expected any) func(any, []any) error {
	return func(src any, actualObject []any) error {
		accessor, ok := actualObject[0].(config.Accessor)
		if !ok {
			return fmt.Errorf("result[0] is not Accessor: %T", actualObject[0])
		}
		got := accessor.Get()
		if !reflect.DeepEqual(got, expected) {
			return fmt.Errorf("get-only: expected %#v, got %#v", expected, got)
		}
		if accessor.Set != nil {
			return fmt.Errorf("get-only: accessor.Set expected nil, got non-nil")
		}
		return nil
	}
}

var sliceStructChangedResultValidator = func(src any, actualObject []any) error {
	srcArray := src.([]any)
	accessor := actualObject[0].(config.Accessor)

	accessor.Set(4) // srcArray:[1,4,3] , accessor:[1,4,3]
	if len(srcArray) != 3 || srcArray[1] != 4 {
		return fmt.Errorf(`set -> src : expect<%d> != actual<%d>`, 4, srcArray[1])
	}
	srcArray = append(srcArray[:1], srcArray[2:]...) // srcArray:[1,3] , accessor:[1,3,3]
	if len(srcArray) != 2 || accessor.Get() != 3.0 { // Go's marshal returns float value
		return fmt.Errorf(`del -> get : expect<%f> != actual<%f>`, 3.0, accessor.Get())
	}
	accessor.Set(5) // srcArray:[1,5] , accessor:[1,5,3]
	if len(srcArray) != 2 || srcArray[1] != 5 {
		return fmt.Errorf(`del -> set -> src : expect<%d> != actual<%d>`, 5, srcArray[1])
	}
	srcArray = append(srcArray[:1], srcArray[2:]...) // srcArray:[1] , accessor:[1,5,3]
	if len(srcArray) != 1 || accessor.Get() != 5 {
		return fmt.Errorf(`del x2 -> get : expect<%d> != actual<%d>`, 5, accessor.Get())
	}
	accessor.Set(6) // srcArray:[1] , accessor:[1,6,3]
	if len(srcArray) != 1 {
		return fmt.Errorf(`del x2 -> set -> len : expect<%d> != actual<%d>`, 1, len(srcArray))
	}
	srcArray = append(srcArray, 7) // srcArray:[1,7] , accessor:[1,7,3]
	if len(srcArray) != 2 || accessor.Get() != 7 {
		return fmt.Errorf(`del x2 -> add -> get : expect<%d> != actual<%d>`, 7, accessor.Get())
	}
	srcArray = append(srcArray, 8) // srcArray:[1,7,8]    , accessor:[1,7,8]
	srcArray = append(srcArray, 9) // srcArray:[1,7,8,9]  , accessor:[1,7,8,9]
	srcArray[1] = 10               // srcArray:[1,10,8,9] , accessor:[1,10,8,9]
	if len(srcArray) != 4 || accessor.Get() != 10 {
		return fmt.Errorf(`del x2 -> add x3 -> update -> get : expect<%d> != actual<%d>`, 10, accessor.Get())
	}
	return nil
}

var mapStructChangedResultValidator = func(src any, actualObject []any) error {
	srcMap := src.(map[string]any)
	accessor := actualObject[0].(config.Accessor)

	accessor.Set(2) // srcMap:{"a":2} , accessor:{"a":2}
	if len(srcMap) != 1 || srcMap[`a`] != 2 {
		return fmt.Errorf(`set -> src : expect<%d> != actual<%d>`, 2, srcMap[`a`])
	}
	delete(srcMap, `a`) // srcMap:{} , accessor:{}
	if accessor.Get() != nil {
		return fmt.Errorf(`del -> get : expect<%v> != actual<%d>`, nil, accessor.Get())
	}
	accessor.Set(3) // srcMap:{"a":3} , accessor:{"a":3}
	if len(srcMap) != 1 || srcMap[`a`] != 3 {
		return fmt.Errorf(`del -> set -> len : expect<%d> != actual<%d>`, 0, len(srcMap))
	}
	delete(srcMap, `a`) // srcMap:{} , accessor:{}
	srcMap[`a`] = 4     // srcMap:{"a":4} , accessor:{"a":4}
	if accessor.Get() != 4 {
		return fmt.Errorf(`del -> update -> get : expect<%v> != actual<%d>`, 4, accessor.Get())
	}
	return nil
}

func maxAggregate(items []any) (any, error) {
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

func twiceFilter(item any) (any, error) {
	if val, ok := item.(float64); ok {
		return val * 2, nil
	}
	return nil, createErrorFunctionFailed("twice", "non-numeric value")
}

func errorFilter(item any) (any, error) {
	return nil, createErrorFunctionFailed("errFilter", "test error")
}
