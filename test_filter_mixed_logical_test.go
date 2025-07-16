package jsonpath

import (
	"testing"
)

func TestFilterLogicalAND_BasicOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a && @.b)]`,
			inputJSON:    `[{"a":1},{"b":2},{"a":3,"b":4}]`,
			expectedJSON: `[{"a":3,"b":4}]`,
		},
		{
			jsonpath:     `$[?(@.a>1 && @.a<3)]`,
			inputJSON:    `[{"a":1},{"a":1.1},{"a":2.9},{"a":3}]`,
			expectedJSON: `[{"a":1.1},{"a":2.9}]`,
		},
		{
			jsonpath:     `$[?(@.a<3 && @.a>1)]`,
			inputJSON:    `[{"a":1},{"a":1.1},{"a":2.9},{"a":3}]`,
			expectedJSON: `[{"a":1.1},{"a":2.9}]`,
		},
		{
			jsonpath:    `$[?((1==2) && @.a>1)]`,
			inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
			expectedErr: createErrorMemberNotExist(`[?((1==2) && @.a>1)]`),
		},
		{
			jsonpath:     `$[?((1==1) && @.a>1)]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
			expectedJSON: `[{"a":2},{"a":3}]`,
		},
		{
			jsonpath:    `$[?(@.a>1 && (1==2))]`,
			inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.a>1 && (1==2))]`),
		},
		{
			jsonpath:     `$[?(@.a>1 && (1==1))]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
			expectedJSON: `[{"a":2},{"a":3}]`,
		},
		{
			jsonpath:    `$[?(@.x && @.b > 2)]`,
			inputJSON:   `[{"a":"a"},{"b":2},{"b":3}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.x && @.b > 2)]`),
		},
		{
			jsonpath:    `$[?(@.b > 2 && @.x)]`,
			inputJSON:   `[{"a":"a"},{"b":2},{"b":3}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.b > 2 && @.x)]`),
		},
		{
			jsonpath:    `$[?(@.x && @.x)]`,
			inputJSON:   `[{"a":"a"},{"b":2},{"b":3}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.x && @.x)]`),
		},
		{
			jsonpath:    `$[?(@.b > 2 && @.b < 2)]`,
			inputJSON:   `[{"b":1},{"b":2},{"b":3}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.b > 2 && @.b < 2)]`),
		},
		{
			jsonpath:     `$.z[?($..x && @.b < 2)]`,
			inputJSON:    `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
			expectedJSON: `[{"b":1}]`,
		},
		{
			jsonpath:    `$.z[?($..xx && @.b < 2)]`,
			inputJSON:   `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
			expectedErr: createErrorMemberNotExist(`[?($..xx && @.b < 2)]`),
		},
	}

	for _, testCase := range testCases {
		runSingleTestCase(t, "TestFilterLogicalAND_BasicOperations", testCase)
	}
}

