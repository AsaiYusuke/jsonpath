package jsonpath

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type TestGroup map[string][]TestCase

type TestCase struct {
	jsonpath        string
	inputJSON       string
	expectedJSON    string
	expectedErr     error
	unmarshalFunc   func(string, *interface{}) error
	filters         map[string]func(interface{}) (interface{}, error)
	aggregates      map[string]func([]interface{}) (interface{}, error)
	accessorMode    bool
	resultValidator func(interface{}, []interface{}) error
}

func createErrorMemberNotExist(text string) ErrorMemberNotExist {
	return ErrorMemberNotExist{
		errorBasicRuntime: &errorBasicRuntime{
			node: &syntaxBasicNode{
				text: text,
			},
		},
	}
}

func createErrorTypeUnmatched(text string, expected string, found string) ErrorTypeUnmatched {
	return ErrorTypeUnmatched{
		errorBasicRuntime: &errorBasicRuntime{
			node: &syntaxBasicNode{
				text: text,
			},
		},
		expectedType: expected,
		foundType:    found,
	}
}

func createErrorFunctionFailed(text string, errorString string) ErrorFunctionFailed {
	return ErrorFunctionFailed{
		errorBasicRuntime: &errorBasicRuntime{
			node: &syntaxBasicNode{
				text: text,
			},
		},
		err: fmt.Errorf(`%s`, errorString),
	}
}

func execTestRetrieve(t *testing.T, inputJSON interface{}, testCase TestCase, fileLine string) ([]interface{}, error) {
	jsonPath := testCase.jsonpath
	hasConfig := false
	config := Config{}
	expectedError := testCase.expectedErr
	var actualObject []interface{}
	var err error
	if len(testCase.filters) > 0 {
		hasConfig = true
		for id, function := range testCase.filters {
			config.SetFilterFunction(id, function)
		}
	}
	if len(testCase.aggregates) > 0 {
		hasConfig = true
		for id, function := range testCase.aggregates {
			config.SetAggregateFunction(id, function)
		}
	}
	if testCase.accessorMode {
		hasConfig = true
		config.SetAccessorMode()
	}
	if hasConfig {
		actualObject, err = Retrieve(jsonPath, inputJSON, config)
	} else {
		actualObject, err = Retrieve(jsonPath, inputJSON)
	}
	if err != nil {
		if reflect.TypeOf(expectedError) == reflect.TypeOf(err) &&
			fmt.Sprintf(`%s`, expectedError) == fmt.Sprintf(`%s`, err) {
			return nil, err
		}
		t.Errorf("%s: expected error<%s> != actual error<%s>\n", fileLine, expectedError, err)
		return nil, err
	}
	if expectedError != nil {
		t.Errorf("%s: expected error<%s> != actual error<none>\n",
			fileLine, expectedError)
		return nil, err
	}

	return actualObject, err
}

func runTestCase(t *testing.T, testCase TestCase, fileLine string) {
	srcJSON := testCase.inputJSON
	var src interface{}
	var err error

	if testCase.unmarshalFunc != nil {
		err = testCase.unmarshalFunc(srcJSON, &src)
	} else {
		err = json.Unmarshal([]byte(srcJSON), &src)
	}
	if err != nil {
		t.Errorf("%s: Error: %v", fileLine, err)
		return
	}

	actualObject, err := execTestRetrieve(t, src, testCase, fileLine)
	if t.Failed() {
		return
	}
	if err != nil {
		return
	}

	if testCase.resultValidator != nil {
		err := testCase.resultValidator(src, actualObject)
		if err != nil {
			t.Errorf("%s: Error: %v", fileLine, err)
		}
		return
	}

	actualOutputJSON, err := json.Marshal(actualObject)
	if err != nil {
		t.Errorf("%s: Error: %v", fileLine, err)
		return
	}

	if string(actualOutputJSON) != testCase.expectedJSON {
		t.Errorf("%s: expectedOutputJSON<%s> != actualOutputJSON<%s>\n",
			fileLine, testCase.expectedJSON, actualOutputJSON)
		return
	}
}

func runTestCases(t *testing.T, testGroupName string, testCases []TestCase) {
	for _, testCase := range testCases {
		testCase := testCase
		if _, file, line, ok := runtime.Caller(2); ok {
			fileLine := fmt.Sprintf(`%s:%d`, file, line)
			t.Run(
				fmt.Sprintf(`%s_<%s>_<%s>`, testGroupName, testCase.jsonpath, testCase.inputJSON),
				func(t *testing.T) {
					t.Parallel()
					runTestCase(t, testCase, fileLine)
				})
		}
	}
}

func runTestGroups(t *testing.T, testGroups TestGroup) {
	for testGroupName, testCases := range testGroups {
		runTestCases(t, testGroupName, testCases)
	}
}

