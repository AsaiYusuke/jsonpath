package jsonpath

import (
	"fmt"
	"testing"
)

func TestSliceNegativeIndex_ComplexCases(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$[-1:-2]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[-1:-2]`),
		},
		{
			jsonpath:    `$[-1:-3]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[-1:-3]`),
		},
		{
			jsonpath:    `$[-1:2]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[-1:2]`),
		},
		{
			jsonpath:     `$[-1:3]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["third"]`,
		},
		{
			jsonpath:    `$[-2:1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[-2:1]`),
		},
		{
			jsonpath:     `$[-2:2]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["second"]`,
		},
		{
			jsonpath:    `$[-3:0]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[-3:0]`),
		},
		{
			jsonpath:     `$[-3:1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first"]`,
		},
		{
			jsonpath:    `$[-4:0]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[-4:0]`),
		},
		{
			jsonpath:     `$[-4:1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first"]`,
		},
		{
			jsonpath:     `$[-4:3]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second","third"]`,
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("ComplexCases_%d", i), test)
	}
}

func TestSliceNegativeIndex_StandardCases(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[0:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second"]`,
		},
		{
			jsonpath:     `$[0:-2]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first"]`,
		},
		{
			jsonpath:    `$[0:-3]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[0:-3]`),
		},
		{
			jsonpath:    `$[0:-4]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[0:-4]`),
		},
		{
			jsonpath:     `$[1:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["second"]`,
		},
		{
			jsonpath:    `$[1:-2]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[1:-2]`),
		},
		{
			jsonpath:    `$[2:-1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[2:-1]`),
		},
		{
			jsonpath:    `$[2:-2]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[2:-2]`),
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("StandardCases_%d", i), test)
	}
}

func TestSliceEmptyStart_Cases(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$[:0]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[:0]`),
		},
		{
			jsonpath:     `$[:1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first"]`,
		},
		{
			jsonpath:     `$[:2]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second"]`,
		},
		{
			jsonpath:     `$[:3]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second","third"]`,
		},
		{
			jsonpath:     `$[:4]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second","third"]`,
		},
		{
			jsonpath:     `$[:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second"]`,
		},
		{
			jsonpath:     `$[:-2]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first"]`,
		},
		{
			jsonpath:    `$[:-3]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[:-3]`),
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("EmptyStart_%d", i), test)
	}
}

func TestSliceEmptyEnd_Cases(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[0:]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second","third"]`,
		},
		{
			jsonpath:     `$[1:]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["second","third"]`,
		},
		{
			jsonpath:     `$[2:]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["third"]`,
		},
		{
			jsonpath:    `$[3:]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[3:]`),
		},
		{
			jsonpath:     `$[-1:]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["third"]`,
		},
		{
			jsonpath:     `$[-2:]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["second","third"]`,
		},
		{
			jsonpath:     `$[-3:]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second","third"]`,
		},
		{
			jsonpath:     `$[-4:]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second","third"]`,
		},
		{
			jsonpath:     `$[:]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second","third"]`,
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("EmptyEnd_%d", i), test)
	}
}

func TestSliceEdgeCases_LargeNumbers(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[-1000000000000000000:1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first"]`,
		},
		{
			jsonpath:    `$[1000000000000000000:1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[1000000000000000000:1]`),
		},
		{
			jsonpath:    `$[1:-1000000000000000000]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[1:-1000000000000000000]`),
		},
		{
			jsonpath:     `$[1:1000000000000000000]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["second","third"]`,
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("LargeNumbers_%d", i), test)
	}
}

func TestSliceTypeErrors_Cases(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$[1:2]`,
			inputJSON:   `{"first":1,"second":2,"third":3}`,
			expectedErr: createErrorTypeUnmatched(`[1:2]`, `array`, `map[string]interface {}`),
		},
		{
			jsonpath:    `$[:]`,
			inputJSON:   `{"first":1,"second":2,"third":3}`,
			expectedErr: createErrorTypeUnmatched(`[:]`, `array`, `map[string]interface {}`),
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("TypeErrors_%d", i), test)
	}
}

func TestSliceSyntaxEdgeCases_Cases(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[+0:+1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first"]`,
		},
		{
			jsonpath:     `$[01:02]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["second"]`,
		},
		{
			jsonpath:    `$[0.0:2]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[0.0:2]`},
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("SyntaxEdgeCases_%d", i), test)
	}
}
