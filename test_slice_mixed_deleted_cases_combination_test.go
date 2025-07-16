package jsonpath

import (
	"testing"
)

func TestSliceAndCombinationDeletedCases(t *testing.T) {
	tests := []struct {
		name     string
		testCase TestCase
	}{
		{
			name: "slice with nested member access",
			testCase: TestCase{
				jsonpath:    `$[0:2].a.b`,
				inputJSON:   `[{"b": 1}]`,
				expectedErr: createErrorMemberNotExist(`.a`),
			},
		},
		{
			name: "recursive descent wildcard on empty array",
			testCase: TestCase{
				jsonpath:    `$..[*]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
		},
		{
			name: "recursive descent filter on empty array",
			testCase: TestCase{
				jsonpath:    `$..[?(@.b)]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
			},
		},
		{
			name: "slice recursive descent on empty array",
			testCase: TestCase{
				jsonpath:    `$[0:2]..a`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
		},
		{
			name: "slice wildcard on empty array",
			testCase: TestCase{
				jsonpath:    `$[0:2].*`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
		},
		{
			name: "slice nested slice on empty array",
			testCase: TestCase{
				jsonpath:    `$[0:2][0:2]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
		},
		{
			name: "slice nested wildcard on empty array",
			testCase: TestCase{
				jsonpath:    `$[0:2][*]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
		},
		{
			name: "slice nested filter on empty array",
			testCase: TestCase{
				jsonpath:    `$[0:2][?(@.b)]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
		},
		{
			name: "wildcard recursive descent on empty array",
			testCase: TestCase{
				jsonpath:    `$[*]..a`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
		},
		{
			name: "wildcard property access on empty array",
			testCase: TestCase{
				jsonpath:    `$[*].*`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
		},
		{
			name: "wildcard nested slice on empty array",
			testCase: TestCase{
				jsonpath:    `$[*][0:2]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
		},
		{
			name: "wildcard nested wildcard on empty array",
			testCase: TestCase{
				jsonpath:    `$[*][*]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
		},
		{
			name: "wildcard nested filter on empty array",
			testCase: TestCase{
				jsonpath:    `$[*][?(@.b)]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
		},
		{
			name: "filter recursive descent type mismatch",
			testCase: TestCase{
				jsonpath:    `$[?(@.b)]..a`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[?(@.b)]`, `object/array`, `string`),
			},
		},
		{
			name: "filter property access type mismatch",
			testCase: TestCase{
				jsonpath:    `$[?(@.b)].*`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[?(@.b)]`, `object/array`, `string`),
			},
		},
		{
			name: "filter slice on empty array",
			testCase: TestCase{
				jsonpath:    `$[?(@)][0:2]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[?(@)]`),
			},
		},
		{
			name: "filter slice type mismatch",
			testCase: TestCase{
				jsonpath:    `$[?(@)][0:2]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[?(@)]`, `object/array`, `string`),
			},
		},
		{
			name: "filter wildcard on empty array",
			testCase: TestCase{
				jsonpath:    `$[?(@)][*]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[?(@)]`),
			},
		},
		{
			name: "filter wildcard type mismatch",
			testCase: TestCase{
				jsonpath:    `$[?(@)][*]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[?(@)]`, `object/array`, `string`),
			},
		},
		{
			name: "nested filter type mismatch",
			testCase: TestCase{
				jsonpath:    `$[?(@.a)][?(@.b)]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[?(@.a)]`, `object/array`, `string`),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestCase(t, tt.testCase, tt.name)
		})
	}
}