func TestRetrieve_dotNotation(t *testing.T) {
	testGroups := TestGroup{
		`root`: []TestCase{
			{
				jsonpath:     `$`,
				inputJSON:    `{"a":"b","c":{"d":"e"}}`,
				expectedJSON: `[{"a":"b","c":{"d":"e"}}]`,
			},
			{
				jsonpath:     `a`,
				inputJSON:    `{"a":"b","c":{"d":"e"}}`,
				expectedJSON: `["b"]`,
			},
			{
				jsonpath:     `\@`,
				inputJSON:    `{"@":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:    `@`,
				inputJSON:   `{"@":1}`,
				expectedErr: ErrorInvalidSyntax{position: 0, reason: `unrecognized input`, near: `@`},
			},
		},
		`child`: []TestCase{
			{
				jsonpath:     `$.a`,
				inputJSON:    `{"a":"b","c":{"d":"e"}}`,
				expectedJSON: `["b"]`,
			},
			{
				jsonpath:     `$.c`,
				inputJSON:    `{"a":"b","c":{"d":"e"}}`,
				expectedJSON: `[{"d":"e"}]`,
			},
			{
				jsonpath:     `$.2`,
				inputJSON:    `{"a":1,"2":2,"3":{"2":1}}`,
				expectedJSON: `[2]`,
			},
			{
				jsonpath:     `$[0].a`,
				inputJSON:    `[{"a":"b","c":{"d":"e"}},{"a":"y"}]`,
				expectedJSON: `["b"]`,
			},
			{
				jsonpath:     `[0].a`,
				inputJSON:    `[{"a":"b","c":{"d":"e"}},{"a":"y"}]`,
				expectedJSON: `["b"]`,
			},
			{
				jsonpath:     `$[2,0].a`,
				inputJSON:    `[{"a":"b","c":{"a":"d"}},{"a":"e"},{"a":"a"}]`,
				expectedJSON: `["a","b"]`,
			},
			{
				jsonpath:     `$[0:2].a`,
				inputJSON:    `[{"a":"b","c":{"d":"e"}},{"a":"a"},{"a":"c"}]`,
				expectedJSON: `["b","a"]`,
			},
			{
				jsonpath:     `$.a.a2`,
				inputJSON:    `{"a":{"a1":"1","a2":"2"},"b":{"b1":"3"}}`,
				expectedJSON: `["2"]`,
			},
			{
				jsonpath:     `.a`,
				inputJSON:    `{"a":"b","c":{"d":"e"}}`,
				expectedJSON: `["b"]`,
			},
		},
		`terms-with-special-meanings`: []TestCase{
			{
				jsonpath:     `$.nil`,
				inputJSON:    `{"nil":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.null`,
				inputJSON:    `{"null":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.true`,
				inputJSON:    `{"true":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.false`,
				inputJSON:    `{"false":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.in`,
				inputJSON:    `{"in":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.length`,
				inputJSON:    `{"length":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:    `$.length`,
				inputJSON:   `["length",1,2]`,
				expectedErr: createErrorTypeUnmatched(`.length`, `object`, `[]interface {}`),
			},
		},
		`character-type::encoded-JSONPath`: []TestCase{
			{
				// 0x20 - 0x2C, 0x2E - 0x2F
				jsonpath:     `$.\ \!\"\#\$\%\&\'\(\)\*\+\,\.\/`,
				inputJSON:    `{" !\"#$%&'()*+,./":1}`,
				expectedJSON: `[1]`,
			},
			{
				// 0x3A - 0x40, 0x5B - 0x5E
				jsonpath:     `$.\:\;\<\=\>\?\@\[\\\]\^`,
				inputJSON:    `{":;<=>?@[\\]^":1}`,
				expectedJSON: `[1]`,
			},
			{
				// 0x60
				jsonpath:     "$.\\`",
				inputJSON:    "{\"`\":1}",
				expectedJSON: `[1]`,
			},
			{
				// 0x7B - 0x7E
				jsonpath:     `$.\{\|\}\~`,
				inputJSON:    `{"{|}~":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:    `$.\a`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.\a`},
			},
			{
				jsonpath:    `$.a\b`,
				inputJSON:   `{"ab":1}`,
				expectedErr: ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: `\b`},
			},
			{
				jsonpath:    `$.\-`,
				inputJSON:   `{"-":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.\-`},
			},
			{
				jsonpath:    `$.\_`,
				inputJSON:   `{"_":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.\_`},
			},
			{
				jsonpath:     `$.\'a\.b\'`,
				inputJSON:    `{"'a.b'":1,"a":{"b":2},"'a'":{"'b'":3},"'a":{"b'":4}}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:    `$.'a\.b'`,
				inputJSON:   `{"'a.b'":1,"a":{"b":2},"'a'":{"'b'":3},"'a":{"b'":4}}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.'a\.b'`},
			},
		},
		`character-type::not-encoded-error`: []TestCase{
			{
				jsonpath:    `$. `,
				inputJSON:   `{" ":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `. `},
			},
			{
				jsonpath:    `$.!`,
				inputJSON:   `{"!":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.!`},
			},
			{
				jsonpath:    `$."`,
				inputJSON:   `{"\"":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `."`},
			},
			{
				jsonpath:    `$.#`,
				inputJSON:   `{"#":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.#`},
			},
			{
				jsonpath:    `$.$`,
				inputJSON:   `{"$":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.$`},
			},
			{
				jsonpath:    `$.%`,
				inputJSON:   `{"%":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.%`},
			},
			{
				jsonpath:    `$.&`,
				inputJSON:   `{"&":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.&`},
			},
			{
				jsonpath:    `$.'`,
				inputJSON:   `{"'":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.'`},
			},
			{
				jsonpath:    `$.(`,
				inputJSON:   `{"(":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.(`},
			},
			{
				jsonpath:    `$.)`,
				inputJSON:   `{")":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.)`},
			},
			{
				jsonpath:    `$.+`,
				inputJSON:   `{"+":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.+`},
			},
			{
				jsonpath:    `$.,`,
				inputJSON:   `{",":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.,`},
			},
			{
				jsonpath:    `$./`,
				inputJSON:   `{"/":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `./`},
			},
			{
				jsonpath:    `$.:`,
				inputJSON:   `{":":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.:`},
			},
			{
				jsonpath:    `$.;`,
				inputJSON:   `{";":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.;`},
			},
			{
				jsonpath:    `$.<`,
				inputJSON:   `{"<":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.<`},
			},
			{
				jsonpath:    `$.=`,
				inputJSON:   `{"=":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.=`},
			},
			{
				jsonpath:    `$.>`,
				inputJSON:   `{">":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.>`},
			},
			{
				jsonpath:    `$.?`,
				inputJSON:   `{"?":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.?`},
			},
			{
				jsonpath:    `$.@`,
				inputJSON:   `{"@":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.@`},
			},
			{
				jsonpath:    `$.\`,
				inputJSON:   `{"\\":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.\`},
			},
			{
				jsonpath:    `$.]`,
				inputJSON:   `{"]":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.]`},
			},
			{
				jsonpath:    `$.^`,
				inputJSON:   `{"^":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.^`},
			},
			{
				jsonpath:    "$.`",
				inputJSON:   "{\"`\":1}",
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: ".`"},
			},
			{
				jsonpath:    `$.{`,
				inputJSON:   `{"{":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.{`},
			},
			{
				jsonpath:    `$.|`,
				inputJSON:   `{"|":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.|`},
			},
			{
				jsonpath:    `$.}`,
				inputJSON:   `{"}":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.}`},
			},
			{
				jsonpath:    `$.~`,
				inputJSON:   `{"~":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.~`},
			},
		},
		`character-type::control-ascii`: []TestCase{
			{
				// TAB
				jsonpath:    `$.` + "\x09",
				inputJSON:   `{"\t":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.` + "\x09"},
			},
			{
				// TAB
				jsonpath:    `$.\` + "\x09",
				inputJSON:   `{"\t":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.\` + "\x09"},
			},
			{
				// CR
				jsonpath:    `$.` + "\x0d",
				inputJSON:   `{"\n":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.` + "\x0d"},
			},
			{
				// CR
				jsonpath:    `$.\` + "\x0d",
				inputJSON:   `{"\n":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.\` + "\x0d"},
			},
			{
				jsonpath:    `$.` + "\x1f",
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.` + "\x1f"},
			},
			{
				jsonpath:    `$.` + "\x7F",
				inputJSON:   `{"` + "\x7F" + `":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.` + "\x7F"},
			},
			{
				jsonpath:    `$.\` + "\x7F",
				inputJSON:   `{"` + "\x7F" + `":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.\` + "\x7F"},
			},
			{
				jsonpath:    `$.a` + "\x7F",
				inputJSON:   `{"a` + "\x7F" + `":1}`,
				expectedErr: ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: "\x7F"},
			},
		},
		`character-type::unicode-character`: []TestCase{
			{
				jsonpath:     `$.ﾃｽﾄソポァゼゾタダＡボマミ①`,
				inputJSON:    `{"ﾃｽﾄソポァゼゾタダＡボマミ①":1}`,
				expectedJSON: `[1]`,
			},
		},
		`not-exist`: []TestCase{
			{
				jsonpath:    `$.d`,
				inputJSON:   `{"a":"b","c":{"d":"e"}}`,
				expectedErr: createErrorMemberNotExist(`.d`),
			},
		},
		`type-unmatched`: []TestCase{
			{
				jsonpath:    `$.2`,
				inputJSON:   `["a","b",{"2":1}]`,
				expectedErr: createErrorTypeUnmatched(`.2`, `object`, `[]interface {}`),
			},
			{
				jsonpath:    `$.-1`,
				inputJSON:   `["a","b",{"2":1}]`,
				expectedErr: createErrorTypeUnmatched(`.-1`, `object`, `[]interface {}`),
			},
			{
				jsonpath:    `$.a.d`,
				inputJSON:   `{"a":"b","c":{"d":"e"}}`,
				expectedErr: createErrorTypeUnmatched(`.d`, `object`, `string`),
			},
			{
				jsonpath:    `$.a.d`,
				inputJSON:   `{"a":123}`,
				expectedErr: createErrorTypeUnmatched(`.d`, `object`, `float64`),
			},
			{
				jsonpath:    `$.a.d`,
				inputJSON:   `{"a":true}`,
				expectedErr: createErrorTypeUnmatched(`.d`, `object`, `bool`),
			},
			{
				jsonpath:    `$.a.d`,
				inputJSON:   `{"a":null}`,
				expectedErr: createErrorTypeUnmatched(`.d`, `object`, `null`),
			},
			{
				jsonpath:    `$.a`,
				inputJSON:   `[1,2]`,
				expectedErr: createErrorTypeUnmatched(`.a`, `object`, `[]interface {}`),
			},
			{
				jsonpath:    `$.a`,
				inputJSON:   `[{"a":1}]`,
				expectedErr: createErrorTypeUnmatched(`.a`, `object`, `[]interface {}`),
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_recursiveDescent(t *testing.T) {
	testGroups := TestGroup{
		`identifier::dotNotation`: []TestCase{
			{
				jsonpath:     `$..a`,
				inputJSON:    `{"a":"b","c":{"a":"d"},"e":["a",{"a":{"a":"h"}}]}`,
				expectedJSON: `["b","d",{"a":"h"},"h"]`,
			},
			{
				jsonpath:     `$.a..b`,
				inputJSON:    `{"a":{"b":1,"c":{"b":2},"d":["b",{"a":3,"b":4}]},"b":5}`,
				expectedJSON: `[1,2,4]`,
			},
			{
				jsonpath:     `$..a.b`,
				inputJSON:    `{"a":{"b":1,"c":{"b":2},"d":["b",{"a":3,"b":4}]},"b":5}`,
				expectedJSON: `[1]`,
			},
		},
		`identifier::bracketNotation`: []TestCase{
			{
				jsonpath:     `$..['a']`,
				inputJSON:    `{"a":"b","c":{"a":"d"},"e":["a",{"a":{"a":"h"}}]}`,
				expectedJSON: `["b","d",{"a":"h"},"h"]`,
			},
			{
				jsonpath:     `$['a']..['b']`,
				inputJSON:    `{"a":{"b":1,"c":{"b":2},"d":["b",{"a":3,"b":4}]},"b":5}`,
				expectedJSON: `[1,2,4]`,
			},
			{
				jsonpath:     `$..['a']['b']`,
				inputJSON:    `{"a":{"b":1,"c":{"b":2},"d":["b",{"a":3,"b":4}]},"b":5}`,
				expectedJSON: `[1]`,
			},
		},
		`qualifier::index`: []TestCase{
			{
				jsonpath:     `$..[1]`,
				inputJSON:    `[{"a":["b",{"c":{"a":"d"}}],"e":["f",{"g":{"a":"h"}}]},0]`,
				expectedJSON: `[0,{"c":{"a":"d"}},{"g":{"a":"h"}}]`,
			},
			{
				jsonpath:     `$..[1].a`,
				inputJSON:    `[{"a":["b",{"a":{"a":"d"}}],"e":["f",{"g":{"a":"h"}}]},0]`,
				expectedJSON: `[{"a":"d"}]`,
			},
			{
				jsonpath:     `$..[1,2]`,
				inputJSON:    `[{"a":["b",{"c":{"a":"d"}}],"e":["f",{"g":{"a":"h"}}]},0,1]`,
				expectedJSON: `[0,1,{"c":{"a":"d"}},{"g":{"a":"h"}}]`,
			},
			{
				jsonpath:     `$..[1,2].a`,
				inputJSON:    `[{"a":["b",{"a":{"a":"d"}}],"e":["f",{"g":{"a":"h"}}]},0,{"a":1}]`,
				expectedJSON: `[1,{"a":"d"}]`,
			},
		},
		`qualifier::slice`: []TestCase{
			{
				jsonpath:     `$..[1:3]`,
				inputJSON:    `[{"a":["b",{"c":{"a":"d"}}],"e":["f",{"g":{"a":"h"}}]},0,1]`,
				expectedJSON: `[0,1,{"c":{"a":"d"}},{"g":{"a":"h"}}]`,
			},
			{
				jsonpath:     `$..[1:3].c`,
				inputJSON:    `[{"a":["b",{"c":{"a":"d"}}],"e":["f",{"g":{"a":"h"}}]},0,1]`,
				expectedJSON: `[{"a":"d"}]`,
			},
			{
				jsonpath:     `$..[1:0:-1]`,
				inputJSON:    `[{"a":["b",{"c":{"a":"d"}}],"e":["f",{"g":{"a":"h"}}]},0,1]`,
				expectedJSON: `[0,{"c":{"a":"d"}},{"g":{"a":"h"}}]`,
			},
		},
		`empty-input`: []TestCase{
			{
				jsonpath:    `$..a`,
				inputJSON:   `{}`,
				expectedErr: createErrorMemberNotExist(`a`),
			},
			{
				jsonpath:    `$..a`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`..`),
			},
		},
		`child-error`: []TestCase{
			{
				jsonpath:    `$..x`,
				inputJSON:   `{"a":"b","c":{"a":"d"},"e":["f",{"g":{"a":"h"}}]}`,
				expectedErr: createErrorMemberNotExist(`x`),
			},
			{
				jsonpath:    `$..a.x`,
				inputJSON:   `{"a":"b","c":{"a":"d"},"e":["f",{"g":{"a":"h"}}]}`,
				expectedErr: createErrorTypeUnmatched(`.x`, `object`, `string`),
			},
			{
				jsonpath:    `$..a.x`,
				inputJSON:   `{"a":"b","c":{"a":"d"},"e":["f",{"g":{"a":{"h":1}}}]}`,
				expectedErr: createErrorMemberNotExist(`.x`),
			},
			{
				// The case where '.x' terminates with an error first
				jsonpath:    `$.x..a`,
				inputJSON:   `{"a":"b","c":{"a":"d"},"e":["f",{"g":{"a":"h"}}]}`,
				expectedErr: createErrorMemberNotExist(`.x`),
			},
		},
		`character-type::Non-alphabet-accepted-in-JSON`: []TestCase{
			{
				jsonpath:     `$..\'a\'`,
				inputJSON:    `{"'a'":1,"b":{"'a'":2},"c":["'a'",{"d":{"'a'":{"'a'":3}}}]}`,
				expectedJSON: `[1,2,{"'a'":3},3]`,
			},
			{
				jsonpath:     `$..\"a\"`,
				inputJSON:    `{"\"a\"":1,"b":{"\"a\"":2},"c":["\"a\"",{"d":{"\"a\"":{"\"a\"":3}}}]}`,
				expectedJSON: `[1,2,{"\"a\"":3},3]`,
			},
		},
		`filter`: []TestCase{
			{
				jsonpath:     `$..[?(@.a)]`,
				inputJSON:    `{"a":1,"b":[{"a":2},{"b":{"a":3}},{"a":{"a":4}}]}`,
				expectedJSON: `[{"a":2},{"a":{"a":4}},{"a":3},{"a":4}]`,
			},
		},
		`multi-identifier`: []TestCase{
			{
				jsonpath:     `$..['a','b']`,
				inputJSON:    `[{"a":1,"b":2,"c":{"a":3}},{"a":4},{"b":5},{"a":6,"b":7},{"d":{"b":8}}]`,
				expectedJSON: `[1,2,3,4,5,6,7,8]`,
			},
		},
		`type-unmatched`: []TestCase{
			{
				jsonpath:    `$..a`,
				inputJSON:   `null`,
				expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `null`),
			},
			{
				jsonpath:    `$..a`,
				inputJSON:   `true`,
				expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `bool`),
			},
			{
				jsonpath:    `$..a`,
				inputJSON:   `"abc"`,
				expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
			},
			{
				jsonpath:    `$..a`,
				inputJSON:   `123`,
				expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `float64`),
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_dotNotation_wildcard(t *testing.T) {
	testGroups := TestGroup{
		`array`: []TestCase{
			{
				jsonpath:     `$.*`,
				inputJSON:    `[[1],[2,3],123,"a",{"b":"c"},[0,1],null]`,
				expectedJSON: `[[1],[2,3],123,"a",{"b":"c"},[0,1],null]`,
			},
			{
				jsonpath:     `$.*[1]`,
				inputJSON:    `[[1],[2,3],[4,[5,6,7]]]`,
				expectedJSON: `[3,[5,6,7]]`,
			},
			{
				jsonpath:     `$.*.a`,
				inputJSON:    `[{"a":1},{"a":[2,3]}]`,
				expectedJSON: `[1,[2,3]]`,
			},
			{
				jsonpath:     `$..*`,
				inputJSON:    `[{"a":1},{"a":[2,3]},null,true]`,
				expectedJSON: `[{"a":1},{"a":[2,3]},null,true,1,[2,3],2,3]`,
			},
			{
				jsonpath:     `$.*['a','b']`,
				inputJSON:    `[{"a":1,"b":2,"c":3},{"a":4,"b":5,"d":6}]`,
				expectedJSON: `[1,2,4,5]`,
			},
			{
				jsonpath:     `*`,
				inputJSON:    `[[1],[2,3],123,"a",{"b":"c"},[0,1],null]`,
				expectedJSON: `[[1],[2,3],123,"a",{"b":"c"},[0,1],null]`,
			},
		},
		`object`: []TestCase{
			{
				jsonpath:     `$.*`,
				inputJSON:    `{"a":[1],"b":[2,3],"c":{"d":4}}`,
				expectedJSON: `[[1],[2,3],{"d":4}]`,
			},
			{
				jsonpath:     `$..*`,
				inputJSON:    `{"a":1,"b":[2,3],"c":{"d":4,"e":[5,6]}}`,
				expectedJSON: `[1,[2,3],{"d":4,"e":[5,6]},2,3,4,[5,6],5,6]`,
			},
			{
				jsonpath:     `$..[*]`,
				inputJSON:    `{"a":1,"b":[2,3],"c":{"d":"e","f":[4,5]}}`,
				expectedJSON: `[1,[2,3],{"d":"e","f":[4,5]},2,3,"e",[4,5],4,5]`,
			},
			{
				jsonpath:     `*`,
				inputJSON:    `{"a":[1],"b":[2,3],"c":{"d":4}}`,
				expectedJSON: `[[1],[2,3],{"d":4}]`,
			},
		},
		`two-wildcards`: []TestCase{
			{
				jsonpath:     `$.*.*`,
				inputJSON:    `[[1,2,3],[4,5,6]]`,
				expectedJSON: `[1,2,3,4,5,6]`,
			},
			{
				jsonpath:     `$.*.a.*`,
				inputJSON:    `[{"a":[1]}]`,
				expectedJSON: `[1]`,
			},
		},
		`empty-input`: []TestCase{
			{
				jsonpath:    `$.*`,
				inputJSON:   `{}`,
				expectedErr: createErrorMemberNotExist(`.*`),
			},
			{
				jsonpath:    `$.*`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`.*`),
			},
		},
		`recursive`: []TestCase{
			{
				jsonpath:    `$..*`,
				inputJSON:   `{}`,
				expectedErr: createErrorMemberNotExist(`*`),
			},
			{
				jsonpath:    `$..*`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`*`),
			},
			{
				jsonpath:     `$..*`,
				inputJSON:    `{"a":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$..*`,
				inputJSON:    `{"b":2,"a":1}`,
				expectedJSON: `[1,2]`,
			},
			{
				jsonpath:     `$..*`,
				inputJSON:    `{"a":{"b":2}}`,
				expectedJSON: `[{"b":2},2]`,
			},
			{
				jsonpath:     `$..*`,
				inputJSON:    `{"a":{"c":3,"b":2}}`,
				expectedJSON: `[{"b":2,"c":3},2,3]`,
			},
			{
				jsonpath:     `$..*`,
				inputJSON:    `[1]`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$..*`,
				inputJSON:    `[2,1]`,
				expectedJSON: `[2,1]`,
			},
			{
				jsonpath:     `$..*`,
				inputJSON:    `[{"a":1}]`,
				expectedJSON: `[{"a":1},1]`,
			},
			{
				jsonpath:     `$..*`,
				inputJSON:    `[{"b":2,"a":1}]`,
				expectedJSON: `[{"a":1,"b":2},1,2]`,
			},
			{
				jsonpath:     `..*`,
				inputJSON:    `{"a":1}`,
				expectedJSON: `[1]`,
			},
		},
		`child-error`: []TestCase{
			{
				jsonpath:    `$.*.a.b`,
				inputJSON:   `{"a":{"b":1}}`,
				expectedErr: createErrorMemberNotExist(`.a`),
			},
			{
				jsonpath:    `$.*.a.b`,
				inputJSON:   `[{"b":1}]`,
				expectedErr: createErrorMemberNotExist(`.a`),
			},
			{
				jsonpath:    `$.*.a.b.c`,
				inputJSON:   `{"a":{"b":1},"b":{"a":2}}`,
				expectedErr: createErrorTypeUnmatched(`.b`, `object`, `float64`),
			},
			{
				jsonpath:    `$.*.a.b.c`,
				inputJSON:   `[{"b":1},{"a":2}]`,
				expectedErr: createErrorTypeUnmatched(`.b`, `object`, `float64`),
			},
			{
				jsonpath:    `$.*.a.b.c`,
				inputJSON:   `{"a":{"a":1},"b":{"a":{"c":2}}}`,
				expectedErr: createErrorMemberNotExist(`.b`),
			},
			{
				jsonpath:    `$.*.a.b.c`,
				inputJSON:   `[{"a":1},{"a":{"c":2}}]`,
				expectedErr: createErrorMemberNotExist(`.b`),
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_bracketNotation(t *testing.T) {
	testGroups := TestGroup{
		`child`: []TestCase{
			{
				jsonpath:     `$['a']`,
				inputJSON:    `{"a":"b","c":{"d":"e"}}`,
				expectedJSON: `["b"]`,
			},
			{
				jsonpath:     `$[0]['a']`,
				inputJSON:    `[{"a":"b","c":{"d":"e"}},{"x":"y"}]`,
				expectedJSON: `["b"]`,
			},
			{
				jsonpath:     `$['a'][0]['b']`,
				inputJSON:    `{"a":[{"b":"x"},"y"],"c":{"d":"e"}}`,
				expectedJSON: `["x"]`,
			},
			{
				jsonpath:     `$[0:2]['b']`,
				inputJSON:    `[{"a":1},{"b":3},{"b":2,"c":4}]`,
				expectedJSON: `[3]`,
			},
			{
				jsonpath:     `$[:]['b']`,
				inputJSON:    `[{"a":1},{"b":3},{"b":2,"c":4}]`,
				expectedJSON: `[3,2]`,
			},
			{
				jsonpath:     `$['a']['a2']`,
				inputJSON:    `{"a":{"a1":"1","a2":"2"},"b":{"b1":"3"}}`,
				expectedJSON: `["2"]`,
			},
		},
		`number-identifier`: []TestCase{
			{
				jsonpath:     `$['0']`,
				inputJSON:    `{"0":1,"a":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['1']`,
				inputJSON:    `{"1":"a","a":2}`,
				expectedJSON: `["a"]`,
			},
		},
		`not-exist`: []TestCase{
			{
				jsonpath:    `$['d']`,
				inputJSON:   `{"a":"b","c":{"d":"e"}}`,
				expectedErr: createErrorMemberNotExist(`['d']`),
			},
		},
		`character-type::single-quoted`: []TestCase{
			{
				jsonpath:     `$['ab']`,
				inputJSON:    `{"ab":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['a\'b']`,
				inputJSON:    `{"a'b":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['ab\'c']`,
				inputJSON:    `{"ab'c":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['a\/b']`,
				inputJSON:    `{"a\/b":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['a\\b']`,
				inputJSON:    `{"a\\b":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['a\bb']`,
				inputJSON:    `{"a\bb":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['a\fb']`,
				inputJSON:    `{"a\fb":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['a\nb']`,
				inputJSON:    `{"a\nb":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['a\rb']`,
				inputJSON:    `{"a\rb":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['a\tb']`,
				inputJSON:    `{"a\tb":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:    `$['a\c']`,
				inputJSON:   `{"ac":1,"b":2}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['a\c']`},
			},
			{
				jsonpath:    `$['a'c']`,
				inputJSON:   `{"ac":1,"b":2}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['a'c']`},
			},
			{
				jsonpath:     `$['a"c']`,
				inputJSON:    `{"a\"c":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['a\uD834\uDD1Ec']`,
				inputJSON:    `{"a\uD834\uDD1Ec":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['a\ud834\udd1ec']`,
				inputJSON:    `{"a\uD834\uDD1Ec":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['\u0000']`,
				inputJSON:    `{"\u0000":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['\uabcd']`,
				inputJSON:    `{"\uabcd":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['\uABCD']`,
				inputJSON:    `{"\uabcd":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:    `$['\uX000']`,
				inputJSON:   `{"a":1,"b":2}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['\uX000']`},
			},
			{
				jsonpath:    `$['\u0X00']`,
				inputJSON:   `{"a":1,"b":2}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['\u0X00']`},
			},
			{
				jsonpath:    `$['\u00X0']`,
				inputJSON:   `{"a":1,"b":2}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['\u00X0']`},
			},
			{
				jsonpath:    `$['\u000X']`,
				inputJSON:   `{"a":1,"b":2}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['\u000X']`},
			},
		},
		`character-type::double-quoted`: []TestCase{
			{
				jsonpath:     `$["ab"]`,
				inputJSON:    `{"ab":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$["a\"b"]`,
				inputJSON:    `{"a\"b":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$["ab\"c"]`,
				inputJSON:    `{"ab\"c":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$["a\/b"]`,
				inputJSON:    `{"a\/b":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$["a\\b"]`,
				inputJSON:    `{"a\\b":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$["a\bb"]`,
				inputJSON:    `{"a\bb":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$["a\fb"]`,
				inputJSON:    `{"a\fb":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$["a\nb"]`,
				inputJSON:    `{"a\nb":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$["a\rb"]`,
				inputJSON:    `{"a\rb":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$["a\tb"]`,
				inputJSON:    `{"a\tb":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:    `$["a\c"]`,
				inputJSON:   `{"ac":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `["a\c"]`},
			},
			{
				jsonpath:    `$["a"b"]`,
				inputJSON:   `{"ab":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `["a"b"]`},
			},
			{
				jsonpath:     `$["a'b"]`,
				inputJSON:    `{"a'b":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$["a\uD834\uDD1Ec"]`,
				inputJSON:    `{"a\uD834\uDD1Ec":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$["a\ud834\udd1ec"]`,
				inputJSON:    `{"a\uD834\uDD1Ec":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$["\u0000"]`,
				inputJSON:    `{"\u0000":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$["\uabcd"]`,
				inputJSON:    `{"\uabcd":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$["\uABCD"]`,
				inputJSON:    `{"\uabcd":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:    `$["\uX000"]`,
				inputJSON:   `{"a":1,"b":2}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `["\uX000"]`},
			},
			{
				jsonpath:    `$["\u0X00"]`,
				inputJSON:   `{"a":1,"b":2}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `["\u0X00"]`},
			},
			{
				jsonpath:    `$["\u00X0"]`,
				inputJSON:   `{"a":1,"b":2}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `["\u00X0"]`},
			},
			{
				jsonpath:    `$["\u000X"]`,
				inputJSON:   `{"a":1,"b":2}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `["\u000X"]`},
			},
		},
		`character-types-like-the-prohibited-dot-notation`: []TestCase{
			{
				jsonpath:     `$['a.b']`,
				inputJSON:    `{"a.b":1,"a":{"b":2}}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$[':']`,
				inputJSON:    `{":":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['[']`,
				inputJSON:    `{"[":1,"]":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$[']']`,
				inputJSON:    `{"[":1,"]":2}`,
				expectedJSON: `[2]`,
			},
			{
				jsonpath:     `$['$']`,
				inputJSON:    `{"$":2}`,
				expectedJSON: `[2]`,
			},
			{
				jsonpath:     `$['@']`,
				inputJSON:    `{"@":2}`,
				expectedJSON: `[2]`,
			},
			{
				jsonpath:     `$['*']`,
				inputJSON:    `{"*":2}`,
				expectedJSON: `[2]`,
			},
			{
				jsonpath:    `$['*']`,
				inputJSON:   `{"a":1,"b":2}`,
				expectedErr: createErrorMemberNotExist(`['*']`),
			},
			{
				jsonpath:     `$['.']`,
				inputJSON:    `{".":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$[',']`,
				inputJSON:    `{",":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['.*']`,
				inputJSON:    `{".*":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['"']`,
				inputJSON:    `{"\"":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$["'"]`,
				inputJSON:    `{"'":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$[':@."$,*\'\\']`,
				inputJSON:    `{":@.\"$,*'\\": 1}`,
				expectedJSON: `[1]`,
			},
		},
		`empty-identifier`: []TestCase{
			{
				jsonpath:     `$['']`,
				inputJSON:    `{"":1, "''":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$[""]`,
				inputJSON:    `{"":1, "''":2,"\"\"":3}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$[''][0]`,
				inputJSON:    `[1,2,3]`,
				expectedJSON: `[1]`,
				expectedErr:  createErrorTypeUnmatched(`['']`, `object`, `[]interface {}`),
			},
		},
		`mixing-bracket-and-dot-notation`: []TestCase{
			{
				jsonpath:     `$['a'].b`,
				inputJSON:    `{"b":2,"a":{"b":1}}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.a['b']`,
				inputJSON:    `{"b":2,"a":{"b":1}}`,
				expectedJSON: `[1]`,
			},
		},
		`empty-input`: []TestCase{
			{
				jsonpath:    `$['a']`,
				inputJSON:   `[]`,
				expectedErr: createErrorTypeUnmatched(`['a']`, `object`, `[]interface {}`),
			},
			{
				jsonpath:    `$['a']`,
				inputJSON:   `{}`,
				expectedErr: createErrorMemberNotExist(`['a']`),
			},
		},
		`type-unmatched`: []TestCase{
			{
				jsonpath:    `$['a']`,
				inputJSON:   `"abc"`,
				expectedErr: createErrorTypeUnmatched(`['a']`, `object`, `string`),
			},
			{
				jsonpath:    `$['a']`,
				inputJSON:   `123`,
				expectedErr: createErrorTypeUnmatched(`['a']`, `object`, `float64`),
			},
			{
				jsonpath:    `$['a']`,
				inputJSON:   `true`,
				expectedErr: createErrorTypeUnmatched(`['a']`, `object`, `bool`),
			},
			{
				jsonpath:    `$['a']`,
				inputJSON:   `null`,
				expectedErr: createErrorTypeUnmatched(`['a']`, `object`, `null`),
			},
			{
				jsonpath:    `$['a']`,
				inputJSON:   `[1,2,3]`,
				expectedErr: createErrorTypeUnmatched(`['a']`, `object`, `[]interface {}`),
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_bracketNotation_multiIdentifiers(t *testing.T) {
	testGroups := TestGroup{
		`identifier-order`: []TestCase{
			{
				jsonpath:     `$['a','b']`,
				inputJSON:    `{"a":1, "b":2}`,
				expectedJSON: `[1,2]`,
			},
			{
				jsonpath:     `$['b','a']`,
				inputJSON:    `{"a":1, "b":2}`,
				expectedJSON: `[2,1]`,
			},
			{
				jsonpath:     `$['b','a']`,
				inputJSON:    `{"b":2,"a":1}`,
				expectedJSON: `[2,1]`,
			},
			{
				jsonpath:     `$['a','b']`,
				inputJSON:    `{"b":2,"a":1}`,
				expectedJSON: `[1,2]`,
			},
			{
				jsonpath:     `$['a','b',*]`,
				inputJSON:    `{"b":2,"a":1,"c":3}`,
				expectedJSON: `[1,2,1,2,3]`,
			},
		},
		`mixing-qualifier-error`: []TestCase{
			{
				jsonpath:    `$['a','b',0]`,
				inputJSON:   `{"b":2,"a":1,"c":3}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['a','b',0]`},
			},
			{
				jsonpath:    `$['a','b',0:1]`,
				inputJSON:   `{"b":2,"a":1,"c":3}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['a','b',0:1]`},
			},
			{
				jsonpath:    `$['a','b',(command)]`,
				inputJSON:   `{"b":2,"a":1,"c":3}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['a','b',(command)]`},
			},
			{
				jsonpath:    `$['a','b',?(@)]`,
				inputJSON:   `{"b":2,"a":1,"c":3}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['a','b',?(@)]`},
			},
		},
		`connecting-child`: []TestCase{
			{
				jsonpath:     `$['a','b'].a`,
				inputJSON:    `{"a":{"a":1}, "b":{"c":2}}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['a','b']['a']`,
				inputJSON:    `{"a":{"a":1}, "b":{"c":2}}`,
				expectedJSON: `[1]`,
			},
		},
		`partial-found`: []TestCase{
			{
				jsonpath:     `$['a','c']`,
				inputJSON:    `{"a":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['b','c']`,
				inputJSON:    `{"a":1,"b":2}`,
				expectedJSON: `[2]`,
			},
			{
				jsonpath:     `$['c','a']`,
				inputJSON:    `{"a":1,"b":2}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['c','b']`,
				inputJSON:    `{"a":1,"b":2}`,
				expectedJSON: `[2]`,
			},
			{
				jsonpath:    `$['c','d']`,
				inputJSON:   `{"a":1,"b":2}`,
				expectedErr: createErrorMemberNotExist(`['c','d']`),
			},
		},
		`two-level-multi-identifiers`: []TestCase{
			{
				jsonpath:     `$['a','b']['a','b']`,
				inputJSON:    `{"a":{"a":1},"b":{"b":2}}`,
				expectedJSON: `[1,2]`,
			},
			{
				jsonpath:     `$['a','b']['a','b']`,
				inputJSON:    `{"a":{"b":1},"b":{"a":2}}`,
				expectedJSON: `[1,2]`,
			},
			{
				jsonpath:     `$['a','b']['a','b']`,
				inputJSON:    `{"a":{"a":1,"b":2},"b":{"c":3}}`,
				expectedJSON: `[1,2]`,
			},
			{
				jsonpath:     `$['a','b']['a','b']`,
				inputJSON:    `{"a":{"b":1},"c":{"a":2}}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:    `$['a','b']['c','d']`,
				inputJSON:   `{"a":{"a":1},"b":{"b":2}}`,
				expectedErr: createErrorMemberNotExist(`['c','d']`),
			},
			{
				jsonpath:    `$['a','b']['c','d']`,
				inputJSON:   `{"a":{"a":1},"c":{"b":2}}`,
				expectedErr: createErrorMemberNotExist(`['c','d']`),
			},
			{
				jsonpath:    `$['a','b']['c','d']`,
				inputJSON:   `{"c":{"a":1},"d":{"b":2}}`,
				expectedErr: createErrorMemberNotExist(`['a','b']`),
			},
			{
				jsonpath:    `$['a','b']['c','d'].e`,
				inputJSON:   `{"a":{"c":1},"b":{"c":2}}`,
				expectedErr: createErrorTypeUnmatched(`.e`, `object`, `float64`),
			},
			{
				jsonpath:    `$['a','b']['c','d'].e`,
				inputJSON:   `{"a":{"a":1},"b":{"c":2}}`,
				expectedErr: createErrorTypeUnmatched(`.e`, `object`, `float64`),
			},
			{
				jsonpath:    `$['a','b','x']['c','d'].e`,
				inputJSON:   `{"a":{"a":1},"b":{"c":2}}`,
				expectedErr: createErrorTypeUnmatched(`.e`, `object`, `float64`),
			},
		},
		`same-identifiers`: []TestCase{
			{
				jsonpath:     `$['a','a']`,
				inputJSON:    `{"b":2,"a":1}`,
				expectedJSON: `[1,1]`,
			},
			{
				jsonpath:     `$['a','a','b','b']`,
				inputJSON:    `{"b":2,"a":1}`,
				expectedJSON: `[1,1,2,2]`,
			},
		},
		`child`: []TestCase{
			{
				jsonpath:     `$[0]['a','b']`,
				inputJSON:    `[{"a":1,"b":2},{"a":3,"b":4},{"a":5,"b":6}]`,
				expectedJSON: `[1,2]`,
			},
			{
				jsonpath:     `$[0]['b','a']`,
				inputJSON:    `[{"a":1,"b":2},{"a":3,"b":4},{"a":5,"b":6}]`,
				expectedJSON: `[2,1]`,
			},
			{
				jsonpath:     `$[0:2]['b','a']`,
				inputJSON:    `[{"a":1,"b":2},{"a":3,"b":4},{"a":5,"b":6}]`,
				expectedJSON: `[2,1,4,3]`,
			},
		},
		`empty-input`: []TestCase{
			{
				jsonpath:    `$['a','b']`,
				inputJSON:   `{}`,
				expectedErr: createErrorMemberNotExist(`['a','b']`),
			},
			{
				jsonpath:    `$['a','b']`,
				inputJSON:   `[]`,
				expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `[]interface {}`),
			},
		},
		`type-unmatched`: []TestCase{
			{
				jsonpath:    `$['a','b']`,
				inputJSON:   `"abc"`,
				expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `string`),
			},
			{
				jsonpath:    `$['a','b']`,
				inputJSON:   `123`,
				expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `float64`),
			},
			{
				jsonpath:    `$['a','b']`,
				inputJSON:   `true`,
				expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `bool`),
			},
			{
				jsonpath:    `$['a','b']`,
				inputJSON:   `null`,
				expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `null`),
			},
			{
				jsonpath:    `$['a','b']`,
				inputJSON:   `[1,2,3]`,
				expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `[]interface {}`),
			},
		},
		`child-error`: []TestCase{
			{
				jsonpath:    `$['a','b'].a.b`,
				inputJSON:   `{"c":{"b":1}}`,
				expectedErr: createErrorMemberNotExist(`['a','b']`),
			},
			{
				jsonpath:    `$['a','b'].a.b`,
				inputJSON:   `{"a":{"b":1}}`,
				expectedErr: createErrorMemberNotExist(`.a`),
			},
			{
				jsonpath:    `$['a','b'].a.b.c`,
				inputJSON:   `{"a":{"b":1},"b":{"a":2}}`,
				expectedErr: createErrorTypeUnmatched(`.b`, `object`, `float64`),
			},
			{
				jsonpath:    `$['a','b'].a.b.c`,
				inputJSON:   `{"a":{"a":1},"b":{"a":{"c":2}}}`,
				expectedErr: createErrorMemberNotExist(`.b`),
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_bracketNotation_wildcard(t *testing.T) {
	testGroups := TestGroup{
		`array`: []TestCase{
			{
				jsonpath:     `$[*]`,
				inputJSON:    `["a",123,true,{"b":"c"},[0,1],null]`,
				expectedJSON: `["a",123,true,{"b":"c"},[0,1],null]`,
			},
		},
		`object`: []TestCase{
			{
				jsonpath:     `$[*]`,
				inputJSON:    `{"a":[1],"b":[2,3]}`,
				expectedJSON: `[[1],[2,3]]`,
			},
		},
		`identifier-union`: []TestCase{
			{
				jsonpath:     `$['a',*]`,
				inputJSON:    `{"a":[1],"b":[2,3]}`,
				expectedJSON: `[[1],[1],[2,3]]`,
			},
			{
				jsonpath:     `$[*,'a']`,
				inputJSON:    `{"a":[1],"b":[2,3]}`,
				expectedJSON: `[[1],[2,3],[1]]`,
			},
			{
				jsonpath:     `$[*,*,*]`,
				inputJSON:    `{"a":[1],"b":[2,3]}`,
				expectedJSON: `[[1],[2,3],[1],[2,3],[1],[2,3]]`,
			},
			{
				jsonpath:     `$['a',*,*]`,
				inputJSON:    `{"a":[1],"b":[2,3]}`,
				expectedJSON: `[[1],[1],[2,3],[1],[2,3]]`,
			},
			{
				jsonpath:     `$[*,'a',*]`,
				inputJSON:    `{"a":[1],"b":[2,3]}`,
				expectedJSON: `[[1],[2,3],[1],[1],[2,3]]`,
			},
			{
				jsonpath:     `$['a','a',*]`,
				inputJSON:    `{"a":[1],"b":[2,3]}`,
				expectedJSON: `[[1],[1],[1],[2,3]]`,
			},
			{
				jsonpath:     `$[*,*,'a']`,
				inputJSON:    `{"a":[1],"b":[2,3]}`,
				expectedJSON: `[[1],[2,3],[1],[2,3],[1]]`,
			},
			{
				jsonpath:     `$['a',*,'a']`,
				inputJSON:    `{"a":[1],"b":[2,3]}`,
				expectedJSON: `[[1],[1],[2,3],[1]]`,
			},
			{
				jsonpath:     `$[*,'a','a']`,
				inputJSON:    `{"a":[1],"b":[2,3]}`,
				expectedJSON: `[[1],[2,3],[1],[1]]`,
			},
			{
				jsonpath:     `$['a','a','a']`,
				inputJSON:    `{"a":[1],"b":[2,3]}`,
				expectedJSON: `[[1],[1],[1]]`,
			},
		},
		`empty-input`: []TestCase{
			{
				jsonpath:    `$[*]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[*]`,
				inputJSON:   `{}`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
		},
		`apply-to-value-group`: []TestCase{
			{
				jsonpath:     `$[0:2][*]`,
				inputJSON:    `[[1,2],[3,4],[5,6]]`,
				expectedJSON: `[1,2,3,4]`,
			},
		},
		`child-after-wildcard`: []TestCase{
			{
				jsonpath:     `$[*].a`,
				inputJSON:    `[{"a":1},{"b":2}]`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$[*].a`,
				inputJSON:    `[{"a":1},{"a":1}]`,
				expectedJSON: `[1,1]`,
			},
			{
				jsonpath:     `$[*].a`,
				inputJSON:    `[{"a":[1,[2]]},{"a":2}]`,
				expectedJSON: `[[1,[2]],2]`,
			},
			{
				jsonpath:     `$[*].a[*]`,
				inputJSON:    `[{"a":[1,[2]]},{"a":2}]`,
				expectedJSON: `[1,[2]]`,
			},
		},
		`child-error`: []TestCase{
			{
				jsonpath:    `$[*].a.b`,
				inputJSON:   `{"a":{"b":1}}`,
				expectedErr: createErrorMemberNotExist(`.a`),
			},
			{
				jsonpath:    `$[*].a.b`,
				inputJSON:   `[{"b":1}]`,
				expectedErr: createErrorMemberNotExist(`.a`),
			},
			{
				jsonpath:    `$[*].a.b.c`,
				inputJSON:   `{"a":{"b":1},"b":{"a":2}}`,
				expectedErr: createErrorTypeUnmatched(`.b`, `object`, `float64`),
			},
			{
				jsonpath:    `$[*].a.b.c`,
				inputJSON:   `[{"b":1},{"a":2}]`,
				expectedErr: createErrorTypeUnmatched(`.b`, `object`, `float64`),
			},
			{
				jsonpath:    `$[*].a.b.c`,
				inputJSON:   `{"a":{"a":1},"b":{"a":{"c":2}}}`,
				expectedErr: createErrorMemberNotExist(`.b`),
			},
			{
				jsonpath:    `$[*].a.b.c`,
				inputJSON:   `[{"a":1},{"a":{"c":2}}]`,
				expectedErr: createErrorMemberNotExist(`.b`),
			},
		}}

	runTestGroups(t, testGroups)
}

func TestRetrieve_valueType(t *testing.T) {
	testGroups := TestGroup{
		`root`: []TestCase{
			{
				jsonpath:     `$`,
				inputJSON:    `"a"`,
				expectedJSON: `["a"]`,
			},
			{
				jsonpath:     `$`,
				inputJSON:    `2`,
				expectedJSON: `[2]`,
			},
			{
				jsonpath:     `$`,
				inputJSON:    `false`,
				expectedJSON: `[false]`,
			},
			{
				jsonpath:     `$`,
				inputJSON:    `true`,
				expectedJSON: `[true]`,
			},
			{
				jsonpath:     `$`,
				inputJSON:    `null`,
				expectedJSON: `[null]`,
			},
			{
				jsonpath:     `$`,
				inputJSON:    `{}`,
				expectedJSON: `[{}]`,
			},
			{
				jsonpath:     `$`,
				inputJSON:    `[]`,
				expectedJSON: `[[]]`,
			},
			{
				jsonpath:     `$`,
				inputJSON:    `[1]`,
				expectedJSON: `[[1]]`,
			},
		},
		`child`: []TestCase{
			{
				jsonpath:     `$.a`,
				inputJSON:    `{"a":"string"}`,
				expectedJSON: `["string"]`,
			},
			{
				jsonpath:     `$.a`,
				inputJSON:    `{"a":123}`,
				expectedJSON: `[123]`,
			},
			{
				jsonpath:     `$.a`,
				inputJSON:    `{"a":-123.456}`,
				expectedJSON: `[-123.456]`,
			},
			{
				jsonpath:     `$.a`,
				inputJSON:    `{"a":true}`,
				expectedJSON: `[true]`,
			},
			{
				jsonpath:     `$.a`,
				inputJSON:    `{"a":false}`,
				expectedJSON: `[false]`,
			},
			{
				jsonpath:     `$.a`,
				inputJSON:    `{"a":null}`,
				expectedJSON: `[null]`,
			},
			{
				jsonpath:     `$.a`,
				inputJSON:    `{"a":{"b":"c"}}`,
				expectedJSON: `[{"b":"c"}]`,
			},
			{
				jsonpath:     `$.a`,
				inputJSON:    `{"a":[1,3,5]}`,
				expectedJSON: `[[1,3,5]]`,
			},
			{
				jsonpath:     `$.a`,
				inputJSON:    `{"a":{}}`,
				expectedJSON: `[{}]`,
			},
			{
				jsonpath:     `$.a`,
				inputJSON:    `{"a":[]}`,
				expectedJSON: `[[]]`,
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_arrayIndex(t *testing.T) {
	testGroups := TestGroup{
		`basic::number-variation-plus`: []TestCase{
			{
				jsonpath:     `$[0]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first"]`,
			},
			{
				jsonpath:     `$[1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second"]`,
			},
			{
				jsonpath:     `$[2]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third"]`,
			},
			{
				jsonpath:    `$[3]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[3]`),
			},
		},
		`basic::number-variation-minus`: []TestCase{
			{
				jsonpath:     `$[-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third"]`,
			},
			{
				jsonpath:     `$[-2]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second"]`,
			},
			{
				jsonpath:     `$[-3]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first"]`,
			},
			{
				jsonpath:    `$[-4]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[-4]`),
			},
		},
		`syntax-check::number`: []TestCase{
			{
				jsonpath:     `$[+1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second"]`,
			},
			{
				jsonpath:     `$[01]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second"]`,
			},
			{
				jsonpath:    `$[1.0]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[1.0]`},
			},
			{
				jsonpath:    `$[0,]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[0,]`},
			},
			{
				jsonpath:    `$[0,a]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[0,a]`},
			},
			{
				jsonpath:    `$[a:1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[a:1]`},
			},
			{
				jsonpath:    `$[0:a]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[0:a]`},
			},
		},
		`levels`: []TestCase{
			{
				jsonpath:     `$[0][1]`,
				inputJSON:    `[["a","b"],["c"],["d"]]`,
				expectedJSON: `["b"]`,
			},
			{
				jsonpath:     `$[0][1][2]`,
				inputJSON:    `[["a",["b","c","d"]],["e"],["f"]]`,
				expectedJSON: `["d"]`,
			},
		},
		`empty-input`: []TestCase{
			{
				jsonpath:    `$[0]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0]`),
			},
			{
				jsonpath:    `$[1]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[1]`),
			},
			{
				jsonpath:    `$[-1]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[-1]`),
			},
		},
		`big-number`: []TestCase{
			{
				jsonpath:    `$[1000000000000000000]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[1000000000000000000]`),
			},
			{
				jsonpath:    `$[-1000000000000000000]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[-1000000000000000000]`),
			},
		},
		`not-array`: []TestCase{
			{
				jsonpath:    `$[0]`,
				inputJSON:   `{"a":1,"b":2}`,
				expectedErr: createErrorTypeUnmatched(`[0]`, `array`, `map[string]interface {}`),
			},
			{
				jsonpath:    `$[0]`,
				inputJSON:   `"abc"`,
				expectedErr: createErrorTypeUnmatched(`[0]`, `array`, `string`),
			},
			{
				jsonpath:    `$[0]`,
				inputJSON:   `123`,
				expectedErr: createErrorTypeUnmatched(`[0]`, `array`, `float64`),
			},
			{
				jsonpath:    `$[0]`,
				inputJSON:   `true`,
				expectedErr: createErrorTypeUnmatched(`[0]`, `array`, `bool`),
			},
			{
				jsonpath:    `$[0]`,
				inputJSON:   `null`,
				expectedErr: createErrorTypeUnmatched(`[0]`, `array`, `null`),
			},
			{
				jsonpath:    `$[0]`,
				inputJSON:   `{}`,
				expectedErr: createErrorTypeUnmatched(`[0]`, `array`, `map[string]interface {}`),
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_arrayUnion(t *testing.T) {
	testGroups := TestGroup{
		`index`: []TestCase{
			{
				jsonpath:     `$[0,0]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","first"]`,
			},
			{
				jsonpath:     `$[0,1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second"]`,
			},
			{
				jsonpath:     `$[0,-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","third"]`,
			},
			{
				jsonpath:     `$[2,0,1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third","first","second"]`,
			},
		},
		`wildcard`: []TestCase{
			{
				jsonpath:     `$[0,*]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","first","second","third"]`,
			},
			{
				jsonpath:     `$[*,0]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third","first"]`,
			},
			{
				jsonpath:     `$[1:2,*]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second","first","second","third"]`,
			},
			{
				jsonpath:     `$[*,1:2]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third","second"]`,
			},
			{
				jsonpath:     `$[*,*]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third","first","second","third"]`,
			},
			{
				jsonpath:     `$[*,*,*]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third","first","second","third","first","second","third"]`,
			},
			{
				jsonpath:     `$[0,*,*]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","first","second","third","first","second","third"]`,
			},
			{
				jsonpath:     `$[*,0,*]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third","first","first","second","third"]`,
			},
			{
				jsonpath:     `$[0,0,*]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","first","first","second","third"]`,
			},
			{
				jsonpath:     `$[*,*,0]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third","first","second","third","first"]`,
			},
			{
				jsonpath:     `$[0,*,0]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","first","second","third","first"]`,
			},
			{
				jsonpath:     `$[*,0,0]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third","first","first"]`,
			},
			{
				jsonpath:     `$[0,0,0]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","first","first"]`,
			},
			{
				jsonpath:     `$[*,*].a`,
				inputJSON:    `[{"a":1},{"b":2}]`,
				expectedJSON: `[1,1]`,
			},
		},
		`slice`: []TestCase{
			{
				jsonpath:     `$[1:2,0]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second","first"]`,
			},
			{
				jsonpath:     `$[:2,0]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","first"]`,
			},
			{
				jsonpath:     `$[1:2,0:2]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second","first","second"]`,
			},
		},
		`not-exist`: []TestCase{
			{
				jsonpath:     `$[0,3]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first"]`,
			},
			{
				jsonpath:    `$[3,3]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[3,3]`),
			},
		},
		`array`: []TestCase{
			{
				jsonpath:     `$[0,1]`,
				inputJSON:    `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
				expectedJSON: `[["11","12","13"],["21","22","23"]]`,
			},
		},
		`child-error`: []TestCase{
			{
				jsonpath:    `$[1,2].a.b`,
				inputJSON:   `[0]`,
				expectedErr: createErrorMemberNotExist(`[1,2]`),
			},
			{
				jsonpath:    `$[0,1].a.b`,
				inputJSON:   `[{"b":1}]`,
				expectedErr: createErrorMemberNotExist(`.a`),
			},
			{
				jsonpath:    `$[0,1].a.b`,
				inputJSON:   `[{"b":1},{"c":2}]`,
				expectedErr: createErrorMemberNotExist(`.a`),
			},
			{
				jsonpath:    `$[0,1].a.b.c`,
				inputJSON:   `[{"a":1},{"b":2}]`,
				expectedErr: createErrorTypeUnmatched(`.b`, `object`, `float64`),
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_arraySlice_StartToEnd(t *testing.T) {
	testGroups := TestGroup{
		`start-zero`: []TestCase{
			{
				jsonpath:    `$[0:0]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[0:0]`),
			},
			{
				jsonpath:     `$[0:1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first"]`,
			},
			{
				jsonpath:     `$[0:2]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second"]`,
			},
			{
				jsonpath:     `$[0:3]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third"]`,
			},
		},
		`start-middle`: []TestCase{
			{
				jsonpath:    `$[1:1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[1:1]`),
			},
			{
				jsonpath:     `$[1:2]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second"]`,
			},
			{
				jsonpath:     `$[1:3]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second","third"]`,
			},
		},
		`start-last-forward`: []TestCase{
			{
				jsonpath:    `$[2:2]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[2:2]`),
			},
			{
				jsonpath:     `$[2:3]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third"]`,
			},
		},
		`start-last-backward`: []TestCase{
			{
				jsonpath:    `$[2:1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[2:1]`),
			},
			{
				jsonpath:    `$[2:0]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[2:0]`),
			},
		},
		`start-after-last`: []TestCase{
			{
				jsonpath:    `$[3:2]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[3:2]`),
			},
			{
				jsonpath:    `$[3:3]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[3:3]`),
			},
			{
				jsonpath:    `$[3:4]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[3:4]`),
			},
		},
		`start-minus-to-minus-forward`: []TestCase{
			{
				jsonpath:    `$[-1:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[-1:-1]`),
			},
			{
				jsonpath:     `$[-2:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second"]`,
			},
			{
				jsonpath:     `$[-3:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second"]`,
			},
		},
		`start-minus-to-minus-backward`: []TestCase{
			{
				jsonpath:    `$[-1:-2]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[-1:-2]`),
			},
			{
				jsonpath:    `$[-1:-3]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[-1:-3]`),
			},
		},
		`start-minus-to-plus`: []TestCase{
			{
				jsonpath:    `$[-1:2]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[-1:2]`),
			},
			{
				jsonpath:     `$[-1:3]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third"]`,
			},
			{
				jsonpath:    `$[-2:1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[-2:1]`),
			},
			{
				jsonpath:     `$[-2:2]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second"]`,
			},
			{
				jsonpath:    `$[-3:0]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[-3:0]`),
			},
			{
				jsonpath:     `$[-3:1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first"]`,
			},
			{
				jsonpath:    `$[-4:0]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[-4:0]`),
			},
			{
				jsonpath:     `$[-4:1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first"]`,
			},
			{
				jsonpath:     `$[-4:3]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third"]`,
			},
		},
		`start-zero-to-minus`: []TestCase{
			{
				jsonpath:     `$[0:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second"]`,
			},
			{
				jsonpath:     `$[0:-2]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first"]`,
			},
			{
				jsonpath:    `$[0:-3]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[0:-3]`),
			},
			{
				jsonpath:    `$[0:-4]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[0:-4]`),
			},
		},
		`start-middle-to-minus`: []TestCase{
			{
				jsonpath:     `$[1:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second"]`,
			},
			{
				jsonpath:    `$[1:-2]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[1:-2]`),
			},
		},
		`start-last-to-minus`: []TestCase{
			{
				jsonpath:    `$[2:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[2:-1]`),
			},
			{
				jsonpath:    `$[2:-2]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[2:-2]`),
			},
		},
		`omitted-start`: []TestCase{
			{
				jsonpath:    `$[:0]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[:0]`),
			},
			{
				jsonpath:     `$[:1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first"]`,
			},
			{
				jsonpath:     `$[:2]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second"]`,
			},
			{
				jsonpath:     `$[:3]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third"]`,
			},
			{
				jsonpath:     `$[:4]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third"]`,
			},
			{
				jsonpath:     `$[:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second"]`,
			},
			{
				jsonpath:     `$[:-2]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first"]`,
			},
			{
				jsonpath:    `$[:-3]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[:-3]`),
			},
		},
		`omitted-last`: []TestCase{
			{
				jsonpath:     `$[0:]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third"]`,
			},
			{
				jsonpath:     `$[1:]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second","third"]`,
			},
			{
				jsonpath:     `$[2:]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third"]`,
			},
			{
				jsonpath:    `$[3:]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[3:]`),
			},
			{
				jsonpath:     `$[-1:]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third"]`,
			},
			{
				jsonpath:     `$[-2:]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second","third"]`,
			},
			{
				jsonpath:     `$[-3:]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third"]`,
			},
			{
				jsonpath:     `$[-4:]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third"]`,
			},
		},
		`omitted-start-and-last`: []TestCase{
			{
				jsonpath:     `$[:]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third"]`,
			},
		},
		`big-number`: []TestCase{
			{
				jsonpath:     `$[-1000000000000000000:1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first"]`,
			},
			{
				jsonpath:    `$[1000000000000000000:1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[1000000000000000000:1]`),
			},
			{
				jsonpath:    `$[1:-1000000000000000000]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[1:-1000000000000000000]`),
			},
			{
				jsonpath:     `$[1:1000000000000000000]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second","third"]`,
			},
		},
		`object`: []TestCase{
			{
				jsonpath:    `$[1:2]`,
				inputJSON:   `{"first":1,"second":2,"third":3}`,
				expectedErr: createErrorTypeUnmatched(`[1:2]`, `array`, `map[string]interface {}`),
			},
			{
				jsonpath:    `$[:]`,
				inputJSON:   `{"first":1,"second":2,"third":3}`,
				expectedErr: createErrorTypeUnmatched(`[:]`, `array`, `map[string]interface {}`),
			},
		},
		`syntax`: []TestCase{
			{
				jsonpath:     `$[+0:+1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first"]`,
			},
			{
				jsonpath:     `$[01:02]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second"]`,
			},
			{
				jsonpath:    `$[0.0:2]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[0.0:2]`},
			},
			{
				jsonpath:    `$[0:2.0]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[0:2.0]`},
			},
		},
		`child-error`: []TestCase{
			{
				jsonpath:    `$[1:2].a.b`,
				inputJSON:   `[0]`,
				expectedErr: createErrorMemberNotExist(`[1:2]`),
			},
			{
				jsonpath:    `$[0:2].a.b`,
				inputJSON:   `[{"b":1}]`,
				expectedErr: createErrorMemberNotExist(`.a`),
			},
			{
				jsonpath:    `$[0:2].a.b`,
				inputJSON:   `[{"b":1},{"c":2}]`,
				expectedErr: createErrorMemberNotExist(`.a`),
			},
			{
				jsonpath:    `$[0:2].a.b.c`,
				inputJSON:   `[{"a":1},{"b":2}]`,
				expectedErr: createErrorTypeUnmatched(`.b`, `object`, `float64`),
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_arraySlice_Step(t *testing.T) {
	testGroups := TestGroup{
		`step-variation`: []TestCase{
			{
				jsonpath:     `$[0:3:1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third"]`,
			},
			{
				jsonpath:     `$[0:3:2]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","third"]`,
			},
			{
				jsonpath:     `$[0:3:3]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first"]`,
			},
			{
				jsonpath:     `$[0:2:2]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first"]`,
			},
			{
				jsonpath:     `$[0:2:3]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first"]`,
			},
			{
				jsonpath:     `$[0:1:3]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first"]`,
			},
		},
		`zero`: []TestCase{
			{
				jsonpath:    `$[0:2:0]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[0:2:0]`),
			},
			{
				jsonpath:    `$[2:0:0]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[2:0:0]`),
			},
		},
		`minus::start-variation`: []TestCase{
			{
				jsonpath:    `$[-3:1:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[-3:1:-1]`),
			},
			{
				jsonpath:    `$[-2:1:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[-2:1:-1]`),
			},
			{
				jsonpath:     `$[-1:1:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third"]`,
			},
			{
				jsonpath:    `$[0:1:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[0:1:-1]`),
			},
			{
				jsonpath:    `$[1:1:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[1:1:-1]`),
			},
			{
				jsonpath:     `$[2:1:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third"]`,
			},
			{
				jsonpath:     `$[3:1:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third"]`,
			},
			{
				jsonpath:     `$[4:1:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third"]`,
			},
		},
		`minus::end-variation::start0`: []TestCase{
			{
				jsonpath:    `$[0:-2:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[0:-2:-1]`),
			},
			{
				jsonpath:    `$[0:-1:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[0:-1:-1]`),
			},
			{
				jsonpath:    `$[0:0:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[0:0:-1]`),
			},
			{
				jsonpath:    `$[0:1:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[0:1:-1]`),
			},
			{
				jsonpath:    `$[0:2:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[0:2:-1]`),
			},
			{
				jsonpath:    `$[0:3:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[0:3:-1]`),
			},
			{
				jsonpath:    `$[0:4:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[0:4:-1]`),
			},
		},
		`minus::end-variation::start1`: []TestCase{
			{
				jsonpath:     `$[1:-5:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second","first"]`,
			},
			{
				jsonpath:     `$[1:-4:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second","first"]`,
			},
			{
				jsonpath:     `$[1:-3:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second"]`,
			},
			{
				jsonpath:    `$[1:-2:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[1:-2:-1]`),
			},
			{
				jsonpath:    `$[1:-1:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[1:-1:-1]`),
			},
			{
				jsonpath:     `$[1:0:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second"]`,
			},
			{
				jsonpath:    `$[1:1:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[1:1:-1]`),
			},
			{
				jsonpath:     `$[1:2:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second","first"]`,
				expectedErr:  createErrorMemberNotExist(`[1:2:-1]`),
			},
			{
				jsonpath:    `$[1:3:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[1:3:-1]`),
			},
		},
		`minus::end-variation::start2`: []TestCase{
			{
				jsonpath:     `$[2:-5:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third","second","first"]`,
			},
			{
				jsonpath:     `$[2:-4:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third","second","first"]`,
			},
			{
				jsonpath:     `$[2:-3:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third","second"]`,
			},
			{
				jsonpath:     `$[2:-2:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third"]`,
			},
			{
				jsonpath:    `$[2:-1:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[2:-1:-1]`),
			},
			{
				jsonpath:     `$[2:0:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third","second"]`,
			},
			{
				jsonpath:     `$[2:1:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third"]`,
			},
			{
				jsonpath:    `$[2:2:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[2:2:-1]`),
			},
			{
				jsonpath:    `$[2:3:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[2:3:-1]`),
			},
			{
				jsonpath:    `$[2:4:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[2:4:-1]`),
			},
			{
				jsonpath:    `$[2:5:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[2:5:-1]`),
			},
		},
		`minus::step-variation`: []TestCase{
			{
				jsonpath:     `$[2:0:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third","second"]`,
			},
			{
				jsonpath:     `$[2:0:-2]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third"]`,
			},
			{
				jsonpath:     `$[2:0:-3]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third"]`,
			},
		},
		`minus::start-end-variation`: []TestCase{
			{
				jsonpath:    `$[2:-1:-2]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: createErrorMemberNotExist(`[2:-1:-2]`),
			},
			{
				jsonpath:     `$[-1:0:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third","second"]`,
			},
		},
		`omitted-number`: []TestCase{
			{
				jsonpath:     `$[0:3:]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third"]`,
			},
			{
				jsonpath:     `$[1::1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second","third"]`,
			},
			{
				jsonpath:     `$[1::-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second","first"]`,
			},
			{
				jsonpath:     `$[:1:1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first"]`,
			},
			{
				jsonpath:     `$[:1:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third"]`,
			},
			{
				jsonpath:     `$[::2]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","third"]`,
			},
			{
				jsonpath:     `$[::-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third","second","first"]`,
			},
			{
				jsonpath:     `$[::]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third"]`,
			},
		},
		`big-number`: []TestCase{
			{
				jsonpath:     `$[1:1000000000000000000:1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second","third"]`,
			},
			{
				jsonpath:     `$[1:-1000000000000000000:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second","first"]`,
			},
			{
				jsonpath:     `$[-1000000000000000000:3:1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third"]`,
			},
			{
				jsonpath:     `$[1000000000000000000:0:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third","second"]`,
			},
			{
				jsonpath:     `$[1:0:-1000000000000000000]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second"]`,
			},
			{
				jsonpath:     `$[0:1:1000000000000000000]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first"]`,
			},
		},
		`syntax`: []TestCase{
			{
				jsonpath:     `$[0:3:+1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third"]`,
			},
			{
				jsonpath:     `$[0:3:01]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third"]`,
			},
			{
				jsonpath:    `$[0:3:1.0]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[0:3:1.0]`},
			},
		},
		`not-array`: []TestCase{
			{
				jsonpath:    `$[2:1:-1]`,
				inputJSON:   `{"first":1,"second":2,"third":3}`,
				expectedErr: createErrorTypeUnmatched(`[2:1:-1]`, `array`, `map[string]interface {}`),
			},
			{
				jsonpath:    `$[::-1]`,
				inputJSON:   `{"first":1,"second":2,"third":3}`,
				expectedErr: createErrorTypeUnmatched(`[::-1]`, `array`, `map[string]interface {}`),
			},
			{
				jsonpath:    `$[2:1:-1]`,
				inputJSON:   `"value"`,
				expectedErr: createErrorTypeUnmatched(`[2:1:-1]`, `array`, `string`),
			},
			{
				jsonpath:    `$[2:1:-1]`,
				inputJSON:   `1`,
				expectedErr: createErrorTypeUnmatched(`[2:1:-1]`, `array`, `float64`),
			},
			{
				jsonpath:    `$[2:1:-1]`,
				inputJSON:   `true`,
				expectedErr: createErrorTypeUnmatched(`[2:1:-1]`, `array`, `bool`),
			},
			{
				jsonpath:    `$[2:1:-1]`,
				inputJSON:   `null`,
				expectedErr: createErrorTypeUnmatched(`[2:1:-1]`, `array`, `null`),
			},
		},
		`child-error`: []TestCase{
			{
				jsonpath:    `$[-1:-1:-1].a.b`,
				inputJSON:   `[0]`,
				expectedErr: createErrorMemberNotExist(`[-1:-1:-1]`),
			},
			{
				jsonpath:    `$[0:-2:-1].a.b`,
				inputJSON:   `[{"b":1}]`,
				expectedErr: createErrorMemberNotExist(`.a`),
			},
			{
				jsonpath:    `$[1:-3:-1].a.b`,
				inputJSON:   `[{"b":1},{"c":2}]`,
				expectedErr: createErrorMemberNotExist(`.a`),
			},
			{
				jsonpath:    `$[1:-3:-1].a.b.c`,
				inputJSON:   `[{"a":1},{"b":2}]`,
				expectedErr: createErrorTypeUnmatched(`.b`, `object`, `float64`),
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_filterExist(t *testing.T) {
	testGroups := TestGroup{
		`current-root`: []TestCase{
			{
				jsonpath:     `$[?(@)]`,
				inputJSON:    `["a","b"]`,
				expectedJSON: `["a","b"]`,
			},
			{
				jsonpath:    `$[?(!@)]`,
				inputJSON:   `["a","b"]`,
				expectedErr: createErrorMemberNotExist(`[?(!@)]`),
			},
		},
		`child`: []TestCase{
			{
				jsonpath:     `$[?(@.a)]`,
				inputJSON:    `[{"b":2},{"a":1},{"a":"value"},{"a":""},{"a":true},{"a":false},{"a":null},{"a":{}},{"a":[]}]`,
				expectedJSON: `[{"a":1},{"a":"value"},{"a":""},{"a":true},{"a":false},{"a":null},{"a":{}},{"a":[]}]`,
			},
			{
				jsonpath:     `$[?(!@.a)]`,
				inputJSON:    `[{"b":2},{"a":1},{"a":"value"},{"a":""},{"a":true},{"a":false},{"a":null},{"a":{}},{"a":[]}]`,
				expectedJSON: `[{"b":2}]`,
			},
			{
				jsonpath:    `$[?(@.c)]`,
				inputJSON:   `[{"a":1},{"b":2}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.c)]`),
			},
			{
				jsonpath:     `$[?(!@.c)]`,
				inputJSON:    `[{"a":1},{"b":2}]`,
				expectedJSON: `[{"a":1},{"b":2}]`,
			},
		},
		`index`: []TestCase{
			{
				jsonpath:     `$[?(@[1])]`,
				inputJSON:    `[[{"a":1}],[{"b":2},{"c":3}],[],{"d":4}]`,
				expectedJSON: `[[{"b":2},{"c":3}]]`,
			},
			{
				jsonpath:     `$[?(!@[1])]`,
				inputJSON:    `[[{"a":1}],[{"b":2},{"c":3}],[],{"d":4}]`,
				expectedJSON: `[[{"a":1}],[],{"d":4}]`,
			},
		},
		`slice`: []TestCase{
			{
				jsonpath:     `$[?(@[1:3])]`,
				inputJSON:    `[[{"a":1}],[{"b":2},{"c":3}],[],{"d":4}]`,
				expectedJSON: `[[{"b":2},{"c":3}]]`,
			},
			{
				jsonpath:     `$[?(!@[1:3])]`,
				inputJSON:    `[[{"a":1}],[{"b":2},{"c":3}],[],{"d":4}]`,
				expectedJSON: `[[{"a":1}],[],{"d":4}]`,
			},
			{
				jsonpath:     `$[?(@[1:3])]`,
				inputJSON:    `[[{"a":1}],[{"b":2},{"c":3},{"e":5}],[],{"d":4}]`,
				expectedJSON: `[[{"b":2},{"c":3},{"e":5}]]`,
			},
			{
				jsonpath:     `$[?(!@[1:3])]`,
				inputJSON:    `[[{"a":1}],[{"b":2},{"c":3},{"e":5}],[],{"d":4}]`,
				expectedJSON: `[[{"a":1}],[],{"d":4}]`,
			},
			{
				jsonpath:     `$[?(@[0:1])]`,
				inputJSON:    `[[{"a":1}],[]]`,
				expectedJSON: `[[{"a":1}]]`,
			},
		},
		`object::root`: []TestCase{
			{
				jsonpath:     `$[?($)]`,
				inputJSON:    `{"a":1,"b":2}`,
				expectedJSON: `[1,2]`,
			},
			{
				jsonpath:    `$[?(!$)]`,
				inputJSON:   `{"a":1,"b":2}`,
				expectedErr: createErrorMemberNotExist(`[?(!$)]`),
			},
		},
		`object::current-root`: []TestCase{
			{
				jsonpath:     `$[?(@)]`,
				inputJSON:    `{"a":1,"b":2}`,
				expectedJSON: `[1,2]`,
			},
			{
				jsonpath:    `$[?(!@)]`,
				inputJSON:   `{"a":1}`,
				expectedErr: createErrorMemberNotExist(`[?(!@)]`),
			},
		},
		`object::child`: []TestCase{
			{
				jsonpath:     `$[?(@.a1)]`,
				inputJSON:    `{"a":{"a1":1},"b":{"b1":2}}`,
				expectedJSON: `[{"a1":1}]`,
			},
			{
				jsonpath:     `$[?(!@.a1)]`,
				inputJSON:    `{"a":{"a1":1},"b":{"b1":2}}`,
				expectedJSON: `[{"b1":2}]`,
			},
		},
		`recursive`: []TestCase{
			{
				jsonpath:     `$[?(@..a)]`,
				inputJSON:    `[{"a":1},{"b":2},{"c":{"a":3}},{"a":{"a":4}}]`,
				expectedJSON: `[{"a":1},{"c":{"a":3}},{"a":{"a":4}}]`,
			},
			{
				jsonpath:     `$[?(!@..a)]`,
				inputJSON:    `[{"a":1},{"b":2},{"c":{"a":3}},{"a":{"a":4}}]`,
				expectedJSON: `[{"b":2}]`,
			},
		},
		`object::index`: []TestCase{
			{
				jsonpath:     `$[?(@[1])]`,
				inputJSON:    `{"a":["a1"],"b":["b1","b2"],"c":[],"d":4}`,
				expectedJSON: `[["b1","b2"]]`,
			},
			{
				jsonpath:     `$[?(!@[1])]`,
				inputJSON:    `{"a":["a1"],"b":["b1","b2"],"c":[],"d":4}`,
				expectedJSON: `[["a1"],[],4]`,
			},
		},
		`object::slice`: []TestCase{
			{
				jsonpath:     `$[?(@[1:3])]`,
				inputJSON:    `{"a":[],"b":[2],"c":[3,4,5,6],"d":4}`,
				expectedJSON: `[[3,4,5,6]]`,
			},
			{
				jsonpath:     `$[?(!@[1:3])]`,
				inputJSON:    `{"a":[],"b":[2],"c":[3,4,5,6],"d":4}`,
				expectedJSON: `[[],[2],4]`,
			},
			{
				jsonpath:     `$[?(@[1:3])]`,
				inputJSON:    `{"a":[],"b":[2],"c":[3,4],"d":4}`,
				expectedJSON: `[[3,4]]`,
			},
			{
				jsonpath:     `$[?(!@[1:3])]`,
				inputJSON:    `{"a":[],"b":[2],"c":[3,4],"d":4}`,
				expectedJSON: `[[],[2],4]`,
			},
		},
		`wildcard-identifier`: []TestCase{
			{
				jsonpath:     `$.*[?(@.a)]`,
				inputJSON:    `[[{"a":1},{"b":2}],[{"c":1},{"d":2}]]`,
				expectedJSON: `[{"a":1}]`,
			},
			{
				jsonpath:    `$.*[?(@.a)]`,
				inputJSON:   `[{"a":1},{"b":2}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.a)]`),
			},
			{
				jsonpath:     `$.*[?(@.a)]`,
				inputJSON:    `{"a":[{"a":1}],"b":[{"b":2}]}`,
				expectedJSON: `[{"a":1}]`,
			},
			{
				jsonpath:    `$.*[?(@.a)]`,
				inputJSON:   `{"a":{"a":1},"b":{"b":2}}`,
				expectedErr: createErrorMemberNotExist(`[?(@.a)]`),
			},
		},
		`root`: []TestCase{
			{
				jsonpath:     `$[?($)]`,
				inputJSON:    `[{"a":1},{"b":2}]`,
				expectedJSON: `[{"a":1},{"b":2}]`,
			},
			{
				jsonpath:     `$[?($[0].a)]`,
				inputJSON:    `[{"a":1},{"b":2}]`,
				expectedJSON: `[{"a":1},{"b":2}]`,
			},
			{
				jsonpath:    `$[?(!$[0].a)]`,
				inputJSON:   `[{"a":1},{"b":2}]`,
				expectedErr: createErrorMemberNotExist(`[?(!$[0].a)]`),
			},
		},
		`multi-identifier`: []TestCase{
			{
				jsonpath:     `$[?(@['a','b'])]`,
				inputJSON:    `[{"a":1},{"b":2}]`,
				expectedJSON: `[{"a":1},{"b":2}]`,
			},
			{
				jsonpath:     `$[?(@['a','c'])]`,
				inputJSON:    `[{"a":1},{"b":2}]`,
				expectedJSON: `[{"a":1}]`,
			},
			{
				jsonpath:    `$[?(@['c','d'])]`,
				inputJSON:   `[{"a":1},{"b":2}]`,
				expectedErr: createErrorMemberNotExist(`[?(@['c','d'])]`),
			},
		},
		`current-wildcard`: []TestCase{
			{
				jsonpath:     `$[?(@.*)]`,
				inputJSON:    `[{"a":1},{"b":2}]`,
				expectedJSON: `[{"a":1},{"b":2}]`,
			},
			{
				jsonpath:    `$[?(@.*)]`,
				inputJSON:   `[1,2]`,
				expectedErr: createErrorMemberNotExist(`[?(@.*)]`),
			},
		},
		`wildcard-qualifier`: []TestCase{
			{
				jsonpath:     `$[?(@[*])]`,
				inputJSON:    `[[{"a":1}],[]]`,
				expectedJSON: `[[{"a":1}]]`,
			},
			{
				jsonpath:     `$[?(@[*])]`,
				inputJSON:    `[{"a":1},{}]`,
				expectedJSON: `[{"a":1}]`,
			},
			{
				jsonpath:    `$[?(@[*])]`,
				inputJSON:   `[1,2]`,
				expectedErr: createErrorMemberNotExist(`[?(@[*])]`),
			},
		},
		`union`: []TestCase{
			{
				jsonpath:     `$[?(@[0,1])]`,
				inputJSON:    `[[{"a":1}],[0,1]]`,
				expectedJSON: `[[{"a":1}],[0,1]]`,
			},
			{
				jsonpath:     `$[?(@[0,1])]`,
				inputJSON:    `[[{"a":1}],[]]`,
				expectedJSON: `[[{"a":1}]]`,
			},
		},
		`filter`: []TestCase{
			{
				jsonpath:     `$[?(@.a[?(@.b)])]`,
				inputJSON:    `[{"a":[{"b":2},{"c":3}]},{"b":4}]`,
				expectedJSON: `[{"a":[{"b":2},{"c":3}]}]`,
			},
			{
				jsonpath:     `$[?(@.a[?(@.b > 1)])]`,
				inputJSON:    `[{"a":[{"b":1},{"c":3}]},{"a":[{"b":2},{"c":5}]},{"b":4}]`,
				expectedJSON: `[{"a":[{"b":2},{"c":5}]}]`,
			},
		},
		`sub-query`: []TestCase{
			{
				jsonpath:     `$[?((@.a>1))]`,
				inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
				expectedJSON: `[{"a":2},{"a":3}]`,
			},
			{
				jsonpath:     `$[?(((@.a>1)))]`,
				inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
				expectedJSON: `[{"a":2},{"a":3}]`,
			},
		},
		`filter-chain`: []TestCase{
			{
				jsonpath:     `$[?(@)][?(@)]`,
				inputJSON:    `[1,[21,[221,[222]]]]`,
				expectedJSON: `[21,[221,[222]]]`,
			},
			{
				jsonpath:     `$[?(@)][?(@)][?(@)]`,
				inputJSON:    `[1,[21,[221,[222]]]]`,
				expectedJSON: `[221,[222]]`,
			},
			{
				jsonpath:     `$[?(@)][?(@)][?(@)][?(@)]`,
				inputJSON:    `[1,[21,[221,[222]]]]`,
				expectedJSON: `[222]`,
			},
		},
		`child-error::array`: []TestCase{
			{
				jsonpath:    `$[?(@.a)].b`,
				inputJSON:   `[{"b":1}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.a)]`),
			},
			{
				jsonpath:    `$[?(@.a)].b.c`,
				inputJSON:   `[{"a":1}]`,
				expectedErr: createErrorMemberNotExist(`.b`),
			},
			{
				jsonpath:    `$[?(@.a)].b.c`,
				inputJSON:   `[{"a":1},{"b":2}]`,
				expectedErr: createErrorMemberNotExist(`.b`),
			},
			{
				jsonpath:    `$[?(@.a)].b.c`,
				inputJSON:   `[{"a":1},{"a":1}]`,
				expectedErr: createErrorMemberNotExist(`.b`),
			},
			{
				jsonpath:    `$[?(@.a)].a.b.c`,
				inputJSON:   `[{"a":1},{"a":{"c":2}}]`,
				expectedErr: createErrorMemberNotExist(`.b`),
			},
			{
				jsonpath:    `$[?(@.a)].a.b.c`,
				inputJSON:   `[{"a":1},{"a":{"c":2}},{"b":3}]`,
				expectedErr: createErrorMemberNotExist(`.b`),
			},
		},
		`child-error::object`: []TestCase{
			{
				jsonpath:    `$[?(@.a)].b`,
				inputJSON:   `{"a":{"b":1}}`,
				expectedErr: createErrorMemberNotExist(`[?(@.a)]`),
			},
			{
				jsonpath:    `$[?(@.a)].b.c`,
				inputJSON:   `{"a":{"a":1}}`,
				expectedErr: createErrorMemberNotExist(`.b`),
			},
			{
				jsonpath:    `$[?(@.a)].b.c`,
				inputJSON:   `{"a":{"a":1},"b":{"b":2}}`,
				expectedErr: createErrorMemberNotExist(`.b`),
			},
			{
				jsonpath:    `$[?(@.a)].b.c`,
				inputJSON:   `{"a":{"a":1},"b":{"a":1}}`,
				expectedErr: createErrorMemberNotExist(`.b`),
			},
			{
				jsonpath:    `$[?(@.a)].a.b.c`,
				inputJSON:   `{"a":{"a":1},"b":{"a":{"c":2}}}`,
				expectedErr: createErrorMemberNotExist(`.b`),
			},
			{
				jsonpath:    `$[?(@.a)].a.b.c`,
				inputJSON:   `{"a":{"a":1},"b":{"a":{"c":2}},"c":{"b":3}}`,
				expectedErr: createErrorMemberNotExist(`.b`),
			},
		},
		`root-value-group`: []TestCase{
			{
				jsonpath:     `$.z[?($..x)]`,
				inputJSON:    `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
				expectedJSON: `[{"b":1},{"b":2},{"b":3}]`,
			},
			{
				jsonpath:     `$.z[?($["x","y"])]`,
				inputJSON:    `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
				expectedJSON: `[{"b":1},{"b":2},{"b":3}]`,
			},
			{
				jsonpath:     `$.z[?($.*)]`,
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
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_filterCompare(t *testing.T) {
	testGroups := TestGroup{
		`eq`: []TestCase{
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
			{
				jsonpath:     `$[?(@.a=='ab')]`,
				inputJSON:    `[{"a":"ab"}]`,
				expectedJSON: `[{"a":"ab"}]`,
			},
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
		},
		`ne`: []TestCase{
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
		},
		`gt`: []TestCase{
			{
				jsonpath:     `$[?(@.a < 1)]`,
				inputJSON:    `[{"a":-9999999},{"a":0.999999},{"a":1.0000000},{"a":1.0000001},{"a":2},{"a":"0.9"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
				expectedJSON: `[{"a":-9999999},{"a":0.999999}]`,
			},
			{
				jsonpath:     `$[?(1 > @.a)]`,
				inputJSON:    `[{"a":-9999999},{"a":0.999999},{"a":1.0000000},{"a":1.0000001},{"a":2},{"a":"0.9"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
				expectedJSON: `[{"a":-9999999},{"a":0.999999}]`,
			},
			{
				jsonpath:    `$[?(1 > @.a)]`,
				inputJSON:   `[{"a":1.0000000},{"a":1.0000001},{"a":2},{"a":"0.9"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
				expectedErr: createErrorMemberNotExist(`[?(1 > @.a)]`),
			},
		},
		`ge`: []TestCase{
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
				jsonpath:    `$[?(1.00001 >= @.a)]`,
				inputJSON:   `[{"a":1.00002},{"a":2,"b":4},{"a":"0.9"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
				expectedErr: createErrorMemberNotExist(`[?(1.00001 >= @.a)]`),
			},
		},
		`lt`: []TestCase{
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
		},
		`le`: []TestCase{
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
			{
				jsonpath:    `$[?(1.000001 <= @.a)]`,
				inputJSON:   `[{"a":0},{"a":1},{"a":1.0000009},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
				expectedErr: createErrorMemberNotExist(`[?(1.000001 <= @.a)]`),
			},
		},
		`syntax-check::string-literal`: []TestCase{
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
				// The character ['&','<','>'] is encoded to [\u0026,\u003c,\u003e] using Go's json.Marshal()
				jsonpath:     `$[?(@.a=="~!@#$%^&*()-_=+[]\\{}|;':\",./<>?")]`,
				inputJSON:    `[{"a":"~!@#$%^&*()-_=+[]\\{}|;':\",./<>?"}]`,
				expectedJSON: `[{"a":"~!@#$%^\u0026*()-_=+[]\\{}|;':\",./\u003c\u003e?"}]`,
			},
			{
				// The character ['&','<','>'] is encoded to [\u0026,\u003c,\u003e] using Go's json.Marshal()
				jsonpath:     `$[?(@.a=='~!@#$%^&*()-_=+[]\\{}|;\':",./<>?')]`,
				inputJSON:    `[{"a":"~!@#$%^&*()-_=+[]\\{}|;':\",./<>?"}]`,
				expectedJSON: `[{"a":"~!@#$%^\u0026*()-_=+[]\\{}|;':\",./\u003c\u003e?"}]`,
			},
			{
				// The character \/ is encoded to / using Go's json.Marshal()
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
				jsonpath:     `$[?(@.a=='a\bb')]`,
				inputJSON:    `[{"a":"a\bb"},{"b":1}]`,
				expectedJSON: `[{"a":"a\bb"}]`,
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
				// The character \uABCD is encoded to ꯍ using Go's json.Marshal()
				jsonpath:     `$[?(@.a=='\uABCD')]`,
				inputJSON:    `[{"a":"\uabcd"},{"b":1}]`,
				expectedJSON: `[{"a":"ꯍ"}]`,
			},
			{
				// The character \uabcd is encoded to ꯍ using Go's json.Marshal()
				jsonpath:     `$[?(@.a=='\uabcd')]`,
				inputJSON:    `[{"a":"\uABCD"},{"b":1}]`,
				expectedJSON: `[{"a":"ꯍ"}]`,
			},
		},
		`syntax-check::number-literal`: []TestCase{
			{
				// The number 5.0 is converted to 5 using Go's json.Marshal().
				jsonpath:     `$[?(@.a==5)]`,
				inputJSON:    `[{"a":4.9},{"a":5.0},{"a":5.1},{"a":5},{"a":-5},{"a":"5"},{"a":"a"},{"a":true},{"a":null},{"a":{}},{"a":[]},{"b":5},{"a":{"a":5}},{"a":[{"a":5}]}]`,
				expectedJSON: `[{"a":5},{"a":5}]`,
			},
			{
				// The number 5.00000 is converted to 5 using Go's json.Marshal().
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
				// The number 1.0 is converted to 1 using Go's json.Marshal().
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
				// The number -0.123e2 is converted to -12.3 using Go's json.Marshal().
				jsonpath:     `$[?(@.a==-0.123e2)]`,
				inputJSON:    `[{"a":-12.3,"b":1},{"a":-0.123e2,"b":2},{"a":-0.123},{"a":-12},{"a":12.3},{"a":2},{"a":"-0.123e2"}]`,
				expectedJSON: `[{"a":-12.3,"b":1},{"a":-12.3,"b":2}]`,
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
				jsonpath:     `$[?(@.a==010)]`,
				inputJSON:    `[{"a":10},{"a":0},{"a":"010"},{"a":"10"}]`,
				expectedJSON: `[{"a":10}]`,
			},
		},
		`syntax-check::bool-literal`: []TestCase{
			{
				jsonpath:     `$[?(@.a==false)]`,
				inputJSON:    `[{"a":null},{"a":false},{"a":true},{"a":0},{"a":1},{"a":"false"},{"b":false}]`,
				expectedJSON: `[{"a":false}]`,
			},
			{
				jsonpath:     `$[?(@.a!=false)]`,
				inputJSON:    `[{"a":null},{"a":false},{"a":true},{"a":0},{"a":1},{"a":"false"},{"b":false}]`,
				expectedJSON: `[{"a":null},{"a":true},{"a":0},{"a":1},{"a":"false"},{"b":false}]`,
			},
			{
				jsonpath:     `$[?(@.a==FALSE)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
			},
			{
				jsonpath:     `$[?(@.a==False)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
			},
			{
				jsonpath:     `$[?(@.a==true)]`,
				inputJSON:    `[{"a":null},{"a":false},{"a":true},{"a":0},{"a":1},{"a":"false"},{"b":true}]`,
				expectedJSON: `[{"a":true}]`,
			},
			{
				jsonpath:     `$[?(@.a!=true)]`,
				inputJSON:    `[{"a":null},{"a":false},{"a":true},{"a":0},{"a":1},{"a":"false"},{"b":false}]`,
				expectedJSON: `[{"a":null},{"a":false},{"a":0},{"a":1},{"a":"false"},{"b":false}]`,
			},
			{
				jsonpath:     `$[?(@.a==TRUE)]`,
				inputJSON:    `[{"a":true}]`,
				expectedJSON: `[{"a":true}]`,
			},
			{
				jsonpath:     `$[?(@.a==True)]`,
				inputJSON:    `[{"a":true}]`,
				expectedJSON: `[{"a":true}]`,
			},
		},
		`syntax-check::null-literal`: []TestCase{
			{
				jsonpath:     `$[?(@.a==null)]`,
				inputJSON:    `[{"a":null},{"a":false},{"a":true},{"a":0},{"a":1},{"a":"false"},{"b":null}]`,
				expectedJSON: `[{"a":null}]`,
			},
			{
				jsonpath:     `$[?(@.a!=null)]`,
				inputJSON:    `[{"a":null},{"a":false},{"a":true},{"a":0},{"a":1},{"a":"false"},{"b":null}]`,
				expectedJSON: `[{"a":false},{"a":true},{"a":0},{"a":1},{"a":"false"},{"b":null}]`,
			},
			{
				jsonpath:     `$[?(@.a==NULL)]`,
				inputJSON:    `[{"a":null}]`,
				expectedJSON: `[{"a":null}]`,
			},
			{
				jsonpath:     `$[?(@.a==Null)]`,
				inputJSON:    `[{"a":null}]`,
				expectedJSON: `[{"a":null}]`,
			},
		},
		`syntax-check::jsonpath`: []TestCase{
			{
				jsonpath:     `$[?(@.a\+10==20)]`,
				inputJSON:    `[{"a":10},{"a":20},{"a":30},{"a+10":20}]`,
				expectedJSON: `[{"a+10":20}]`,
			},
			{
				jsonpath:     `$[?(@.a-10==20)]`,
				inputJSON:    `[{"a":10},{"a":20},{"a":30},{"a-10":20}]`,
				expectedJSON: `[{"a-10":20}]`,
			},
			{
				// The number 11.0 is converted to 11 using Go's json.Marshal().
				jsonpath:     `$[?(@.a\*2==11)]`,
				inputJSON:    `[{"a":6},{"a":5},{"a":5.5},{"a":-5},{"a*2":10.999},{"a*2":11.0},{"a*2":11.1},{"a*2":5},{"a*2":"11"}]`,
				expectedJSON: `[{"a*2":11}]`,
			},
			{
				jsonpath:     `$[?(@.a\/10==5)]`,
				inputJSON:    `[{"a":60},{"a":50},{"a":51},{"a":-50},{"a/10":5},{"a/10":"5"}]`,
				expectedJSON: `[{"a/10":5}]`,
			},
			{
				jsonpath:     `$[?(@['a']<2.1)]`,
				inputJSON:    `[{"a":1.9},{"a":2},{"a":2.1},{"a":3},{"a":"test"}]`,
				expectedJSON: `[{"a":1.9},{"a":2}]`,
			},
			{
				jsonpath:     `$[?(@['$a']<2.1)]`,
				inputJSON:    `[{"$a":1.9},{"a":2},{"a":2.1},{"a":3},{"$a":"test"}]`,
				expectedJSON: `[{"$a":1.9}]`,
			},
			{
				jsonpath:     `$[?(@['@a']<2.1)]`,
				inputJSON:    `[{"@a":1.9},{"a":2},{"a":2.1},{"a":3},{"@a":"test"}]`,
				expectedJSON: `[{"@a":1.9}]`,
			},
			{
				jsonpath:     `$[?(@['a==b']<2.1)]`,
				inputJSON:    `[{"a==b":1.9},{"a":2},{"a":2.1},{"b":3},{"a==b":"test"}]`,
				expectedJSON: `[{"a==b":1.9}]`,
			},
			{
				// The character '<' is encoded to \u003c using Go's json.Marshal()
				jsonpath:     `$[?(@['a<=b']<2.1)]`,
				inputJSON:    `[{"a<=b":1.9},{"a":2},{"a":2.1},{"b":3},{"a<=b":"test"}]`,
				expectedJSON: `[{"a\u003c=b":1.9}]`,
			},
			{
				jsonpath:     `$[?(@[-1]==2)]`,
				inputJSON:    `[[0,1],[0,2],[2],["2"],["a","b"],["b"]]`,
				expectedJSON: `[[0,2],[2]]`,
			},
			{
				jsonpath:     `$[?(@.a.b == 1)]`,
				inputJSON:    `[{"a":1},{"a":{"b":1}},{"a":{"a":1}}]`,
				expectedJSON: `[{"a":{"b":1}}]`,
			},
		},
		`jsonpath::start-from-root`: []TestCase{
			{
				jsonpath:     `$[?(@.a == $[2].b)]`,
				inputJSON:    `[{"a":0},{"a":1},{"a":2,"b":1},{"b":1}]`,
				expectedJSON: `[{"a":1}]`,
			},
			{
				jsonpath:     `$[?($[2].b == @.a)]`,
				inputJSON:    `[{"a":0},{"a":1},{"a":2,"b":1},{"b":1}]`,
				expectedJSON: `[{"a":1}]`,
			},
			{
				jsonpath:     `$.a[?(@.b==$.c)]`,
				inputJSON:    `{"a":[{"b":123},{"b":123.456},{"b":"123.456"}],"c":123.456}`,
				expectedJSON: `[{"b":123.456}]`,
			},
			{
				jsonpath:    `$.x[?(@[*]>=$.y[*])]`,
				inputJSON:   `{"x":[[1,2],[3,4],[5,6]],"y":[3,4,5]}`,
				expectedErr: ErrorInvalidSyntax{position: 6, reason: `JSONPath that returns a value group is prohibited`, near: `@[*]>=$.y[*])]`},
			},
			{
				jsonpath:    `$.x[?(@[*]>=$.y.a[0:1])]`,
				inputJSON:   `{"x":[[1,2],[3,4],[5,6]],"y":{"a":[3,4,5]}}`,
				expectedErr: ErrorInvalidSyntax{position: 6, reason: `JSONPath that returns a value group is prohibited`, near: `@[*]>=$.y.a[0:1])]`},
			},
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
			{
				jsonpath:    `$[?(@.b == $[0].a)]`,
				inputJSON:   `[{"a":1},{"a":2}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.b == $[0].a)]`),
			},
			{
				jsonpath:    `$[?($[0].a == @.b)]`,
				inputJSON:   `[{"a":1},{"a":2}]`,
				expectedErr: createErrorMemberNotExist(`[?($[0].a == @.b)]`),
			},
		},
		`child-after-filter`: []TestCase{
			{
				jsonpath:     `$[?(@.a == 2)].b`,
				inputJSON:    `[{"a":0},{"a":1},{"a":2,"b":4}]`,
				expectedJSON: `[4]`,
			},
		},
		`filter-after-recursive`: []TestCase{
			{
				jsonpath:     `$..[?(@.a==2)]`,
				inputJSON:    `{"a":2,"x":[{"a":2},{"b":{"a":2}},{"a":{"a":2}},[{"a":2}]]}`,
				expectedJSON: `[{"a":2},{"a":2},{"a":2},{"a":2}]`,
			},
			{
				jsonpath:    `$..*[?(@.a>2)]`,
				inputJSON:   `[{"b":"1","a":1},{"c":"2","a":2},{"d":"3","a":3}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.a>2)]`),
			},
			{
				jsonpath:     `$..*[?(@.a>2)]`,
				inputJSON:    `{"z":[{"d":"1","a":1},{"c":"2","a":2},{"b":"3","a":3}],"y":{"b":"4","a":4}}`,
				expectedJSON: `[{"a":3,"b":"3"}]`,
			},
			{
				jsonpath:     `$..*[?(@.a>2)]`,
				inputJSON:    `{"x":{"z":[{"x":"1","a":1},{"z":"2","a":2},{"y":"3","a":3}],"y":{"b":"4","a":4}}}`,
				expectedJSON: `[{"a":4,"b":"4"},{"a":3,"y":"3"}]`,
			},
			{
				jsonpath:     `$..*[?(@.a>2)]`,
				inputJSON:    `[{"x":{"z":[{"b":"1","a":1},{"b":"2","a":2},{"b":"3","a":3},{"b":"6","a":6}],"y":{"b":"4","a":4}}},{"b":"5","a":5}]`,
				expectedJSON: `[{"a":4,"b":"4"},{"a":3,"b":"3"},{"a":6,"b":"6"}]`,
			},
		},
		`both-number`: []TestCase{
			{
				jsonpath:     `$[?(10==10)]`,
				inputJSON:    `[{"a":10},{"a":20},{"a":30},{"a+10":20}]`,
				expectedJSON: `[{"a":10},{"a":20},{"a":30},{"a+10":20}]`,
			},
			{
				jsonpath:    `$[?(10==20)]`,
				inputJSON:   `[{"a":10},{"a":20},{"a":30},{"a+10":20}]`,
				expectedErr: createErrorMemberNotExist(`[?(10==20)]`),
			},
		},
		`both-jsonpath`: []TestCase{
			{
				jsonpath:    `$[?(@.a==@.a)]`,
				inputJSON:   `[{"a":10},{"a":20},{"a":30},{"a+10":20}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `comparison between two current nodes is prohibited`, near: `@.a==@.a)]`},
			},
		},
		`value-group-jsonpath::slice-qualifier`: []TestCase{
			{
				jsonpath:    `$[?(@[0:1]==1)]`,
				inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0:1]==1)]`},
			},
			{
				jsonpath:    `$[?(@[0:2]==1)]`,
				inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0:2]==1)]`},
			},
			{
				jsonpath:    `$[?(@[0:2].a==1)]`,
				inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0:2].a==1)]`},
			},
			{
				jsonpath:    `$[?(@.a[0:2]==1)]`,
				inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[0:2]==1)]`},
			},
		},
		`value-group-jsonpath::wildcard-qualifier`: []TestCase{
			{
				jsonpath:    `$[?(@[*]==1)]`,
				inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[*]==1)]`},
			},
			{
				jsonpath:    `$[?(@[*].a==1)]`,
				inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[*].a==1)]`},
			},
			{
				jsonpath:    `$[?(@.a[*]==1)]`,
				inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[*]==1)]`},
			},
		},
		`value-group-jsonpath::union-qualifier`: []TestCase{
			{
				jsonpath:    `$[?(@[0,1]==1)]`,
				inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0,1]==1)]`},
			},
			{
				jsonpath:    `$[?(@[0,1:2]==1)]`,
				inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0,1:2]==1)]`},
			},
			{
				jsonpath:    `$[?(@[0,1].a==1)]`,
				inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0,1].a==1)]`},
			},
			{
				jsonpath:    `$[?(@.a[0,1]==1)]`,
				inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[0,1]==1)]`},
			},
		},
		`value-group-jsonpath::recursive`: []TestCase{
			{
				jsonpath:    `$[?(@..a==123)]`,
				inputJSON:   `[{"a":"123"},{"a":123}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@..a==123)]`},
			},
			{
				jsonpath:    `$[?(@..a.b==123)]`,
				inputJSON:   `[{"a":"123"},{"a":123}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@..a.b==123)]`},
			},
			{
				jsonpath:    `$[?(@.a..b==123)]`,
				inputJSON:   `[{"a":"123"},{"a":123}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a..b==123)]`},
			},
			{
				jsonpath:    `$[?(@..a..b==123)]`,
				inputJSON:   `[{"a":"123"},{"a":123}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@..a..b==123)]`},
			},
		},
		`value-group-jsonpath::multi-identifier`: []TestCase{
			{
				jsonpath:    `$[?(@['a','b']==123)]`,
				inputJSON:   `[{"a":"123"},{"a":123}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@['a','b']==123)]`},
			},
			{
				jsonpath:    `$[?(@['a','b','c']==123)]`,
				inputJSON:   `[{"a":"123"},{"a":123}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@['a','b','c']==123)]`},
			},
			{
				jsonpath:    `$[?(@['a','b']['a']==123)]`,
				inputJSON:   `[{"a":"123"},{"a":123}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@['a','b']['a']==123)]`},
			},
			{
				jsonpath:    `$[?(@['a']['a','b']==123)]`,
				inputJSON:   `[{"a":"123"},{"a":123}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@['a']['a','b']==123)]`},
			},
		},
		`value-group-jsonpath::wildcard-dot-child-identifier`: []TestCase{
			{
				jsonpath:    `$[?(@.*==2)]`,
				inputJSON:   `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.*==2)]`},
			},
			{
				jsonpath:    `$[?(@.*[0]==2)]`,
				inputJSON:   `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.*[0]==2)]`},
			},
			{
				jsonpath:    `$[?(@.*.a==2)]`,
				inputJSON:   `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.*.a==2)]`},
			},
			{
				jsonpath:    `$[?(@.a.*==2)]`,
				inputJSON:   `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a.*==2)]`},
			},
		},
		`match-array`: []TestCase{
			{
				jsonpath:     `$[?(@==$[1])]`,
				inputJSON:    `[[1],[2],[2],[3]]`,
				expectedJSON: `[[2],[2]]`,
			},
			{
				jsonpath:     `$[?(@==$[1])]`,
				inputJSON:    `[[1],[2,2],[2,2],[3]]`,
				expectedJSON: `[[2,2],[2,2]]`,
			},
			{
				jsonpath:     `$[?(@==$[1])]`,
				inputJSON:    `[[1],[2,{"2":2}],[2,{"2":2}],[3]]`,
				expectedJSON: `[[2,{"2":2}],[2,{"2":2}]]`,
			},
			{
				jsonpath:     `$.*[?(@==1)]`,
				inputJSON:    `[[1],{"b":2}]`,
				expectedJSON: `[1]`,
			},
		},
		`match-object`: []TestCase{
			{
				jsonpath:     `$[?(@==$[1])]`,
				inputJSON:    `[{"a":[1]},{"a":[2]},{"a":[2]},{"a":[3]}]`,
				expectedJSON: `[{"a":[2]},{"a":[2]}]`,
			},
			{
				jsonpath:     `$[?(@==$[1])]`,
				inputJSON:    `[{"a":[1]},{"a":[2,2]},{"a":[2,2]},{"a":[3]}]`,
				expectedJSON: `[{"a":[2,2]},{"a":[2,2]}]`,
			},
			{
				jsonpath:     `$[?(@==$[1])]`,
				inputJSON:    `[{"a":[1]},{"a":[2,{"2":2}]},{"a":[2,{"2":2}]},{"a":[3]}]`,
				expectedJSON: `[{"a":[2,{"2":2}]},{"a":[2,{"2":2}]}]`,
			},
			{
				jsonpath:     `$.*[?(@==1)]`,
				inputJSON:    `[{"a":1},{"b":2}]`,
				expectedJSON: `[1]`,
			},
		},
		`object-input`: []TestCase{
			{
				jsonpath:     `$[?(@)]`,
				inputJSON:    `{"a":1,"b":null}`,
				expectedJSON: `[1,null]`,
			},
			{
				jsonpath:     `$[?(@.a)]`,
				inputJSON:    `{"a":{"a":1},"b":{"b":2}}`,
				expectedJSON: `[{"a":1}]`,
			},
			{
				jsonpath:     `$[?(@[1]=="b")]`,
				inputJSON:    `{"a":["a","b"],"b":["b"]}`,
				expectedJSON: `[["a","b"]]`,
			},
		},
		`filter-chain`: []TestCase{
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
			{
				jsonpath:    `$[?(@[1][0]>1)][?(@[1][0]>1)][?(@[1]>1)]`,
				inputJSON:   `[1,[21,[221,[222]]]]`,
				expectedErr: createErrorMemberNotExist(`[?(@[1]>1)]`),
			},
		},
		`found-path-and-not-found-root-path`: []TestCase{
			{
				jsonpath:    `$[?(@.a == $.b)]`,
				inputJSON:   `[{"a":0},{"a":1}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.a == $.b)]`),
			},
			{
				jsonpath:     `$[?(@.a != $.b)]`,
				inputJSON:    `[{"a":0},{"a":1}]`,
				expectedJSON: `[{"a":0},{"a":1}]`,
			},
			{
				jsonpath:    `$[?(@.a < $.b)]`,
				inputJSON:   `[{"a":0},{"a":1}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.a < $.b)]`),
			},
			{
				jsonpath:    `$[?(@.a <= $.b)]`,
				inputJSON:   `[{"a":0},{"a":1}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.a <= $.b)]`),
			},
			{
				jsonpath:    `$[?(@.a > $.b)]`,
				inputJSON:   `[{"a":0},{"a":1}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.a > $.b)]`),
			},
			{
				jsonpath:    `$[?(@.a >= $.b)]`,
				inputJSON:   `[{"a":0},{"a":1}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.a >= $.b)]`),
			}, {
				jsonpath:    `$[?($.b == @.a)]`,
				inputJSON:   `[{"a":0},{"a":1}]`,
				expectedErr: createErrorMemberNotExist(`[?($.b == @.a)]`),
			},
			{
				jsonpath:     `$[?($.b != @.a)]`,
				inputJSON:    `[{"a":0},{"a":1}]`,
				expectedJSON: `[{"a":0},{"a":1}]`,
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
		},
		`not-found-path-and-found-root-path`: []TestCase{
			{
				jsonpath:    `$[?(@.b == $[0].a)]`,
				inputJSON:   `[{"a":0},{"a":1}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.b == $[0].a)]`),
			},
			{
				jsonpath:     `$[?(@.b != $[0].a)]`,
				inputJSON:    `[{"a":0},{"a":1}]`,
				expectedJSON: `[{"a":0},{"a":1}]`,
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
				jsonpath:    `$[?(@.b >= $[0].a)]`,
				inputJSON:   `[{"a":0},{"a":1}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.b >= $[0].a)]`),
			}, {
				jsonpath:    `$[?($[0].a == @.b)]`,
				inputJSON:   `[{"a":0},{"a":1}]`,
				expectedErr: createErrorMemberNotExist(`[?($[0].a == @.b)]`),
			},
			{
				jsonpath:     `$[?($[0].a != @.b)]`,
				inputJSON:    `[{"a":0},{"a":1}]`,
				expectedJSON: `[{"a":0},{"a":1}]`,
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
		},
		`not-found-path-and-not-found-root-path`: []TestCase{
			{
				jsonpath:     `$[?(@.b == $.b)]`,
				inputJSON:    `[{"a":0},{"a":1}]`,
				expectedJSON: `[{"a":0},{"a":1}]`,
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
				jsonpath:     `$[?($.b == @.b)]`,
				inputJSON:    `[{"a":0},{"a":1}]`,
				expectedJSON: `[{"a":0},{"a":1}]`,
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
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_filterSubFilter(t *testing.T) {
	testGroups := TestGroup{
		`allowed`: []TestCase{
			{
				jsonpath:     `$[?(@.a[?(@.b>1)])]`,
				inputJSON:    `[{"a":[{"b":1},{"b":2}]},{"a":[{"b":1}]}]`,
				expectedJSON: `[{"a":[{"b":1},{"b":2}]}]`,
			},
		},
		`prohibited`: []TestCase{
			{
				jsonpath:    `$[?(@.a[?(@.b)] > 1)]`,
				inputJSON:   `[{"a":[{"b":1},{"b":2}]},{"a":[{"b":1}]}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[?(@.b)] > 1)]`},
			},
			{
				jsonpath:    `$[?(@.a[?(@.b)] > 1)]`,
				inputJSON:   `[{"a":[{"b":2}]},{"a":[{"b":1}]}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[?(@.b)] > 1)]`},
			},
			{
				jsonpath:    `$[?(@.a[?(@.b)] > 1)]`,
				inputJSON:   `[{"a":[{"c":2}]},{"a":[{"d":1}]}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[?(@.b)] > 1)]`},
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_filterRegex(t *testing.T) {
	testGroups := TestGroup{
		`Go-regex-syntax`: []TestCase{
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
				jsonpath:    `$[?(@.a=~/(?x)CASE/)]`,
				inputJSON:   `[{"a":"case"},{"a":"CASE"},{"a":"Case"},{"a":"abc"}]`,
				expectedErr: ErrorInvalidArgument{argument: `(?x)CASE`, err: fmt.Errorf("error parsing regexp: invalid or unsupported Perl syntax: `(?x`")},
			},
			{
				jsonpath:    `$[?(@.a.b=~/abc/)]`,
				inputJSON:   `[{"a":"abc"}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.a.b=~/abc/)]`),
			},
		},
		`jsonpath::index-qualifier`: []TestCase{
			{
				jsonpath:     `$[?(@[0]=~/123/)]`,
				inputJSON:    `[["123"],["456"]]`,
				expectedJSON: `[["123"]]`,
			},
		},
		`value-group-jsonpath::slice-qualifier`: []TestCase{
			{
				jsonpath:    `$[?(@[0:1]=~/123/)]`,
				inputJSON:   `[{"b":["123"]},{"a":["123"]}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0:1]=~/123/)]`},
			},
			{
				jsonpath:    `$[?(@[0:2]=~/123/)]`,
				inputJSON:   `[{"b":["123"]},{"a":["123"]}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0:2]=~/123/)]`},
			},
			{
				jsonpath:    `$[?(@[0:2].a=~/123/)]`,
				inputJSON:   `[{"b":["123"]},{"a":["123"]}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0:2].a=~/123/)]`},
			},
			{
				jsonpath:    `$[?(@.a[0:2]=~/123/)]`,
				inputJSON:   `[{"b":["123"]},{"a":["123"]}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[0:2]=~/123/)]`},
			},
		},
		`value-group-jsonpath::wildcard-qualifier`: []TestCase{
			{
				jsonpath:    `$[?(@[*]=~/123/)]`,
				inputJSON:   `[{"b":"123"},{"a":"123"}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[*]=~/123/)]`},
			},
			{
				jsonpath:    `$[?(@[*].a=~/123/)]`,
				inputJSON:   `[{"b":"123"},{"a":"123"}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[*].a=~/123/)]`},
			},
			{
				jsonpath:    `$[?(@.a[*]=~/123/)]`,
				inputJSON:   `[{"b":"123"},{"a":"123"}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[*]=~/123/)]`},
			},
		},
		`value-group-jsonpath::union-qualifier`: []TestCase{
			{
				jsonpath:    `$[?(@[0,1]=~/123/)]`,
				inputJSON:   `[{"b":["123"]},{"a":[123,"123"]}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0,1]=~/123/)]`},
			},
			{
				jsonpath:    `$[?(@[0,1:2]=~/123/)]`,
				inputJSON:   `[{"b":["123"]},{"a":[123,"123"]}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0,1:2]=~/123/)]`},
			},
			{
				jsonpath:    `$[?(@[0,1].a=~/123/)]`,
				inputJSON:   `[{"b":["123"]},{"a":[123,"123"]}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@[0,1].a=~/123/)]`},
			},
			{
				jsonpath:    `$[?(@.a[0,1]=~/123/)]`,
				inputJSON:   `[{"b":["123"]},{"a":[123,"123"]}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[0,1]=~/123/)]`},
			},
		},
		`value-group-jsonpath::recursive`: []TestCase{
			{
				jsonpath:    `$[?($..a=~/123/)]`,
				inputJSON:   `[{"a":"123"},{"a":123}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `$..a=~/123/)]`},
			},
			{
				jsonpath:    `$[?($..a.b=~/123/)]`,
				inputJSON:   `[{"a":"123"},{"a":123}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `$..a.b=~/123/)]`},
			},
			{
				jsonpath:    `$[?($.a..b=~/123/)]`,
				inputJSON:   `[{"a":"123"},{"a":123}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `$.a..b=~/123/)]`},
			},
			{
				jsonpath:    `$[?($..a..b=~/123/)]`,
				inputJSON:   `[{"a":"123"},{"a":123}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `$..a..b=~/123/)]`},
			},
		},
		`value-group-jsonpath::multi-identifier`: []TestCase{
			{
				jsonpath:    `$[?(@['a','b']=~/123/)]`,
				inputJSON:   `[{"b":"123"},{"a":"123"}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@['a','b']=~/123/)]`},
			},
			{
				jsonpath:    `$[?(@['a','b','c']=~/123/)]`,
				inputJSON:   `[{"b":"123"},{"a":"123"}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@['a','b','c']=~/123/)]`},
			},
			{
				jsonpath:    `$[?(@['a','b']['a']=~/123/)]`,
				inputJSON:   `[{"b":"123"},{"a":"123"}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@['a','b']['a']=~/123/)]`},
			},
			{
				jsonpath:    `$[?(@['a']['a','b']=~/123/)]`,
				inputJSON:   `[{"b":"123"},{"a":"123"}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@['a']['a','b']=~/123/)]`},
			},
		},
		`value-group-jsonpath::wildcard-dot-child-identifier`: []TestCase{
			{
				jsonpath:    `$[?(@.*=~/123/)]`,
				inputJSON:   `[{"b":"123"},{"a":"123"}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.*=~/123/)]`},
			},
			{
				jsonpath:    `$[?(@.*[0]=~/123/)]`,
				inputJSON:   `[{"b":"123"},{"a":"123"}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.*[0]=~/123/)]`},
			},
			{
				jsonpath:    `$[?(@.*.a=~/123/)]`,
				inputJSON:   `[{"b":"123"},{"a":"123"}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.*.a=~/123/)]`},
			},
			{
				jsonpath:    `$[?(@.a.*=~/123/)]`,
				inputJSON:   `[{"b":"123"},{"a":"123"}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a.*=~/123/)]`},
			},
		},
		`value-group-jsonpath::sub-filter`: []TestCase{
			{
				jsonpath:    `$[?(@.a[?(@.b)]=~/123/)]`,
				inputJSON:   `[{"b":"123"},{"a":"123"}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[?(@.b)]=~/123/)]`},
			},
			{
				jsonpath:    `$[?(@.a[?(@.b>1)]=~/123/)]`,
				inputJSON:   `[{"a":[{"b":1},{"b":2}]},{"a":[{"b":1}]}]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `JSONPath that returns a value group is prohibited`, near: `@.a[?(@.b>1)]=~/123/)]`},
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_filterLogicalCombination(t *testing.T) {
	testGroups := TestGroup{
		`logical OR`: []TestCase{
			{
				jsonpath:     `$[?(@.a || @.b)]`,
				inputJSON:    `[{"a":1},{"b":2},{"c":3}]`,
				expectedJSON: `[{"a":1},{"b":2}]`,
			},
			{
				jsonpath:     `$[?(@.a>2 || @.a<2)]`,
				inputJSON:    `[{"a":1},{"a":1.9},{"a":2},{"a":2.1},{"a":3}]`,
				expectedJSON: `[{"a":1},{"a":1.9},{"a":2.1},{"a":3}]`,
			},
			{
				jsonpath:     `$[?(@.a<2 || @.a>2)]`,
				inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
				expectedJSON: `[{"a":1},{"a":3}]`,
			},
			{
				jsonpath:     `$[?((1==2) || @.a>1)]`,
				inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
				expectedJSON: `[{"a":2},{"a":3}]`,
			},
			{
				jsonpath:     `$[?((1==1) || @.a>1)]`,
				inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
				expectedJSON: `[{"a":1},{"a":2},{"a":3}]`,
			},
			{
				jsonpath:     `$[?(@.a>1 || (1==2))]`,
				inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
				expectedJSON: `[{"a":2},{"a":3}]`,
			},
			{
				jsonpath:     `$[?(@.a>1 || (1==1))]`,
				inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
				expectedJSON: `[{"a":1},{"a":2},{"a":3}]`,
			},
			{
				jsonpath:     `$[?(@.x || @.b > 2)]`,
				inputJSON:    `[{"a":"a"},{"b":2},{"b":3}]`,
				expectedJSON: `[{"b":3}]`,
			},
			{
				jsonpath:     `$[?(@.b > 2 || @.x)]`,
				inputJSON:    `[{"a":"a"},{"b":2},{"b":3}]`,
				expectedJSON: `[{"b":3}]`,
			},
			{
				jsonpath:    `$[?(@.x || @.x)]`,
				inputJSON:   `[{"a":"a"},{"b":2},{"b":3}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.x || @.x)]`),
			},
			{
				jsonpath:     `$[?(@.b > 2 || @.b < 2)]`,
				inputJSON:    `[{"b":1},{"b":2},{"b":3}]`,
				expectedJSON: `[{"b":1},{"b":3}]`,
			},
			{
				jsonpath:     `$.z[?($..x || @.b < 2)]`,
				inputJSON:    `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
				expectedJSON: `[{"b":1},{"b":2},{"b":3}]`,
			},
			{
				jsonpath:     `$.z[?($..xx || @.b < 2)]`,
				inputJSON:    `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
				expectedJSON: `[{"b":1}]`,
			},
		},
		`logical AND`: []TestCase{
			{
				jsonpath:     `$[?(@.a && @.b)]`,
				inputJSON:    `[{"a":1},{"b":2},{"a":3,"b":4}]`,
				expectedJSON: `[{"a":3,"b":4}]`,
			},
			{
				jsonpath:     `$[?(@.a>1 && @.a<3)]`,
				inputJSON:    `[{"a":1},{"a":1.1},{"a":2.9},{"a":3}]`,
				expectedJSON: `[{"a":1.1},{"a":2.9}]`,
			},
			{
				jsonpath:     `$[?(@.a<3 && @.a>1)]`,
				inputJSON:    `[{"a":1},{"a":1.1},{"a":2.9},{"a":3}]`,
				expectedJSON: `[{"a":1.1},{"a":2.9}]`,
			},
			{
				jsonpath:     `$[?((1==2) && @.a>1)]`,
				inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
				expectedJSON: `[]`,
				expectedErr:  createErrorMemberNotExist(`[?((1==2) && @.a>1)]`),
			},
			{
				jsonpath:     `$[?((1==1) && @.a>1)]`,
				inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
				expectedJSON: `[{"a":2},{"a":3}]`,
			},
			{
				jsonpath:     `$[?(@.a>1 && (1==2))]`,
				inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
				expectedJSON: `[]`,
				expectedErr:  createErrorMemberNotExist(`[?(@.a>1 && (1==2))]`),
			},
			{
				jsonpath:     `$[?(@.a>1 && (1==1))]`,
				inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
				expectedJSON: `[{"a":2},{"a":3}]`,
			},
			{
				jsonpath:    `$[?(@.x && @.b > 2)]`,
				inputJSON:   `[{"a":"a"},{"b":2},{"b":3}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.x && @.b > 2)]`),
			},
			{
				jsonpath:    `$[?(@.b > 2 && @.x)]`,
				inputJSON:   `[{"a":"a"},{"b":2},{"b":3}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.b > 2 && @.x)]`),
			},
			{
				jsonpath:    `$[?(@.x && @.x)]`,
				inputJSON:   `[{"a":"a"},{"b":2},{"b":3}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.x && @.x)]`),
			},
			{
				jsonpath:    `$[?(@.b > 2 && @.b < 2)]`,
				inputJSON:   `[{"b":1},{"b":2},{"b":3}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.b > 2 && @.b < 2)]`),
			},
			{
				jsonpath:     `$.z[?($..x && @.b < 2)]`,
				inputJSON:    `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
				expectedJSON: `[{"b":1}]`,
			},
			{
				jsonpath:    `$.z[?($..xx && @.b < 2)]`,
				inputJSON:   `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
				expectedErr: createErrorMemberNotExist(`[?($..xx && @.b < 2)]`),
			},
		},
		`logical NOT`: []TestCase{
			{
				jsonpath:     `$[?(!@.a)]`,
				inputJSON:    `[{"a":1},{"b":2},{"a":3,"b":4}]`,
				expectedJSON: `[{"b":2}]`,
			},
			{
				jsonpath:     `$[?(!@.c)]`,
				inputJSON:    `[{"a":1},{"b":2},{"a":3,"b":4}]`,
				expectedJSON: `[{"a":1},{"b":2},{"a":3,"b":4}]`,
			},
			{
				jsonpath:    `$.z[?(!$..x)]`,
				inputJSON:   `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
				expectedErr: createErrorMemberNotExist(`[?(!$..x)]`),
			},
			{
				jsonpath:     `$.z[?(!$..xx)]`,
				inputJSON:    `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
				expectedJSON: `[{"b":1},{"b":2},{"b":3}]`,
			},
		},
		`priority`: []TestCase{
			{
				jsonpath:     `$[?(@.a && @.b || @.c)]`,
				inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
				expectedJSON: `[{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6}]`,
			},
			{
				jsonpath:     `$[?(@.a && (@.b || @.c))]`,
				inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
				expectedJSON: `[{"a":2,"b":2},{"a":3,"b":3,"c":3},{"a":5,"c":5}]`,
			},
			{
				jsonpath:     `$[?((@.a && @.b) || @.c)]`,
				inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
				expectedJSON: `[{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6}]`,
			},
			{
				jsonpath:     `$[?(@.a || @.b && @.c)]`,
				inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
				expectedJSON: `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5}]`,
			},
		},
		`not-combination`: []TestCase{
			{
				jsonpath:     `$[?(!@.a && @.b || @.c)]`,
				inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
				expectedJSON: `[{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
			},
			{
				jsonpath:     `$[?(@.a && !@.b || @.c)]`,
				inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
				expectedJSON: `[{"a":1},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6}]`,
			},
			{
				jsonpath:     `$[?(!@.a && !@.b || @.c)]`,
				inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
				expectedJSON: `[{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6}]`,
			},
			{
				jsonpath:     `$[?(@.a && @.b || !@.c)]`,
				inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
				expectedJSON: `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":7}]`,
			},
			{
				jsonpath:     `$[?(!@.a && @.b || !@.c)]`,
				inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
				expectedJSON: `[{"a":1},{"a":2,"b":2},{"b":4,"c":4},{"b":7}]`,
			},
			{
				jsonpath:     `$[?(@.a && !@.b || !@.c)]`,
				inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
				expectedJSON: `[{"a":1},{"a":2,"b":2},{"a":5,"c":5},{"b":7}]`,
			},
			{
				jsonpath:     `$[?(!@.a && !@.b || !@.c)]`,
				inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
				expectedJSON: `[{"a":1},{"a":2,"b":2},{"c":6},{"b":7}]`,
			},
		},
		`comparator`: []TestCase{
			{
				jsonpath:     `$[?(@.a || @.b > 2)]`,
				inputJSON:    `[{"a":"a"},{"b":2},{"b":3}]`,
				expectedJSON: `[{"a":"a"},{"b":3}]`,
			},
			{
				jsonpath:     `$[?(@.b > 2 || @.a)]`,
				inputJSON:    `[{"a":"a"},{"b":2},{"b":3}]`,
				expectedJSON: `[{"a":"a"},{"b":3}]`,
			},
			{
				jsonpath:     `$[?(@.a =~ /a/ && @.b == 2)]`,
				inputJSON:    `[{"a":"a"},{"a":"a","b":2}]`,
				expectedJSON: `[{"a":"a","b":2}]`,
			},
			{
				jsonpath:     `$[?(@.b == 2 && @.a =~ /a/)]`,
				inputJSON:    `[{"a":"a"},{"a":"a","b":2}]`,
				expectedJSON: `[{"a":"a","b":2}]`,
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_valueGroupCombination_Recursive_descent(t *testing.T) {
	testGroups := TestGroup{
		`Recursive-descent`: []TestCase{
			{
				jsonpath:     `$..a..b`,
				inputJSON:    `[{"a":{"a":{"b":1},"c":2}},{"b":{"a":{"d":3,"b":4}}}]`,
				expectedJSON: `[1,1,4]`,
			},
			{
				jsonpath:    `$..a..b`,
				inputJSON:   `[{"a":{"a":{"x":1},"c":2}},{"b":{"a":{"d":3,"x":4}}}]`,
				expectedErr: createErrorMemberNotExist(`b`),
			},
			{
				jsonpath:    `$..a..b`,
				inputJSON:   `[{"a":"b"},{"b":{"a":"b"}}]`,
				expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
			},
			{
				jsonpath:    `$..a..b`,
				inputJSON:   `[{"x":{"x":{"b":1},"c":2}},{"b":{"x":{"d":3,"b":4}}}]`,
				expectedErr: createErrorMemberNotExist(`a`),
			},
			{
				jsonpath:    `$..a..b`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
			},
		},
		`Multiple-identifier`: []TestCase{
			{
				jsonpath:     `$..['a','b']`,
				inputJSON:    `[{"a":1,"c":2},{"d":3,"b":4}]`,
				expectedJSON: `[1,4]`,
			},
			{
				jsonpath:    `$..['a','b']`,
				inputJSON:   `[{"x":1,"c":2},{"d":3,"x":4}]`,
				expectedErr: createErrorMemberNotExist(`['a','b']`),
			},
			{
				jsonpath:    `$..['a','b']`,
				inputJSON:   `{}`,
				expectedErr: createErrorMemberNotExist(`['a','b']`),
			},
			{
				jsonpath:    `$..['a','b']`,
				inputJSON:   `[]`,
				expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `[]interface {}`),
			},
			{
				jsonpath:    `$..['a','b']`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
			},
		},
		`Wildcard-identifier`: []TestCase{
			{
				jsonpath:     `$..*`,
				inputJSON:    `[{"a":1,"c":2},{"d":3,"b":4}]`,
				expectedJSON: `[{"a":1,"c":2},{"b":4,"d":3},1,2,4,3]`,
			},
			{
				jsonpath:    `$..*`,
				inputJSON:   `{}`,
				expectedErr: createErrorMemberNotExist(`*`),
			},
			{
				jsonpath:    `$..*`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
			},
		},
		`Slice-qualifier`: []TestCase{
			{
				jsonpath:     `$..[0:2]`,
				inputJSON:    `{"a":[1,3,2],"b":{"a":[4,6,5]}}`,
				expectedJSON: `[1,3,4,6]`,
			},
			{
				jsonpath:    `$..[0:2]`,
				inputJSON:   `{"a":[],"b":{"a":[]}}`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
			{
				jsonpath:    `$..[0:2]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
			},
		},
		`Wildcard-qualifier`: []TestCase{
			{
				jsonpath:     `$..[*]`,
				inputJSON:    `[[1,3,2],[4,6,5]]`,
				expectedJSON: `[[1,3,2],[4,6,5],1,3,2,4,6,5]`,
			},
			{
				jsonpath:    `$..[*]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$..a[*]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
			},
		},
		`Union-in-qualifier`: []TestCase{
			{
				jsonpath:     `$..[0,1]`,
				inputJSON:    `[[1,3,2],[4,6,5]]`,
				expectedJSON: `[[1,3,2],[4,6,5],1,3,4,6]`,
			},
			{
				jsonpath:    `$..[0,1]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$..[0,1]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
			},
		},
		`Filter-qualifier`: []TestCase{
			{
				jsonpath:     `$..[?(@.b)]`,
				inputJSON:    `[[{"a":1},{"b":2}],[{"a":3},{"b":4}]]`,
				expectedJSON: `[{"b":2},{"b":4}]`,
			},
			{
				jsonpath:    `$..[?(@.b)]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
			},
			{
				jsonpath:    `$..[?(@.b)]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_valueGroupCombination_Multiple_identifier(t *testing.T) {
	testGroups := TestGroup{
		`Recursive-descent`: []TestCase{
			{
				jsonpath:     `$['a','b']..a`,
				inputJSON:    `{"a":{"a":1,"c":2},"b":{"a":{"d":3,"a":4}}}`,
				expectedJSON: `[1,{"a":4,"d":3},4]`,
			},
			{
				jsonpath:    `$['a','b']..a`,
				inputJSON:   `{"a":{"x":1,"c":2},"b":{"x":{"d":3,"x":4}}}`,
				expectedErr: createErrorMemberNotExist(`a`),
			},
			{
				jsonpath:    `$['a','b']..a`,
				inputJSON:   `{"a":"a","b":"a"}`,
				expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
			},
			{
				jsonpath:    `$['a','b']..a`,
				inputJSON:   `{"x":{"x":1,"c":2},"y":{"x":{"d":3,"x":4}}}`,
				expectedErr: createErrorMemberNotExist(`['a','b']`),
			},
			{
				jsonpath:    `$['a','b']..a`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `string`),
			},
		},
		`Multiple-identifier`: []TestCase{
			{
				jsonpath:     `$['a','b']['c','d']`,
				inputJSON:    `{"a":{"a":1,"c":2},"b":{"d":3,"a":4}}`,
				expectedJSON: `[2,3]`,
			},
			{
				jsonpath:    `$['a','b']['c','d']`,
				inputJSON:   `{"a":{"a":1,"x":2},"b":{"x":3,"a":4}}`,
				expectedErr: createErrorMemberNotExist(`['c','d']`),
			},
			{
				jsonpath:    `$['a','b']['c','d']`,
				inputJSON:   `{"x":{"a":1,"c":2},"x":{"d":3,"a":4}}`,
				expectedErr: createErrorMemberNotExist(`['a','b']`),
			},
			{
				jsonpath:    `$['a','b']['c','d']`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `string`),
			},
		},
		`Wildcard-identifier`: []TestCase{
			{
				jsonpath:     `$['a','b'].*`,
				inputJSON:    `{"a":{"a":1,"c":2},"b":{"d":3,"a":4}}`,
				expectedJSON: `[1,2,4,3]`,
			},
			{
				jsonpath:    `$['a','b'].*`,
				inputJSON:   `{"a":{},"b":{}}`,
				expectedErr: createErrorMemberNotExist(`.*`),
			},
			{
				jsonpath:    `$['a','b'].*`,
				inputJSON:   `{"x":[1,3,2],"y":[4,6,5]}`,
				expectedErr: createErrorMemberNotExist(`['a','b']`),
			},
			{
				jsonpath:    `$['a','b'].*`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `string`),
			},
		},
		`Slice-qualifier`: []TestCase{
			{
				jsonpath:     `$['a','b'][0:2]`,
				inputJSON:    `{"a":[1,3,2],"b":[4,6,5]}`,
				expectedJSON: `[1,3,4,6]`,
			},
			{
				jsonpath:    `$['a','b'][0:2]`,
				inputJSON:   `{"a":[],"b":[]}`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
			{
				jsonpath:    `$['a','b'][0:2]`,
				inputJSON:   `{"x":[1,3,2],"y":[4,6,5]}`,
				expectedErr: createErrorMemberNotExist(`['a','b']`),
			},
			{
				jsonpath:    `$['a','b'][0:2]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `string`),
			},
		},
		`Wildcard-qualifier`: []TestCase{
			{
				jsonpath:     `$['a','b'][*]`,
				inputJSON:    `{"a":[1,3,2],"b":[4,6,5]}`,
				expectedJSON: `[1,3,2,4,6,5]`,
			},
			{
				jsonpath:    `$['a','b'][*]`,
				inputJSON:   `{"a":[],"b":[]}`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$['a','b'][*]`,
				inputJSON:   `{"x":[1,3,2],"y":[4,6,5]}`,
				expectedErr: createErrorMemberNotExist(`['a','b']`),
			},
			{
				jsonpath:    `$['a','b'][*]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `string`),
			},
		},
		`Union-in-qualifier`: []TestCase{
			{
				jsonpath:     `$['a','b'][0,1]`,
				inputJSON:    `{"a":[1,3,2],"b":[4,6,5]}`,
				expectedJSON: `[1,3,4,6]`,
			},
			{
				jsonpath:    `$['a','b'][0,1]`,
				inputJSON:   `{"a":[],"b":[]}`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$['a','b'][0,1]`,
				inputJSON:   `{"x":[1,3,2],"y":[4,6,5]}`,
				expectedErr: createErrorMemberNotExist(`['a','b']`),
			},
			{
				jsonpath:    `$['a','b'][0,1]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `string`),
			},
		},
		`Filter-qualifier`: []TestCase{
			{
				jsonpath:     `$['a','b'][?(@.b)]`,
				inputJSON:    `{"a":[{"a":1},{"b":2}],"b":[{"a":3},{"b":4}]}`,
				expectedJSON: `[{"b":2},{"b":4}]`,
			},
			{
				jsonpath:    `$['a','b'][?(@.b)]`,
				inputJSON:   `{"a":[{"a":1},{"x":2}],"b":[{"a":3},{"x":4}]}`,
				expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
			},
			{
				jsonpath:    `$['a','b'][?(@.b)]`,
				inputJSON:   `{"x":[{"a":1},{"b":2}],"y":[{"a":3},{"b":4}]}`,
				expectedErr: createErrorMemberNotExist(`['a','b']`),
			},
			{
				jsonpath:    `$['a','b'][?(@.b)]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`['a','b']`, `object`, `string`),
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_valueGroupCombination_Wildcard_identifier(t *testing.T) {
	testGroups := TestGroup{
		`Recursive-descent`: []TestCase{
			{
				jsonpath:     `$.*..a`,
				inputJSON:    `{"a":{"a":1,"c":2},"b":{"d":{"e":3,"a":4}}}`,
				expectedJSON: `[1,4]`,
			},
			{
				jsonpath:    `$.*..a`,
				inputJSON:   `{"x":{"x":1,"c":2},"b":{"d":{"e":3,"x":4}}}`,
				expectedErr: createErrorMemberNotExist(`a`),
			},
			{
				jsonpath:    `$.*..a`,
				inputJSON:   `{"a":"a","b":"b"}`,
				expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
			},
			{
				jsonpath:    `$.*..a`,
				inputJSON:   `{}`,
				expectedErr: createErrorMemberNotExist(`.*`),
			},
			{
				jsonpath:    `$.*..a`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`.*`, `object/array`, `string`),
			},
		},
		`Multiple-identifier`: []TestCase{
			{
				jsonpath:     `$.*['a','b']`,
				inputJSON:    `{"a":{"a":1},"c":{"c":3},"b":{"b":2}}`,
				expectedJSON: `[1,2]`,
			},
			{
				jsonpath:    `$.*['a','b']`,
				inputJSON:   `{"a":{"x":1},"c":{"c":3},"b":{"x":2}}`,
				expectedErr: createErrorMemberNotExist(`['a','b']`),
			},
			{
				jsonpath:    `$.*['a','b']`,
				inputJSON:   `{}`,
				expectedErr: createErrorMemberNotExist(`.*`),
			},
			{
				jsonpath:    `$.*['a','b']`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`.*`, `object/array`, `string`),
			},
		},
		`Wildcard-identifier`: []TestCase{
			{
				jsonpath:     `$.*.*`,
				inputJSON:    `{"a":{"a":1,"c":2},"b":{"d":3,"a":4}}`,
				expectedJSON: `[1,2,4,3]`,
			},
			{
				jsonpath:    `$.*.*`,
				inputJSON:   `{"a":{},"b":{}}`,
				expectedErr: createErrorMemberNotExist(`.*`),
			},
			{
				jsonpath:    `$.*.*`,
				inputJSON:   `{}`,
				expectedErr: createErrorMemberNotExist(`.*`),
			},
			{
				jsonpath:    `$.*.*`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`.*`, `object/array`, `string`),
			},
		},
		`Slice-qualifier`: []TestCase{
			{
				jsonpath:     `$.*[0:2]`,
				inputJSON:    `{"a":[1,3,2],"b":[4,6,5]}`,
				expectedJSON: `[1,3,4,6]`,
			},
			{
				jsonpath:    `$.*[0:2]`,
				inputJSON:   `{"a":[],"b":[]}`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
			{
				jsonpath:    `$.*[0:2]`,
				inputJSON:   `{}`,
				expectedErr: createErrorMemberNotExist(`.*`),
			},
			{
				jsonpath:    `$.*[0:2]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`.*`, `object/array`, `string`),
			},
		},
		`Wildcard-qualifier`: []TestCase{
			{
				jsonpath:     `$.*[*]`,
				inputJSON:    `{"a":[1,3,2],"b":[4,6,5]}`,
				expectedJSON: `[1,3,2,4,6,5]`,
			},
			{
				jsonpath:    `$.*[*]`,
				inputJSON:   `{"a":[],"b":[]}`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$.*[*]`,
				inputJSON:   `{}`,
				expectedErr: createErrorMemberNotExist(`.*`),
			},
			{
				jsonpath:    `$.*[*]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`.*`, `object/array`, `string`),
			},
		},
		`Union-in-qualifier`: []TestCase{
			{
				jsonpath:     `$.*[0,1]`,
				inputJSON:    `{"a":[1,3,2],"b":[4,6,5]}`,
				expectedJSON: `[1,3,4,6]`,
			},
			{
				jsonpath:    `$.*[0,1]`,
				inputJSON:   `{"a":[],"b":[]}`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$.*[0,1]`,
				inputJSON:   `{}`,
				expectedErr: createErrorMemberNotExist(`.*`),
			},
			{
				jsonpath:    `$.*[0,1]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`.*`, `object/array`, `string`),
			},
		},
		`Filter-qualifier`: []TestCase{
			{
				jsonpath:     `$.*[?(@.b)]`,
				inputJSON:    `{"a":[{"a":1},{"b":2}],"b":[{"a":3},{"b":4}]}`,
				expectedJSON: `[{"b":2},{"b":4}]`,
			},
			{
				jsonpath:    `$.*[?(@.b)]`,
				inputJSON:   `{"a":[{"a":1},{"x":2}],"b":[{"a":3},{"x":4}]}`,
				expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
			},
			{
				jsonpath:    `$.*[?(@.b)]`,
				inputJSON:   `{}`,
				expectedErr: createErrorMemberNotExist(`.*`),
			},
			{
				jsonpath:    `$.*[?(@.b)]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`.*`, `object/array`, `string`),
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_valueGroupCombination_Slice_qualifier(t *testing.T) {
	testGroups := TestGroup{
		`Recursive-descent`: []TestCase{
			{
				jsonpath:     `$[0:2]..a`,
				inputJSON:    `[{"a":1},{"b":{"a":2}},{"a":3}]`,
				expectedJSON: `[1,2]`,
			},
			{
				jsonpath:    `$[0:2]..a`,
				inputJSON:   `[{"x":1},{"b":{"x":2}},{"a":3}]`,
				expectedErr: createErrorMemberNotExist(`a`),
			},
			{
				jsonpath:    `$[0:2]..a`,
				inputJSON:   `["a","b",{"a":3}]`,
				expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
			},
			{
				jsonpath:    `$[0:2]..a`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
			{
				jsonpath:    `$[0:2]..a`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[0:2]`, `array`, `string`),
			},
		},
		`Multiple-identifier`: []TestCase{
			{
				jsonpath:     `$[0:2]['a','b']`,
				inputJSON:    `[{"a":1},{"b":2}]`,
				expectedJSON: `[1,2]`,
			},
			{
				jsonpath:    `$[0:2]['a','b']`,
				inputJSON:   `[{"x":1},{"x":2}]`,
				expectedErr: createErrorMemberNotExist(`['a','b']`),
			},
			{
				jsonpath:    `$[0:2]['a','b']`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
			{
				jsonpath:    `$[0:2]['a','b']`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[0:2]`, `array`, `string`),
			},
		},
		`Wildcard-identifier`: []TestCase{
			{
				jsonpath:     `$[0:2].*`,
				inputJSON:    `[{"a":1,"c":2},{"d":3,"b":4}]`,
				expectedJSON: `[1,2,4,3]`,
			},
			{
				jsonpath:    `$[0:2].*`,
				inputJSON:   `[[],[]]`,
				expectedErr: createErrorMemberNotExist(`.*`),
			},
			{
				jsonpath:    `$[0:2].*`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
			{
				jsonpath:    `$[0:2].*`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[0:2]`, `array`, `string`),
			},
		},
		`Slice-qualifier`: []TestCase{
			{
				jsonpath:     `$[0:2][0:2]`,
				inputJSON:    `[[1,2,3],[4,5,6]]`,
				expectedJSON: `[1,2,4,5]`,
			},
			{
				jsonpath:    `$[0:2][0:2]`,
				inputJSON:   `[[],[]]`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
			{
				jsonpath:    `$[0:2][0:2]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
			{
				jsonpath:    `$[0:2][0:2]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[0:2]`, `array`, `string`),
			},
		},
		`Wildcard-qualifier`: []TestCase{
			{
				jsonpath:     `$[0:2][*]`,
				inputJSON:    `[{"a":1,"c":3},{"d":4,"b":2},{"e":5}]`,
				expectedJSON: `[1,3,2,4]`,
			},
			{
				jsonpath:    `$[0:2][*]`,
				inputJSON:   `[{},{}]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[0:2][*]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
			{
				jsonpath:    `$[0:2][*]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[0:2]`, `array`, `string`),
			},
		},
		`Union-in-qualifier`: []TestCase{
			{
				jsonpath:     `$[0:2][0,1]`,
				inputJSON:    `[[1,3,2],[4,6,5],[7]]`,
				expectedJSON: `[1,3,4,6]`,
			},
			{
				jsonpath:    `$[0:2][0,1]`,
				inputJSON:   `[[],[],[7]]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$[0:2][0,1]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
			{
				jsonpath:    `$[0:2][0,1]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[0:2]`, `array`, `string`),
			},
		},
		`Filter-qualifier`: []TestCase{
			{
				jsonpath:     `$[0:2][?(@.b)]`,
				inputJSON:    `[[{"a":1},{"b":2}],[{"a":3},{"b":4}]]`,
				expectedJSON: `[{"b":2},{"b":4}]`,
			},
			{
				jsonpath:    `$[0:2][?(@.b)]`,
				inputJSON:   `[[{"a":1},{"x":2}],[{"a":3},{"x":4}]]`,
				expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
			},
			{
				jsonpath:    `$[0:2][?(@.b)]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
			{
				jsonpath:    `$[0:2][?(@.b)]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[0:2]`, `array`, `string`),
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_valueGroupCombination_Wildcard_qualifier(t *testing.T) {
	testGroups := TestGroup{
		`Recursive-descent`: []TestCase{
			{
				jsonpath:     `$[*]..a`,
				inputJSON:    `[{"a":1},{"b":{"a":2}},{"c":3}]`,
				expectedJSON: `[1,2]`,
			},
			{
				jsonpath:    `$[*]..a`,
				inputJSON:   `[{"x":1},{"b":{"x":2}},{"c":3}]`,
				expectedErr: createErrorMemberNotExist(`a`),
			},
			{
				jsonpath:    `$[*]..a`,
				inputJSON:   `["a","b","c"]`,
				expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
			},
			{
				jsonpath:    `$[*]..a`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[*]..a`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[*]`, `object/array`, `string`),
			},
		},
		`Multiple-identifier`: []TestCase{
			{
				jsonpath:     `$[*]['a','b']`,
				inputJSON:    `[{"c":4},{"b":2,"a":1},{"a":3}]`,
				expectedJSON: `[1,2,3]`,
			},
			{
				jsonpath:    `$[*]['a','b']`,
				inputJSON:   `[{"c":4},{"x":2},{"x":1}]`,
				expectedErr: createErrorMemberNotExist(`['a','b']`),
			},
			{
				jsonpath:    `$[*]['a','b']`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[*]['a','b']`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[*]`, `object/array`, `string`),
			},
		},
		`Wildcard-identifier`: []TestCase{
			{
				jsonpath:     `$[*].*`,
				inputJSON:    `[{"c":4},{"b":2,"a":1},{"a":3}]`,
				expectedJSON: `[4,1,2,3]`,
			},
			{
				jsonpath:    `$[*].*`,
				inputJSON:   `[{},{},{}]`,
				expectedErr: createErrorMemberNotExist(`.*`),
			},
			{
				jsonpath:    `$[*].*`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[*].*`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[*]`, `object/array`, `string`),
			},
		},
		`Slice-qualifier`: []TestCase{
			{
				jsonpath:     `$[*][0:2]`,
				inputJSON:    `[[1,2,3],[4,5],[6]]`,
				expectedJSON: `[1,2,4,5,6]`,
			},
			{
				jsonpath:    `$[*][0:2]`,
				inputJSON:   `[[],[],[]]`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
			{
				jsonpath:    `$[*][0:2]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[*][0:2]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[*]`, `object/array`, `string`),
			},
		},
		`Wildcard-qualifier`: []TestCase{
			{
				jsonpath:     `$[*][*]`,
				inputJSON:    `[[1,2,3],[4,5],[6]]`,
				expectedJSON: `[1,2,3,4,5,6]`,
			},
			{
				jsonpath:    `$[*][*]`,
				inputJSON:   `[[],[],[]]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[*][*]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[*][*]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[*]`, `object/array`, `string`),
			},
		},
		`Union-in-qualifier`: []TestCase{
			{
				jsonpath:     `$[*][0,1]`,
				inputJSON:    `[[1,3,2],[4,6,5],[7]]`,
				expectedJSON: `[1,3,4,6,7]`,
			},
			{
				jsonpath:    `$[*][0,1]`,
				inputJSON:   `[[],[],[]]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$[*][0,1]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[*][0,1]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[*]`, `object/array`, `string`),
			},
		},
		`Filter-qualifier`: []TestCase{
			{
				jsonpath:     `$[*][?(@.b)]`,
				inputJSON:    `[[{"a":1},{"b":2}],[{"a":3},{"b":4}]]`,
				expectedJSON: `[{"b":2},{"b":4}]`,
			},
			{
				jsonpath:    `$[*][?(@.b)]`,
				inputJSON:   `[[{"a":1},{"x":2}],[{"a":3},{"x":4}]]`,
				expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
			},
			{
				jsonpath:    `$[*][?(@.b)]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[*][?(@.b)]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[*]`, `object/array`, `string`),
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_valueGroupCombination_Union_in_qualifier(t *testing.T) {
	testGroups := TestGroup{
		`Recursive-descent`: []TestCase{
			{
				jsonpath:     `$[0,1]..a`,
				inputJSON:    `[{"a":1},{"b":{"a":2}},{"a":3}]`,
				expectedJSON: `[1,2]`,
			},
			{
				jsonpath:    `$[0,1]..a`,
				inputJSON:   `[{"x":1},{"b":{"x":2}},{"a":3}]`,
				expectedErr: createErrorMemberNotExist(`a`),
			},
			{
				jsonpath:    `$[0,1]..a`,
				inputJSON:   `["a","b",{"a":3}]`,
				expectedErr: createErrorTypeUnmatched(`..`, `object/array`, `string`),
			},
			{
				jsonpath:    `$[0,1]..a`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$[0,1]..a`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
			},
		},
		`Multiple-identifier`: []TestCase{
			{
				jsonpath:     `$[0,1]['a','b']`,
				inputJSON:    `[{"a":1},{"b":2}]`,
				expectedJSON: `[1,2]`,
			},
			{
				jsonpath:    `$[0,1]['a','b']`,
				inputJSON:   `[{"x":1},{"x":2}]`,
				expectedErr: createErrorMemberNotExist(`['a','b']`),
			},
			{
				jsonpath:    `$[0,1]['a','b']`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$[0,1]['a','b']`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
			},
		},
		`Wildcard-identifier`: []TestCase{
			{
				jsonpath:     `$[0,1].*`,
				inputJSON:    `[{"a":1,"c":2},{"d":3,"b":4}]`,
				expectedJSON: `[1,2,4,3]`,
			},
			{
				jsonpath:    `$[0,1].*`,
				inputJSON:   `[[],[]]`,
				expectedErr: createErrorMemberNotExist(`.*`),
			},
			{
				jsonpath:    `$[0,1].*`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$[0,1].*`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
			},
		},
		`Slice-qualifier`: []TestCase{
			{
				jsonpath:     `$[0,1][0:2]`,
				inputJSON:    `[[1,2,3],[4,5,6]]`,
				expectedJSON: `[1,2,4,5]`,
			},
			{
				jsonpath:    `$[0,1][0:2]`,
				inputJSON:   `[[],[]]`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
			{
				jsonpath:    `$[0,1][0:2]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$[0,1][0:2]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
			},
		},
		`Wildcard-qualifier`: []TestCase{
			{
				jsonpath:     `$[0,1][*]`,
				inputJSON:    `[{"a":1,"c":3},{"d":4,"b":2},{"e":5}]`,
				expectedJSON: `[1,3,2,4]`,
			},
			{
				jsonpath:    `$[0,1][*]`,
				inputJSON:   `[{},{}]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[0,1][*]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$[0,1][*]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
			},
		},
		`Union-in-qualifier`: []TestCase{
			{
				jsonpath:     `$[0,1][0,1]`,
				inputJSON:    `[[1,3,2],[4,6,5],[7]]`,
				expectedJSON: `[1,3,4,6]`,
			},
			{
				jsonpath:    `$[0,1][0,1]`,
				inputJSON:   `[[],[],[7]]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$[0,1][0,1]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$[0,1][0,1]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
			},
		},
		`Filter-qualifier`: []TestCase{
			{
				jsonpath:     `$[0,1][?(@.b)]`,
				inputJSON:    `[[{"a":1},{"b":2}],[{"a":3},{"b":4}]]`,
				expectedJSON: `[{"b":2},{"b":4}]`,
			},
			{
				jsonpath:    `$[0,1][?(@.b)]`,
				inputJSON:   `[[{"a":1},{"x":2}],[{"a":3},{"x":4}]]`,
				expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
			},
			{
				jsonpath:    `$[0,1][?(@.b)]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$[0,1][?(@.b)]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[0,1]`, `array`, `string`),
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_valueGroupCombination_Filter_qualifier(t *testing.T) {
	testGroups := TestGroup{
		`Recursive-descent`: []TestCase{
			{
				jsonpath:     `$[?(@.b)]..a`,
				inputJSON:    `[{"a":1},{"b":{"a":2}},{"c":3},{"b":[{"a":4}]}]`,
				expectedJSON: `[2,4]`,
			},
			{
				jsonpath:    `$[?(@.b)]..a`,
				inputJSON:   `[{"a":1},{"b":{"x":2}},{"c":3},{"b":[{"x":4}]}]`,
				expectedErr: createErrorMemberNotExist(`a`),
			},
			{
				jsonpath:    `$[?(@.b)]..a`,
				inputJSON:   `[{"a":1},{"b":"a"},{"c":3},{"b":"a"}]`,
				expectedErr: createErrorMemberNotExist(`a`),
			},
			{
				jsonpath:    `$[?(@.b)]..a`,
				inputJSON:   `[{"a":1},{"x":{"a":2}},{"c":3},{"x":[{"a":4}]}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
			},
			{
				jsonpath:    `$[?(@.b)]..a`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[?(@.b)]`, `object/array`, `string`),
			},
		},
		`Multiple-identifier`: []TestCase{
			{
				jsonpath:     `$[?(@.b)]['a','c']`,
				inputJSON:    `[{"a":1},{"b":2},{"a":3,"b":4},{"c":5},{"a":6,"c":7},{"b":8,"c":9},{"a":10,"b":11,"c":12}]`,
				expectedJSON: `[3,9,10,12]`,
			},
			{
				jsonpath:    `$[?(@.b)]['a','c']`,
				inputJSON:   `[{"a":1},{"b":2},{"x":3,"b":4},{"c":5},{"a":6,"c":7},{"b":8,"z":9},{"x":10,"b":11,"z":12}]`,
				expectedErr: createErrorMemberNotExist(`['a','c']`),
			},
			{
				jsonpath:    `$[?(@.b)]['a','c']`,
				inputJSON:   `[{"a":1},{"x":2},{"a":3,"x":4},{"c":5},{"a":6,"c":7},{"x":8,"c":9},{"a":10,"x":11,"c":12}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
			},
			{
				jsonpath:    `$[?(@.b)]['a','c']`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[?(@.b)]`, `object/array`, `string`),
			},
		},
		`Wildcard-identifier`: []TestCase{
			{
				jsonpath:     `$[?(@.b)].*`,
				inputJSON:    `[{"a":1},{"b":2},{"a":3,"b":4},{"c":5},{"a":6,"c":7},{"b":8,"c":9},{"a":10,"b":11,"c":12}]`,
				expectedJSON: `[2,3,4,8,9,10,11,12]`,
			},
			{
				jsonpath:    `$[?(@.b)].*`,
				inputJSON:   `[{"a":1},{"x":2},{"a":3,"x":4},{"c":5},{"a":6,"c":7},{"x":8,"c":9},{"a":10,"x":11,"c":12}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
			},
			{
				jsonpath:    `$[?(@.b)].*`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[?(@.b)]`, `object/array`, `string`),
			},
		},
		`Slice-qualifier`: []TestCase{
			{
				jsonpath:     `$[?(@)][0:2]`,
				inputJSON:    `[[1,2],[3,4],[5,6,7]]`,
				expectedJSON: `[1,2,3,4,5,6]`,
			},
			{
				jsonpath:    `$[?(@)][0:2]`,
				inputJSON:   `[[],[],[]]`,
				expectedErr: createErrorMemberNotExist(`[0:2]`),
			},
			{
				jsonpath:    `$[?(@)][0:2]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[?(@)]`),
			},
			{
				jsonpath:    `$[?(@)][0:2]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[?(@)]`, `object/array`, `string`),
			},
		},
		`Wildcard-qualifier`: []TestCase{
			{
				jsonpath:     `$[?(@)][*]`,
				inputJSON:    `[[1,2],[3,4],[5,6,7]]`,
				expectedJSON: `[1,2,3,4,5,6,7]`,
			},
			{
				jsonpath:    `$[?(@)][*]`,
				inputJSON:   `[[],[],[]]`,
				expectedErr: createErrorMemberNotExist(`[*]`),
			},
			{
				jsonpath:    `$[?(@)][*]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[?(@)]`),
			},
			{
				jsonpath:    `$[?(@)][*]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[?(@)]`, `object/array`, `string`),
			},
		},
		`Union-in-qualifier`: []TestCase{
			{
				jsonpath:     `$[?(@)][0,1]`,
				inputJSON:    `[[1,2],[3,4],[5,6,7]]`,
				expectedJSON: `[1,2,3,4,5,6]`,
			},
			{
				jsonpath:    `$[?(@)][0,1]`,
				inputJSON:   `[[],[],[]]`,
				expectedErr: createErrorMemberNotExist(`[0,1]`),
			},
			{
				jsonpath:    `$[?(@)][0,1]`,
				inputJSON:   `[]`,
				expectedErr: createErrorMemberNotExist(`[?(@)]`),
			},
			{
				jsonpath:    `$[?(@)][0,1]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[?(@)]`, `object/array`, `string`),
			},
		},
		`Filter-qualifier`: []TestCase{
			{
				jsonpath:     `$[?(@.a)][?(@.b)]`,
				inputJSON:    `[{"a":{"b":2}},{"b":{"a":1}},{"a":{"a":3}}]`,
				expectedJSON: `[{"b":2}]`,
			},
			{
				jsonpath:    `$[?(@.a)][?(@.b)]`,
				inputJSON:   `[{"a":{"x":2}},{"b":{"a":1}},{"a":{"a":3}}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.b)]`),
			},
			{
				jsonpath:    `$[?(@.a)][?(@.b)]`,
				inputJSON:   `[{"x":{"b":2}},{"b":{"a":1}},{"x":{"a":3}}]`,
				expectedErr: createErrorMemberNotExist(`[?(@.a)]`),
			},
			{
				jsonpath:    `$[?(@.a)][?(@.b)]`,
				inputJSON:   `"x"`,
				expectedErr: createErrorTypeUnmatched(`[?(@.a)]`, `object/array`, `string`),
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_space(t *testing.T) {
	testGroups := TestGroup{
		`Space`: []TestCase{
			{
				jsonpath:     ` $.a `,
				inputJSON:    `{"a":123}`,
				expectedJSON: `[123]`,
			},
			{
				jsonpath:    "\t" + `$.a` + "\t",
				inputJSON:   `{"a":123}`,
				expectedErr: ErrorInvalidSyntax{position: 0, reason: `unrecognized input`, near: "\t" + `$.a` + "\t"},
			},
			{
				jsonpath:    `$.a` + "\n",
				inputJSON:   `{"a":123}`,
				expectedErr: ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: "\n"},
			},
			{
				jsonpath:     `$[ "a" , "c" ]`,
				inputJSON:    `{"a":1,"b":2,"c":3}`,
				expectedJSON: `[1,3]`,
			},
			{
				jsonpath:     `$[ 0 , 2 : 4 , * ]`,
				inputJSON:    `[1,2,3,4,5]`,
				expectedJSON: `[1,3,4,1,2,3,4,5]`,
			},
			{
				jsonpath:     `$[ ?( @.a == 1 ) ]`,
				inputJSON:    `[{"a":1}]`,
				expectedJSON: `[{"a":1}]`,
			},
			{
				jsonpath:     `$[ ?( @.a != 1 ) ]`,
				inputJSON:    `[{"a":2}]`,
				expectedJSON: `[{"a":2}]`,
			},
			{
				jsonpath:     `$[ ?( @.a <= 1 ) ]`,
				inputJSON:    `[{"a":1}]`,
				expectedJSON: `[{"a":1}]`,
			},
			{
				jsonpath:     `$[ ?( @.a < 1 ) ]`,
				inputJSON:    `[{"a":0}]`,
				expectedJSON: `[{"a":0}]`,
			},
			{
				jsonpath:     `$[ ?( @.a >= 1 ) ]`,
				inputJSON:    `[{"a":1}]`,
				expectedJSON: `[{"a":1}]`,
			},
			{
				jsonpath:     `$[ ?( @.a > 1 ) ]`,
				inputJSON:    `[{"a":2}]`,
				expectedJSON: `[{"a":2}]`,
			},
			{
				jsonpath:     `$[ ?( @.a =~ /a/ ) ]`,
				inputJSON:    `[{"a":"abc"}]`,
				expectedJSON: `[{"a":"abc"}]`,
			},
			{
				jsonpath:     `$[ ?( @.a == 1 && @.b == 2 ) ]`,
				inputJSON:    `[{"a":1,"b":2}]`,
				expectedJSON: `[{"a":1,"b":2}]`,
			},
			{
				jsonpath:     `$[ ?( @.a == 1 || @.b == 2 ) ]`,
				inputJSON:    `[{"a":1},{"b":2}]`,
				expectedJSON: `[{"a":1},{"b":2}]`,
			},
			{
				jsonpath:     `$[ ?( ! @.a ) ]`,
				inputJSON:    `[{"a":1},{"b":2}]`,
				expectedJSON: `[{"b":2}]`,
			},
			{
				jsonpath:     `$[ ?( ( @.a ) ) ]`,
				inputJSON:    `[{"a":1},{"b":2}]`,
				expectedJSON: `[{"a":1}]`,
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_invalidSyntax(t *testing.T) {
	testGroups := TestGroup{
		`root-identifier`: []TestCase{
			{
				jsonpath:    ``,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 0, reason: `unrecognized input`, near: ``},
			},
			{
				jsonpath:    `$$`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `$`},
			},
		},
		`root-less-identifier`: []TestCase{
			{
				jsonpath:    `a.`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.`},
			},
			{
				jsonpath:    `b.`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.`},
			},
		},
		`dot-child-identifier`: []TestCase{
			{
				jsonpath:    `$a`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `a`},
			},
			{
				jsonpath:    `.`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 0, reason: `unrecognized input`, near: `.`},
			},
			{
				jsonpath:    `$.`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.`},
			},
			{
				jsonpath:    `..`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 0, reason: `unrecognized input`, near: `..`},
			},
			{
				jsonpath:    `$..`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `..`},
			},
			{
				jsonpath:    `$.a.`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: `.`},
			},
			{
				jsonpath:    `$.a..`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: `..`},
			},
			{
				jsonpath:    `$..a.`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `unrecognized input`, near: `.`},
			},
			{
				jsonpath:    `$..a..`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `unrecognized input`, near: `..`},
			},
			{
				jsonpath:    `$...a`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `...a`},
			},
		},
		`bracket-child-identifier`: []TestCase{
			{
				jsonpath:    `$['a]`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['a]`},
			},
			{
				jsonpath:    `$["a]`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `["a]`},
			},
			{
				jsonpath:    `$[a']`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[a']`},
			},
			{
				jsonpath:    `$[a"]`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[a"]`},
			},
			{
				jsonpath:    `$[a]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[a]`},
			},
			{
				jsonpath:    `$.['a']`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.['a']`},
			},
			{
				jsonpath:    `$.["a"]`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.["a"]`},
			},
			{
				jsonpath:    `$.[a]`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.[a]`},
			},
			{
				jsonpath:    `$['a'.'b']`,
				inputJSON:   `["a"]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['a'.'b']`},
			},
			{
				jsonpath:    `$[a.b]`,
				inputJSON:   `[{"a":{"b":1}}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[a.b]`},
			},
			{
				jsonpath:    `$['a'b']`,
				inputJSON:   `["a"]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['a'b']`},
			},
			{
				jsonpath:    `$['a\\'b']`,
				inputJSON:   `["a"]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['a\\'b']`},
			},
			{
				jsonpath:    `$['ab\']`,
				inputJSON:   `["a"]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['ab\']`},
			},
		},
		`qualifier::index`: []TestCase{
			{
				jsonpath:    `$[0].[1]`,
				inputJSON:   `[["a","b"],["c"],["d"]]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `unrecognized input`, near: `.[1]`},
			},
			{
				jsonpath:    `$[0,1].[1]`,
				inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
				expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `.[1]`},
			},
			{
				jsonpath:    `$[0:2].[1]`,
				inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
				expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `.[1]`},
			},
		},
		`qualifier::union`: []TestCase{
			{
				jsonpath:    `$[0].[1,2]`,
				inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `unrecognized input`, near: `.[1,2]`},
			},
			{
				jsonpath:    `$[0,1].[1,2]`,
				inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
				expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `.[1,2]`},
			},
			{
				jsonpath:    `$[0:2].[1,2]`,
				inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
				expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `.[1,2]`},
			},
		},
		`qualifier::slice`: []TestCase{
			{
				jsonpath:    `$[0].[1:3]`,
				inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `unrecognized input`, near: `.[1:3]`},
			},
			{
				jsonpath:    `$[0,1].[1:3]`,
				inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
				expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `.[1:3]`},
			},
			{
				jsonpath:    `$[0:1].[1:3]`,
				inputJSON:   `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
				expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `.[1:3]`},
			},
		},
		`qualifier::filter`: []TestCase{
			{
				jsonpath:    `$[?(@.a),?(@.b)]`,
				inputJSON:   `{}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a),?(@.b)]`},
			},
		},
		`qualifier::empty`: []TestCase{
			{
				jsonpath:    `$[]`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[]`},
			},
			{
				jsonpath:    `$.a[]`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: `[]`},
			},
			{
				jsonpath:    `$.a.b[]`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 5, reason: `unrecognized input`, near: `[]`},
			},
			{
				jsonpath:    `$[]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[]`},
			},
			{
				jsonpath:    `$[?()]`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?()]`},
			},
			{
				jsonpath:    `$[()]`,
				inputJSON:   `{}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[()]`},
			},
		},
		`qualifier::brace`: []TestCase{
			{
				jsonpath:    `$()`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `()`},
			},
			{
				jsonpath:    `$(a)`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `(a)`},
			},
			{
				jsonpath:    `$[`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[`},
			},
			{
				jsonpath:    `$[(`,
				inputJSON:   `{}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[(`},
			},
			{
				jsonpath:    `$[(]`,
				inputJSON:   `{}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[(]`},
			},
			{
				jsonpath:    `$[0`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[0`},
			},
			{
				jsonpath:    `$[?@a]`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?@a]`},
			},
			{
				jsonpath:    `$[?(@.a=='abc`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=='abc`},
			},
			{
				jsonpath:    `$[?(@.a=="abc`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=="abc`},
			},
			{
				jsonpath:    `$[?((@.a>1 )]`,
				inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?((@.a>1 )]`},
			},
			{
				jsonpath:    `$[?((@.a>1`,
				inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?((@.a>1`},
			},
		},
		`qualifier::big-number`: []TestCase{
			{
				jsonpath:    `$[0,10000000000000000000,]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[0,10000000000000000000,]`},
			},
			{
				jsonpath:    `$[0:10000000000000000000:a]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[0:10000000000000000000:a]`},
			},
		},
		`qualifier::comparator`: []TestCase{
			{
				jsonpath:    `$[?(@.a!!=1)]`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a!!=1)]`},
			},
			{
				jsonpath:    `$[?(@.a==)]`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==)]`},
			},
			{
				jsonpath:    `$[?(@.a!=)]`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a!=)]`},
			},
			{
				jsonpath:    `$[?(@.a<=)]`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a<=)]`},
			},
			{
				jsonpath:    `$[?(@.a<)]`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a<)]`},
			},
			{
				jsonpath:    `$[?(@.a>=)]`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a>=)]`},
			},
			{
				jsonpath:    `$[?(@.a>)]`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a>)]`},
			},
			{
				jsonpath:    `$[?(==@.a)]`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(==@.a)]`},
			},
			{
				jsonpath:    `$[?(!=@.a)]`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(!=@.a)]`},
			},
			{
				jsonpath:    `$[?(<=@.a)]`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(<=@.a)]`},
			},
			{
				jsonpath:    `$[?(<@.a)]`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(<@.a)]`},
			},
			{
				jsonpath:    `$[?(>=@.a)]`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(>=@.a)]`},
			},
			{
				jsonpath:    `$[?(>@.a)]`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(>@.a)]`},
			},
			{
				jsonpath:    `$[?(@.a===1)]`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a===1)]`},
			},
			{
				jsonpath:    `$[?(@.a=2)]`,
				inputJSON:   `[{"a":2}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=2)]`},
			},
			{
				jsonpath:    `$[?(@.a<>2)]`,
				inputJSON:   `[{"a":2}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a<>2)]`},
			},
			{
				jsonpath:    `$[?(@.a=<2)]`,
				inputJSON:   `[{"a":2}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=<2)]`},
			},
		},
		`qualifier::literal`: []TestCase{
			{
				jsonpath:     `$[?(false)]`,
				inputJSON:    `[0,1,false,true,null,{},[]]`,
				expectedJSON: `[]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(false)]`},
			},
			{
				jsonpath:     `$[?(true)]`,
				inputJSON:    `[0,1,false,true,null,{},[]]`,
				expectedJSON: `[]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(true)]`},
			},
			{
				jsonpath:     `$[?(null)]`,
				inputJSON:    `[0,1,false,true,null,{},[]]`,
				expectedJSON: `[]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(null)]`},
			},
			{
				jsonpath:    `$[?(@.a==["b"])]`,
				inputJSON:   `[{"a":["b"]}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==["b"])]`},
			},
			{
				jsonpath:    `$[?(@[0:1]==[1])]`,
				inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@[0:1]==[1])]`},
			},
			{
				jsonpath:    `$[?(@.*==[1,2])]`,
				inputJSON:   `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.*==[1,2])]`},
			},
			{
				jsonpath:    `$[?(@.*==['1','2'])]`,
				inputJSON:   `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.*==['1','2'])]`},
			},
			{
				jsonpath:    `$[?(@=={"k":"v"})]`,
				inputJSON:   `{}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@=={"k":"v"})]`},
			},
			{
				jsonpath:     `$[?(@.a==fAlse)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==fAlse)]`},
			},
			{
				jsonpath:     `$[?(@.a==faLse)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==faLse)]`},
			},
			{
				jsonpath:     `$[?(@.a==falSe)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==falSe)]`},
			},
			{
				jsonpath:     `$[?(@.a==falsE)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==falsE)]`},
			},
			{
				jsonpath:     `$[?(@.a==FaLse)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==FaLse)]`},
			},
			{
				jsonpath:     `$[?(@.a==FalSe)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==FalSe)]`},
			},
			{
				jsonpath:     `$[?(@.a==FalsE)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==FalsE)]`},
			},
			{
				jsonpath:     `$[?(@.a==FaLSE)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==FaLSE)]`},
			},
			{
				jsonpath:     `$[?(@.a==FAlSE)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==FAlSE)]`},
			},
			{
				jsonpath:     `$[?(@.a==FALsE)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==FALsE)]`},
			},
			{
				jsonpath:     `$[?(@.a==FALSe)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==FALSe)]`},
			},
			{
				jsonpath:     `$[?(@.a==tRue)]`,
				inputJSON:    `[{"a":true}]`,
				expectedJSON: `[{"a":true}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==tRue)]`},
			},
			{
				jsonpath:     `$[?(@.a==trUe)]`,
				inputJSON:    `[{"a":true}]`,
				expectedJSON: `[{"a":true}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==trUe)]`},
			},
			{
				jsonpath:     `$[?(@.a==truE)]`,
				inputJSON:    `[{"a":true}]`,
				expectedJSON: `[{"a":true}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==truE)]`},
			},
			{
				jsonpath:     `$[?(@.a==TrUe)]`,
				inputJSON:    `[{"a":true}]`,
				expectedJSON: `[{"a":true}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==TrUe)]`},
			},
			{
				jsonpath:     `$[?(@.a==TruE)]`,
				inputJSON:    `[{"a":true}]`,
				expectedJSON: `[{"a":true}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==TruE)]`},
			},
			{
				jsonpath:     `$[?(@.a==TrUE)]`,
				inputJSON:    `[{"a":true}]`,
				expectedJSON: `[{"a":true}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==TrUE)]`},
			},
			{
				jsonpath:     `$[?(@.a==TRuE)]`,
				inputJSON:    `[{"a":true}]`,
				expectedJSON: `[{"a":true}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==TRuE)]`},
			},
			{
				jsonpath:     `$[?(@.a==TRUe)]`,
				inputJSON:    `[{"a":true}]`,
				expectedJSON: `[{"a":true}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==TRUe)]`},
			},
			{
				jsonpath:     `$[?(@.a==nUll)]`,
				inputJSON:    `[{"a":null}]`,
				expectedJSON: `[{"a":null}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==nUll)]`},
			},
			{
				jsonpath:     `$[?(@.a==nuLl)]`,
				inputJSON:    `[{"a":null}]`,
				expectedJSON: `[{"a":null}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==nuLl)]`},
			},
			{
				jsonpath:     `$[?(@.a==nulL)]`,
				inputJSON:    `[{"a":null}]`,
				expectedJSON: `[{"a":null}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==nulL)]`},
			},
			{
				jsonpath:     `$[?(@.a==NuLl)]`,
				inputJSON:    `[{"a":null}]`,
				expectedJSON: `[{"a":null}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==NuLl)]`},
			},
			{
				jsonpath:     `$[?(@.a==NulL)]`,
				inputJSON:    `[{"a":null}]`,
				expectedJSON: `[{"a":null}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==NulL)]`},
			},
			{
				jsonpath:     `$[?(@.a==NuLL)]`,
				inputJSON:    `[{"a":null}]`,
				expectedJSON: `[{"a":null}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==NuLL)]`},
			},
			{
				jsonpath:     `$[?(@.a==NUlL)]`,
				inputJSON:    `[{"a":null}]`,
				expectedJSON: `[{"a":null}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==NUlL)]`},
			},
			{
				jsonpath:     `$[?(@.a==NULl)]`,
				inputJSON:    `[{"a":null}]`,
				expectedJSON: `[{"a":null}]`,
				expectedErr:  ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a==NULl)]`},
			},
		},
		`qualifier::sub-expression`: []TestCase{
			{
				jsonpath:    `$[?((@.a<2)==false)]`,
				inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?((@.a<2)==false)]`},
			},
			{
				jsonpath:    `$[?((@.a<2)==true)]`,
				inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?((@.a<2)==true)]`},
			},
			{
				jsonpath:    `$[?((@.a<2)==1)]`,
				inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?((@.a<2)==1)]`},
			},
		},
		`qualifier::logical-operator`: []TestCase{
			{
				jsonpath:    `$[?(@.a>1 && )]`,
				inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a>1 && )]`},
			},
			{
				jsonpath:    `$[?(@.a>1 || )]`,
				inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a>1 || )]`},
			},
			{
				jsonpath:    `$[?( && @.a>1 )]`,
				inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?( && @.a>1 )]`},
			},
			{
				jsonpath:    `$[?( || @.a>1 )]`,
				inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?( || @.a>1 )]`},
			},
			{
				jsonpath:    `$[?(@.a>1 && false)]`,
				inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a>1 && false)]`},
			},
			{
				jsonpath:    `$[?(@.a>1 && true)]`,
				inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a>1 && true)]`},
			},
			{
				jsonpath:    `$[?(@.a>1 || false)]`,
				inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a>1 || false)]`},
			},
			{
				jsonpath:    `$[?(@.a>1 || true)]`,
				inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a>1 || true)]`},
			},
			{
				jsonpath:    `$[?(@.a>1 && ())]`,
				inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a>1 && ())]`},
			},
			{
				jsonpath:    `$[?(@.a & @.b)]`,
				inputJSON:   `{}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a & @.b)]`},
			},
			{
				jsonpath:    `$[?(@.a | @.b)]`,
				inputJSON:   `{}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a | @.b)]`},
			},
			{
				jsonpath:    `$[?(!(@.a==2))]`,
				inputJSON:   `[{"a":1.9999},{"a":2},{"a":2.0001},{"a":"2"},{"a":true},{"a":{}},{"a":[]},{"a":["b"]},{"a":{"a":"value"}},{"b":"value"}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(!(@.a==2))]`},
			},
			{
				jsonpath:    `$[?(!(@.a<2))]`,
				inputJSON:   `[{"a":1.9999},{"a":2},{"a":2.0001},{"a":"2"},{"a":true},{"a":{}},{"a":[]},{"a":["b"]},{"a":{"a":"value"}},{"b":"value"}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(!(@.a<2))]`},
			},
		},
		`qualifier::regular-expression`: []TestCase{
			{
				jsonpath:    `$[?(@.a=~/abc)]`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=~/abc)]`},
			},
			{
				jsonpath:    `$[?(@.a=~///)]`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=~///)]`},
			},
			{
				jsonpath:    `$[?(@.a=~s/a/b/)]`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=~s/a/b/)]`},
			},
			{
				jsonpath:    `$[?(@.a=~@abc@)]`,
				inputJSON:   `[]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(@.a=~@abc@)]`},
			},
			{
				jsonpath:    `$[?(a=~/123/)]`,
				inputJSON:   `[{"a":"123"},{"a":123}]`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[?(a=~/123/)]`},
			},
		},
		`function`: []TestCase{
			{
				jsonpath:    `$.func(`,
				inputJSON:   `{}`,
				expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `(`},
			},
			{
				jsonpath:    `$.func(a`,
				inputJSON:   `{}`,
				expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `(a`},
			},
			{
				jsonpath:    `$.func(a)`,
				inputJSON:   `{}`,
				expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `(a)`},
			},
			{
				jsonpath:    `$.func()(`,
				inputJSON:   `{}`,
				expectedErr: ErrorFunctionNotFound{function: `.func()`},
			},
			{
				jsonpath:    `$.func(){}`,
				inputJSON:   `{}`,
				expectedErr: ErrorFunctionNotFound{function: `.func()`},
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_invalidArgument(t *testing.T) {
	testGroups := TestGroup{
		`number-overflow`: []TestCase{
			{
				jsonpath:  `$[10000000000000000000]`,
				inputJSON: `["first","second","third"]`,
				expectedErr: ErrorInvalidArgument{
					argument: `10000000000000000000`,
					err:      fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
				},
			},
			{
				jsonpath:  `$[0,10000000000000000000]`,
				inputJSON: `["first","second","third"]`,
				expectedErr: ErrorInvalidArgument{
					argument: `10000000000000000000`,
					err:      fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
				},
			},
			{
				jsonpath:  `$[10000000000000000000:1]`,
				inputJSON: `["first","second","third"]`,
				expectedErr: ErrorInvalidArgument{
					argument: `10000000000000000000`,
					err:      fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
				},
			},
			{
				jsonpath:  `$[1:10000000000000000000]`,
				inputJSON: `["first","second","third"]`,
				expectedErr: ErrorInvalidArgument{
					argument: `10000000000000000000`,
					err:      fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
				},
			},
			{
				jsonpath:  `$[0:3:10000000000000000000]`,
				inputJSON: `["first","second","third"]`,
				expectedErr: ErrorInvalidArgument{
					argument: `10000000000000000000`,
					err:      fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
				},
			},
		},
		`number-syntax`: []TestCase{
			{
				jsonpath:  `$[?(@.a==1e1abc)]`,
				inputJSON: `{}`,
				expectedErr: ErrorInvalidArgument{
					argument: `1e1abc`,
					err:      fmt.Errorf(`strconv.ParseFloat: parsing "1e1abc": invalid syntax`),
				},
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieve_notSupported(t *testing.T) {
	testGroups := TestGroup{
		`Not supported`: []TestCase{
			{
				jsonpath:    `$[(command)]`,
				inputJSON:   `{}`,
				expectedErr: ErrorNotSupported{feature: `script`, path: `[(command)]`},
			},
			{
				jsonpath:    `$[( command )]`,
				inputJSON:   `{}`,
				expectedErr: ErrorNotSupported{feature: `script`, path: `[(command)]`},
			},
		},
	}

	runTestGroups(t, testGroups)
}

var useJSONNumberDecoderFunction = func(srcJSON string, src *interface{}) error {
	reader := strings.NewReader(srcJSON)
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()
	return decoder.Decode(src)
}

func TestRetrieve_jsonNumber(t *testing.T) {
	testGroups := TestGroup{
		`filter`: []TestCase{
			{
				jsonpath:      `$[?(@.a > 123)].a`,
				inputJSON:     `[{"a":123.456}]`,
				expectedJSON:  `[123.456]`,
				unmarshalFunc: useJSONNumberDecoderFunction,
			},
			{
				jsonpath:      `$[?(@.a > 123.46)].a`,
				inputJSON:     `[{"a":123.456}]`,
				expectedJSON:  `[]`,
				expectedErr:   createErrorMemberNotExist(`[?(@.a > 123.46)]`),
				unmarshalFunc: useJSONNumberDecoderFunction,
			},
			{
				jsonpath:      `$[?(@.a > 122)].a`,
				inputJSON:     `[{"a":123}]`,
				expectedJSON:  `[123]`,
				unmarshalFunc: useJSONNumberDecoderFunction,
			},
			{
				jsonpath:      `$[?(123 < @.a)].a`,
				inputJSON:     `[{"a":123.456}]`,
				expectedJSON:  `[123.456]`,
				unmarshalFunc: useJSONNumberDecoderFunction,
			},
			{
				jsonpath:      `$[?(@.a==-0.123e2)]`,
				inputJSON:     `[{"a":-12.3,"b":1},{"a":-0.123e2,"b":2},{"a":-0.123},{"a":-12},{"a":12.3},{"a":2},{"a":"-0.123e2"}]`,
				expectedJSON:  `[{"a":-12.3,"b":1},{"a":-0.123e2,"b":2}]`,
				unmarshalFunc: useJSONNumberDecoderFunction,
			},
			{
				jsonpath:      `$[?(@.a==11)]`,
				inputJSON:     `[{"a":10.999},{"a":11.00},{"a":11.10}]`,
				expectedJSON:  `[{"a":11.00}]`,
				unmarshalFunc: useJSONNumberDecoderFunction,
			},
		},
	}

	runTestGroups(t, testGroups)
}

var twiceFunc = func(param interface{}) (interface{}, error) {
	if input, ok := param.(float64); ok {
		return input * 2, nil
	}
	return nil, fmt.Errorf(`type error`)
}
var quarterFunc = func(param interface{}) (interface{}, error) {
	if input, ok := param.(float64); ok {
		return input / 4, nil
	}
	return nil, fmt.Errorf(`type error`)
}
var maxFunc = func(param []interface{}) (interface{}, error) {
	var result float64
	for _, value := range param {
		if result < value.(float64) {
			result = value.(float64)
		}
	}
	return result, nil
}
var minFunc = func(param []interface{}) (interface{}, error) {
	var result float64 = 999
	for _, value := range param {
		if result > value.(float64) {
			result = value.(float64)
		}
	}
	return result, nil
}
var errAggregateFunc = func(param []interface{}) (interface{}, error) {
	return nil, fmt.Errorf(`aggregate error`)
}
var errFilterFunc = func(param interface{}) (interface{}, error) {
	return nil, fmt.Errorf(`filter error`)
}

func TestRetrieve_configFunction(t *testing.T) {
	testGroups := TestGroup{
		`filter-function`: []TestCase{
			{
				jsonpath:     `$.*.twice()`,
				inputJSON:    `[123.456,256]`,
				expectedJSON: `[246.912,512]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`twice`: twiceFunc,
				},
			},
			{
				jsonpath:     `$.*.twice().twice()`,
				inputJSON:    `[123.456,256]`,
				expectedJSON: `[493.824,1024]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`twice`: twiceFunc,
				},
			},
			{
				jsonpath:     `$.*.twice().quarter()`,
				inputJSON:    `[123.456,256]`,
				expectedJSON: `[61.728,128]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`twice`:   twiceFunc,
					`quarter`: quarterFunc,
				},
			},
			{
				jsonpath:     `$.*.quarter().twice()`,
				inputJSON:    `[123.456,256]`,
				expectedJSON: `[61.728,128]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`twice`:   twiceFunc,
					`quarter`: quarterFunc,
				},
			},
			{
				jsonpath:     `$[?(@.twice())]`,
				inputJSON:    `[123.456,256]`,
				expectedJSON: `[123.456,256]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`twice`: twiceFunc,
				},
			},
			{
				jsonpath:     `$[?(@.twice() == 512)]`,
				inputJSON:    `[123.456,256]`,
				expectedJSON: `[256]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`twice`: twiceFunc,
				},
			},
			{
				jsonpath:     `$[?(512 != @.twice())]`,
				inputJSON:    `[123.456,256]`,
				expectedJSON: `[123.456]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`twice`: twiceFunc,
				},
			},
			{
				jsonpath:     `$[?(@.twice() == $[0].twice())]`,
				inputJSON:    `[123.456,256]`,
				expectedJSON: `[123.456]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`twice`: twiceFunc,
				},
			},
		},
		`aggregate-function`: []TestCase{
			{
				jsonpath:     `$.*.max()`,
				inputJSON:    `[122.345,123.45,123.456]`,
				expectedJSON: `[123.456]`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxFunc,
				},
			},
			{
				jsonpath:     `$.*.max().max()`,
				inputJSON:    `[122.345,123.45,123.456]`,
				expectedJSON: `[123.456]`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxFunc,
				},
			},
			{
				jsonpath:     `$.*.max().min()`,
				inputJSON:    `[122.345,123.45,123.456]`,
				expectedJSON: `[123.456]`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxFunc,
					`min`: minFunc,
				},
			},
			{
				jsonpath:     `$.*.min().max()`,
				inputJSON:    `[122.345,123.45,123.456]`,
				expectedJSON: `[122.345]`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxFunc,
					`min`: minFunc,
				},
			},
			{
				jsonpath:     `$[?(@.max())]`,
				inputJSON:    `[122.345,123.45,123.456]`,
				expectedJSON: `[122.345,123.45,123.456]`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxFunc,
				},
			},
			{
				jsonpath:     `$[?(@.max() == 123.45)]`,
				inputJSON:    `[122.345,123.45,123.456]`,
				expectedJSON: `[123.45]`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxFunc,
				},
			},
			{
				jsonpath:     `$[?(123.45 != @.max())]`,
				inputJSON:    `[122.345,123.45,123.456]`,
				expectedJSON: `[122.345,123.456]`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxFunc,
				},
			},
			{
				jsonpath:     `$[?(@.max() != 123.45)]`,
				inputJSON:    `[[122.345,123.45,123.456],[122.345,123.45]]`,
				expectedJSON: `[[122.345,123.45,123.456]]`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxFunc,
				},
			},
			{
				jsonpath:     `$[?(@.max() == $[1].max())]`,
				inputJSON:    `[[122.345,123.45,123.456],[122.345,123.45]]`,
				expectedJSON: `[[122.345,123.45]]`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxFunc,
				},
			},
		},
		`aggregate-filter-mix`: []TestCase{
			{
				jsonpath:     `$.*.max().twice()`,
				inputJSON:    `[122.345,123.45,123.456]`,
				expectedJSON: `[246.912]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`twice`: twiceFunc,
				},
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxFunc,
				},
			},
			{
				jsonpath:     `$.*.twice().max()`,
				inputJSON:    `[122.345,123.45,123.456]`,
				expectedJSON: `[246.912]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`twice`: twiceFunc,
				},
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxFunc,
				},
			},
		},
		`filter-error`: []TestCase{
			{
				jsonpath:  `$.errFilter()`,
				inputJSON: `[122.345,123.45,123.456]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter`: errFilterFunc,
				},

				expectedErr: createErrorFunctionFailed(`.errFilter()`, `filter error`),
			},
			{
				jsonpath:  `$.*.errFilter()`,
				inputJSON: `[122.345,123.45,123.456]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter`: errFilterFunc,
				},

				expectedErr: createErrorFunctionFailed(`.errFilter()`, `filter error`),
			},
			{
				jsonpath:  `$.*.max().errFilter()`,
				inputJSON: `[122.345,123.45,123.456]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter`: errFilterFunc,
				},
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxFunc,
				},
				expectedErr: createErrorFunctionFailed(`.errFilter()`, `filter error`),
			},
			{
				jsonpath:  `$.*.twice().errFilter()`,
				inputJSON: `[122.345,123.45,123.456]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter`: errFilterFunc,
					`twice`:     twiceFunc,
				},

				expectedErr: createErrorFunctionFailed(`.errFilter()`, `filter error`),
			}, {
				jsonpath:  `$.errFilter().twice()`,
				inputJSON: `[122.345,123.45,123.456]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter`: errFilterFunc,
					`twice`:     twiceFunc,
				},

				expectedErr: createErrorFunctionFailed(`.errFilter()`, `filter error`),
			},
			{
				jsonpath:  `$.*.errFilter().twice()`,
				inputJSON: `[122.345,123.45,123.456]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter`: errFilterFunc,
					`twice`:     twiceFunc,
				},

				expectedErr: createErrorFunctionFailed(`.errFilter()`, `filter error`),
			},
			{
				jsonpath:  `$.*.max().errFilter().twice()`,
				inputJSON: `[122.345,123.45,123.456]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter`: errFilterFunc,
					`twice`:     twiceFunc,
				},
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxFunc,
				},
				expectedErr: createErrorFunctionFailed(`.errFilter()`, `filter error`),
			},
		},
		`aggregate-error`: []TestCase{
			{
				jsonpath:  `$.*.errAggregate()`,
				inputJSON: `[122.345,123.45,123.456]`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`errAggregate`: errAggregateFunc,
				},
				expectedErr: createErrorFunctionFailed(`.errAggregate()`, `aggregate error`),
			},
			{
				jsonpath:  `$.*.max().errAggregate()`,
				inputJSON: `[122.345,123.45,123.456]`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`errAggregate`: errAggregateFunc,
					`max`:          maxFunc,
				},
				expectedErr: createErrorFunctionFailed(`.errAggregate()`, `aggregate error`),
			},
			{
				jsonpath:  `$.*.twice().errAggregate()`,
				inputJSON: `[122.345,123.45,123.456]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`twice`: twiceFunc,
				},
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`errAggregate`: errAggregateFunc,
				},
				expectedErr: createErrorFunctionFailed(`.errAggregate()`, `aggregate error`),
			},
			{
				jsonpath:  `$.*.errAggregate().twice()`,
				inputJSON: `[122.345,123.45,123.456]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`twice`: twiceFunc,
				},
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`errAggregate`: errAggregateFunc,
				},
				expectedErr: createErrorFunctionFailed(`.errAggregate()`, `aggregate error`),
			},
			{
				jsonpath:  `$.*.max().errAggregate().twice()`,
				inputJSON: `[122.345,123.45,123.456]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`twice`: twiceFunc,
				},
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`errAggregate`: errAggregateFunc,
					`max`:          maxFunc,
				},
				expectedErr: createErrorFunctionFailed(`.errAggregate()`, `aggregate error`),
			},
			{
				jsonpath:  `$.a.max()`,
				inputJSON: `{}`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxFunc,
				},
				expectedErr: createErrorMemberNotExist(`.a`),
			},
		},
		`jsonpath-error-with-filter`: []TestCase{
			{
				jsonpath:  `$.x.*.errAggregate()`,
				inputJSON: `{"a":[122.345,123.45,123.456]}`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`errAggregate`: errAggregateFunc,
				},
				expectedErr: createErrorMemberNotExist(`.x`),
			},
			{
				jsonpath:  `$.*.a.b.c.errFilter()`,
				inputJSON: `[{"a":{"b":1}},{"a":2}]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter`: errFilterFunc,
				},

				expectedErr: createErrorTypeUnmatched(`.c`, `object`, `float64`),
			},
			{
				jsonpath:  `$.*.a.b.c.errFilter1().errFilter2()`,
				inputJSON: `[{"a":{"b":1}},{"a":2}]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter1`: errFilterFunc,
					`errFilter2`: errFilterFunc,
				},

				expectedErr: createErrorTypeUnmatched(`.c`, `object`, `float64`),
			},
			{
				jsonpath:  `$.*.a.b.c.errAggregate()`,
				inputJSON: `[{"a":{"b":1}},{"a":2}]`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`errAggregate`: errAggregateFunc,
				},
				expectedErr: createErrorTypeUnmatched(`.c`, `object`, `float64`),
			},
			{
				jsonpath:  `$.*.a.b.c.errAggregate1().errAggregate2()`,
				inputJSON: `[{"a":{"b":1}},{"a":2}]`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`errAggregate1`: errAggregateFunc,
					`errAggregate2`: errAggregateFunc,
				},
				expectedErr: createErrorTypeUnmatched(`.c`, `object`, `float64`),
			},
			{
				jsonpath:  `$.*.a.b.c.errAggregate().errFilter()`,
				inputJSON: `[{"a":{"b":1}},{"a":2}]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter`: twiceFunc,
				},
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`errAggregate`: errAggregateFunc,
				},
				expectedErr: createErrorTypeUnmatched(`.c`, `object`, `float64`),
			},
			{
				jsonpath:  `$.*.a.b.c.errFilter().errAggregate()`,
				inputJSON: `[{"a":{"b":1}},{"a":2}]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter`: twiceFunc,
				},
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`errAggregate`: errAggregateFunc,
				},
				expectedErr: createErrorTypeUnmatched(`.c`, `object`, `float64`),
			},
		},
		`function-syntax-check`: []TestCase{
			{
				jsonpath:     `$.*.TWICE()`,
				inputJSON:    `[123.456,256]`,
				expectedJSON: `[246.912,512]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`TWICE`: twiceFunc,
				},
			},
			{
				jsonpath:     `$.*.123()`,
				inputJSON:    `[123.456,256]`,
				expectedJSON: `[246.912,512]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`123`: twiceFunc,
				},
			},
			{
				jsonpath:     `$.*.--()`,
				inputJSON:    `[123.456,256]`,
				expectedJSON: `[246.912,512]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`--`: twiceFunc,
				},
			},
			{
				jsonpath:     `$.*.__()`,
				inputJSON:    `[123.456,256]`,
				expectedJSON: `[246.912,512]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`__`: twiceFunc,
				},
			},
			{
				jsonpath:    `$.*.unknown()`,
				inputJSON:   `[123.456,256]`,
				expectedErr: ErrorFunctionNotFound{function: `.unknown()`},
			},
		},
	}

	runTestGroups(t, testGroups)
}

func createAccessorModeValidator(
	resultIndex int,
	expectedValue1, expectedValue2, expectedValue3 interface{},
	srcGetter func(interface{}) interface{},
	srcSetter func(interface{}, interface{})) func(interface{}, []interface{}) error {
	return func(src interface{}, actualObject []interface{}) error {
		accessor := actualObject[resultIndex].(Accessor)

		getValue := accessor.Get()
		if getValue != expectedValue1 {
			return fmt.Errorf(`Get : expect<%f> != actual<%f>`, expectedValue1, getValue)
		}

		accessor.Set(expectedValue2)

		newSrcValue := srcGetter(src)

		if newSrcValue != expectedValue2 {
			return fmt.Errorf(`Set : expect<%f> != actual<%f>`, expectedValue2, newSrcValue)
		}

		getValue = accessor.Get()
		if getValue != expectedValue2 {
			return fmt.Errorf(`Set -> Get : expect<%f> != actual<%f>`, expectedValue2, getValue)
		}

		srcSetter(src, expectedValue3)

		getValue = accessor.Get()
		if getValue != expectedValue3 {
			return fmt.Errorf(`Src -> Get : expect<%f> != actual<%f>`, expectedValue3, getValue)
		}

		return nil
	}
}

var getOnlyValidator = func(src interface{}, actualObject []interface{}) error {
	accessor := actualObject[0].(Accessor)

	if accessor.Set != nil {
		return fmt.Errorf(`Set != nil`)
	}

	if !reflect.DeepEqual(accessor.Get(), src) {
		return fmt.Errorf(`Get != src`)
	}

	return nil
}

var echoAggregateFunc = func(param []interface{}) (interface{}, error) {
	return param, nil
}

var sliceStructChangedResultValidator = func(src interface{}, actualObject []interface{}) error {
	srcArray := src.([]interface{})
	accessor := actualObject[0].(Accessor)

	accessor.Set(4) // srcArray:[1,4,3] , accessor:[1,4,3]
	if len(srcArray) != 3 || srcArray[1] != 4 {
		return fmt.Errorf(`Set -> Src : expect<%d> != actual<%d>`, 4, srcArray[1])
	}
	srcArray = append(srcArray[:1], srcArray[2:]...) // srcArray:[1,3] , accessor:[1,3,3]
	if len(srcArray) != 2 || accessor.Get() != 3.0 { // Go's marshal returns float value
		return fmt.Errorf(`Del -> Get : expect<%f> != actual<%f>`, 3.0, accessor.Get())
	}
	accessor.Set(5) // srcArray:[1,5] , accessor:[1,5,3]
	if len(srcArray) != 2 || srcArray[1] != 5 {
		return fmt.Errorf(`Del -> Set -> Src : expect<%d> != actual<%d>`, 5, srcArray[1])
	}
	srcArray = append(srcArray[:1], srcArray[2:]...) // srcArray:[1] , accessor:[1,5,3]
	if len(srcArray) != 1 || accessor.Get() != 5 {
		return fmt.Errorf(`Del x2 -> Get : expect<%d> != actual<%d>`, 5, accessor.Get())
	}
	accessor.Set(6) // srcArray:[1] , accessor:[1,6,3]
	if len(srcArray) != 1 {
		return fmt.Errorf(`Del x2 -> Set -> Len : expect<%d> != actual<%d>`, 1, len(srcArray))
	}
	srcArray = append(srcArray, 7) // srcArray:[1,7] , accessor:[1,7,3]
	if len(srcArray) != 2 || accessor.Get() != 7 {
		return fmt.Errorf(`Del x2 -> Add -> Get : expect<%d> != actual<%d>`, 7, accessor.Get())
	}
	srcArray = append(srcArray, 8) // srcArray:[1,7,8]    , accessor:[1,7,8]
	srcArray = append(srcArray, 9) // srcArray:[1,7,8,9]  , accessor:[1,7,8,9]
	srcArray[1] = 10               // srcArray:[1,10,8,9] , accessor:[1,10,8,9]
	if len(srcArray) != 4 || accessor.Get() != 10 {
		return fmt.Errorf(`Del x2 -> Add x3 -> Update -> Get : expect<%d> != actual<%d>`, 10, accessor.Get())
	}
	return nil
}

var mapStructChangedResultValidator = func(src interface{}, actualObject []interface{}) error {
	srcMap := src.(map[string]interface{})
	accessor := actualObject[0].(Accessor)

	accessor.Set(2) // srcMap:{"a":2} , accessor:{"a":2}
	if len(srcMap) != 1 || srcMap[`a`] != 2 {
		return fmt.Errorf(`Set -> Src : expect<%d> != actual<%d>`, 2, srcMap[`a`])
	}
	delete(srcMap, `a`) // srcMap:{} , accessor:{}
	if accessor.Get() != nil {
		return fmt.Errorf(`Del -> Get : expect<%v> != actual<%d>`, nil, accessor.Get())
	}
	accessor.Set(3) // srcMap:{"a":3} , accessor:{"a":3}
	if len(srcMap) != 1 || srcMap[`a`] != 3 {
		return fmt.Errorf(`Del -> Set -> Len : expect<%d> != actual<%d>`, 0, len(srcMap))
	}
	delete(srcMap, `a`) // srcMap:{} , accessor:{}
	srcMap[`a`] = 4     // srcMap:{"a":4} , accessor:{"a":4}
	if accessor.Get() != 4 {
		return fmt.Errorf(`Del -> Update -> Get : expect<%v> != actual<%d>`, 4, accessor.Get())
	}
	return nil
}

func TestRetrieve_configAccessorMode(t *testing.T) {
	testGroups := TestGroup{
		`getter-setter`: []TestCase{
			{
				jsonpath:     `$.b`,
				inputJSON:    `{"a":11,"b":22,"c":33}`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					0, 22.0, 33.0, 44.0,
					func(src interface{}) interface{} {
						return src.(map[string]interface{})[`b`]
					},
					func(src, value interface{}) {
						src.(map[string]interface{})[`b`] = value
					}),
			},
			{
				jsonpath:     `$['b']`,
				inputJSON:    `{"a":123,"b":456,"c":789}`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					0, 456.0, 246.0, 369.0,
					func(src interface{}) interface{} {
						return src.(map[string]interface{})[`b`]
					},
					func(src, value interface{}) {
						src.(map[string]interface{})[`b`] = value
					}),
			},
			{
				jsonpath:     `b`,
				inputJSON:    `{"a":11,"b":22,"c":33}`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					0, 22.0, 33.0, 44.0,
					func(src interface{}) interface{} {
						return src.(map[string]interface{})[`b`]
					},
					func(src, value interface{}) {
						src.(map[string]interface{})[`b`] = value
					}),
			},
			{
				jsonpath:     `$['a','b','c']`,
				inputJSON:    `{"a":11,"b":22}`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					1, 22.0, 44.0, 55.0,
					func(src interface{}) interface{} {
						return src.(map[string]interface{})[`b`]
					},
					func(src, value interface{}) {
						src.(map[string]interface{})[`b`] = value
					}),
			},
			{
				jsonpath:     `$.*`,
				inputJSON:    `{"a":11,"b":22,"c":33}`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					1, 22.0, 33.0, 44.0,
					func(src interface{}) interface{} {
						return src.(map[string]interface{})[`b`]
					},
					func(src, value interface{}) {
						src.(map[string]interface{})[`b`] = value
					}),
			},
			{
				jsonpath:     `$.*`,
				inputJSON:    `[11,22,33]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					1, 22.0, 33.0, 44.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[1]
					},
					func(src, value interface{}) {
						src.([]interface{})[1] = value
					}),
			},
			{
				jsonpath:     `$..a`,
				inputJSON:    `{"b":{"a":11}, "c":66, "a":77}`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					1, 11.0, 22.0, 44.0,
					func(src interface{}) interface{} {
						return src.(map[string]interface{})[`b`].(map[string]interface{})[`a`]
					},
					func(src, value interface{}) {
						src.(map[string]interface{})[`b`].(map[string]interface{})[`a`] = value
					}),
			},
			{
				jsonpath:     `$[1]`,
				inputJSON:    `[123.456,256,789]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					0, 256.0, 512.0, 1024.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[1]
					},
					func(src, value interface{}) {
						src.([]interface{})[1] = value
					}),
			},
			{
				jsonpath:     `$[2,1]`,
				inputJSON:    `[123.456,256,789]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					1, 256.0, 512.0, 1024.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[1]
					},
					func(src, value interface{}) {
						src.([]interface{})[1] = value
					}),
			},
			{
				jsonpath:     `$[0:2]`,
				inputJSON:    `[11,22,33,44]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					0, 11.0, 22.0, 44.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[0]
					},
					func(src, value interface{}) {
						src.([]interface{})[0] = value
					}),
			},
			{
				jsonpath:     `$[*]`,
				inputJSON:    `[11,22]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					0, 11.0, 22.0, 44.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[0]
					},
					func(src, value interface{}) {
						src.([]interface{})[0] = value
					}),
			},
			{
				jsonpath:     `$[0:2].a`,
				inputJSON:    `[{"a":11},{"a":22},{"a":33}]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					0, 11.0, 22.0, 44.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[0].(map[string]interface{})[`a`]
					},
					func(src, value interface{}) {
						src.([]interface{})[0].(map[string]interface{})[`a`] = value
					}),
			},
			{
				jsonpath:     `$[?(@==11||@==33)]`,
				inputJSON:    `{"a":11,"b":22,"c":33}`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					1, 33.0, 22.0, 44.0,
					func(src interface{}) interface{} {
						return src.(map[string]interface{})[`c`]
					},
					func(src, value interface{}) {
						src.(map[string]interface{})[`c`] = value
					}),
			},
			{
				jsonpath:     `$[?(@==11||@==33)]`,
				inputJSON:    `[11,22,33]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					1, 33.0, 22.0, 44.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[2]
					},
					func(src, value interface{}) {
						src.([]interface{})[2] = value
					}),
			},
			{
				jsonpath:     `$[?(@.a==11||@.a==33)].a`,
				inputJSON:    `[{"a":11},{"a":22},{"a":33}]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					1, 33.0, 22.0, 44.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[2].(map[string]interface{})[`a`]
					},
					func(src, value interface{}) {
						src.([]interface{})[2].(map[string]interface{})[`a`] = value
					}),
			},
			{
				jsonpath:     `$[?(@==$[1])]`,
				inputJSON:    `[11,22,33]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					0, 22.0, 33.0, 44.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[1]
					},
					func(src, value interface{}) {
						src.([]interface{})[1] = value
					}),
			},
		},
		`get-only`: []TestCase{
			{
				jsonpath:        `$`,
				inputJSON:       `[1,2,3]`,
				accessorMode:    true,
				resultValidator: getOnlyValidator,
			},
			{
				jsonpath:     `$.echo()`,
				inputJSON:    `[122.345,123.45,123.456]`,
				accessorMode: true,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`echo`: echoAggregateFunc,
				},
				resultValidator: getOnlyValidator,
			},
		},
		`convert-srcJSON`: []TestCase{
			{
				jsonpath:        `$[1]`,
				inputJSON:       `[1,2,3]`,
				accessorMode:    true,
				resultValidator: sliceStructChangedResultValidator,
			},
			{
				jsonpath:        `$.a`,
				inputJSON:       `{"a":1}`,
				accessorMode:    true,
				resultValidator: mapStructChangedResultValidator,
			},
		},
	}

	runTestGroups(t, testGroups)
}

