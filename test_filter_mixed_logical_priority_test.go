package jsonpath

import "testing"

func TestRetrieve_filterLogicalCombination_priority_AND_OR(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a && @.b || @.c)]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_priority_AND_OR")
}

func TestRetrieve_filterLogicalCombination_priority_AND_parentheses_OR(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a && (@.b || @.c))]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":2,"b":2},{"a":3,"b":3,"c":3},{"a":5,"c":5}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_priority_AND_parentheses_OR")
}

func TestRetrieve_filterLogicalCombination_priority_parentheses_AND_OR(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?((@.a && @.b) || @.c)]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_priority_parentheses_AND_OR")
}

func TestRetrieve_filterLogicalCombination_priority_OR_AND(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a || @.b && @.c)]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_priority_OR_AND")
}
