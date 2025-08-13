package tests

import (
	"testing"
)

func TestFilterComparison_NumericComparisons(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.a < 1)]`,
			inputJSON:    `[{"a":-9999999},{"a":0.999999},{"a":1.0000000},{"a":1.0000001},{"a":2},{"a":"0.9"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
			expectedJSON: `[{"a":-9999999},{"a":0.999999}]`,
		},
		{
			jsonpath:    `$[?(1 > @.a)]`,
			inputJSON:   `[{"a":1.0000000},{"a":1.0000001},{"a":2},{"a":"0.9"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
			expectedErr: createErrorMemberNotExist(`[?(1 > @.a)]`),
		},
		{
			jsonpath:     `$[?(@.a <= 1.00001)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":1.00001},{"a":1.00002},{"a":2,"b":4},{"a":"0.9"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
			expectedJSON: `[{"a":0},{"a":1},{"a":1.00001}]`,
		},
		{
			jsonpath:     `$[?(1.00001 >= @.a)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":1.00001},{"a":1.00002},{"a":2,"b":4},{"a":"0.9"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
			expectedJSON: `[{"a":0},{"a":1},{"a":1.00001}]`,
		},
		{
			jsonpath:     `$[?(@.a > 1)]`,
			inputJSON:    `[{"a":0},{"a":0.9999},{"a":1},{"a":1.000001},{"a":2,"b":4},{"a":9999999999},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
			expectedJSON: `[{"a":1.000001},{"a":2,"b":4},{"a":9999999999}]`,
		},
		{
			jsonpath:     `$[?(1 < @.a)]`,
			inputJSON:    `[{"a":0},{"a":0.9999},{"a":1},{"a":1.000001},{"a":2,"b":4},{"a":9999999999},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
			expectedJSON: `[{"a":1.000001},{"a":2,"b":4},{"a":9999999999}]`,
		},
		{
			jsonpath:    `$[?(1 < @.a)]`,
			inputJSON:   `[{"a":0},{"a":0.9999},{"a":1},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
			expectedErr: createErrorMemberNotExist(`[?(1 < @.a)]`),
		},
		{
			jsonpath:     `$[?(@.a >= 1.000001)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":1.000001},{"a":1.0000009},{"a":1.001},{"a":2,"b":4},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
			expectedJSON: `[{"a":1.000001},{"a":1.001},{"a":2,"b":4}]`,
		},
		{
			jsonpath:     `$[?(1.000001 <= @.a)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":1.000001},{"a":1.0000009},{"a":1.001},{"a":2,"b":4},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
			expectedJSON: `[{"a":1.000001},{"a":1.001},{"a":2,"b":4}]`,
		},
	}

	runTestCases(t, "TestFilterComparison_NumericComparisons", tests)
}

func TestFilterComparison_SpecialCharacters(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a=='a\rb')]`,
			inputJSON:    `[{"a":"a\rb"},{"b":1}]`,
			expectedJSON: `[{"a":"a\rb"}]`,
		},
		{
			jsonpath:     `$[?(@.a=='a\tb')]`,
			inputJSON:    `[{"a":"a\tb"},{"b":1}]`,
			expectedJSON: `[{"a":"a\tb"}]`,
		},
		{
			jsonpath:     `$[?(@.a=='\u0000')]`,
			inputJSON:    `[{"a":"\u0000"},{"b":1}]`,
			expectedJSON: `[{"a":"\u0000"}]`,
		},
		{
			jsonpath:     `$[?(@.a=='\uABCD')]`,
			inputJSON:    `[{"a":"\uabcd"},{"b":1}]`,
			expectedJSON: `[{"a":"ꯍ"}]`,
		},
		{
			jsonpath:     `$[?(@.a=='\uabcd')]`,
			inputJSON:    `[{"a":"\uABCD"},{"b":1}]`,
			expectedJSON: `[{"a":"ꯍ"}]`,
		},
	}

	runTestCases(t, "TestFilterComparison_SpecialCharacters", testCases)
}

func TestFilterComparison_NumericPrecisionComparisons(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a==5)]`,
			inputJSON:    `[{"a":4.9},{"a":5.0},{"a":5.1},{"a":5},{"a":-5},{"a":"5"},{"a":"a"},{"a":true},{"a":null},{"a":{}},{"a":[]},{"b":5},{"a":{"a":5}},{"a":[{"a":5}]}]`,
			expectedJSON: `[{"a":5},{"a":5}]`,
		},
		{
			jsonpath:     `$[?(@==5)]`,
			inputJSON:    `[4.999999,5.00000,5.00001,5,-5,"5","a",null,{},[],{"a":5},[5]]`,
			expectedJSON: `[5,5]`,
		},
		{
			jsonpath:    `$[?(@.a==5)]`,
			inputJSON:   `[{"a":4.9},{"a":5.1},{"a":-5},{"a":"5"},{"a":"a"},{"a":true},{"a":null},{"a":{}},{"a":[]},{"b":5},{"a":{"a":5}},{"a":[{"a":5}]}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.a==5)]`),
		},
		{
			jsonpath:     `$[?(@.a==1)]`,
			inputJSON:    `{"a":{"a":0.999999},"b":{"a":1.0},"c":{"a":1.00001},"d":{"a":1},"e":{"a":-1},"f":{"a":"1"},"g":{"a":[1]}}`,
			expectedJSON: `[{"a":1},{"a":1}]`,
		},
		{
			jsonpath:    `$[?(@.a==1)]`,
			inputJSON:   `{"a":1}`,
			expectedErr: createErrorMemberNotExist(`[?(@.a==1)]`),
		},
		{
			jsonpath:     `$[?(@.a==-0.123E2)]`,
			inputJSON:    `[{"a":-12.3}]`,
			expectedJSON: `[{"a":-12.3}]`,
		},
		{
			jsonpath:     `$[?(@.a==+0.123e+2)]`,
			inputJSON:    `[{"a":-12.3},{"a":12.3}]`,
			expectedJSON: `[{"a":12.3}]`,
		},
		{
			jsonpath:     `$[?(@.a==-1.23e-1)]`,
			inputJSON:    `[{"a":-12.3},{"a":-1.23},{"a":-0.123}]`,
			expectedJSON: `[{"a":-0.123}]`,
		},
		{
			jsonpath:     `$[?(@.a==null)]`,
			inputJSON:    `[{"a":null},{"a":false},{"a":true},{"a":0},{"a":1},{"a":"false"},{"b":null}]`,
			expectedJSON: `[{"a":null}]`,
		},
	}

	runTestCases(t, "TestFilterComparison_NumericPrecisionComparisons", testCases)
}

