package jsonpath

import (
	"fmt"
	"testing"
)

func TestInvalidDotNotationSpecialCharacters(t *testing.T) {
	tests := []TestCase{
		// Root symbol (@) at start
		{
			jsonpath:    `@`,
			inputJSON:   `{"@":1}`,
			expectedErr: ErrorInvalidSyntax{position: 0, reason: `unrecognized input`, near: `@`},
		},

		// Backslash escape sequences
		{
			jsonpath:    `$.\a`,
			inputJSON:   `{"a":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.\a`},
		},
		{
			jsonpath:    `$.a\b`,
			inputJSON:   `{"ab":1}`,
			expectedErr: ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: `\b`},
		},
		{
			jsonpath:    `$.\-`,
			inputJSON:   `{"-":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.\-`},
		},
		{
			jsonpath:    `$.\_`,
			inputJSON:   `{"_":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.\_`},
		},

		// Quoted key with backslash
		{
			jsonpath:    `$.'a\.b'`,
			inputJSON:   `{"'a.b'":1,"a":{"b":2},"'a'":{"'b'":3},"'a":{"b'":4}}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.'a\.b'`},
		},

		// Special characters after dot
		{
			jsonpath:    `$. `,
			inputJSON:   `{" ":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `. `},
		},
		{
			jsonpath:    `$.!`,
			inputJSON:   `{"!":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.!`},
		},
		{
			jsonpath:    `$."`,
			inputJSON:   `{"\"":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `."`},
		},
		{
			jsonpath:    `$.#`,
			inputJSON:   `{"#":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.#`},
		},
		{
			jsonpath:    `$.$`,
			inputJSON:   `{"$":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.$`},
		},
		{
			jsonpath:    `$.%`,
			inputJSON:   `{"%":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.%`},
		},
		{
			jsonpath:    `$.&`,
			inputJSON:   `{"&":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.&`},
		},
		{
			jsonpath:    `$.'`,
			inputJSON:   `{"'":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.'`},
		},
		{
			jsonpath:    `$.(`,
			inputJSON:   `{"(":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.(`},
		},
		{
			jsonpath:    `$.)`,
			inputJSON:   `{")":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.)`},
		},
		{
			jsonpath:    `$.+`,
			inputJSON:   `{"+":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.+`},
		},
		{
			jsonpath:    "$.`",
			inputJSON:   "{\"`\":1}",
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: ".`"},
		},
		{
			jsonpath:    `$./`,
			inputJSON:   `{"/":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `./`},
		},
		{
			jsonpath:    `$.,`,
			inputJSON:   `{",":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.,`},
		},
		{
			jsonpath:    `$.:`,
			inputJSON:   `{":":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.:`},
		},
		{
			jsonpath:    `$.;`,
			inputJSON:   `{";":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.;`},
		},
		{
			jsonpath:    `$.<`,
			inputJSON:   `{"<":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.<`},
		},
		{
			jsonpath:    `$.=`,
			inputJSON:   `{"=":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.=`},
		},
		{
			jsonpath:    `$.>`,
			inputJSON:   `{">":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.>`},
		},
		{
			jsonpath:    `$.?`,
			inputJSON:   `{"?":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.?`},
		},
		{
			jsonpath:    `$.@`,
			inputJSON:   `{"@":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.@`},
		},
		{
			jsonpath:    `$.\`,
			inputJSON:   `{"\\":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.\`},
		},
		{
			jsonpath:    `$.]`,
			inputJSON:   `{"]":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.]`},
		},
		{
			jsonpath:    `$.^`,
			inputJSON:   `{"^":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.^`},
		},
		{
			jsonpath:    `$.{`,
			inputJSON:   `{"{":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.{`},
		},
		{
			jsonpath:    `$.|`,
			inputJSON:   `{"|":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.|`},
		},
		{
			jsonpath:    `$.}`,
			inputJSON:   `{"}":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.}`},
		},
		{
			jsonpath:    `$.~`,
			inputJSON:   `{"~":1}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.~`},
		},
	}

	for i, tc := range tests {
		runTestCase(t, tc, fmt.Sprintf("TestInvalidDotNotationSpecialCharacters_case_%d", i))
	}
}

func TestInvalidDotNotationSpecialCharacters_additional_brace(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.{`,
		inputJSON:   `{"{":1}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.{`},
	}
	runTestCase(t, testCase, "TestInvalidDotNotationSpecialCharacters_additional_brace")
}

func TestInvalidDotNotationSpecialCharacters_additional_pipe(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.|`,
		inputJSON:   `{"|":1}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.|`},
	}
	runTestCase(t, testCase, "TestInvalidDotNotationSpecialCharacters_additional_pipe")
}
