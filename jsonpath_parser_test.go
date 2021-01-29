package jsonpath

import (
	"fmt"
	"reflect"
	"testing"
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
		{input: `"`, expectError: ErrorInvalidArgument{argument: `"`, err: fmt.Errorf(`invalid character '"' after top-level value`)}},
		{input: `\u`, expectError: ErrorInvalidArgument{argument: `\u`, err: fmt.Errorf(`invalid character '"' in \u hexadecimal character escape`)}},
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
		{input: `\u`, expectError: ErrorInvalidArgument{argument: `\u`, err: fmt.Errorf(`invalid character '"' in \u hexadecimal character escape`)}},
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
		actualError := recover().(error)
		if reflect.TypeOf(expectedError) != reflect.TypeOf(actualError) ||
			fmt.Sprintf(`%s`, expectedError) != fmt.Sprintf(`%s`, actualError) {
			t.Errorf("expect<%s> != actual<%s>\n", expectedError, actualError)
		}

	}()
	function(text)
	t.Errorf("expect panic")
}