func TestFilterComparison_RootPathComparisons(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$.a[?(@.b==$.c)]`,
			inputJSON:    `{"a":[{"b":123},{"b":123.456},{"b":"123.456"}],"c":123.456}`,
			expectedJSON: `[{"b":123.456}]`,
		},
		{
			jsonpath:    `$[?(@.a < $.b)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.a < $.b)]`),
		},
		{
			jsonpath:    `$[?(@.a > $.b)]`,
			inputJSON:   `[{"a":0},{"a":1}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.a > $.b)]`),
		},
	}

	runTestCases(t, "TestFilterComparison_RootPathComparisons", testCases)
}

func TestFilterComparison_ArrayIndexComparisons(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@[1]=="b")]`,
			inputJSON:    `{"a":["a","b"],"b":["b"]}`,
			expectedJSON: `[["a","b"]]`,
		},
	}

	runTestCases(t, "TestFilterComparison_ArrayIndexComparisons", testCases)
}

func TestFilterComparison_ArrayIndexWithQuotes(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@[1]=="b")]`,
			inputJSON:    `[[0,1],[0,2],[2],["2"],["a","b"],["b"]]`,
			expectedJSON: `[["a","b"]]`,
		},
		{
			jsonpath:     `$[?(@[1]=="a\"b")]`,
			inputJSON:    `[[0,1],[2],["a","a\"b"],["a\"b"]]`,
			expectedJSON: `[["a","a\"b"]]`,
		},
		{
			jsonpath:     `$[?(@[1]=='b')]`,
			inputJSON:    `[[0,1],[2],["a","b"],["b"]]`,
			expectedJSON: `[["a","b"]]`,
		},
		{
			jsonpath:     `$[?(@[1]=='a\'b')]`,
			inputJSON:    `[[0,1],[2],["a","a'b"],["a'b"]]`,
			expectedJSON: `[["a","a'b"]]`,
		},
	}

	runTestCases(t, "TestFilterComparison_ArrayIndexWithQuotes", testCases)
}

func TestFilterComparison_EscapeSequences(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a=='a\'b')]`,
			inputJSON:    `[{"a":"a'b"},{"b":1}]`,
			expectedJSON: `[{"a":"a'b"}]`,
		},
		{
			jsonpath:    `$[?(@.a=='a\b')]`,
			inputJSON:   `[{"a":"ab"}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.a=='a\b')]`),
		},
		{
			jsonpath:    `$[?(@.a=="a\b")]`,
			inputJSON:   `[{"a":"ab"}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.a=="a\b")]`),
		},
		{
			jsonpath:     `$[?(@.a=='a\/b')]`,
			inputJSON:    `[{"a":"a\/b"},{"b":1}]`,
			expectedJSON: `[{"a":"a/b"}]`,
		},
		{
			jsonpath:     `$[?(@.a=='a\\b')]`,
			inputJSON:    `[{"a":"a\\b"},{"b":1}]`,
			expectedJSON: `[{"a":"a\\b"}]`,
		},
		{
			jsonpath:     `$[?(@.a=='a\fb')]`,
			inputJSON:    `[{"a":"a\fb"},{"b":1}]`,
			expectedJSON: `[{"a":"a\fb"}]`,
		},
		{
			jsonpath:     `$[?(@.a=='a\nb')]`,
			inputJSON:    `[{"a":"a\nb"},{"b":1}]`,
			expectedJSON: `[{"a":"a\nb"}]`,
		},
	}

	runTestCases(t, "TestFilterComparison_EscapeSequences", testCases)
}

