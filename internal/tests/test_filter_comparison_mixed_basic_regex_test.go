package tests

import (
	"testing"
)

func TestComparisonRegex_BasicOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a =~ /ab/)]`,
			inputJSON:    `[{"a":"abc"},{"a":"def"},{"a":"abef"}]`,
			expectedJSON: `[{"a":"abc"},{"a":"abef"}]`,
		},
		{
			jsonpath:     `$[?(@.a =~ /123/)]`,
			inputJSON:    `[{"a":"abc123"},{"a":"def"},{"a":"123abc"}]`,
			expectedJSON: `[{"a":"abc123"},{"a":"123abc"}]`,
		},
		{
			jsonpath:     `$[?(@.a=~/テスト/)]`,
			inputJSON:    `[{"a":"テストデータ"},{"a":"サンプル"},{"a":"テスト"}]`,
			expectedJSON: `[{"a":"テストデータ"},{"a":"テスト"}]`,
		},
	}

	runTestCases(t, "TestComparisonRegex_BasicOperations", testCases)
}

func TestComparisonRegex_ComplexPatterns(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a=~/^\d+[a-d]\/\\$/)]`,
			inputJSON:    `[{"a":"012b/\\"},{"a":"ab/\\"},{"a":"1b\\"},{"a":"1b//"},{"a":"1b/\""}]`,
			expectedJSON: `[{"a":"012b/\\"}]`,
		},
		{
			jsonpath:     `$[?(@.a=~/(?i)CASE/)]`,
			inputJSON:    `[{"a":"case"},{"a":"CASE"},{"a":"Case"},{"a":"abc"}]`,
			expectedJSON: `[{"a":"case"},{"a":"CASE"},{"a":"Case"}]`,
		},
	}

	runTestCases(t, "TestComparisonRegex_ComplexPatterns", testCases)
}

func TestComparisonRegex_ArrayElements(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@[0]=~/123/)]`,
			inputJSON:    `[["123"],["456"]]`,
			expectedJSON: `[["123"]]`,
		},
	}

	runTestCases(t, "TestComparisonRegex_ArrayElements", testCases)
}

func TestComparisonRegex_MixedTypes(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a =~ /ab/)]`,
			inputJSON:    `[{"a":"abc"},{"a":1},{"a":"def"}]`,
			expectedJSON: `[{"a":"abc"}]`,
		},
		{
			jsonpath:     `$[?(@.a =~ /123/)]`,
			inputJSON:    `[{"a":123},{"a":"123"},{"a":"12"},{"a":"23"},{"a":"0123"},{"a":"1234"}]`,
			expectedJSON: `[{"a":"123"},{"a":"0123"},{"a":"1234"}]`,
		},
		{
			jsonpath:     `$[?(@.a=~/テスト/)]`,
			inputJSON:    `[{"a":"123テストabc"}]`,
			expectedJSON: `[{"a":"123テストabc"}]`,
		},
	}

	runTestCases(t, "TestComparisonRegex_MixedTypes", testCases)
}
