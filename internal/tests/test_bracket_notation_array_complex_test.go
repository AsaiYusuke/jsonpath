package tests

import (
	"testing"
)

func TestArrayAccess_ComplexIndexCombinations(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[*,*,*]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second","third","first","second","third","first","second","third"]`,
		},
		{
			jsonpath:     `$[0,*,*]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","first","second","third","first","second","third"]`,
		},
		{
			jsonpath:     `$[0,0,*]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","first","first","second","third"]`,
		},
		{
			jsonpath:     `$[*,*,0]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second","third","first","second","third","first"]`,
		},
		{
			jsonpath:     `$[*,0,0]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second","third","first","first"]`,
		},
		{
			jsonpath:     `$[0,0,0]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","first","first"]`,
		},
	}

	runTestCases(t, "TestArrayAccess_ComplexIndexCombinations", tests)
}

func TestArrayAccess_WildcardNestedAccess(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[*,*].a`,
			inputJSON:    `[{"a":1},{"b":2}]`,
			expectedJSON: `[1,1]`,
		},
	}

	runTestCases(t, "TestArrayAccess_WildcardNestedAccess", tests)
}
