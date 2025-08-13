package tests

import (
	"testing"
)

func TestFilterTest_PropertyExistence(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.name)]`,
			inputJSON:    `[{"name":"John","age":30},{"age":25},{"name":"Jane"}]`,
			expectedJSON: `[{"age":30,"name":"John"},{"name":"Jane"}]`,
		},
		{
			jsonpath:     `$[?(@.value)]`,
			inputJSON:    `[{"value":0},{"value":""},{"other":"data"},{"value":null}]`,
			expectedJSON: `[{"value":0},{"value":""},{"value":null}]`,
		},
		{
			jsonpath:     `$[?(@.active)]`,
			inputJSON:    `[{"active":true},{"active":false},{"disabled":true}]`,
			expectedJSON: `[{"active":true},{"active":false}]`,
		},
	}

	runTestCases(t, "TestFilterTest_PropertyExistence", tests)
}

func TestFilterTest_PropertyNonExistence(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(!@.name)]`,
			inputJSON:    `[{"name":"John","age":30},{"age":25},{"name":"Jane"}]`,
			expectedJSON: `[{"age":25}]`,
		},
		{
			jsonpath:     `$[?(!@.active)]`,
			inputJSON:    `[{"active":true},{"disabled":false},{"active":false}]`,
			expectedJSON: `[{"disabled":false}]`,
		},
	}

	runTestCases(t, "TestFilterTest_PropertyNonExistence", tests)
}

func TestFilterTest_NestedProperties(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.user.name)]`,
			inputJSON:    `[{"user":{"name":"John"}},{"user":{}},{"other":"data"}]`,
			expectedJSON: `[{"user":{"name":"John"}}]`,
		},
		{
			jsonpath:     `$[?(@.data.value)]`,
			inputJSON:    `[{"data":{"value":42}},{"data":{"other":1}},{"different":"structure"}]`,
			expectedJSON: `[{"data":{"value":42}}]`,
		},
	}

	runTestCases(t, "TestFilterTest_NestedProperties", tests)
}

func TestFilterTest_ArrayElements(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@[0])]`,
			inputJSON:    `[[1,2,3],[],{"0":"value"}]`,
			expectedJSON: `[[1,2,3]]`,
		},
		{
			jsonpath:     `$[?(@[-1])]`,
			inputJSON:    `[[1,2,3],[4],["single"]]`,
			expectedJSON: `[[1,2,3],[4],["single"]]`,
		},
	}

	runTestCases(t, "TestFilterTest_ArrayElements", tests)
}

func TestFilterTest_MissingProperties(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$[?(@.missing)]`,
			inputJSON:   `[{"name":"John"},{"age":30}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.missing)]`),
		},
		{
			jsonpath:    `$[?(@.user.missing)]`,
			inputJSON:   `[{"user":{"name":"John"}},{"other":"data"}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.user.missing)]`),
		},
	}

	runTestCases(t, "TestFilterTest_MissingProperties", tests)
}
