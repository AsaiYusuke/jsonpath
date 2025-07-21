package tests

import (
	"testing"
)

func TestFilter_RootNodeComparison(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@==$[1])]`,
			inputJSON:    `[[1],[2],[2],[3]]`,
			expectedJSON: `[[2],[2]]`,
		},
		{
			jsonpath:     `$[?(@==$[1])]`,
			inputJSON:    `[[1],[2,2],[2,2],[3]]`,
			expectedJSON: `[[2,2],[2,2]]`,
		},
		{
			jsonpath:     `$[?(@==$[1])]`,
			inputJSON:    `[[1],[2,{"2":2}],[2,{"2":2}],[3]]`,
			expectedJSON: `[[2,{"2":2}],[2,{"2":2}]]`,
		},
		{
			jsonpath:     `$[?(@==$[1])]`,
			inputJSON:    `[{"a":[1]},{"a":[2]},{"a":[2]},{"a":[3]}]`,
			expectedJSON: `[{"a":[2]},{"a":[2]}]`,
		},
		{
			jsonpath:     `$[?(@==$[1])]`,
			inputJSON:    `[{"a":[1]},{"a":[2,2]},{"a":[2,2]},{"a":[3]}]`,
			expectedJSON: `[{"a":[2,2]},{"a":[2,2]}]`,
		},
		{
			jsonpath:     `$[?(@==$[1])]`,
			inputJSON:    `[{"a":[1]},{"a":[2,{"2":2}]},{"a":[2,{"2":2}]},{"a":[3]}]`,
			expectedJSON: `[{"a":[2,{"2":2}]},{"a":[2,{"2":2}]}]`,
		},
	}

	runTestCases(t, "TestFilter_RootNodeComparison", testCases)
}

func TestFilter_ScientificNotation(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a==-0.123e2)]`,
		inputJSON:    `[{"a":-12.3,"b":1},{"a":-0.123e2,"b":2},{"a":-0.123},{"a":-12},{"a":12.3},{"a":2},{"a":"-0.123e2"}]`,
		expectedJSON: `[{"a":-12.3,"b":1},{"a":-12.3,"b":2}]`,
	}
	runTestCase(t, testCase, "TestFilter_ScientificNotation")
}

func TestFilter_FalseConditions(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(10==20)]`,
		inputJSON:    `[{"a":10},{"a":20},{"a":30},{"a+10":20}]`,
		expectedJSON: ``,
		expectedErr:  createErrorMemberNotExist(`[?(10==20)]`),
	}
	runTestCase(t, testCase, "TestFilter_FalseConditions")
}
