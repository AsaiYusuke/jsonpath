package jsonpath

import (
	"fmt"
	"testing"
)

func TestValueGroupCombinationMultipleIdentifier_Recursive(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$['a','b']..a`,
			inputJSON:    `{"a":{"a":1,"c":2},"b":{"a":{"d":3,"a":4}}}`,
			expectedJSON: `[1,{"a":4,"d":3},4]`,
		},
		{
			jsonpath:    `$['a','b']..a`,
			inputJSON:   `{"a":{"x":1,"c":2},"b":{"x":{"d":3,"x":4}}}`,
			expectedErr: createErrorMemberNotExist(`a`),
		},
		{
			jsonpath:    `$['a','b']..a`,
			inputJSON:   `{"a":"a","b":"a"}`,
			expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
		},
		{
			jsonpath:    `$['a','b']..a`,
			inputJSON:   `{"x":{"x":1,"c":2},"y":{"x":{"d":3,"x":4}}}`,
			expectedErr: createErrorMemberNotExist(`['a','b']`),
		},
		{
			jsonpath:    `$['a','b']..a`,
			inputJSON:   `"x"`,
			expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `string`),
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("Recursive_%d", i), test)
	}
}

func TestValueGroupCombinationMultipleIdentifier_MultipleIdentifier(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$['a','b']['c','d']`,
			inputJSON:    `{"a":{"a":1,"c":2},"b":{"d":3,"a":4}}`,
			expectedJSON: `[2,3]`,
		},
		{
			jsonpath:    `$['a','b']['c','d']`,
			inputJSON:   `{"a":{"a":1,"x":2},"b":{"x":3,"a":4}}`,
			expectedErr: createErrorMemberNotExist(`['c','d']`),
		},
		{
			jsonpath:    `$['a','b']['c','d']`,
			inputJSON:   `{"x":{"a":1,"c":2},"x":{"d":3,"a":4}}`,
			expectedErr: createErrorMemberNotExist(`['a','b']`),
		},
		{
			jsonpath:    `$['a','b']['c','d']`,
			inputJSON:   `"x"`,
			expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `string`),
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("MultipleIdentifier_%d", i), test)
	}
}

func TestValueGroupCombinationMultipleIdentifier_WildcardIdentifier(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$['a','b'].*`,
			inputJSON:    `{"a":{"a":1,"c":2},"b":{"d":3,"a":4}}`,
			expectedJSON: `[1,2,4,3]`,
		},
		{
			jsonpath:    `$['a','b'].*`,
			inputJSON:   `{"a":{},"b":{}}`,
			expectedErr: createErrorMemberNotExist(`.*`),
		},
		{
			jsonpath:    `$['a','b'].*`,
			inputJSON:   `{"x":[1,3,2],"y":[4,6,5]}`,
			expectedErr: createErrorMemberNotExist(`['a','b']`),
		},
		{
			jsonpath:    `$['a','b'].*`,
			inputJSON:   `"x"`,
			expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `string`),
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("WildcardIdentifier_%d", i), test)
	}
}

func TestValueGroupCombinationMultipleIdentifier_SliceQualifier(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$['a','b'][0:2]`,
			inputJSON:    `{"a":[1,3,2],"b":[4,6,5]}`,
			expectedJSON: `[1,3,4,6]`,
		},
		{
			jsonpath:    `$['a','b'][0:2]`,
			inputJSON:   `{"a":[],"b":[]}`,
			expectedErr: createErrorMemberNotExist(`[0:2]`),
		},
		{
			jsonpath:    `$['a','b'][0:2]`,
			inputJSON:   `{"x":[1,3,2],"y":[4,6,5]}`,
			expectedErr: createErrorMemberNotExist(`['a','b']`),
		},
		{
			jsonpath:    `$['a','b'][0:2]`,
			inputJSON:   `"x"`,
			expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `string`),
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("SliceQualifier_%d", i), test)
	}
}

func TestValueGroupCombinationMultipleIdentifier_WildcardQualifier(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$['a','b'][*]`,
			inputJSON:    `{"a":[1,3,2],"b":[4,6,5]}`,
			expectedJSON: `[1,3,2,4,6,5]`,
		},
		{
			jsonpath:    `$['a','b'][*]`,
			inputJSON:   `{"a":[],"b":[]}`,
			expectedErr: createErrorMemberNotExist(`[*]`),
		},
		{
			jsonpath:    `$['a','b'][*]`,
			inputJSON:   `{"x":[1,3,2],"y":[4,6,5]}`,
			expectedErr: createErrorMemberNotExist(`['a','b']`),
		},
		{
			jsonpath:    `$['a','b'][*]`,
			inputJSON:   `"x"`,
			expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `string`),
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("WildcardQualifier_%d", i), test)
	}
}

func TestValueGroupCombinationMultipleIdentifier_UnionInQualifier(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$['a','b'][0,1]`,
			inputJSON:    `{"a":[1,3,2],"b":[4,6,5]}`,
			expectedJSON: `[1,3,4,6]`,
		},
		{
			jsonpath:    `$['a','b'][0,1]`,
			inputJSON:   `{"a":[],"b":[]}`,
			expectedErr: createErrorMemberNotExist(`[0,1]`),
		},
		{
			jsonpath:    `$['a','b'][0,1]`,
			inputJSON:   `{"x":[1,3,2],"y":[4,6,5]}`,
			expectedErr: createErrorMemberNotExist(`['a','b']`),
		},
		{
			jsonpath:    `$['a','b'][0,1]`,
			inputJSON:   `"x"`,
			expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `string`),
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("UnionInQualifier_%d", i), test)
	}
}

func TestValueGroupCombinationMultipleIdentifier_FilterQualifier(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$['a','b'][?(@.b)]`,
			inputJSON:    `{"a":[{"a":1},{"b":2}],"b":[{"a":3},{"b":4}]}`,
			expectedJSON: `[{"b":2},{"b":4}]`,
		},
		{
			jsonpath:    `$['a','b'][?(@.b)]`,
			inputJSON:   `{"a":[{"a":1},{"x":2}],"b":[{"a":3},{"x":4}]}`,
			expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
		},
		{
			jsonpath:    `$['a','b'][?(@.b)]`,
			inputJSON:   `{"x":[{"a":1},{"b":2}],"y":[{"a":3},{"b":4}]}`,
			expectedErr: createErrorMemberNotExist(`['a','b']`),
		},
		{
			jsonpath:    `$['a','b'][?(@.b)]`,
			inputJSON:   `"x"`,
			expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `string`),
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("FilterQualifier_%d", i), test)
	}
}
