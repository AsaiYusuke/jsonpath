package jsonpath

import "testing"

func TestRetrieve_filterLogicalCombination_comparator_OR_exists_comparison(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a || @.b > 2)]`,
		inputJSON:    `[{"a":"a"},{"b":2},{"b":3}]`,
		expectedJSON: `[{"a":"a"},{"b":3}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_comparator_OR_exists_comparison")
}

func TestRetrieve_filterLogicalCombination_comparator_OR_comparison_exists(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.b > 2 || @.a)]`,
		inputJSON:    `[{"a":"a"},{"b":2},{"b":3}]`,
		expectedJSON: `[{"a":"a"},{"b":3}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_comparator_OR_comparison_exists")
}

func TestRetrieve_filterLogicalCombination_comparator_AND_regex_equality(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a =~ /a/ && @.b == 2)]`,
		inputJSON:    `[{"a":"a"},{"a":"a","b":2}]`,
		expectedJSON: `[{"a":"a","b":2}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_comparator_AND_regex_equality")
}

func TestRetrieve_filterLogicalCombination_comparator_AND_equality_regex(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.b == 2 && @.a =~ /a/)]`,
		inputJSON:    `[{"a":"a"},{"a":"a","b":2}]`,
		expectedJSON: `[{"a":"a","b":2}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterLogicalCombination_comparator_AND_equality_regex")
}
