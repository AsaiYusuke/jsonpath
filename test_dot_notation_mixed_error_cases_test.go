package jsonpath

import (
	"fmt"
	"testing"
)

func TestDotNotation_MemberNotExistErrors(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$.d`,
			inputJSON:   `{"a":"b","c":{"d":"e"}}`,
			expectedErr: createErrorMemberNotExist(`.d`),
		},
	}

	for i, tc := range tests {
		runTestCase(t, tc, fmt.Sprintf("TestDotNotation_MemberNotExistErrors_case_%d", i))
	}
}

func TestDotNotation_ArrayTypeErrors(t *testing.T) {
	tests := []TestCase{
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
		runTestCase(t, tc, fmt.Sprintf("TestDotNotation_ArrayTypeErrors_case_%d", i))
	}
}

func TestDotNotation_NestedTypeErrors(t *testing.T) {
	tests := []TestCase{
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
		runTestCase(t, tc, fmt.Sprintf("TestDotNotation_NestedTypeErrors_case_%d", i))
	}
}
