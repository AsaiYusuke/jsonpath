package jsonpath

import (
	"testing"
)

func TestRetrieve_space_basic_space(t *testing.T) {
	testCase := TestCase{
		jsonpath:     ` $.a `,
		inputJSON:    `{"a":123}`,
		expectedJSON: `[123]`,
	}
	runTestCase(t, testCase, "TestRetrieve_space_basic_space")
}

func TestRetrieve_space_tab_invalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:    "\t" + `$.a` + "\t",
		inputJSON:   `{"a":123}`,
		expectedErr: ErrorInvalidSyntax{position: 0, reason: `unrecognized input`, near: "\t" + `$.a` + "\t"},
	}
	runTestCase(t, testCase, "TestRetrieve_space_tab_invalid")
}

func TestRetrieve_space_newline_invalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.a` + "\n",
		inputJSON:   `{"a":123}`,
		expectedErr: ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: "\n"},
	}
	runTestCase(t, testCase, "TestRetrieve_space_newline_invalid")
}

func TestRetrieve_space_bracket_notation_multi_identifier(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ "a" , "c" ]`,
		inputJSON:    `{"a":1,"b":2,"c":3}`,
		expectedJSON: `[1,3]`,
	}
	runTestCase(t, testCase, "TestRetrieve_space_bracket_notation_multi_identifier")
}

func TestRetrieve_space_bracket_notation_slice_union(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ 0 , 2 : 4 , * ]`,
		inputJSON:    `[1,2,3,4,5]`,
		expectedJSON: `[1,3,4,1,2,3,4,5]`,
	}
	runTestCase(t, testCase, "TestRetrieve_space_bracket_notation_slice_union")
}

func TestRetrieve_space_filter_equal(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( @.a == 1 ) ]`,
		inputJSON:    `[{"a":1}]`,
		expectedJSON: `[{"a":1}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_space_filter_equal")
}

func TestRetrieve_space_filter_not_equal(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( @.a != 1 ) ]`,
		inputJSON:    `[{"a":2}]`,
		expectedJSON: `[{"a":2}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_space_filter_not_equal")
}

func TestRetrieve_space_filter_less_equal(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( @.a <= 1 ) ]`,
		inputJSON:    `[{"a":1}]`,
		expectedJSON: `[{"a":1}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_space_filter_less_equal")
}

func TestRetrieve_space_filter_less_than(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( @.a < 1 ) ]`,
		inputJSON:    `[{"a":0}]`,
		expectedJSON: `[{"a":0}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_space_filter_less_than")
}

func TestRetrieve_space_filter_greater_equal(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( @.a >= 1 ) ]`,
		inputJSON:    `[{"a":1}]`,
		expectedJSON: `[{"a":1}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_space_filter_greater_equal")
}

func TestRetrieve_space_filter_greater_than(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( @.a > 1 ) ]`,
		inputJSON:    `[{"a":2}]`,
		expectedJSON: `[{"a":2}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_space_filter_greater_than")
}

func TestRetrieve_space_filter_regex(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( @.a =~ /a/ ) ]`,
		inputJSON:    `[{"a":"abc"}]`,
		expectedJSON: `[{"a":"abc"}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_space_filter_regex")
}

func TestRetrieve_space_filter_logical_and(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( @.a == 1 && @.b == 2 ) ]`,
		inputJSON:    `[{"a":1,"b":2}]`,
		expectedJSON: `[{"a":1,"b":2}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_space_filter_logical_and")
}

func TestRetrieve_space_filter_logical_or(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( @.a == 1 || @.b == 2 ) ]`,
		inputJSON:    `[{"a":1},{"b":2}]`,
		expectedJSON: `[{"a":1},{"b":2}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_space_filter_logical_or")
}

func TestRetrieve_space_filter_logical_not(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( ! @.a ) ]`,
		inputJSON:    `[{"a":1},{"b":2}]`,
		expectedJSON: `[{"b":2}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_space_filter_logical_not")
}

func TestRetrieve_space_filter_parentheses(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( ( @.a ) ) ]`,
		inputJSON:    `[{"a":1},{"b":2}]`,
		expectedJSON: `[{"a":1}]`,
	}
	runTestCase(t, testCase, "TestRetrieve_space_filter_parentheses")
}

// Control character tests
func TestRetrieve_space_tab_character_invalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.` + "\x09",
		inputJSON:   `{"\t":1}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.` + "\x09"},
	}
	runTestCase(t, testCase, "TestRetrieve_space_tab_character_invalid")
}

func TestRetrieve_space_escaped_tab_character_invalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.\` + "\x09",
		inputJSON:   `{"\t":1}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.\` + "\x09"},
	}
	runTestCase(t, testCase, "TestRetrieve_space_escaped_tab_character_invalid")
}

func TestRetrieve_space_cr_character_invalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.` + "\x0d",
		inputJSON:   `{"\n":1}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.` + "\x0d"},
	}
	runTestCase(t, testCase, "TestRetrieve_space_cr_character_invalid")
}

func TestRetrieve_space_escaped_cr_character_invalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.\` + "\x0d",
		inputJSON:   `{"\n":1}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.\` + "\x0d"},
	}
	runTestCase(t, testCase, "TestRetrieve_space_escaped_cr_character_invalid")
}

func TestRetrieve_space_control_character_1f_invalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.` + "\x1f",
		inputJSON:   `{"a":1}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.` + "\x1f"},
	}
	runTestCase(t, testCase, "TestRetrieve_space_control_character_1f_invalid")
}

func TestRetrieve_space_del_character_invalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.` + "\x7F",
		inputJSON:   `{"` + "\x7F" + `":1}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.` + "\x7F"},
	}
	runTestCase(t, testCase, "TestRetrieve_space_del_character_invalid")
}

func TestRetrieve_space_escaped_del_character_invalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.\` + "\x7F",
		inputJSON:   `{"` + "\x7F" + `":1}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.\` + "\x7F"},
	}
	runTestCase(t, testCase, "TestRetrieve_space_escaped_del_character_invalid")
}

func TestRetrieve_space_property_with_del_character_invalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.a` + "\x7F",
		inputJSON:   `{"a` + "\x7F" + `":1}`,
		expectedErr: ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: "\x7F"},
	}
	runTestCase(t, testCase, "TestRetrieve_space_property_with_del_character_invalid")
}