func TestFilterComparison_InvalidSyntaxCases(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$[?(@.a == $.b)]`,
			inputJSON:   `[{"a":1},{"a":2}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.a == $.b)]`),
		},
		{
			jsonpath:    `$[?($.b == @.a)]`,
			inputJSON:   `[{"a":1},{"a":2}]`,
			expectedErr: createErrorMemberNotExist(`[?($.b == @.a)]`),
		},
	}

	runTestCases(t, "TestFilterComparison_InvalidSyntaxCases", tests)
}

func TestFilterComparison_InvalidSyntaxRegexCases(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$[?(@[0:1]==1)]`,
			inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
			expectedErr: createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@[0:1]==1)]`),
		},
		{
			jsonpath:    `$[?(@[0:2]==1)]`,
			inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
			expectedErr: createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@[0:2]==1)]`),
		},
		{
			jsonpath:    `$[?(@[0:2].a==1)]`,
			inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
			expectedErr: createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@[0:2].a==1)]`),
		},
		{
			jsonpath:    `$[?(@.a[0:2]==1)]`,
			inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
			expectedErr: createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@.a[0:2]==1)]`),
		},
		{
			jsonpath:    `$[?(@[0,1]==1)]`,
			inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
			expectedErr: createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@[0,1]==1)]`),
		},
		{
			jsonpath:    `$[?(@[0,1].a==1)]`,
			inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
			expectedErr: createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@[0,1].a==1)]`),
		},
		{
			jsonpath:    `$[?(@.a[0,1]==1)]`,
			inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
			expectedErr: createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@.a[0,1]==1)]`),
		},
		{
			jsonpath:    `$[?(@..a==123)]`,
			inputJSON:   `[{"a":"123"},{"a":123}]`,
			expectedErr: createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@..a==123)]`),
		},
		{
			jsonpath:    `$[?(@..a.b==123)]`,
			inputJSON:   `[{"a":"123"},{"a":123}]`,
			expectedErr: createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@..a.b==123)]`),
		},
		{
			jsonpath:    `$[?(@.a..b==123)]`,
			inputJSON:   `[{"a":"123"},{"a":123}]`,
			expectedErr: createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@.a..b==123)]`),
		},
		{
			jsonpath:    `$[?(@..a..b==123)]`,
			inputJSON:   `[{"a":"123"},{"a":123}]`,
			expectedErr: createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@..a..b==123)]`),
		},
		{
			jsonpath:    `$[?(@['a','b']==123)]`,
			inputJSON:   `[{"a":"123"},{"a":123}]`,
			expectedErr: createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@['a','b']==123)]`),
		},
		{
			jsonpath:    `$[?(@['a','b','c']==123)]`,
			inputJSON:   `[{"a":"123"},{"a":123}]`,
			expectedErr: createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@['a','b','c']==123)]`),
		},
		{
			jsonpath:    `$[?(@['a','b']['a']==123)]`,
			inputJSON:   `[{"a":"123"},{"a":123}]`,
			expectedErr: createErrorInvalidSyntax(4, `JSONPath that returns a value group is prohibited`, `@['a','b']['a']==123)]`),
		},
	}

	runTestCases(t, "TestFilterComparison_InvalidSyntaxRegexCases", tests)
}
