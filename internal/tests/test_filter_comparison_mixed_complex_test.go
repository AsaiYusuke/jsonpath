package tests

import (
	"testing"
)

func TestFilterComparison_ComplexNestedArrayAccess(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@[1][0]>1)]`,
			inputJSON:    `[1,[21,[221,[222]]]]`,
			expectedJSON: `[[21,[221,[222]]]]`,
		},
		{
			jsonpath:     `$[?(@[1][0]>1)][?(@[1][0]>1)]`,
			inputJSON:    `[1,[21,[221,[222]]]]`,
			expectedJSON: `[[221,[222]]]`,
		},
		{
			jsonpath:     `$[?(@[1][0]>1)][?(@[1][0]>1)][?(@[0]>1)]`,
			inputJSON:    `[1,[21,[221,[222]]]]`,
			expectedJSON: `[[222]]`,
		},
	}

	runTestCases(t, "TestFilterComparison_ComplexNestedArrayAccess", testCases)
}

func TestFilterComparison_ComplexNestedArrayErrors(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@[1][0]>1)][?(@[1][0]>1)][?(@[1]>1)]`,
			inputJSON:   `[1,[21,[221,[222]]]]`,
			expectedErr: createErrorMemberNotExist(`[?(@[1]>1)]`),
		},
	}

	runTestCases(t, "TestFilterComparison_ComplexNestedArrayErrors", testCases)
}

func TestFilterComparison_RootReferenceErrors(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@.a == $.b)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.a == $.b)]`),
		},
		{
			jsonpath:    `$[?(@.a <= $.b)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.a <= $.b)]`),
		},
		{
			jsonpath:    `$[?($.b == @.a)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?($.b == @.a)]`),
		},
		{
			jsonpath:    `$[?($.b < @.a)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?($.b < @.a)]`),
		},
		{
			jsonpath:    `$[?($.b <= @.a)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?($.b <= @.a)]`),
		},
		{
			jsonpath:    `$[?($.b > @.a)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?($.b > @.a)]`),
		},
		{
			jsonpath:    `$[?($.b >= @.a)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?($.b >= @.a)]`),
		},
		{
			jsonpath:    `$[?(@.b == $[0].a)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.b == $[0].a)]`),
		},
		{
			jsonpath:    `$[?(@.b < $[0].a)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.b < $[0].a)]`),
		},
		{
			jsonpath:    `$[?(@.b <= $[0].a)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.b <= $[0].a)]`),
		},
		{
			jsonpath:    `$[?(@.b > $[0].a)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.b > $[0].a)]`),
		},
		{
			jsonpath:    `$[?($[0].a == @.b)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?($[0].a == @.b)]`),
		},
		{
			jsonpath:    `$[?($[0].a < @.b)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?($[0].a < @.b)]`),
		},
		{
			jsonpath:    `$[?($[0].a <= @.b)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?($[0].a <= @.b)]`),
		},
		{
			jsonpath:    `$[?($[0].a > @.b)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?($[0].a > @.b)]`),
		},
		{
			jsonpath:    `$[?($[0].a >= @.b)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?($[0].a >= @.b)]`),
		},
		{
			jsonpath:    `$[?(@.b != $.b)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.b != $.b)]`),
		},
		{
			jsonpath:    `$[?(@.b < $.b)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.b < $.b)]`),
		},
		{
			jsonpath:    `$[?(@.b <= $.b)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.b <= $.b)]`),
		},
		{
			jsonpath:    `$[?(@.b > $.b)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.b > $.b)]`),
		},
		{
			jsonpath:    `$[?(@.b >= $.b)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.b >= $.b)]`),
		},
		{
			jsonpath:    `$[?($.b != @.b)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?($.b != @.b)]`),
		},
		{
			jsonpath:    `$[?($.b < @.b)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?($.b < @.b)]`),
		},
		{
			jsonpath:    `$[?($.b <= @.b)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?($.b <= @.b)]`),
		},
		{
			jsonpath:    `$[?($.b > @.b)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?($.b > @.b)]`),
		},
		{
			jsonpath:    `$[?($.b >= @.b)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?($.b >= @.b)]`),
		},
	}

	runTestCases(t, "TestFilterComparison_RootReferenceErrors", testCases)
}

func TestFilterComparison_RootReferenceSuccessCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a != $.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: `[{"a":0},{"a":1}]`,
		},
		{
			jsonpath:     `$[?($.b != @.a)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: `[{"a":0},{"a":1}]`,
		},
		{
			jsonpath:     `$[?(@.b != $[0].a)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: `[{"a":0},{"a":1}]`,
		},
		{
			jsonpath:     `$[?($[0].a != @.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: `[{"a":0},{"a":1}]`,
		},
		{
			jsonpath:     `$[?(@.b == $.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: `[{"a":0},{"a":1}]`,
		},
		{
			jsonpath:     `$[?($.b == @.b)]`,
			inputJSON:    `[{"a":0},{"a":1}]`,
			expectedJSON: `[{"a":0},{"a":1}]`,
		},
	}

	runTestCases(t, "TestFilterComparison_RootReferenceSuccessCases", testCases)
}

func TestFilterComparison_ComplexRootSelectors(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$.z[?($["x","y"])]`,
			inputJSON:    `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
			expectedJSON: `[{"b":1},{"b":2},{"b":3}]`,
		},
		{
			jsonpath:     `$[1].z[?($[0:1])]`,
			inputJSON:    `[0,{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}]`,
			expectedJSON: `[{"b":1},{"b":2},{"b":3}]`,
		},
		{
			jsonpath:     `$.z[?($[*])]`,
			inputJSON:    `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
			expectedJSON: `[{"b":1},{"b":2},{"b":3}]`,
		},
		{
			jsonpath:     `$[1].z[?($[0,1])]`,
			inputJSON:    `[0,{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}]`,
			expectedJSON: `[{"b":1},{"b":2},{"b":3}]`,
		},
		{
			jsonpath:     `$.z[?($[?(@.x)])]`,
			inputJSON:    `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
			expectedJSON: `[{"b":1},{"b":2},{"b":3}]`,
		},
	}

	runTestCases(t, "TestFilterComparison_ComplexRootSelectors", tests)
}