func TestRetrieveExecTwice(t *testing.T) {
	jsonpath1 := `$.a`
	srcJSON1 := `{"a":123}`
	expectedOutput1 := "[123]"
	jsonpath2 := `$[1].b`
	srcJSON2 := `[123,{"b":456}]`
	expectedOutput2 := "[456]"

	var src1 interface{}
	if err := json.Unmarshal([]byte(srcJSON1), &src1); err != nil {
		t.Errorf("%s", err)
		return
	}
	var src2 interface{}
	if err := json.Unmarshal([]byte(srcJSON2), &src2); err != nil {
		t.Errorf("%s", err)
		return
	}

	actualObject1, err := Retrieve(jsonpath1, src1)
	if err != nil {
		t.Errorf("expected error<nil> != actual error<%s>\n", err)
		return
	}
	actualObject2, err := Retrieve(jsonpath2, src2)
	if err != nil {
		t.Errorf("expected error<nil> != actual error<%s>\n", err)
		return
	}

	actualOutputJSON1, err := json.Marshal(actualObject1)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	actualOutputJSON2, err := json.Marshal(actualObject2)
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	if string(actualOutputJSON1) != string(expectedOutput1) || string(actualOutputJSON2) != string(expectedOutput2) {
		t.Errorf("actualOutputJSON1<%s> != expectedOutput1<%s> || actualOutputJSON2<%s> != expectedOutput2<%s>\n",
			string(actualOutputJSON1), string(expectedOutput1), string(actualOutputJSON2), string(expectedOutput2))
		return
	}
}

