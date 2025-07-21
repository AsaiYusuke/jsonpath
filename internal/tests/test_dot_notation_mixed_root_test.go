package tests

import (
	"testing"
)

func TestBasicAccess_RootAccess(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$`,
			inputJSON:    `"a"`,
			expectedJSON: `["a"]`,
		},
		{
			jsonpath:     `$`,
			inputJSON:    `2`,
			expectedJSON: `[2]`,
		},
		{
			jsonpath:     `$`,
			inputJSON:    `false`,
			expectedJSON: `[false]`,
		},
		{
			jsonpath:     `$`,
			inputJSON:    `true`,
			expectedJSON: `[true]`,
		},
		{
			jsonpath:     `$`,
			inputJSON:    `null`,
			expectedJSON: `[null]`,
		},
		{
			jsonpath:     `$`,
			inputJSON:    `{}`,
			expectedJSON: `[{}]`,
		},
		{
			jsonpath:     `$`,
			inputJSON:    `[]`,
			expectedJSON: `[[]]`,
		},
		{
			jsonpath:     `$`,
			inputJSON:    `[1]`,
			expectedJSON: `[[1]]`,
		},
	}

	runTestCases(t, "TestBasicAccess_RootAccess", tests)
}

func TestBasicAccess_SimplePropertyAccess(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$.a`,
			inputJSON:    `{"a":"string"}`,
			expectedJSON: `["string"]`,
		},
		{
			jsonpath:     `$.a`,
			inputJSON:    `{"a":123}`,
			expectedJSON: `[123]`,
		},
		{
			jsonpath:     `$.a`,
			inputJSON:    `{"a":-123.456}`,
			expectedJSON: `[-123.456]`,
		},
		{
			jsonpath:     `$.a`,
			inputJSON:    `{"a":true}`,
			expectedJSON: `[true]`,
		},
		{
			jsonpath:     `$.a`,
			inputJSON:    `{"a":false}`,
			expectedJSON: `[false]`,
		},
		{
			jsonpath:     `$.a`,
			inputJSON:    `{"a":null}`,
			expectedJSON: `[null]`,
		},
		{
			jsonpath:     `$.a`,
			inputJSON:    `{"a":{"b":"c"}}`,
			expectedJSON: `[{"b":"c"}]`,
		},
		{
			jsonpath:     `$.a`,
			inputJSON:    `{"a":[1,3,5]}`,
			expectedJSON: `[[1,3,5]]`,
		},
		{
			jsonpath:     `$.a`,
			inputJSON:    `{"a":{}}`,
			expectedJSON: `[{}]`,
		},
		{
			jsonpath:     `$.a`,
			inputJSON:    `{"a":[]}`,
			expectedJSON: `[[]]`,
		},
	}

	runTestCases(t, "TestBasicAccess_SimplePropertyAccess", tests)
}
