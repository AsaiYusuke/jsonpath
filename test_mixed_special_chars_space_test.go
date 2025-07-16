package jsonpath

import (
	"testing"
)

func TestSpecialChars_BasicSpaceHandling(t *testing.T) {
	testCase := TestCase{
		jsonpath:     ` $.a `,
		inputJSON:    `{"a":123}`,
		expectedJSON: `[123]`,
	}
	runTestCase(t, testCase, "TestSpecialChars_basic_space")
}

func TestSpecialChars_InvalidTabCharacters(t *testing.T) {
	testCase := TestCase{
		jsonpath:    "\t" + `$.a` + "\t",
		inputJSON:   `{"a":123}`,
		expectedErr: ErrorInvalidSyntax{position: 0, reason: `unrecognized input`, near: "\t" + `$.a` + "\t"},
	}
	runTestCase(t, testCase, "TestSpecialChars_tab_invalid")
}

func TestSpecialChars_InvalidNewlineCharacters(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.a` + "\n",
		inputJSON:   `{"a":123}`,
		expectedErr: ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: "\n"},
	}
	runTestCase(t, testCase, "TestSpecialChars_newline_invalid")
}

func TestSpecialChars_BracketNotationMultiIdentifier(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ "a" , "c" ]`,
		inputJSON:    `{"a":1,"b":2,"c":3}`,
		expectedJSON: `[1,3]`,
	}
	runTestCase(t, testCase, "TestSpecialChars_BracketNotationMultiIdentifier")
}

func TestSpecialChars_BracketNotationSliceUnion(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ 0 , 2 : 4 , * ]`,
		inputJSON:    `[1,2,3,4,5]`,
		expectedJSON: `[1,3,4,1,2,3,4,5]`,
	}
	runTestCase(t, testCase, "TestSpecialChars_BracketNotationSliceUnion")
}

func TestSpecialChars_FilterEqual(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( @.a == 1 ) ]`,
		inputJSON:    `[{"a":1}]`,
		expectedJSON: `[{"a":1}]`,
	}
	runTestCase(t, testCase, "TestSpecialChars_FilterEqual")
}

func TestSpecialChars_FilterNotEqual(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( @.a != 1 ) ]`,
		inputJSON:    `[{"a":2}]`,
		expectedJSON: `[{"a":2}]`,
	}
	runTestCase(t, testCase, "TestSpecialChars_FilterNotEqual")
}

func TestSpecialChars_FilterLessEqual(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( @.a <= 1 ) ]`,
		inputJSON:    `[{"a":1}]`,
		expectedJSON: `[{"a":1}]`,
	}
	runTestCase(t, testCase, "TestSpecialChars_FilterLessEqual")
}

func TestSpecialChars_FilterLessThan(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( @.a < 1 ) ]`,
		inputJSON:    `[{"a":0}]`,
		expectedJSON: `[{"a":0}]`,
	}
	runTestCase(t, testCase, "TestSpecialChars_FilterLessThan")
}

func TestSpecialChars_FilterGreaterEqual(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( @.a >= 1 ) ]`,
		inputJSON:    `[{"a":1}]`,
		expectedJSON: `[{"a":1}]`,
	}
	runTestCase(t, testCase, "TestSpecialChars_FilterGreaterEqual")
}

func TestSpecialChars_FilterGreaterThan(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( @.a > 1 ) ]`,
		inputJSON:    `[{"a":2}]`,
		expectedJSON: `[{"a":2}]`,
	}
	runTestCase(t, testCase, "TestSpecialChars_FilterGreaterThan")
}

func TestSpecialChars_FilterRegex(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( @.a =~ /a/ ) ]`,
		inputJSON:    `[{"a":"abc"}]`,
		expectedJSON: `[{"a":"abc"}]`,
	}
	runTestCase(t, testCase, "TestSpecialChars_FilterRegex")
}

func TestSpecialChars_FilterLogicalAnd(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( @.a == 1 && @.b == 2 ) ]`,
		inputJSON:    `[{"a":1,"b":2}]`,
		expectedJSON: `[{"a":1,"b":2}]`,
	}
	runTestCase(t, testCase, "TestSpecialChars_FilterLogicalAnd")
}

func TestSpecialChars_FilterLogicalOr(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( @.a == 1 || @.b == 2 ) ]`,
		inputJSON:    `[{"a":1},{"b":2}]`,
		expectedJSON: `[{"a":1},{"b":2}]`,
	}
	runTestCase(t, testCase, "TestSpecialChars_FilterLogicalOr")
}

func TestSpecialChars_FilterLogicalNot(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( ! @.a ) ]`,
		inputJSON:    `[{"a":1},{"b":2}]`,
		expectedJSON: `[{"b":2}]`,
	}
	runTestCase(t, testCase, "TestSpecialChars_FilterLogicalNot")
}

func TestSpecialChars_FilterParentheses(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[ ?( ( @.a ) ) ]`,
		inputJSON:    `[{"a":1},{"b":2}]`,
		expectedJSON: `[{"a":1}]`,
	}
	runTestCase(t, testCase, "TestSpecialChars_FilterParentheses")
}

func TestSpecialChars_TabCharacterInvalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.` + "\x09",
		inputJSON:   `{"\t":1}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.` + "\x09"},
	}
	runTestCase(t, testCase, "TestSpecialChars_TabCharacterInvalid")
}

func TestSpecialChars_EscapedTabCharacterInvalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.\` + "\x09",
		inputJSON:   `{"\t":1}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.\` + "\x09"},
	}
	runTestCase(t, testCase, "TestSpecialChars_EscapedTabCharacterInvalid")
}

func TestSpecialChars_CrCharacterInvalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.` + "\x0d",
		inputJSON:   `{"\n":1}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.` + "\x0d"},
	}
	runTestCase(t, testCase, "TestSpecialChars_CrCharacterInvalid")
}

func TestSpecialChars_EscapedCrCharacterInvalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.\` + "\x0d",
		inputJSON:   `{"\n":1}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.\` + "\x0d"},
	}
	runTestCase(t, testCase, "TestSpecialChars_EscapedCrCharacterInvalid")
}

func TestSpecialChars_ControlCharacter1fInvalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.` + "\x1f",
		inputJSON:   `{"a":1}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.` + "\x1f"},
	}
	runTestCase(t, testCase, "TestSpecialChars_ControlCharacter1fInvalid")
}

func TestSpecialChars_DelCharacterInvalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.` + "\x7F",
		inputJSON:   `{"` + "\x7F" + `":1}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.` + "\x7F"},
	}
	runTestCase(t, testCase, "TestSpecialChars_DelCharacterInvalid")
}

func TestSpecialChars_EscapedDelCharacterInvalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.\` + "\x7F",
		inputJSON:   `{"` + "\x7F" + `":1}`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.\` + "\x7F"},
	}
	runTestCase(t, testCase, "TestSpecialChars_EscapedDelCharacterInvalid")
}

func TestSpecialChars_PropertyWithDelCharacterInvalid(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.a` + "\x7F",
		inputJSON:   `{"a` + "\x7F" + `":1}`,
		expectedErr: ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: "\x7F"},
	}
	runTestCase(t, testCase, "TestSpecialChars_PropertyWithDelCharacterInvalid")
}
