package jsonpath

import (
	"testing"
)

func TestRetrieve_dotNotation_root_basic(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$`,
		inputJSON:    `{"a":"b","c":{"d":"e"}}`,
		expectedJSON: `[{"a":"b","c":{"d":"e"}}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_root_basic")
}

func TestRetrieve_dotNotation_root_simple_property(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `a`,
		inputJSON:    `{"a":"b","c":{"d":"e"}}`,
		expectedJSON: `["b"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_root_simple_property")
}

func TestRetrieve_dotNotation_root_escaped_at_symbol(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `\@`,
		inputJSON:    `{"@":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_root_escaped_at_symbol")
}

func TestRetrieve_dotNotation_child_simple_a(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.a`,
		inputJSON:    `{"a":"b","c":{"d":"e"}}`,
		expectedJSON: `["b"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_child_simple_a")
}

func TestRetrieve_dotNotation_child_simple_c(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.c`,
		inputJSON:    `{"a":"b","c":{"d":"e"}}`,
		expectedJSON: `[{"d":"e"}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_child_simple_c")
}

func TestRetrieve_dotNotation_child_numeric_key(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.2`,
		inputJSON:    `{"a":1,"2":2,"3":{"2":1}}`,
		expectedJSON: `[2]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_child_numeric_key")
}

func TestRetrieve_dotNotation_child_array_index_and_property(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0].a`,
		inputJSON:    `[{"a":"b","c":{"d":"e"}},{"a":"y"}]`,
		expectedJSON: `["b"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_child_array_index_and_property")
}

func TestRetrieve_dotNotation_child_array_without_dollar(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `[0].a`,
		inputJSON:    `[{"a":"b","c":{"d":"e"}},{"a":"y"}]`,
		expectedJSON: `["b"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_child_array_without_dollar")
}

func TestRetrieve_dotNotation_child_array_multi_index(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[2,0].a`,
		inputJSON:    `[{"a":"b","c":{"a":"d"}},{"a":"e"},{"a":"a"}]`,
		expectedJSON: `["a","b"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_child_array_multi_index")
}

func TestRetrieve_dotNotation_child_array_slice(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0:2].a`,
		inputJSON:    `[{"a":"b","c":{"d":"e"}},{"a":"a"},{"a":"c"}]`,
		expectedJSON: `["b","a"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_child_array_slice")
}

func TestRetrieve_dotNotation_child_nested_properties(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.a.a2`,
		inputJSON:    `{"a":{"a1":"1","a2":"2"},"b":{"b1":"3"}}`,
		expectedJSON: `["2"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_child_nested_properties")
}

func TestRetrieve_dotNotation_child_without_dollar_prefix(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `.a`,
		inputJSON:    `{"a":"b","c":{"d":"e"}}`,
		expectedJSON: `["b"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_child_without_dollar_prefix")
}

func TestRetrieve_dotNotation_special_nil(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.nil`,
		inputJSON:    `{"nil":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_special_nil")
}

func TestRetrieve_dotNotation_special_null(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.null`,
		inputJSON:    `{"null":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_special_null")
}

func TestRetrieve_dotNotation_special_true(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.true`,
		inputJSON:    `{"true":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_special_true")
}

func TestRetrieve_dotNotation_special_false(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.false`,
		inputJSON:    `{"false":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_special_false")
}

func TestRetrieve_dotNotation_special_in(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.in`,
		inputJSON:    `{"in":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_special_in")
}

func TestRetrieve_dotNotation_special_length_object(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.length`,
		inputJSON:    `{"length":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_special_length_object")
}

func TestRetrieve_dotNotation_special_length_array_error(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.length`,
		inputJSON:   `["length",1,2]`,
		expectedErr: createErrorTypeUnmatched(`.length`, `object`, `[]interface {}`),
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_special_length_array_error")
}

func TestRetrieve_dotNotation_mixed_bracket_dot_notation(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$['a'].b`,
		inputJSON:    `{"b":2,"a":{"b":1}}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_mixed_bracket_dot_notation")
}

func TestRetrieve_dotNotation_escaped_special_characters_part1(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.\ \!\"\#\$\%\&\'\(\)\*\+\,\.\/`,
		inputJSON:    `{" !\"#$%&'()*+,./":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_escaped_special_characters_part1")
}

func TestRetrieve_dotNotation_escaped_special_characters_part2(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.\:\;\<\=\>\?\@\[\\\]\^`,
		inputJSON:    `{":;<=>?@[\\]^":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_escaped_special_characters_part2")
}

func TestRetrieve_dotNotation_escaped_special_characters_part3(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.\{\|\}\~`,
		inputJSON:    `{"{|}~":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_escaped_special_characters_part3")
}

func TestRetrieve_dotNotation_escaped_quoted_property_with_backslash(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.\'a\.b\'`,
		inputJSON:    `{"'a.b'":1,"a":{"b":2},"'a'":{"'b'":3},"'a":{"b'":4}}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_escaped_quoted_property_with_backslash")
}

func TestRetrieve_dotNotation_backtick_property(t *testing.T) {
	testCase := TestCase{
		jsonpath:     "$.\\`",
		inputJSON:    "{\"`\":1}",
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_backtick_property")
}

func TestRetrieve_dotNotation_unicode_property(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.ﾃｽﾄソポァゼゾタダＡボマミ①`,
		inputJSON:    `{"ﾃｽﾄソポァゼゾタダＡボマミ①":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_unicode_property")
}

func TestRetrieve_dotNotation_union_property_chained_missing_property1(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[1,2].a.b`,
		inputJSON:   `[0]`,
		expectedErr: createErrorMemberNotExist(`[1,2]`),
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_union_property_chained_missing_property1")
}

func TestRetrieve_dotNotation_union_property_chained_missing_property2(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1].a.b`,
		inputJSON:   `[{"b":1}]`,
		expectedErr: createErrorMemberNotExist(`.a`),
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_union_property_chained_missing_property2")
}

func TestRetrieve_dotNotation_union_property_chained_missing_property3(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1].a.b`,
		inputJSON:   `[{"b":1},{"c":2}]`,
		expectedErr: createErrorMemberNotExist(`.a`),
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_union_property_chained_missing_property3")
}

func TestRetrieve_dotNotation_union_property_deeply_nested_missing(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0,1].a.b.c`,
		inputJSON:   `[{"a":1},{"b":2}]`,
		expectedErr: createErrorTypeUnmatched(`.b`, `object`, `float64`),
	}
	runTestCase(t, testCase, "TestRetrieve_dotNotation_union_property_deeply_nested_missing")
}
