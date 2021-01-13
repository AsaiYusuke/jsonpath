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
	}
	parser := jsonPathParser{}

	for _, testcase := range testcases {
		if testcase.expectError != nil {
			expectPanic(t, testcase.input, testcase.expectError)
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
	}
	parser := jsonPathParser{}

	for _, testcase := range testcases {
		if testcase.expectError != nil {
			expectPanic(t, testcase.input, testcase.expectError)
			continue
		}
		actual := parser.unescapeSingleQuotedString(testcase.input)
		if testcase.expectOut != actual {
			t.Errorf("expect<%s> != actual<%s>\n", testcase.expectOut, actual)
			return
		}

	}
}

func expectPanic(t *testing.T, text string, expectedError error) {
	defer func() {
		actualError := recover().(error)
		if reflect.TypeOf(expectedError) != reflect.TypeOf(actualError) ||
			fmt.Sprintf(`%s`, expectedError) != fmt.Sprintf(`%s`, actualError) {
			t.Errorf("expect<%s> != actual<%s>\n", expectedError, actualError)
		}

	}()
	parser.unescapeDoubleQuotedString(text)
	t.Errorf("expect panic")
}
