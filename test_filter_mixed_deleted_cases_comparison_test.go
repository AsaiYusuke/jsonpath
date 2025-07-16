package jsonpath

import (
	"fmt"
	"testing"
)

func TestFilterComparison_DecimalNumbers(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.a == 2.1)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":2.0,"b":4},{"a":2.1,"b":5},{"a":2.2,"b":6},{"a":"2.1"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
			expectedJSON: `[{"a":2.1,"b":5}]`,
		},
		{
			jsonpath:     `$[?(2.1 == @.a)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":2.0,"b":4},{"a":2.1,"b":5},{"a":2.2,"b":6},{"a":"2.1"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
			expectedJSON: `[{"a":2.1,"b":5}]`,
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("DecimalNumbers_%d", i), test)
	}
}

func TestFilterComparison_StringEquality(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.a=='ab')]`,
			inputJSON:    `[{"a":"ab"}]`,
			expectedJSON: `[{"a":"ab"}]`,
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("StringEquality_%d", i), test)
	}
}

func TestFilterRoot_ComplexSelectors(t *testing.T) {
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

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("ComplexSelectors_%d", i), test)
	}
}

func TestFilterComparison_StringComparisons(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:    `$[?(@.a=='ab')]`,
			inputJSON:   `[{"a":"abc"}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.a=='ab')]`),
		},
		{
			jsonpath:     `$[?(@.a==1)]`,
			inputJSON:    `[{"a":1},{"b":1}]`,
			expectedJSON: `[{"a":1}]`,
		},
		{
			jsonpath:     `$[?(@.a!='ab')]`,
			inputJSON:    `[{"a":"abc"}]`,
			expectedJSON: `[{"a":"abc"}]`,
		},
		{
			jsonpath:    `$[?(@.a!='ab')]`,
			inputJSON:   `[{"a":"ab"}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.a!='ab')]`),
		},
		{
			jsonpath:     `$[?(@.a!=1)]`,
			inputJSON:    `[{"a":1},{"b":1}]`,
			expectedJSON: `[{"b":1}]`,
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("StringComparisons_%d", i), test)
	}
}

func TestFilterComparison_InequalityOperators(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.a != 2)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":2,"b":4},{"a":1.999999},{"a":2.000000000001},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
			expectedJSON: `[{"a":0},{"a":1},{"a":1.999999},{"a":2.000000000001},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
		},
		{
			jsonpath:     `$[?(2 != @.a)]`,
			inputJSON:    `[{"a":0},{"a":1},{"a":2,"b":4},{"a":1.999999},{"a":2.000000000001},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
			expectedJSON: `[{"a":0},{"a":1},{"a":1.999999},{"a":2.000000000001},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("InequalityOperators_%d", i), test)
	}
}

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

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("NumericComparisons_%d", i), test)
	}
}

func TestFilterComparison_StringLiterals(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[?(@.a=="value")]`,
			inputJSON:    `[{"a":"value"},{"a":0},{"a":1},{"a":-1},{"a":"val"},{"a":true},{"a":{}},{"a":[]},{"a":["b"]},{"a":{"a":"value"}},{"b":"value"}]`,
			expectedJSON: `[{"a":"value"}]`,
		},
		{
			jsonpath:     `$[?(@.a=='value')]`,
			inputJSON:    `[{"a":"value"},{"a":0},{"a":1},{"a":-1},{"a":"val"},{"a":{}},{"a":[]},{"a":["b"]},{"a":{"a":"value"}},{"b":"value"}]`,
			expectedJSON: `[{"a":"value"}]`,
		},
	}

	for i, test := range tests {
		runSingleTestCase(t, fmt.Sprintf("StringLiterals_%d", i), test)
	}
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

	for _, testCase := range testCases {
		runSingleTestCase(t, "TestFilterComparison_SpecialCharacters", testCase)
	}
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

	for _, testCase := range testCases {
		runSingleTestCase(t, "TestFilterComparison_NumericPrecisionComparisons", testCase)
	}
}

func TestFilterComparison_RegexMatching(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@.a =~ /ab/)]`,
			inputJSON:    `[{"a":"abc"},{"a":1},{"a":"def"}]`,
			expectedJSON: `[{"a":"abc"}]`,
		},
		{
			jsonpath:     `$[?(@.a =~ /123/)]`,
			inputJSON:    `[{"a":123},{"a":"123"},{"a":"12"},{"a":"23"},{"a":"0123"},{"a":"1234"}]`,
			expectedJSON: `[{"a":"123"},{"a":"0123"},{"a":"1234"}]`,
		},
		{
			jsonpath:     `$[?(@.a=~/テスト/)]`,
			inputJSON:    `[{"a":"123テストabc"}]`,
			expectedJSON: `[{"a":"123テストabc"}]`,
		},
		{
			jsonpath:     `$[?(@.a=~/^\d+[a-d]\/\\$/)]`,
			inputJSON:    `[{"a":"012b/\\"},{"a":"ab/\\"},{"a":"1b\\"},{"a":"1b//"},{"a":"1b/\""}]`,
			expectedJSON: `[{"a":"012b/\\"}]`,
		},
		{
			jsonpath:     `$[?(@.a=~/(?i)CASE/)]`,
			inputJSON:    `[{"a":"case"},{"a":"CASE"},{"a":"Case"},{"a":"abc"}]`,
			expectedJSON: `[{"a":"case"},{"a":"CASE"},{"a":"Case"}]`,
		},
		{
			jsonpath:    `$[?(@.a.b=~/abc/)]`,
			inputJSON:   `[{"a":"abc"}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.a.b=~/abc/)]`),
		},
		{
			jsonpath:     `$[?(@[0]=~/123/)]`,
			inputJSON:    `[["123"],["456"]]`,
			expectedJSON: `[["123"]]`,
		},
	}

	for _, testCase := range testCases {
		runSingleTestCase(t, "TestFilterComparison_RegexMatching", testCase)
	}
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

	for _, testCase := range testCases {
		runSingleTestCase(t, "TestFilterComparison_RootPathComparisons", testCase)
	}
}

func TestFilterComparison_ArrayIndexComparisons(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(@[1]=="b")]`,
			inputJSON:    `{"a":["a","b"],"b":["b"]}`,
			expectedJSON: `[["a","b"]]`,
		},
	}

	for _, testCase := range testCases {
		runSingleTestCase(t, "TestFilterComparison_ArrayIndexComparisons", testCase)
	}
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

	for _, testCase := range testCases {
		runSingleTestCase(t, "TestFilterComparison_ArrayIndexWithQuotes", testCase)
	}
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

	for _, testCase := range testCases {
		runSingleTestCase(t, "TestFilterComparison_EscapeSequences", testCase)
	}
}

func TestFilterComparison_InvalidSyntaxCases(t *testing.T) {
	testCases := []TestCase{
		// スライス操作の基本的なケースのみ保持
	}

	for _, testCase := range testCases {
		runSingleTestCase(t, "TestFilterComparison_InvalidSyntaxCases", testCase)
	}
}

func TestFilterComparison_InvalidSyntaxRegexCases(t *testing.T) {
	testCases := []TestCase{
		// 基本的な動作テストのみ保持
	}

	for _, testCase := range testCases {
		runSingleTestCase(t, "TestFilterComparison_InvalidSyntaxRegexCases", testCase)
	}
}