func TestFilterLogicalOR_BasicOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a || @.b)]`,
			inputJSON:    `[{"a":1},{"b":2},{"c":3}]`,
			expectedJSON: `[{"a":1},{"b":2}]`,
		},
		{
			jsonpath:     `$[?(@.a>2 || @.a<2)]`,
			inputJSON:    `[{"a":1},{"a":1.9},{"a":2},{"a":2.1},{"a":3}]`,
			expectedJSON: `[{"a":1},{"a":1.9},{"a":2.1},{"a":3}]`,
		},
		{
			jsonpath:     `$[?(@.a<2 || @.a>2)]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
			expectedJSON: `[{"a":1},{"a":3}]`,
		},
		{
			jsonpath:     `$[?((1==2) || @.a>1)]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
			expectedJSON: `[{"a":2},{"a":3}]`,
		},
		{
			jsonpath:     `$[?((1==1) || @.a>1)]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
			expectedJSON: `[{"a":1},{"a":2},{"a":3}]`,
		},
		{
			jsonpath:     `$[?(@.a>1 || (1==2))]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
			expectedJSON: `[{"a":2},{"a":3}]`,
		},
		{
			jsonpath:     `$[?(@.a>1 || (1==1))]`,
			inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
			expectedJSON: `[{"a":1},{"a":2},{"a":3}]`,
		},
		{
			jsonpath:     `$[?(@.x || @.b > 2)]`,
			inputJSON:    `[{"a":"a"},{"b":2},{"b":3}]`,
			expectedJSON: `[{"b":3}]`,
		},
		{
			jsonpath:     `$[?(@.b > 2 || @.x)]`,
			inputJSON:    `[{"a":"a"},{"b":2},{"b":3}]`,
			expectedJSON: `[{"b":3}]`,
		},
		{
			jsonpath:    `$[?(@.x || @.x)]`,
			inputJSON:   `[{"a":"a"},{"b":2},{"b":3}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.x || @.x)]`),
		},
		{
			jsonpath:     `$[?(@.b > 2 || @.b < 2)]`,
			inputJSON:    `[{"b":1},{"b":2},{"b":3}]`,
			expectedJSON: `[{"b":1},{"b":3}]`,
		},
		{
			jsonpath:     `$.z[?($..x || @.b < 2)]`,
			inputJSON:    `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
			expectedJSON: `[{"b":1},{"b":2},{"b":3}]`,
		},
		{
			jsonpath:     `$.z[?($..xx || @.b < 2)]`,
			inputJSON:    `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
			expectedJSON: `[{"b":1}]`,
		},
	}

	for _, testCase := range testCases {
		runSingleTestCase(t, "TestFilterLogicalOR_BasicOperations", testCase)
	}
}

func TestFilterLogicalNOT_BasicOperations(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@.a)]`,
		inputJSON:    `[{"a":1},{"b":2},{"a":3,"b":4}]`,
		expectedJSON: `[{"b":2}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_NOT_basic")
}

func TestFilterLogicalNOT_MissingField(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@.c)]`,
		inputJSON:    `[{"a":1},{"b":2},{"a":3,"b":4}]`,
		expectedJSON: `[{"a":1},{"b":2},{"a":3,"b":4}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_NOT_missing_field")
}

func TestFilterLogicalNOT_RecursiveDescentExistsError(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.z[?(!$..x)]`,
		inputJSON:   `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
		expectedErr: createErrorMemberNotExist(`[?(!$..x)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_NOT_recursive_descent_exists_error")
}

func TestFilterLogicalNOT_RecursiveDescentMissing(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.z[?(!$..xx)]`,
		inputJSON:    `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
		expectedJSON: `[{"b":1},{"b":2},{"b":3}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_NOT_recursive_descent_missing")
}

func TestFilterLogicalNOT_RootReferenceError(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(!$)]`,
		inputJSON:   `{"a":1,"b":2}`,
		expectedErr: createErrorMemberNotExist(`[?(!$)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_NOT_root_reference_error")
}

func TestFilterLogicalNOT_CurrentNodeError(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(!@)]`,
		inputJSON:   `{"a":1}`,
		expectedErr: createErrorMemberNotExist(`[?(!@)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_NOT_current_node_error")
}

func TestFilterLogicalCombination_ComplexOperations(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a || @.b > 2)]`,
		inputJSON:    `[{"a":"a"},{"b":2},{"b":3}]`,
		expectedJSON: `[{"a":"a"},{"b":3}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_comparator_OR_exists_comparison")
}

func TestFilterLogicalCombination_ComparisonExists(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.b > 2 || @.a)]`,
		inputJSON:    `[{"a":"a"},{"b":2},{"b":3}]`,
		expectedJSON: `[{"a":"a"},{"b":3}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_comparator_OR_comparison_exists")
}

func TestFilterLogicalCombination_RegexEquality(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a =~ /a/ && @.b == 2)]`,
		inputJSON:    `[{"a":"a"},{"a":"a","b":2}]`,
		expectedJSON: `[{"a":"a","b":2}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_comparator_AND_regex_equality")
}

func TestFilterLogicalCombination_EqualityRegex(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.b == 2 && @.a =~ /a/)]`,
		inputJSON:    `[{"a":"a"},{"a":"a","b":2}]`,
		expectedJSON: `[{"a":"a","b":2}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_comparator_AND_equality_regex")
}

func TestFilterLogicalPriority_AND_OR(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a && @.b || @.c)]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_priority_AND_OR")
}

func TestFilterLogicalPriority_AND_Parentheses_OR(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a && (@.b || @.c))]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":2,"b":2},{"a":3,"b":3,"c":3},{"a":5,"c":5}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_priority_AND_parentheses_OR")
}

func TestFilterLogicalPriority_Parentheses_AND_OR(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?((@.a && @.b) || @.c)]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_priority_parentheses_AND_OR")
}

func TestFilterLogicalPriority_OR_AND(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a || @.b && @.c)]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_priority_OR_AND")
}

func TestFilterLogicalNOT_AND_OR_1(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@.a && @.b || @.c)]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_NOT_AND_OR_1")
}

func TestFilterLogicalNOT_AND_NOT_OR(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a && !@.b || @.c)]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":1},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_AND_NOT_OR")
}

func TestFilterLogicalNOT_NOT_AND_NOT_OR(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@.a && !@.b || @.c)]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_NOT_AND_NOT_OR")
}

func TestFilterLogicalNOT_AND_OR_NOT(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a && @.b || !@.c)]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":7}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_AND_OR_NOT")
}

func TestFilterLogicalNOT_NOT_AND_OR_NOT(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@.a && @.b || !@.c)]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":1},{"a":2,"b":2},{"b":4,"c":4},{"b":7}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_NOT_AND_OR_NOT")
}

func TestFilterLogicalNOT_AND_NOT_OR_NOT(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a && !@.b || !@.c)]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":1},{"a":2,"b":2},{"a":5,"c":5},{"b":7}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_AND_NOT_OR_NOT")
}

func TestFilterLogicalNOT_NOT_AND_NOT_OR_NOT(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@.a && !@.b || !@.c)]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":1},{"a":2,"b":2},{"c":6},{"b":7}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_NOT_AND_NOT_OR_NOT")
}
