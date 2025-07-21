package tests

import (
	"testing"
)

func TestInvalidDotNotationSpecialCharacters(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `@`,
			inputJSON:   `{"@":1}`,
			expectedErr: createErrorInvalidSyntax(0, `unrecognized input`, `@`),
		},

		{
			jsonpath:    `$.\a`,
			inputJSON:   `{"a":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.\a`),
		},
		{
			jsonpath:    `$.a\b`,
			inputJSON:   `{"ab":1}`,
			expectedErr: createErrorInvalidSyntax(3, `unrecognized input`, `\b`),
		},
		{
			jsonpath:    `$.\-`,
			inputJSON:   `{"-":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.\-`),
		},
		{
			jsonpath:    `$.\_`,
			inputJSON:   `{"_":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.\_`),
		},

		{
			jsonpath:    `$.'a\.b'`,
			inputJSON:   `{"'a.b'":1,"a":{"b":2},"'a'":{"'b'":3},"'a":{"b'":4}}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.'a\.b'`),
		},

		{
			jsonpath:    `$. `,
			inputJSON:   `{" ":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `. `),
		},
		{
			jsonpath:    `$.!`,
			inputJSON:   `{"!":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.!`),
		},
		{
			jsonpath:    `$."`,
			inputJSON:   `{"\"":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `."`),
		},
		{
			jsonpath:    `$.#`,
			inputJSON:   `{"#":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.#`),
		},
		{
			jsonpath:    `$.$`,
			inputJSON:   `{"$":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.$`),
		},
		{
			jsonpath:    `$.%`,
			inputJSON:   `{"%":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.%`),
		},
		{
			jsonpath:    `$.&`,
			inputJSON:   `{"&":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.&`),
		},
		{
			jsonpath:    `$.'`,
			inputJSON:   `{"'":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.'`),
		},
		{
			jsonpath:    `$.(`,
			inputJSON:   `{"(":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.(`),
		},
		{
			jsonpath:    `$.)`,
			inputJSON:   `{")":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.)`),
		},
		{
			jsonpath:    `$.+`,
			inputJSON:   `{"+":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.+`),
		},
		{
			jsonpath:    "$.`",
			inputJSON:   "{\"`\":1}",
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, ".`"),
		},
		{
			jsonpath:    `$./`,
			inputJSON:   `{"/":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `./`),
		},
		{
			jsonpath:    `$.,`,
			inputJSON:   `{",":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.,`),
		},
		{
			jsonpath:    `$.:`,
			inputJSON:   `{":":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.:`),
		},
		{
			jsonpath:    `$.;`,
			inputJSON:   `{";":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.;`),
		},
		{
			jsonpath:    `$.<`,
			inputJSON:   `{"<":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.<`),
		},
		{
			jsonpath:    `$.=`,
			inputJSON:   `{"=":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.=`),
		},
		{
			jsonpath:    `$.>`,
			inputJSON:   `{">":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.>`),
		},
		{
			jsonpath:    `$.?`,
			inputJSON:   `{"?":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.?`),
		},
		{
			jsonpath:    `$.@`,
			inputJSON:   `{"@":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.@`),
		},
		{
			jsonpath:    `$.\`,
			inputJSON:   `{"\\":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.\`),
		},
		{
			jsonpath:    `$.]`,
			inputJSON:   `{"]":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.]`),
		},
		{
			jsonpath:    `$.^`,
			inputJSON:   `{"^":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.^`),
		},
		{
			jsonpath:    `$.{`,
			inputJSON:   `{"{":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.{`),
		},
		{
			jsonpath:    `$.|`,
			inputJSON:   `{"|":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.|`),
		},
		{
			jsonpath:    `$.}`,
			inputJSON:   `{"}":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.}`),
		},
		{
			jsonpath:    `$.~`,
			inputJSON:   `{"~":1}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `.~`),
		},
	}

	runTestCases(t, "TestInvalidDotNotationSpecialCharacters", tests)
}
