package tests

import (
	"testing"
)

func TestLiteralNumber_BasicOperations(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.value == 42)]`,
			inputJSON:    `[{"value":42},{"value":43}]`,
			expectedJSON: `[{"value":42}]`,
		},
		{
			jsonpath:     `$[?(42 == @.value)]`,
			inputJSON:    `[{"value":42},{"value":43}]`,
			expectedJSON: `[{"value":42}]`,
		},
		{
			jsonpath:     `$[?(@.price == 3.14)]`,
			inputJSON:    `[{"price":3.14},{"price":2.71}]`,
			expectedJSON: `[{"price":3.14}]`,
		},
		{
			jsonpath:     `$[?(3.14 == @.price)]`,
			inputJSON:    `[{"price":3.14},{"price":2.71}]`,
			expectedJSON: `[{"price":3.14}]`,
		},
	}

	runTestCases(t, "TestLiteralNumber_BasicOperations", tests)
}

func TestLiteralNumber_NegativeNumbers(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.value == -10)]`,
			inputJSON:    `[{"value":-10},{"value":10}]`,
			expectedJSON: `[{"value":-10}]`,
		},
		{
			jsonpath:     `$[?(-10 == @.value)]`,
			inputJSON:    `[{"value":-10},{"value":10}]`,
			expectedJSON: `[{"value":-10}]`,
		},
		{
			jsonpath:     `$[?(@.value == -3.14)]`,
			inputJSON:    `[{"value":-3.14},{"value":3.14}]`,
			expectedJSON: `[{"value":-3.14}]`,
		},
	}

	runTestCases(t, "TestLiteralNumber_NegativeNumbers", tests)
}

func TestLiteralNumber_Zero(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.value == 0)]`,
			inputJSON:    `[{"value":0},{"value":1}]`,
			expectedJSON: `[{"value":0}]`,
		},
		{
			jsonpath:     `$[?(0 == @.value)]`,
			inputJSON:    `[{"value":0},{"value":1}]`,
			expectedJSON: `[{"value":0}]`,
		},
		{
			jsonpath:     `$[?(@.value == 0.0)]`,
			inputJSON:    `[{"value":0.0},{"value":1.0}]`,
			expectedJSON: `[{"value":0}]`,
		},
	}

	runTestCases(t, "TestLiteralNumber_Zero", tests)
}

func TestLiteralNumber_ScientificNotation(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.value == 1e2)]`,
			inputJSON:    `[{"value":100},{"value":1e2}]`,
			expectedJSON: `[{"value":100},{"value":100}]`,
		},
		{
			jsonpath:     `$[?(@.value == -1.23e2)]`,
			inputJSON:    `[{"value":-123},{"value":-1.23e2}]`,
			expectedJSON: `[{"value":-123},{"value":-123}]`,
		},
	}

	runTestCases(t, "TestLiteralNumber_ScientificNotation", tests)
}