func TestParserFuncExecTwice(t *testing.T) {
	jsonpath := `$.a`
	srcJSON1 := `{"a":1}`
	expectedOutput1 := "[1]"
	srcJSON2 := `{"a":2}`
	expectedOutput2 := "[2]"

	var src1 interface{}
	if err := json.Unmarshal([]byte(srcJSON1), &src1); err != nil {
		t.Errorf("%s", err)
		return
	}
	var src2 interface{}
	if err := json.Unmarshal([]byte(srcJSON2), &src2); err != nil {
		t.Errorf("%s", err)
		return
	}

	parserFunc, err := Parse(jsonpath)
	if err != nil {
		t.Errorf("expected error<nil> != actual error<%s>\n", err)
		return
	}

	actualObject1, err := parserFunc(src1)
	if err != nil {
		t.Errorf("expected error<nil> != actual error<%s>\n", err)
		return
	}
	actualObject2, err := parserFunc(src2)
	if err != nil {
		t.Errorf("expected error<nil> != actual error<%s>\n", err)
		return
	}

	actualOutputJSON1, err := json.Marshal(actualObject1)
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	actualOutputJSON2, err := json.Marshal(actualObject2)
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	if string(actualOutputJSON1) != string(expectedOutput1) || string(actualOutputJSON2) != string(expectedOutput2) {
		t.Errorf("actualOutputJSON1<%s> != expectedOutput1<%s> || actualOutputJSON2<%s> != expectedOutput2<%s>\n",
			string(actualOutputJSON1), string(expectedOutput1), string(actualOutputJSON2), string(expectedOutput2))
		return
	}
}

