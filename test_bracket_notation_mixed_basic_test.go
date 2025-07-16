package jsonpath

import (
	"fmt"
	"testing"
)

func TestBracketNotationBasic_SimpleAccess(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$['a']`,
			inputJSON:    `{"a":"b","c":{"d":"e"}}`,
			expectedJSON: `["b"]`,
		},
		{
			jsonpath:     `$[0]['a']`,
			inputJSON:    `[{"a":"b","c":{"d":"e"}},{"x":"y"}]`,
			expectedJSON: `["b"]`,
		},
		{
			jsonpath:     `$['a'][0]['b']`,
			inputJSON:    `{"a":[{"b":"x"},"y"],"c":{"d":"e"}}`,
			expectedJSON: `["x"]`,
		},
		{
			jsonpath:     `$[0:2]['b']`,
			inputJSON:    `[{"a":1},{"b":3},{"b":2,"c":4}]`,
			expectedJSON: `[3]`,
		},
		{
			jsonpath:     `$[:]['b']`,
			inputJSON:    `[{"a":1},{"b":3},{"b":2,"c":4}]`,
			expectedJSON: `[3,2]`,
		},
		{
			jsonpath:     `$['a']['a2']`,
			inputJSON:    `{"a":{"a1":"1","a2":"2"},"b":{"b1":"3"}}`,
			expectedJSON: `["2"]`,
		},
		{
			jsonpath:     `$['0']`,
			inputJSON:    `{"0":1,"a":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$['1']`,
			inputJSON:    `{"1":"a","a":2}`,
			expectedJSON: `["a"]`,
		},
		{
			jsonpath:    `$['d']`,
			inputJSON:   `{"a":"b","c":{"d":"e"}}`,
			expectedErr: createErrorMemberNotExist(`['d']`),
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("SimpleAccess_%d", i), test)
	}
}

func TestBracketNotationBasic_SpecialCharacters(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$['ab']`,
			inputJSON:    `{"ab":1,"b":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$['a\'b']`,
			inputJSON:    `{"a'b":1,"b":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$['ab\'c']`,
			inputJSON:    `{"ab'c":1,"b":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$['a\/b']`,
			inputJSON:    `{"a\/b":1,"b":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$['a\\b']`,
			inputJSON:    `{"a\\b":1,"b":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$['a\bb']`,
			inputJSON:    `{"a\bb":1,"b":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$['a\fb']`,
			inputJSON:    `{"a\fb":1,"b":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$['a\nb']`,
			inputJSON:    `{"a\nb":1,"b":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$['a\rb']`,
			inputJSON:    `{"a\rb":1,"b":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$['a\tb']`,
			inputJSON:    `{"a\tb":1,"b":2}`,
			expectedJSON: `[1]`,
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("SpecialCharacters_%d", i), test)
	}
}

func TestBracketNotationBasic_UnicodeEscapes(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$["a\"b"]`,
			inputJSON:    `{"a\"b":1}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$["ab\"c"]`,
			inputJSON:    `{"ab\"c":1}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$["a\/b"]`,
			inputJSON:    `{"a\/b":1}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$["a\\b"]`,
			inputJSON:    `{"a\\b":1}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$["a\bb"]`,
			inputJSON:    `{"a\bb":1}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$["a\fb"]`,
			inputJSON:    `{"a\fb":1}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$["a\nb"]`,
			inputJSON:    `{"a\nb":1}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$["a\rb"]`,
			inputJSON:    `{"a\rb":1}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$["a\tb"]`,
			inputJSON:    `{"a\tb":1}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$["a\uD834\uDD1Ec"]`,
			inputJSON:    `{"a\uD834\uDD1Ec":1,"b":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$["a\ud834\udd1ec"]`,
			inputJSON:    `{"a\uD834\uDD1Ec":1,"b":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$["\u0000"]`,
			inputJSON:    `{"\u0000":1,"b":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$["\uabcd"]`,
			inputJSON:    `{"\uabcd":1,"b":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$["\uABCD"]`,
			inputJSON:    `{"\uabcd":1,"b":2}`,
			expectedJSON: `[1]`,
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("UnicodeEscapes_%d", i), test)
	}
}

func TestBracketNotationBasic_InvalidUnicode(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$['\uX000']`,
			inputJSON:   `{"a":1,"b":2}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['\uX000']`},
		},
		{
			jsonpath:    `$['\u0X00']`,
			inputJSON:   `{"a":1,"b":2}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['\u0X00']`},
		},
		{
			jsonpath:    `$['\u00X0']`,
			inputJSON:   `{"a":1,"b":2}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['\u00X0']`},
		},
		{
			jsonpath:    `$['\u000X']`,
			inputJSON:   `{"a":1,"b":2}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['\u000X']`},
		},
		{
			jsonpath:    `$["\uX000"]`,
			inputJSON:   `{"a":1,"b":2}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `["\uX000"]`},
		},
		{
			jsonpath:    `$["\u0X00"]`,
			inputJSON:   `{"a":1,"b":2}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `["\u0X00"]`},
		},
		{
			jsonpath:    `$["\u00X0"]`,
			inputJSON:   `{"a":1,"b":2}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `["\u00X0"]`},
		},
		{
			jsonpath:    `$["\u000X"]`,
			inputJSON:   `{"a":1,"b":2}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `["\u000X"]`},
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("InvalidUnicode_%d", i), test)
	}
}

