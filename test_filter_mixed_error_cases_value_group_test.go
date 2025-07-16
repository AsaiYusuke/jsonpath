package jsonpath

import (
	"fmt"
	"testing"
)

// TestRetrieve_valueGroupFilterInvalidSyntax tests invalid syntax where value group is prohibited in filter
func TestRetrieve_valueGroupFilterInvalidSyntax(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@['a']['a','b']==123)]`,
		inputJSON:    `[{"a":"123"},{"a":123}]`,
		expectedJSON: ``,
		expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@['a']['a','b']==123)]`},
	}
	runTestCase(t, testCase, "TestRetrieve_valueGroupFilterInvalidSyntax")
}

// TestRetrieve_filterRecursiveDescentDeleted tests deleted recursive descent filter cases
func TestRetrieve_filterRecursiveDescentDeleted(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.z[?($.*)]`,
		inputJSON:    `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
		expectedJSON: `[{"b":1},{"b":2},{"b":3}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_filterRecursiveDescentDeleted")
}

// TestRetrieve_valueGroupFilterInvalidSyntaxDeletedCases tests deleted value group invalid syntax cases
func TestRetrieve_valueGroupFilterInvalidSyntaxDeletedCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$.x[?(@[*]>=$.y.a[0:1])]`,
			inputJSON:    `{"x":[[1,2],[3,4],[5,6]],"y":{"a":[3,4,5]}}`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 6, reason: `JSONPath that returns a value group is prohibited`, near: `@[*]>=$.y.a[0:1])]`},
		},
		{
			jsonpath:     `$[?(@[0,1:2]==1)]`,
			inputJSON:    `[[1,2,3],[1],[2,3],1,2]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0,1:2]==1)]`},
		},
		{
			jsonpath:     `$[?(@[0:1]=~/123/)]`,
			inputJSON:    `[{"b":["123"]},{"a":["123"]}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0:1]=~/123/)]`},
		},
		{
			jsonpath:     `$[?(@[0:2]=~/123/)]`,
			inputJSON:    `[{"b":["123"]},{"a":["123"]}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0:2]=~/123/)]`},
		},
		{
			jsonpath:     `$[?(@[0:2].a=~/123/)]`,
			inputJSON:    `[{"b":["123"]},{"a":["123"]}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0:2].a=~/123/)]`},
		},
		{
			jsonpath:     `$[?(@.a[0:2]=~/123/)]`,
			inputJSON:    `[{"b":["123"]},{"a":["123"]}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[0:2]=~/123/)]`},
		},
		{
			jsonpath:     `$[?(@[*]=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[*]=~/123/)]`},
		},
		{
			jsonpath:     `$[?(@[*].a=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[*].a=~/123/)]`},
		},
		{
			jsonpath:     `$[?(@.a[*]=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[*]=~/123/)]`},
		},
		{
			jsonpath:     `$[?(@[0,1]=~/123/)]`,
			inputJSON:    `[{"b":["123"]},{"a":[123,"123"]}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0,1]=~/123/)]`},
		},
		{
			jsonpath:     `$[?(@[0,1:2]=~/123/)]`,
			inputJSON:    `[{"b":["123"]},{"a":[123,"123"]}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0,1:2]=~/123/)]`},
		},
		{
			jsonpath:     `$[?(@[0,1].a=~/123/)]`,
			inputJSON:    `[{"b":["123"]},{"a":[123,"123"]}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0,1].a=~/123/)]`},
		},
		{
			jsonpath:     `$[?(@.a[0,1]=~/123/)]`,
			inputJSON:    `[{"b":["123"]},{"a":[123,"123"]}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[0,1]=~/123/)]`},
		},
		{
			jsonpath:     `$[?($..a=~/123/)]`,
			inputJSON:    `[{"a":"123"},{"a":123}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `$..a=~/123/)]`},
		},
		{
			jsonpath:     `$[?($..a.b=~/123/)]`,
			inputJSON:    `[{"a":"123"},{"a":123}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `$..a.b=~/123/)]`},
		},
		{
			jsonpath:     `$[?($.a..b=~/123/)]`,
			inputJSON:    `[{"a":"123"},{"a":123}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `$.a..b=~/123/)]`},
		},
		{
			jsonpath:     `$[?($..a..b=~/123/)]`,
			inputJSON:    `[{"a":"123"},{"a":123}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `$..a..b=~/123/)]`},
		},
		{
			jsonpath:     `$[?(@['a','b']=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@['a','b']=~/123/)]`},
		},
		{
			jsonpath:     `$[?(@['a','b','c']=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@['a','b','c']=~/123/)]`},
		},
		{
			jsonpath:     `$[?(@['a','b']['a']=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@['a','b']['a']=~/123/)]`},
		},
		{
			jsonpath:     `$[?(@['a']['a','b']=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@['a']['a','b']=~/123/)]`},
		},
		// Value group errors with wildcard operators
		{
			jsonpath:     `$[?(@.*=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.*=~/123/)]`},
		},
		{
			jsonpath:     `$[?(@.*[0]=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.*[0]=~/123/)]`},
		},
		{
			jsonpath:     `$[?(@.*.a=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.*.a=~/123/)]`},
		},
		{
			jsonpath:     `$[?(@.a.*=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a.*=~/123/)]`},
		},

		// Value group errors with sub-filters
		{
			jsonpath:     `$[?(@.a[?(@.b)]=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[?(@.b)]=~/123/)]`},
		},
		{
			jsonpath:     `$[?(@.a[?(@.b>1)]=~/123/)]`,
			inputJSON:    `[{"a":[{"b":1},{"b":2}]},{"a":[{"b":1}]}]`,
			expectedJSON: ``,
			expectedErr:  ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[?(@.b>1)]=~/123/)]`},
		},
	}

	for i, testCase := range testCases {
		runSingleTestCase(t, fmt.Sprintf("ValueGroupInvalidSyntaxDeleted_%d", i), testCase)
	}
}
