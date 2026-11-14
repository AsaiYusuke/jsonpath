package syntax

import (
	"testing"

	"github.com/AsaiYusuke/jsonpath/v2/errors"
)

type TestCaseUnescapeJSONString struct {
	input       string
	expectOut   string
	expectError error
}

func TestUnescapeDoubleQuotedString(t *testing.T) {
	testcases := []TestCaseUnescapeJSONString{
		{input: `\"`, expectOut: `"`},
		{input: `\\`, expectOut: `\`},
		{input: `\/`, expectOut: `/`},
		{input: `\b`, expectOut: "\b"},
		{input: `\f`, expectOut: "\f"},
		{input: `\n`, expectOut: "\n"},
		{input: `\r`, expectOut: "\r"},
		{input: `\t`, expectOut: "\t"},
		{input: `\uD834\uDD1E`, expectOut: `𝄞`},
		{input: `"`, expectError: errors.ErrorInvalidArgument{ArgumentName: `"`}},
		{input: `\u`, expectError: errors.ErrorInvalidArgument{ArgumentName: `\u`}},
	}
	parser := jsonPathParser{}

	for _, testcase := range testcases {
		if testcase.expectError != nil {
			expectPanic(t, testcase.input, testcase.expectError, parser.unescapeDoubleQuotedString)
			continue
		}
		actual := parser.unescapeDoubleQuotedString(testcase.input)
		if testcase.expectOut != actual {
			t.Errorf("expect<%s> != actual<%s>\n", testcase.expectOut, actual)
			return
		}

	}
}

func TestUnescapeSingleQuotedString(t *testing.T) {
	testcases := []TestCaseUnescapeJSONString{
		{input: `\'`, expectOut: `'`},
		{input: `"`, expectOut: `"`},
		{input: `\\`, expectOut: `\`},
		{input: `\/`, expectOut: `/`},
		{input: `\b`, expectOut: "\b"},
		{input: `\f`, expectOut: "\f"},
		{input: `\n`, expectOut: "\n"},
		{input: `\r`, expectOut: "\r"},
		{input: `\t`, expectOut: "\t"},
		{input: `\uD834\uDD1E`, expectOut: `𝄞`},
		{input: `'`, expectOut: `'`},
		{input: `\u`, expectError: errors.ErrorInvalidArgument{ArgumentName: `\u`}},
	}
	parser := jsonPathParser{}

	for _, testcase := range testcases {
		if testcase.expectError != nil {
			expectPanic(t, testcase.input, testcase.expectError, parser.unescapeSingleQuotedString)
			continue
		}
		actual := parser.unescapeSingleQuotedString(testcase.input)
		if testcase.expectOut != actual {
			t.Errorf("expect<%s> != actual<%s>\n", testcase.expectOut, actual)
			return
		}

	}
}

func expectPanic(t *testing.T, text string, expectedError error, function func(string) string) {
	defer func() {
		actual := recover()
		actualError, ok := actual.(errors.ErrorInvalidArgument)
		expected, expectedOK := expectedError.(errors.ErrorInvalidArgument)
		if !ok || !expectedOK || expected.ArgumentName != actualError.ArgumentName || actualError.Err == nil {
			t.Errorf("input %q: expected %T with argument %q, got %T (%v)", text, expectedError, expected.ArgumentName, actual, actual)
		}
	}()
	function(text)
	t.Errorf("expect panic")
}