func TestBracketNotationBasic_SpecialQuoting(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$['"']`,
			inputJSON:    `{"\"":1}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$[':@."$,*\'\\']`,
			inputJSON:    `{":@.\"$,*'\\": 1}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$[""]`,
			inputJSON:    `{"":1, "''":2,"\"\"":3}`,
			expectedJSON: `[1]`,
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("SpecialQuoting_%d", i), test)
	}
}

func TestBracketNotationBasic_InvalidSyntax(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$.['a']`,
			inputJSON:   `{"a":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.['a']`},
		},
		{
			jsonpath:    `$['a\\'b']`,
			inputJSON:   `["a"]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['a\\'b']`},
		},
		{
			jsonpath:    `$['ab\']`,
			inputJSON:   `["a"]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['ab\']`},
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("InvalidSyntax_%d", i), test)
	}
}

func TestBracketNotationBasic_MoreSpecialCharacters(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$['a"c']`,
			inputJSON:    `{"a\"c":1,"b":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$['a\uD834\uDD1Ec']`,
			inputJSON:    `{"a\uD834\uDD1Ec":1,"b":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$['a\ud834\udd1ec']`,
			inputJSON:    `{"a\uD834\uDD1Ec":1,"b":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$['\u0000']`,
			inputJSON:    `{"\u0000":1,"b":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$['\uabcd']`,
			inputJSON:    `{"\uabcd":1,"b":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$['\uABCD']`,
			inputJSON:    `{"\uabcd":1,"b":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$["ab"]`,
			inputJSON:    `{"ab":1}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$["a'b"]`,
			inputJSON:    `{"a'b":1}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$['a.b']`,
			inputJSON:    `{"a.b":1,"a":{"b":2}}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$[':']`,
			inputJSON:    `{":":1,"b":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$['[']`,
			inputJSON:    `{"[":1,"]":2}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$[']']`,
			inputJSON:    `{"[":1,"]":2}`,
			expectedJSON: `[2]`,
		},
		{
			jsonpath:     `$['$']`,
			inputJSON:    `{"$":2}`,
			expectedJSON: `[2]`,
		},
		{
			jsonpath:     `$['@']`,
			inputJSON:    `{"@":2}`,
			expectedJSON: `[2]`,
		},
		{
			jsonpath:     `$['*']`,
			inputJSON:    `{"*":2}`,
			expectedJSON: `[2]`,
		},
		{
			jsonpath:    `$['*']`,
			inputJSON:   `{"a":1,"b":2}`,
			expectedErr: createErrorMemberNotExist(`['*']`),
		},
		{
			jsonpath:     `$['.']`,
			inputJSON:    `{".":1}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$[',']`,
			inputJSON:    `{",":1}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$['.*']`,
			inputJSON:    `{".*":1}`,
			expectedJSON: `[1]`,
		},
		{
			jsonpath:     `$["'"]`,
			inputJSON:    `{"'":1}`,
			expectedJSON: `[1]`,
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("MoreSpecialCharacters_%d", i), test)
	}
}
