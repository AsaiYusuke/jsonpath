package jsonpath

import "testing"

func TestRetrieve_filterLogicalCombination_NOT_AND_OR_1(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@.a && @.b || @.c)]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_NOT_AND_OR_1")
}

func TestRetrieve_filterLogicalCombination_AND_NOT_OR(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a && !@.b || @.c)]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":1},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_AND_NOT_OR")
}

func TestRetrieve_filterLogicalCombination_NOT_AND_NOT_OR(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@.a && !@.b || @.c)]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_NOT_AND_NOT_OR")
}

func TestRetrieve_filterLogicalCombination_AND_OR_NOT(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a && @.b || !@.c)]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":7}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_AND_OR_NOT")
}

func TestRetrieve_filterLogicalCombination_NOT_AND_OR_NOT(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@.a && @.b || !@.c)]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":1},{"a":2,"b":2},{"b":4,"c":4},{"b":7}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_NOT_AND_OR_NOT")
}

func TestRetrieve_filterLogicalCombination_AND_NOT_OR_NOT(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a && !@.b || !@.c)]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":1},{"a":2,"b":2},{"a":5,"c":5},{"b":7}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_AND_NOT_OR_NOT")
}

func TestRetrieve_filterLogicalCombination_NOT_AND_NOT_OR_NOT(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(!@.a && !@.b || !@.c)]`,
		inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
		expectedJSON: `[{"a":1},{"a":2,"b":2},{"c":6},{"b":7}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_NOT_AND_NOT_OR_NOT")
}
