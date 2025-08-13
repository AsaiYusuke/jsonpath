package tests

import (
	"testing"
)

func TestBracketNotationMultiple_BasicMultipleSelection(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$['a','b']`,
			inputJSON:    `{"a":1, "b":2}`,
			expectedJSON: `[1,2]`,
		},
		{
			jsonpath:     `$['b','a']`,
			inputJSON:    `{"a":1, "b":2}`,
			expectedJSON: `[2,1]`,
		},
		{
			jsonpath:     `$['b','a']`,
			inputJSON:    `{"b":2,"a":1}`,
			expectedJSON: `[2,1]`,
		},
		{
			jsonpath:     `$['a','b']`,
			inputJSON:    `{"b":2,"a":1}`,
			expectedJSON: `[1,2]`,
		},
		{
			jsonpath:     `$['a','b',*]`,
			inputJSON:    `{"b":2,"a":1,"c":3}`,
			expectedJSON: `[1,2,1,2,3]`,
		},
		{
			jsonpath:     `$['b','c']`,
			inputJSON:    `{"a":1,"b":2}`,
			expectedJSON: `[2]`,
		},
		{
			jsonpath:     `$['c','a']`,
			inputJSON:    `{"a":1,"b":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$['c','b']`,
			inputJSON:    `{"a":1,"b":2}`,
			expectedJSON: `[2]`,
		},
		{
			jsonpath:    `$['c','d']`,
			inputJSON:   `{"a":1,"b":2}`,
			expectedErr: createErrorMemberNotExist(`['c','d']`),
		},
	}

	runTestCases(t, "TestBracketNotationMultiple_BasicMultipleSelection", tests)
}

func TestBracketNotationMultiple_InvalidSyntax(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$['a','b',0]`,
			inputJSON:   `{"b":2,"a":1,"c":3}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `['a','b',0]`),
		},
		{
			jsonpath:    `$['a','b',0:1]`,
			inputJSON:   `{"b":2,"a":1,"c":3}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `['a','b',0:1]`),
		},
		{
			jsonpath:    `$['a','b',(command)]`,
			inputJSON:   `{"b":2,"a":1,"c":3}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `['a','b',(command)]`),
		},
		{
			jsonpath:    `$['a','b',?(@)]`,
			inputJSON:   `{"b":2,"a":1,"c":3}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `['a','b',?(@)]`),
		},
	}

	runTestCases(t, "TestBracketNotationMultiple_InvalidSyntax", tests)
}

func TestBracketNotationMultiple_NestedAccess(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$['a','b'].a`,
			inputJSON:    `{"a":{"a":1}, "b":{"c":2}}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$['a','b']['a']`,
			inputJSON:    `{"a":{"a":1}, "b":{"c":2}}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$['a','c']`,
			inputJSON:    `{"a":1,"b":2}`,
			expectedJSON: `[1]`,
		},
	}

	runTestCases(t, "TestBracketNotationMultiple_NestedAccess", tests)
}

func TestBracketNotationMultiple_DoubleNestedAccess(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$['a','b']['a','b']`,
			inputJSON:    `{"a":{"a":1},"b":{"b":2}}`,
			expectedJSON: `[1,2]`,
		},
		{
			jsonpath:     `$['a','b']['a','b']`,
			inputJSON:    `{"a":{"b":1},"b":{"a":2}}`,
			expectedJSON: `[1,2]`,
		},
		{
			jsonpath:     `$['a','b']['a','b']`,
			inputJSON:    `{"a":{"a":1,"b":2},"b":{"c":3}}`,
			expectedJSON: `[1,2]`,
		},
		{
			jsonpath:     `$['a','b']['a','b']`,
			inputJSON:    `{"a":{"b":1},"c":{"a":2}}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:    `$['a','b']['c','d']`,
			inputJSON:   `{"a":{"a":1},"b":{"b":2}}`,
			expectedErr: createErrorMemberNotExist(`['c','d']`),
		},
		{
			jsonpath:    `$['a','b']['c','d']`,
			inputJSON:   `{"a":{"a":1},"c":{"b":2}}`,
			expectedErr: createErrorMemberNotExist(`['c','d']`),
		},
		{
			jsonpath:    `$['a','b']['c','d']`,
			inputJSON:   `{"c":{"a":1},"d":{"b":2}}`,
			expectedErr: createErrorMemberNotExist(`['a','b']`),
		},
		{
			jsonpath:    `$['a','b']['c','d'].e`,
			inputJSON:   `{"a":{"c":1},"b":{"c":2}}`,
			expectedErr: createErrorTypeUnmatched(`.e`, `object`, `float64`),
		},
		{
			jsonpath:    `$['a','b']['c','d'].e`,
			inputJSON:   `{"a":{"a":1},"b":{"c":2}}`,
			expectedErr: createErrorTypeUnmatched(`.e`, `object`, `float64`),
		},
	}

	runTestCases(t, "TestBracketNotationMultiple_DoubleNestedAccess", tests)
}

func TestBracketNotationMultiple_EmptyStringKey(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$['']`,
			inputJSON:    `{"":1, "''":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:    `$[''][0]`,
			inputJSON:   `[1,2,3]`,
			expectedErr: createErrorTypeUnmatched(`['']`, `object`, `[]interface {}`),
		},
	}

	runTestCases(t, "TestBracketNotationMultiple_EmptyStringKey", tests)
}

func TestBracketNotationMultiple_TypeMismatchErrors(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$['a']`,
			inputJSON:   `"abc"`,
			expectedErr: createErrorTypeUnmatched(`['a']`, `object`, `string`),
		},
		{
			jsonpath:    `$['a']`,
			inputJSON:   `[1,2,3]`,
			expectedErr: createErrorTypeUnmatched(`['a']`, `object`, `[]interface {}`),
		},
	}

	runTestCases(t, "TestBracketNotationMultiple_TypeMismatchErrors", tests)
}

