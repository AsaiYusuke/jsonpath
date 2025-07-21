package tests

import (
	"testing"
)

func TestLiteralString_BasicOperations(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.name == 'John')]`,
			inputJSON:    `[{"name":"John","age":30},{"name":"Jane","age":25}]`,
			expectedJSON: `[{"age":30,"name":"John"}]`,
		},
		{
			jsonpath:     `$[?(@.name == "John")]`,
			inputJSON:    `[{"name":"John","age":30},{"name":"Jane","age":25}]`,
			expectedJSON: `[{"age":30,"name":"John"}]`,
		},
		{
			jsonpath:     `$[?('text' == @.value)]`,
			inputJSON:    `[{"value":"text"},{"value":"other"}]`,
			expectedJSON: `[{"value":"text"}]`,
		},
		{
			jsonpath:     `$[?("text" == @.value)]`,
			inputJSON:    `[{"value":"text"},{"value":"other"}]`,
			expectedJSON: `[{"value":"text"}]`,
		},
	}

	runTestCases(t, "TestLiteralString_BasicOperations", tests)
}

func TestLiteralString_SpecialCharacters(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.text == 'hello world')]`,
			inputJSON:    `[{"text":"hello world"},{"text":"hello"}]`,
			expectedJSON: `[{"text":"hello world"}]`,
		},
		{
			jsonpath:     `$[?(@.text == "hello\nworld")]`,
			inputJSON:    `[{"text":"hello\nworld"},{"text":"hello"}]`,
			expectedJSON: `[{"text":"hello\nworld"}]`,
		},
		{
			jsonpath:     `$[?(@.text == 'with "quotes"')]`,
			inputJSON:    `[{"text":"with \"quotes\""},{"text":"other"}]`,
			expectedJSON: `[{"text":"with \"quotes\""}]`,
		},
	}

	runTestCases(t, "TestLiteralString_SpecialCharacters", tests)
}

func TestLiteralString_EmptyStrings(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.value == '')]`,
			inputJSON:    `[{"value":""},{"value":"text"}]`,
			expectedJSON: `[{"value":""}]`,
		},
		{
			jsonpath:     `$[?(@.value == "")]`,
			inputJSON:    `[{"value":""},{"value":"text"}]`,
			expectedJSON: `[{"value":""}]`,
		},
	}

	runTestCases(t, "TestLiteralString_EmptyStrings", tests)
}
