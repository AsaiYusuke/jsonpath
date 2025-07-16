package jsonpath

import (
	"fmt"
	"testing"
)

func TestDotNotationMemberNotExist(t *testing.T) {
	tests := []TestCase{
		// Member does not exist in object
		{
			jsonpath:    `$.d`,
			inputJSON:   `{"a":"b","c":{"d":"e"}}`,
			expectedErr: createErrorMemberNotExist(`.d`),
		},
	}

	for i, tc := range tests {
		runTestCase(t, tc, fmt.Sprintf("TestDotNotationMemberNotExist_case_%d", i))
	}
}

func TestDotNotationArrayTypeError(t *testing.T) {
	tests := []TestCase{
		// Dot notation on array should fail
		{
			jsonpath:    `$.2`,
			inputJSON:   `["a","b",{"2":1}]`,
			expectedErr: createErrorTypeUnmatched(`.2`, `object`, `[]interface {}`),
		},
		{
			jsonpath:    `$.-1`,
			inputJSON:   `["a","b",{"2":1}]`,
			expectedErr: createErrorTypeUnmatched(`.-1`, `object`, `[]interface {}`),
		},
		{
			jsonpath:    `$.a`,
			inputJSON:   `[1,2]`,
			expectedErr: createErrorTypeUnmatched(`.a`, `object`, `[]interface {}`),
		},
		{
			jsonpath:    `$.a`,
			inputJSON:   `[{"a":1}]`,
			expectedErr: createErrorTypeUnmatched(`.a`, `object`, `[]interface {}`),
		},
	}

	for i, tc := range tests {
		runTestCase(t, tc, fmt.Sprintf("TestDotNotationArrayTypeError_case_%d", i))
	}
}

func TestDotNotationNestedTypeError(t *testing.T) {
	tests := []TestCase{
		// Dot notation on non-object types
		{
			jsonpath:    `$.a.d`,
			inputJSON:   `{"a":"b","c":{"d":"e"}}`,
			expectedErr: createErrorTypeUnmatched(`.d`, `object`, `string`),
		},
		{
			jsonpath:    `$.a.d`,
			inputJSON:   `{"a":123}`,
			expectedErr: createErrorTypeUnmatched(`.d`, `object`, `float64`),
		},
		{
			jsonpath:    `$.a.d`,
			inputJSON:   `{"a":true}`,
			expectedErr: createErrorTypeUnmatched(`.d`, `object`, `bool`),
		},
		{
			jsonpath:    `$.a.d`,
			inputJSON:   `{"a":null}`,
			expectedErr: createErrorTypeUnmatched(`.d`, `object`, `null`),
		},
	}

	for i, tc := range tests {
		runTestCase(t, tc, fmt.Sprintf("TestDotNotationNestedTypeError_case_%d", i))
	}
}
