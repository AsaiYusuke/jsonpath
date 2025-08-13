package tests

import (
	"testing"
)

func TestFilter_ValueGroupInvalidSyntax(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@['a']['a','b']==123)]`,
		inputJSON:    `[{"a":"123"},{"a":123}]`,
		expectedJSON: ``,
		expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@['a']['a','b']==123)]`),
	}
	runTestCase(t, testCase, "TestFilter_ValueGroupInvalidSyntax")
}

func TestFilter_RecursiveDescentDeleted(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.z[?($.*)]`,
		inputJSON:    `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
		expectedJSON: `[{"b":1},{"b":2},{"b":3}]`,
	}
	runTestCase(t, testCase, "TestFilter_RecursiveDescentDeleted")
}

func TestFilter_ValueGroupInvalidSyntaxDeletedCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$.x[?(@[*]>=$.y.a[0:1])]`,
			inputJSON:    `{"x":[[1,2],[3,4],[5,6]],"y":{"a":[3,4,5]}}`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(6, `JSONPath that returns a value group is prohibited`, `@[*]>=$.y.a[0:1])]`),
		},
		{
			jsonpath:     `$[?(@[0,1:2]==1)]`,
			inputJSON:    `[[1,2,3],[1],[2,3],1,2]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@[0,1:2]==1)]`),
		},
		{
			jsonpath:     `$[?(@[0:1]=~/123/)]`,
			inputJSON:    `[{"b":["123"]},{"a":["123"]}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@[0:1]=~/123/)]`),
		},
		{
			jsonpath:     `$[?(@[0:2]=~/123/)]`,
			inputJSON:    `[{"b":["123"]},{"a":["123"]}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@[0:2]=~/123/)]`),
		},
		{
			jsonpath:     `$[?(@[0:2].a=~/123/)]`,
			inputJSON:    `[{"b":["123"]},{"a":["123"]}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@[0:2].a=~/123/)]`),
		},
		{
			jsonpath:     `$[?(@.a[0:2]=~/123/)]`,
			inputJSON:    `[{"b":["123"]},{"a":["123"]}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@.a[0:2]=~/123/)]`),
		},
		{
			jsonpath:     `$[?(@[*]=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@[*]=~/123/)]`),
		},
		{
			jsonpath:     `$[?(@[*].a=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@[*].a=~/123/)]`),
		},
		{
			jsonpath:     `$[?(@.a[*]=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@.a[*]=~/123/)]`),
		},
		{
			jsonpath:     `$[?(@[0,1]=~/123/)]`,
			inputJSON:    `[{"b":["123"]},{"a":[123,"123"]}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@[0,1]=~/123/)]`),
		},
		{
			jsonpath:     `$[?(@[0,1:2]=~/123/)]`,
			inputJSON:    `[{"b":["123"]},{"a":[123,"123"]}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@[0,1:2]=~/123/)]`),
		},
		{
			jsonpath:     `$[?(@[0,1].a=~/123/)]`,
			inputJSON:    `[{"b":["123"]},{"a":[123,"123"]}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@[0,1].a=~/123/)]`),
		},
		{
			jsonpath:     `$[?(@.a[0,1]=~/123/)]`,
			inputJSON:    `[{"b":["123"]},{"a":[123,"123"]}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@.a[0,1]=~/123/)]`),
		},
		{
			jsonpath:     `$[?($..a=~/123/)]`,
			inputJSON:    `[{"a":"123"},{"a":123}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `$..a=~/123/)]`),
		},
		{
			jsonpath:     `$[?($..a.b=~/123/)]`,
			inputJSON:    `[{"a":"123"},{"a":123}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `$..a.b=~/123/)]`),
		},
		{
			jsonpath:     `$[?($.a..b=~/123/)]`,
			inputJSON:    `[{"a":"123"},{"a":123}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `$.a..b=~/123/)]`),
		},
		{
			jsonpath:     `$[?($..a..b=~/123/)]`,
			inputJSON:    `[{"a":"123"},{"a":123}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `$..a..b=~/123/)]`),
		},
		{
			jsonpath:     `$[?(@['a','b']=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@['a','b']=~/123/)]`),
		},
		{
			jsonpath:     `$[?(@['a','b','c']=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@['a','b','c']=~/123/)]`),
		},
		{
			jsonpath:     `$[?(@['a','b']['a']=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@['a','b']['a']=~/123/)]`),
		},
		{
			jsonpath:     `$[?(@['a']['a','b']=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@['a']['a','b']=~/123/)]`),
		},
		{
			jsonpath:     `$[?(@.*=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@.*=~/123/)]`),
		},
		{
			jsonpath:     `$[?(@.*[0]=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@.*[0]=~/123/)]`),
		},
		{
			jsonpath:     `$[?(@.*.a=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@.*.a=~/123/)]`),
		},
		{
			jsonpath:     `$[?(@.a.*=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@.a.*=~/123/)]`),
		},
		{
			jsonpath:     `$[?(@.a[?(@.b)]=~/123/)]`,
			inputJSON:    `[{"b":"123"},{"a":"123"}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@.a[?(@.b)]=~/123/)]`),
		},
		{
			jsonpath:     `$[?(@.a[?(@.b>1)]=~/123/)]`,
			inputJSON:    `[{"a":[{"b":1},{"b":2}]},{"a":[{"b":1}]}]`,
			expectedJSON: ``,
			expectedErr:  createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@.a[?(@.b>1)]=~/123/)]`),
		},
	}

	runTestCases(t, "TestFilter_ValueGroupInvalidSyntaxDeletedCases", testCases)
}
