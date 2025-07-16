package jsonpath

import (
	"testing"
)

// TestRetrieve_filterRootNodeComparisonDeleted tests deleted filter cases with root node comparisons
func TestRetrieve_filterRootNodeComparisonDeleted(t *testing.T) {
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

	for i, testCase := range testCases {
		runTestCase(t, testCase, "TestRetrieve_filterRootNodeComparisonDeleted_case_"+string(rune('A'+i)))
	}
}

// TestRetrieve_filterExistsDeleted tests deleted filter cases for existence checks
func TestRetrieve_filterExistsDeleted(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@)]`,
			inputJSON:    `{"a":1,"b":null}`,
			expectedJSON: `[1,null]`,
		},
		{
			jsonpath:     `$[?(@.a)]`,
			inputJSON:    `{"a":{"a":1},"b":{"b":2}}`,
			expectedJSON: `[{"a":1}]`,
		},
	}

	for i, testCase := range testCases {
		runTestCase(t, testCase, "TestRetrieve_filterExistsDeleted_case_"+string(rune('A'+i)))
	}
}

// TestRetrieve_filterNumericScientificNotationDeleted tests deleted filter cases with scientific notation
func TestRetrieve_filterNumericScientificNotationDeleted(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.a==-0.123e2)]`,
		inputJSON:    `[{"a":-12.3,"b":1},{"a":-0.123e2,"b":2},{"a":-0.123},{"a":-12},{"a":12.3},{"a":2},{"a":"-0.123e2"}]`,
		expectedJSON: `[{"a":-12.3,"b":1},{"a":-12.3,"b":2}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterNumericScientificNotationDeleted")
}

// TestRetrieve_filterFalseConditionsDeleted tests deleted filter cases that always return false
func TestRetrieve_filterFalseConditionsDeleted(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(10==20)]`,
		inputJSON:    `[{"a":10},{"a":20},{"a":30},{"a+10":20}]`,
		expectedJSON: ``,
		expectedErr:  createErrorMemberNotExist(`[?(10==20)]`),
	}
	runTestCase(t, testCase, "TestRetrieve_filterFalseConditionsDeleted")
}