func TestBracketNotationMultiple_DuplicateKeys(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$['a','a']`,
			inputJSON:    `{"b":2,"a":1}`,
			expectedJSON: `[1,1]`,
		},
		{
			jsonpath:     `$['a','a','b','b']`,
			inputJSON:    `{"b":2,"a":1}`,
			expectedJSON: `[1,1,2,2]`,
		},
	}

	runTestCases(t, "TestBracketNotationMultiple_DuplicateKeys", tests)
}

func TestBracketNotationMultiple_ArrayIndexCombination(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[0]['a','b']`,
			inputJSON:    `[{"a":1,"b":2},{"a":3,"b":4},{"a":5,"b":6}]`,
			expectedJSON: `[1,2]`,
		},
		{
			jsonpath:     `$[0]['b','a']`,
			inputJSON:    `[{"a":1,"b":2},{"a":3,"b":4},{"a":5,"b":6}]`,
			expectedJSON: `[2,1]`,
		},
		{
			jsonpath:     `$[0:2]['b','a']`,
			inputJSON:    `[{"a":1,"b":2},{"a":3,"b":4},{"a":5,"b":6}]`,
			expectedJSON: `[2,1,4,3]`,
		},
	}

	runTestCases(t, "TestBracketNotationMultiple_ArrayIndexCombination", tests)
}

func TestBracketNotationMultiple_AdditionalTypeErrors(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$['a','b']`,
			inputJSON:   `{}`,
			expectedErr: createErrorMemberNotExist(`['a','b']`),
		},
		{
			jsonpath:    `$['a','b']`,
			inputJSON:   `[]`,
			expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `[]interface {}`),
		},
		{
			jsonpath:    `$['a','b']`,
			inputJSON:   `"abc"`,
			expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `string`),
		},
		{
			jsonpath:    `$['a','b']`,
			inputJSON:   `123`,
			expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `float64`),
		},
		{
			jsonpath:    `$['a','b']`,
			inputJSON:   `true`,
			expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `bool`),
		},
		{
			jsonpath:    `$['a','b']`,
			inputJSON:   `null`,
			expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `null`),
		},
		{
			jsonpath:    `$['a','b']`,
			inputJSON:   `[1,2,3]`,
			expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `[]interface {}`),
		},
	}

	runTestCases(t, "TestBracketNotationMultiple_AdditionalTypeErrors", tests)
}

func TestBracketNotationMultiple_DeepChainErrors(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$['a','b'].a.b`,
			inputJSON:   `{"c":{"b":1}}`,
			expectedErr: createErrorMemberNotExist(`['a','b']`),
		},
		{
			jsonpath:    `$['a','b'].a.b`,
			inputJSON:   `{"a":{"b":1}}`,
			expectedErr: createErrorMemberNotExist(`.a`),
		},
		{
			jsonpath:    `$['a','b'].a.b.c`,
			inputJSON:   `{"a":{"b":1},"b":{"a":2}}`,
			expectedErr: createErrorTypeUnmatched(`.b`, `object`, `float64`),
		},
		{
			jsonpath:    `$['a','b'].a.b.c`,
			inputJSON:   `{"a":{"a":1},"b":{"a":{"c":2}}}`,
			expectedErr: createErrorMemberNotExist(`.b`),
		},
		{
			jsonpath:    `$['a','b','x']['c','d'].e`,
			inputJSON:   `{"a":{"a":1},"b":{"c":2}}`,
			expectedErr: createErrorTypeUnmatched(`.e`, `object`, `float64`),
		},
	}

	runTestCases(t, "TestBracketNotationMultiple_DeepChainErrors", tests)
}
