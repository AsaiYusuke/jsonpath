package tests

import (
	"testing"
)

func TestValueGroupCombinationRecursiveDescent_Recursive(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$..a..b`,
			inputJSON:    `[{"a":{"a":{"b":1},"c":2}},{"b":{"a":{"d":3,"b":4}}}]`,
			expectedJSON: `[1,1,4]`,
		},
		{
			jsonpath:    `$..a..b`,
			inputJSON:   `[{"a":{"a":{"x":1},"c":2}},{"b":{"a":{"d":3,"x":4}}}]`,
			expectedErr: createErrorMemberNotExist(`b`),
		},
		{
			jsonpath:    `$..a..b`,
			inputJSON:   `[{"a":"b"},{"b":{"a":"b"}}]`,
			expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
		},
		{
			jsonpath:    `$..a..b`,
			inputJSON:   `[{"x":{"x":{"b":1},"c":2}},{"b":{"x":{"d":3,"b":4}}}]`,
			expectedErr: createErrorMemberNotExist(`a`),
		},
		{
			jsonpath:    `$..a..b`,
			inputJSON:   `"x"`,
			expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
		},
	}

	runTestCases(t, "TestValueGroupCombinationRecursiveDescent_Recursive", tests)
}

func TestValueGroupCombinationRecursiveDescent_MultipleIdentifier(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$..['a','b']`,
			inputJSON:    `[{"a":1,"c":2},{"d":3,"b":4}]`,
			expectedJSON: `[1,4]`,
		},
		{
			jsonpath:    `$..['a','b']`,
			inputJSON:   `[{"x":1,"c":2},{"d":3,"x":4}]`,
			expectedErr: createErrorMemberNotExist(`['a','b']`),
		},
		{
			jsonpath:    `$..['a','b']`,
			inputJSON:   `{}`,
			expectedErr: createErrorMemberNotExist(`['a','b']`),
		},
		{
			jsonpath:    `$..['a','b']`,
			inputJSON:   `[]`,
			expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `[]interface {}`),
		},
		{
			jsonpath:    `$..['a','b']`,
			inputJSON:   `"x"`,
			expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
		},
	}

	runTestCases(t, "TestValueGroupCombinationRecursiveDescent_MultipleIdentifier", tests)
}

func TestValueGroupCombinationRecursiveDescent_WildcardIdentifier(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$..*`,
			inputJSON:    `[{"a":1,"c":2},{"d":3,"b":4}]`,
			expectedJSON: `[{"a":1,"c":2},{"b":4,"d":3},1,2,4,3]`,
		},
		{
			jsonpath:    `$..*`,
			inputJSON:   `"x"`,
			expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
		},
	}

	runTestCases(t, "TestValueGroupCombinationRecursiveDescent_WildcardIdentifier", tests)
}

func TestValueGroupCombinationRecursiveDescent_SliceQualifier(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$..[0:2]`,
			inputJSON:    `{"a":[1,3,2],"b":{"a":[4,6,5]}}`,
			expectedJSON: `[1,3,4,6]`,
		},
		{
			jsonpath:    `$..[0:2]`,
			inputJSON:   `{"a":[],"b":{"a":[]}}`,
			expectedErr: createErrorMemberNotExist(`[0:2]`),
		},
		{
			jsonpath:    `$..[0:2]`,
			inputJSON:   `"x"`,
			expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
		},
	}

	runTestCases(t, "TestValueGroupCombinationRecursiveDescent_SliceQualifier", tests)
}

func TestValueGroupCombinationRecursiveDescent_WildcardQualifier(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$..[*]`,
			inputJSON:    `[[1,3,2],[4,6,5]]`,
			expectedJSON: `[[1,3,2],[4,6,5],1,3,2,4,6,5]`,
		},
		{
			jsonpath:    `$..a[*]`,
			inputJSON:   `"x"`,
			expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
		},
	}

	runTestCases(t, "TestValueGroupCombinationRecursiveDescent_WildcardQualifier", tests)
}

func TestValueGroupCombinationRecursiveDescent_UnionInQualifier(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$..[0,1]`,
			inputJSON:    `[[1,3,2],[4,6,5]]`,
			expectedJSON: `[[1,3,2],[4,6,5],1,3,4,6]`,
		},
		{
			jsonpath:    `$..[0,1]`,
			inputJSON:   `[]`,
			expectedErr: createErrorMemberNotExist(`[0,1]`),
		},
		{
			jsonpath:    `$..[0,1]`,
			inputJSON:   `"x"`,
			expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
		},
	}

	runTestCases(t, "TestValueGroupCombinationRecursiveDescent_UnionInQualifier", tests)
}

func TestValueGroupCombinationRecursiveDescent_FilterQualifier(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$..[?(@.b)]`,
			inputJSON:    `[[{"a":1},{"b":2}],[{"a":3},{"b":4}]]`,
			expectedJSON: `[{"b":2},{"b":4}]`,
		},
		{
			jsonpath:    `$..[?(@.b)]`,
			inputJSON:   `[]`,
			expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
		},
		{
			jsonpath:    `$..[?(@.b)]`,
			inputJSON:   `"x"`,
			expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
		},
	}

	runTestCases(t, "TestValueGroupCombinationRecursiveDescent_FilterQualifier", tests)
}
