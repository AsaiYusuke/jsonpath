package jsonpath

import (
	"fmt"
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

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("RootAccess_%d", i), test)
	}
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
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("SimplePropertyAccess_%d", i), test)
	}
}

func TestBasicAccess_ComplexPropertyAccess(t *testing.T) {
	tests := []TestCase{
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

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("ComplexPropertyAccess_%d", i), test)
	}
}