type UnsupportedStructChild struct {
	B string
	C int
}

type UnsupportedStructParent struct {
	A UnsupportedStructChild
}

func TestRetrieve_unsupportedStruct(t *testing.T) {
	inputJSON := UnsupportedStructParent{A: UnsupportedStructChild{B: `test`, C: 123}}
	jsonpath := `$.A.B`
	expectedError := createErrorTypeUnmatched(`.A`, `object`, `jsonpath.UnsupportedStructParent`)
	_, err := Retrieve(jsonpath, inputJSON)

	if reflect.TypeOf(expectedError) != reflect.TypeOf(err) ||
		fmt.Sprintf(`%s`, expectedError) != fmt.Sprintf(`%s`, err) {
		t.Errorf("expected error<%s> != actual error<%s>\n",
			expectedError, err)
	}
}

func TestPegParserExecuteFunctions(t *testing.T) {
	stdoutBackup := os.Stdout
	os.Stdout = nil

	parser := pegJSONPathParser{Buffer: `$`}
	parser.Init()
	parser.Parse()
	parser.Execute()

	parser.Print()
	parser.Reset()
	parser.PrintSyntaxTree()
	parser.SprintSyntaxTree()

	err := parseError{p: &parser, max: token32{begin: 0, end: 1}}
	_ = err.Error()

	parser.buffer = []rune{'\n'}
	_ = err.Error()

	parser.Parse(1)
	parser.Parse(3)

	Pretty(true)(&parser)
	parser.PrintSyntaxTree()

	_ = err.Error()

	Size(10)(&parser)

	parser.Init(func(p *pegJSONPathParser) error {
		return fmt.Errorf(`test error`)
	})

	parser.Buffer = ``
	parser.PrintSyntaxTree()

	memoizeFunc := DisableMemoize()
	memoizeFunc(&parser)

	os.Stdout = stdoutBackup
}
