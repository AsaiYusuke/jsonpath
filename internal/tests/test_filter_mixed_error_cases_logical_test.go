package tests

import (
	"testing"
)

func TestFilterBooleanOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@.a==nulL)]`,
			inputJSON:   `[{"a":null}]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a==nulL)]`),
		},
		{
			jsonpath:    `$[?(@.a==NulL)]`,
			inputJSON:   `[{"a":null}]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a==NulL)]`),
		},
		{
			jsonpath:    `$[?(@.a==NuLL)]`,
			inputJSON:   `[{"a":null}]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a==NuLL)]`),
		},
		{
			jsonpath:    `$[?(@.a==NUlL)]`,
			inputJSON:   `[{"a":null}]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a==NUlL)]`),
		},
		{
			jsonpath:    `$[?(@.a==NULl)]`,
			inputJSON:   `[{"a":null}]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a==NULl)]`),
		},

		{
			jsonpath:    `$[?((@.a<2)==false)]`,
			inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?((@.a<2)==false)]`),
		},
		{
			jsonpath:    `$[?((@.a<2)==true)]`,
			inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?((@.a<2)==true)]`),
		},
		{
			jsonpath:    `$[?((@.a<2)==1)]`,
			inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?((@.a<2)==1)]`),
		},

		{
			jsonpath:    `$[?(@.a & @.b)]`,
			inputJSON:   `{}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a & @.b)]`),
		},
		{
			jsonpath:    `$[?(@.a | @.b)]`,
			inputJSON:   `{}`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a | @.b)]`),
		},

		{
			jsonpath:    `$[?(!(@.a==2))]`,
			inputJSON:   `[{"a":1.9999},{"a":2},{"a":2.0001},{"a":"2"},{"a":true},{"a":{}},{"a":[]},{"a":["b"]},{"a":{"a":"value"}},{"b":"value"}]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(!(@.a==2))]`),
		},
		{
			jsonpath:    `$[?(!(@.a<2))]`,
			inputJSON:   `[{"a":1.9999},{"a":2},{"a":2.0001},{"a":"2"},{"a":true},{"a":{}},{"a":[]},{"a":["b"]},{"a":{"a":"value"}},{"b":"value"}]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(!(@.a<2))]`),
		},
	}

	runTestCases(t, "TestFilterBooleanOperations", testCases)
}

func TestFilterBooleanValueVariations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@.a==falSe)]`,
			inputJSON:   `[{"a":false}]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a==falSe)]`),
		},
		{
			jsonpath:    `$[?(@.a==FaLSE)]`,
			inputJSON:   `[{"a":false}]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a==FaLSE)]`),
		},
		{
			jsonpath:    `$[?(@.a==trUe)]`,
			inputJSON:   `[{"a":true}]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a==trUe)]`),
		},
		{
			jsonpath:    `$[?(@.a==NuLl)]`,
			inputJSON:   `[{"a":null}]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[?(@.a==NuLl)]`),
		},
	}

	runTestCases(t, "TestFilterBooleanValueVariations", testCases)
}
