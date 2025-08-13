package tests

import (
	"testing"
)

func TestLiteralNull_EqualityComparisons(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.value == null)]`,
			inputJSON:    `[{"value":null},{"value":"text"},{"value":0}]`,
			expectedJSON: `[{"value":null}]`,
		},
		{
			jsonpath:     `$[?(null == @.value)]`,
			inputJSON:    `[{"value":null},{"value":"text"},{"value":0}]`,
			expectedJSON: `[{"value":null}]`,
		},
		{
			jsonpath:     `$[?(@.data == null)]`,
			inputJSON:    `[{"data":null,"id":1},{"data":"content","id":2},{"id":3}]`,
			expectedJSON: `[{"data":null,"id":1}]`,
		},
	}

	runTestCases(t, "TestLiteralNull_EqualityComparisons", tests)
}

func TestLiteralNull_InequalityComparisons(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.value != null)]`,
			inputJSON:    `[{"value":null},{"value":"text"},{"value":0}]`,
			expectedJSON: `[{"value":"text"},{"value":0}]`,
		},
		{
			jsonpath:     `$[?(null != @.value)]`,
			inputJSON:    `[{"value":null},{"value":"text"},{"value":0}]`,
			expectedJSON: `[{"value":"text"},{"value":0}]`,
		},
		{
			jsonpath:     `$[?(@.data != null)]`,
			inputJSON:    `[{"data":null,"id":1},{"data":"content","id":2},{"id":3}]`,
			expectedJSON: `[{"data":"content","id":2},{"id":3}]`,
		},
	}

	runTestCases(t, "TestLiteralNull_InequalityComparisons", tests)
}

func TestLiteralNull_MissingProperty(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.missing == null)]`,
			inputJSON:    `[{"value":"text"},{"missing":null}]`,
			expectedJSON: `[{"missing":null}]`,
		},
		{
			jsonpath:     `$[?(@.missing != null)]`,
			inputJSON:    `[{"value":"text"},{"missing":null}]`,
			expectedJSON: `[{"value":"text"}]`,
		},
	}

	runTestCases(t, "TestLiteralNull_MissingProperty", tests)
}
