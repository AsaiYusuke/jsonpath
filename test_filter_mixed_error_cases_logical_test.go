package jsonpath

import (
	"fmt"
	"testing"
)

func TestFilterBooleanOperations(t *testing.T) {
	testCases := []TestCase{
		// Case variations of null - these are invalid syntax
		{
			jsonpath:    `$[?(@.a==nulL)]`,
			inputJSON:   `[{"a":null}]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==nulL)]`},
		},
		{
			jsonpath:    `$[?(@.a==NulL)]`,
			inputJSON:   `[{"a":null}]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==NulL)]`},
		},
		{
			jsonpath:    `$[?(@.a==NuLL)]`,
			inputJSON:   `[{"a":null}]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==NuLL)]`},
		},
		{
			jsonpath:    `$[?(@.a==NUlL)]`,
			inputJSON:   `[{"a":null}]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==NUlL)]`},
		},
		{
			jsonpath:    `$[?(@.a==NULl)]`,
			inputJSON:   `[{"a":null}]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==NULl)]`},
		},

		// Parenthesized boolean expression comparisons
		{
			jsonpath:    `$[?((@.a<2)==false)]`,
			inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?((@.a<2)==false)]`},
		},
		{
			jsonpath:    `$[?((@.a<2)==true)]`,
			inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?((@.a<2)==true)]`},
		},
		{
			jsonpath:    `$[?((@.a<2)==1)]`,
			inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?((@.a<2)==1)]`},
		},

		// Bitwise operators (invalid in JSONPath filters)
		{
			jsonpath:    `$[?(@.a & @.b)]`,
			inputJSON:   `{}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a & @.b)]`},
		},
		{
			jsonpath:    `$[?(@.a | @.b)]`,
			inputJSON:   `{}`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a | @.b)]`},
		},

		// Negation with complex expressions - these also have invalid syntax
		{
			jsonpath:    `$[?(!(@.a==2))]`,
			inputJSON:   `[{"a":1.9999},{"a":2},{"a":2.0001},{"a":"2"},{"a":true},{"a":{}},{"a":[]},{"a":["b"]},{"a":{"a":"value"}},{"b":"value"}]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(!(@.a==2))]`},
		},
		{
			jsonpath:    `$[?(!(@.a<2))]`,
			inputJSON:   `[{"a":1.9999},{"a":2},{"a":2.0001},{"a":"2"},{"a":true},{"a":{}},{"a":[]},{"a":["b"]},{"a":{"a":"value"}},{"b":"value"}]`,
			expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(!(@.a<2))]`},
		},
	}

	for i, tc := range testCases {
		runTestCase(t, tc, fmt.Sprintf("TestFilterBooleanOperations_case_%d", i))
	}
}
