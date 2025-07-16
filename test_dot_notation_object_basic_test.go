package jsonpath

import (
	"testing"
)

func TestDotNotation_RootBasic(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$`,
		inputJSON:    `{"a":"b","c":{"d":"e"}}`,
		expectedJSON: `[{"a":"b","c":{"d":"e"}}]`,
	}
	runTestCase(t, testCase, "TestDotNotation_RootBasic")
}

func TestDotNotation_RootSimpleProperty(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `a`,
		inputJSON:    `{"a":"b","c":{"d":"e"}}`,
		expectedJSON: `["b"]`,
	}
	runTestCase(t, testCase, "TestDotNotation_RootSimpleProperty")
}

func TestDotNotation_RootEscapedAtSymbol(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `\@`,
		inputJSON:    `{"@":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestDotNotation_RootEscapedAtSymbol")
}

func TestDotNotation_ChildSimpleA(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.a`,
		inputJSON:    `{"a":"b","c":{"d":"e"}}`,
		expectedJSON: `["b"]`,
	}
	runTestCase(t, testCase, "TestDotNotation_ChildSimpleA")
}

func TestDotNotation_ChildSimpleC(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.c`,
		inputJSON:    `{"a":"b","c":{"d":"e"}}`,
		expectedJSON: `[{"d":"e"}]`,
	}
	runTestCase(t, testCase, "TestDotNotation_ChildSimpleC")
}

func TestDotNotation_ChildNumericKey(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.2`,
		inputJSON:    `{"a":1,"2":2,"3":{"2":1}}`,
		expectedJSON: `[2]`,
	}
	runTestCase(t, testCase, "TestDotNotation_ChildNumericKey")
}

func TestDotNotation_ChildArrayIndexAndProperty(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0].a`,
		inputJSON:    `[{"a":"b","c":{"d":"e"}},{"a":"y"}]`,
		expectedJSON: `["b"]`,
	}
	runTestCase(t, testCase, "TestDotNotation_ChildArrayIndexAndProperty")
}

func TestDotNotation_ChildArrayWithoutDollar(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `[0].a`,
		inputJSON:    `[{"a":"b","c":{"d":"e"}},{"a":"y"}]`,
		expectedJSON: `["b"]`,
	}
	runTestCase(t, testCase, "TestDotNotation_ChildArrayWithoutDollar")
}

func TestDotNotation_ChildArrayMultiIndex(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[2,0].a`,
		inputJSON:    `[{"a":"b","c":{"a":"d"}},{"a":"e"},{"a":"a"}]`,
		expectedJSON: `["a","b"]`,
	}
	runTestCase(t, testCase, "TestDotNotation_ChildArrayMultiIndex")
}

func TestDotNotation_ChildArraySlice(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0:2].a`,
		inputJSON:    `[{"a":"b","c":{"d":"e"}},{"a":"a"},{"a":"c"}]`,
		expectedJSON: `["b","a"]`,
	}
	runTestCase(t, testCase, "TestDotNotation_ChildArraySlice")
}

func TestDotNotation_ChildNestedProperties(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.a.a2`,
		inputJSON:    `{"a":{"a1":"1","a2":"2"},"b":{"b1":"3"}}`,
		expectedJSON: `["2"]`,
	}
	runTestCase(t, testCase, "TestDotNotation_ChildNestedProperties")
}

func TestDotNotation_ChildWithoutDollarPrefix(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `.a`,
		inputJSON:    `{"a":"b","c":{"d":"e"}}`,
		expectedJSON: `["b"]`,
	}
	runTestCase(t, testCase, "TestDotNotation_ChildWithoutDollarPrefix")
}

func TestDotNotation_SpecialNil(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.nil`,
		inputJSON:    `{"nil":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestDotNotation_SpecialNil")
}

func TestDotNotation_SpecialNull(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.null`,
		inputJSON:    `{"null":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestDotNotation_SpecialNull")
}

func TestDotNotation_SpecialTrue(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.true`,
		inputJSON:    `{"true":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestDotNotation_SpecialTrue")
}

func TestDotNotation_SpecialFalse(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.false`,
		inputJSON:    `{"false":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestDotNotation_SpecialFalse")
}

func TestDotNotation_SpecialIn(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.in`,
		inputJSON:    `{"in":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestDotNotation_SpecialIn")
}

func TestDotNotation_SpecialLengthObject(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.length`,
		inputJSON:    `{"length":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestDotNotation_SpecialLengthObject")
}

func TestDotNotation_SpecialLengthArrayError(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.length`,
		inputJSON:   `["length",1,2]`,
		expectedErr: createErrorTypeUnmatched(`.length`, `object`, `[]interface {}`),
	}
	runTestCase(t, testCase, "TestDotNotation_SpecialLengthArrayError")
}

func TestDotNotation_MixedBracketDotNotation(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$['a'].b`,
		inputJSON:    `{"b":2,"a":{"b":1}}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestDotNotation_MixedBracketDotNotation")
}

func TestDotNotation_escaped_special_characters_part1(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.\ \!\"\#\$\%\&\'\(\)\*\+\,\.\/`,
		inputJSON:    `{" !\"#$%&'()*+,./":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestDotNotation_escaped_special_characters_part1")
}

func TestDotNotation_escaped_special_characters_part2(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.\:\;\<\=\>\?\@\[\\\]\^`,
		inputJSON:    `{":;<=>?@[\\]^":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestDotNotation_escaped_special_characters_part2")
}

func TestDotNotation_escaped_special_characters_part3(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.\{\|\}\~`,
		inputJSON:    `{"{|}~":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestDotNotation_escaped_special_characters_part3")
}

func TestDotNotation_escaped_quoted_property_with_backslash(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.\'a\.b\'`,
		inputJSON:    `{"'a.b'":1,"a":{"b":2},"'a'":{"'b'":3},"'a":{"b'":4}}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestDotNotation_escaped_quoted_property_with_backslash")
}

func TestDotNotation_backtick_property(t *testing.T) {
	testCase := TestCase{
		jsonpath:     "$.\\`",
		inputJSON:    "{\"`\":1}",
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestDotNotation_backtick_property")
}

func TestDotNotation_unicode_property(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.ﾃｽﾄソポァゼゾタダＡボマミ①`,
		inputJSON:    `{"ﾃｽﾄソポァゼゾタダＡボマミ①":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestDotNotation_unicode_property")
}

func TestDotNotation_union_property_chained_missing_property1(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[1,2].a.b`,
		inputJSON:   `[0]`,
		expectedErr: createErrorMemberNotExist(`[1,2]`),
	}
	runTestCase(t, testCase, "TestDotNotation_union_property_chained_missing_property1")
}

func TestDotNotation_union_property_chained_missing_property2(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1].a.b`,
		inputJSON:   `[{"b":1}]`,
		expectedErr: createErrorMemberNotExist(`.a`),
	}
	runTestCase(t, testCase, "TestDotNotation_union_property_chained_missing_property2")
}

func TestDotNotation_union_property_chained_missing_property3(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1].a.b`,
		inputJSON:   `[{"b":1},{"c":2}]`,
		expectedErr: createErrorMemberNotExist(`.a`),
	}
	runTestCase(t, testCase, "TestDotNotation_union_property_chained_missing_property3")
}

func TestDotNotation_union_property_deeply_nested_missing(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1].a.b.c`,
		inputJSON:   `[{"a":1},{"b":2}]`,
		expectedErr: createErrorTypeUnmatched(`.b`, `object`, `float64`),
	}
	runTestCase(t, testCase, "TestDotNotation_union_property_deeply_nested_missing")
}
