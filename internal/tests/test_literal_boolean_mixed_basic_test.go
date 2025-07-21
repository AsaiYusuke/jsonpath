package tests

import (
	"testing"
)

func TestLiteralBoolean_TrueComparisons(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.active == true)]`,
			inputJSON:    `[{"active":true},{"active":false}]`,
			expectedJSON: `[{"active":true}]`,
		},
		{
			jsonpath:     `$[?(true == @.active)]`,
			inputJSON:    `[{"active":true},{"active":false}]`,
			expectedJSON: `[{"active":true}]`,
		},
		{
			jsonpath:     `$[?(@.enabled == true)]`,
			inputJSON:    `[{"enabled":true,"name":"test1"},{"enabled":false,"name":"test2"}]`,
			expectedJSON: `[{"enabled":true,"name":"test1"}]`,
		},
	}

	runTestCases(t, "TestLiteralBoolean_TrueComparisons", tests)
}

func TestLiteralBoolean_FalseComparisons(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.active == false)]`,
			inputJSON:    `[{"active":true},{"active":false}]`,
			expectedJSON: `[{"active":false}]`,
		},
		{
			jsonpath:     `$[?(false == @.active)]`,
			inputJSON:    `[{"active":true},{"active":false}]`,
			expectedJSON: `[{"active":false}]`,
		},
		{
			jsonpath:     `$[?(@.enabled == false)]`,
			inputJSON:    `[{"enabled":true,"name":"test1"},{"enabled":false,"name":"test2"}]`,
			expectedJSON: `[{"enabled":false,"name":"test2"}]`,
		},
	}

	runTestCases(t, "TestLiteralBoolean_FalseComparisons", tests)
}

func TestLiteralBoolean_MixedComparisons(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.active != true)]`,
			inputJSON:    `[{"active":true},{"active":false},{"active":"true"}]`,
			expectedJSON: `[{"active":false},{"active":"true"}]`,
		},
		{
			jsonpath:     `$[?(@.active != false)]`,
			inputJSON:    `[{"active":true},{"active":false},{"active":"false"}]`,
			expectedJSON: `[{"active":true},{"active":"false"}]`,
		},
	}

	runTestCases(t, "TestLiteralBoolean_MixedComparisons", tests)
}
