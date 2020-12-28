package jsonpath

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
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

func execTestRetrieve(t *testing.T, inputJSON interface{}, testCase TestCase) []interface{} {
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
			return nil
		}
		t.Errorf("expected error<%s> != actual error<%s>\n",
			expectedError, err)
		return nil
	}
	if expectedError != nil {
		t.Errorf("expected error<%w> != actual error<none>\n", expectedError)
		return nil
	}

	return actualObject
}

func execTestRetrieveTestGroups(t *testing.T, testGroup TestGroup) {
	for testGroupName, testCases := range testGroup {
		for _, testCase := range testCases {
			testCase := testCase
			jsonPath := testCase.jsonpath
			srcJSON := testCase.inputJSON
			expectedOutputJSON := testCase.expectedJSON

			t.Run(
				fmt.Sprintf(`%s <%s> <%s>`, testGroupName, jsonPath, srcJSON),
				func(t *testing.T) {
					t.Parallel()

					var src interface{}
					var err error
					if testCase.unmarshalFunc != nil {
						err = testCase.unmarshalFunc(srcJSON, &src)
					} else {
						err = json.Unmarshal([]byte(srcJSON), &src)
					}
					if err != nil {
						t.Errorf("%w", err)
						return
					}

					actualObject := execTestRetrieve(t, src, testCase)
					if t.Failed() {
						return
					}

					if actualObject == nil {
						return
					}

					if testCase.resultValidator != nil {
						err := testCase.resultValidator(src, actualObject)
						if err != nil {
							t.Errorf("%w", err)
						}
						return
					}

					actualOutputJSON, err := json.Marshal(actualObject)
					if err != nil {
						t.Errorf("%w", err)
						return
					}

					if string(actualOutputJSON) != expectedOutputJSON {
						t.Errorf("expectedOutputJSON<%s> != actualOutputJSON<%s>\n",
							expectedOutputJSON, actualOutputJSON)
						return
					}

				})
		}
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
				expectedErr: ErrorTypeUnmatched{expectedType: `object`, foundType: `[]interface {}`, path: `.length`},
			},
		},
		`character-type::Non-ASCII-syntax-accepted-in-JSON`: []TestCase{
			{
				jsonpath:     `$.a-b`,
				inputJSON:    `{"a-b":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.a:b`,
				inputJSON:    `{"a:b":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.$`,
				inputJSON:    `{"$":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.@`,
				inputJSON:    `{"@":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.'a'`,
				inputJSON:    `{"'a'":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$."a"`,
				inputJSON:    `{"\"a\"":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.'a.b'`,
				inputJSON:    `{"'a.b'":1,"a":{"b":2},"'a'":{"'b'":3},"'a":{"b'":4}}`,
				expectedJSON: `[4]`,
			},
		},
		`character-type::encoded-JSONPath`: []TestCase{
			{
				jsonpath:     `$.'a\.b'`,
				inputJSON:    `{"'a.b'":1,"a":{"b":2},"'a'":{"'b'":3},"'a":{"b'":4}}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.\\`,
				inputJSON:    `{"\\":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.\.`,
				inputJSON:    `{".":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.\[`,
				inputJSON:    `{"[":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.\(`,
				inputJSON:    `{"(":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.\)`,
				inputJSON:    `{")":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.\=`,
				inputJSON:    `{"=":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.\!`,
				inputJSON:    `{"!":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.\>`,
				inputJSON:    `{">":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.\<`,
				inputJSON:    `{"<":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.\ `,
				inputJSON:    `{" ":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:    `$.\` + "\t",
				inputJSON:   `{"":123}`,
				expectedErr: ErrorMemberNotExist{path: `.\` + "\t"},
			},
			{
				jsonpath:    `$.\` + "\r",
				inputJSON:   `{"":123}`,
				expectedErr: ErrorMemberNotExist{path: `.\` + "\r"},
			},
			{
				jsonpath:    `$.\` + "\n",
				inputJSON:   `{"":123}`,
				expectedErr: ErrorMemberNotExist{path: `.\` + "\n"},
			},
			{
				jsonpath:    `$.\a`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.\a`},
			},
			{
				jsonpath:     `$.a\\b`,
				inputJSON:    `{"a\\b":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.a\.b`,
				inputJSON:    `{"a.b":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.a\[b`,
				inputJSON:    `{"a[b":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.a\(b`,
				inputJSON:    `{"a(b":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.a\)b`,
				inputJSON:    `{"a)b":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.a\=b`,
				inputJSON:    `{"a=b":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.a\!b`,
				inputJSON:    `{"a!b":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.a\>b`,
				inputJSON:    `{"a>b":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.a\<b`,
				inputJSON:    `{"a<b":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$.a\ b`,
				inputJSON:    `{"a b":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:    `$.a\` + "\t" + `b`,
				inputJSON:   `{"ab":123}`,
				expectedErr: ErrorMemberNotExist{path: `.a\` + "\t" + `b`},
			},
			{
				jsonpath:    `$.a\` + "\r" + `b`,
				inputJSON:   `{"ab":123}`,
				expectedErr: ErrorMemberNotExist{path: `.a\` + "\r" + `b`},
			},
			{
				jsonpath:    `$.a\` + "\n" + `b`,
				inputJSON:   `{"ab":123}`,
				expectedErr: ErrorMemberNotExist{path: `.a\` + "\n" + `b`},
			},
			{
				jsonpath:    `$.a\a`,
				inputJSON:   `{"aa":1}`,
				expectedErr: ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: `\a`},
			},
		},
		`character-type::not-encoded-error`: []TestCase{
			{
				jsonpath:    `$.\`,
				inputJSON:   `{"\\":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.\`},
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
				jsonpath:    `$.=`,
				inputJSON:   `{"=":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.=`},
			},
			{
				jsonpath:    `$.!`,
				inputJSON:   `{"!":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.!`},
			},
			{
				jsonpath:    `$.>`,
				inputJSON:   `{">":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.>`},
			},
			{
				jsonpath:    `$.<`,
				inputJSON:   `{"<":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.<`},
			},
			{
				jsonpath:    `$. `,
				inputJSON:   `{" ":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `. `},
			},
			{
				jsonpath:    `$.` + "\t",
				inputJSON:   `{"":123}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.` + "\t"},
			},
			{
				jsonpath:    `$.` + "\r",
				inputJSON:   `{"":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.` + "\r"},
			},
			{
				jsonpath:    `$.` + "\n",
				inputJSON:   `{"":123}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.` + "\n"},
			},
			{
				jsonpath:    `$.a\b`,
				inputJSON:   `{"a\\b":1}`,
				expectedErr: ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: `\b`},
			},
			{
				jsonpath:    `$.a(b`,
				inputJSON:   `{"(":1}`,
				expectedErr: ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: `(b`},
			},
			{
				jsonpath:    `$.a)b`,
				inputJSON:   `{")":1}`,
				expectedErr: ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: `)b`},
			},
			{
				jsonpath:    `$.a=b`,
				inputJSON:   `{"=":1}`,
				expectedErr: ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: `=b`},
			},
			{
				jsonpath:    `$.a!b`,
				inputJSON:   `{"!":1}`,
				expectedErr: ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: `!b`},
			},
			{
				jsonpath:    `$.a>b`,
				inputJSON:   `{">":1}`,
				expectedErr: ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: `>b`},
			},
			{
				jsonpath:    `$.a<b`,
				inputJSON:   `{"<":1}`,
				expectedErr: ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: `<b`},
			},
			{
				jsonpath:    `$.a b`,
				inputJSON:   `{" ":1}`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `unrecognized input`, near: `b`},
			},
			{
				jsonpath:    `$.a` + "\t" + `b`,
				inputJSON:   `{"":123}`,
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `unrecognized input`, near: `b`},
			},
			{
				jsonpath:    `$.a` + "\r" + `b`,
				inputJSON:   `{"":1}`,
				expectedErr: ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: "\r" + `b`},
			},
			{
				jsonpath:    `$.a` + "\n" + `b`,
				inputJSON:   `{"":123}`,
				expectedErr: ErrorInvalidSyntax{position: 3, reason: `unrecognized input`, near: "\n" + `b`},
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
				expectedErr: ErrorMemberNotExist{path: `.d`},
			},
		},
		`type-unmatched`: []TestCase{
			{
				jsonpath:    `$.2`,
				inputJSON:   `["a","b",{"2":1}]`,
				expectedErr: ErrorTypeUnmatched{expectedType: `object`, foundType: `[]interface {}`, path: `.2`},
			},
			{
				jsonpath:    `$.-1`,
				inputJSON:   `["a","b",{"2":1}]`,
				expectedErr: ErrorTypeUnmatched{expectedType: `object`, foundType: `[]interface {}`, path: `.-1`},
			},
			{
				jsonpath:    `$.a.d`,
				inputJSON:   `{"a":"b","c":{"d":"e"}}`,
				expectedErr: ErrorTypeUnmatched{expectedType: `object/array`, foundType: `string`, path: `.d`},
			},
			{
				jsonpath:    `$.a.d`,
				inputJSON:   `{"a":123}`,
				expectedErr: ErrorTypeUnmatched{expectedType: `object/array`, foundType: `float64`, path: `.d`},
			},
			{
				jsonpath:    `$.a.d`,
				inputJSON:   `{"a":true}`,
				expectedErr: ErrorTypeUnmatched{expectedType: `object/array`, foundType: `bool`, path: `.d`},
			},
			{
				jsonpath:    `$.a.d`,
				inputJSON:   `{"a":null}`,
				expectedErr: ErrorTypeUnmatched{expectedType: `object/array`, foundType: `null`, path: `.d`},
			},
			{
				jsonpath:    `$.a`,
				inputJSON:   `[1,2]`,
				expectedErr: ErrorTypeUnmatched{expectedType: `object`, foundType: `[]interface {}`, path: `.a`},
			},
			{
				jsonpath:    `$.a`,
				inputJSON:   `[{"a":1}]`,
				expectedErr: ErrorTypeUnmatched{expectedType: `object`, foundType: `[]interface {}`, path: `.a`},
			},
		},
	}

	execTestRetrieveTestGroups(t, testGroups)
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
		`none-matched`: []TestCase{
			{
				jsonpath:    `$..x`,
				inputJSON:   `{"a":"b","c":{"a":"d"},"e":["f",{"g":{"a":"h"}}]}`,
				expectedErr: ErrorNoneMatched{path: `..x`},
			},
			{
				jsonpath:    `$..a.x`,
				inputJSON:   `{"a":"b","c":{"a":"d"},"e":["f",{"g":{"a":"h"}}]}`,
				expectedErr: ErrorNoneMatched{path: `..a.x`},
			},
			{
				// The case where '.x' terminates with an error first
				jsonpath:    `$.x..a`,
				inputJSON:   `{"a":"b","c":{"a":"d"},"e":["f",{"g":{"a":"h"}}]}`,
				expectedErr: ErrorMemberNotExist{path: `.x`},
			},
		},
		`character-type::Non-ASCII-syntax-accepted-in-JSON`: []TestCase{
			{
				jsonpath:     `$..'a'`,
				inputJSON:    `{"'a'":1,"b":{"'a'":2},"c":["'a'",{"d":{"'a'":{"'a'":3}}}]}`,
				expectedJSON: `[1,2,{"'a'":3},3]`,
			},
			{
				jsonpath:     `$.."a"`,
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
	}

	execTestRetrieveTestGroups(t, testGroups)
}

func TestRetrieve_dotNotation_asterisk(t *testing.T) {
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
		},
		`two-asterisks`: []TestCase{
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
		`empty`: []TestCase{
			{
				jsonpath:    `$.*`,
				inputJSON:   `{}`,
				expectedErr: ErrorNoneMatched{path: `.*`},
			},
			{
				jsonpath:    `$.*`,
				inputJSON:   `[]`,
				expectedErr: ErrorNoneMatched{path: `.*`},
			},
		},
		`recursive`: []TestCase{
			{
				jsonpath:    `$..*`,
				inputJSON:   `"a"`,
				expectedErr: ErrorNoneMatched{path: `..*`},
			},
			{
				jsonpath:    `$..*`,
				inputJSON:   `true`,
				expectedErr: ErrorNoneMatched{path: `..*`},
			},
			{
				jsonpath:    `$..*`,
				inputJSON:   `1`,
				expectedErr: ErrorNoneMatched{path: `..*`},
			},
		},
	}

	execTestRetrieveTestGroups(t, testGroups)
}

func TestRetrieve_bracketNotation(t *testing.T) {
	testGroups := TestGroup{
		`bracket-notation`: []TestCase{
			{
				jsonpath:     `$['a']`,
				inputJSON:    `{"a":"b","c":{"d":"e"}}`,
				expectedJSON: `["b"]`,
			},
			{
				jsonpath:    `$['d']`,
				inputJSON:   `{"a":"b","c":{"d":"e"}}`,
				expectedErr: ErrorMemberNotExist{path: `['d']`},
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
			{
				jsonpath:     `$['0']`,
				inputJSON:    `{"0":1,"a":2}`,
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
				jsonpath:    `$['a\c']`,
				inputJSON:   `{"ac":1,"b":2}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['a\c']`},
			},
			{
				jsonpath:    `$["a\c"]`,
				inputJSON:   `{"ac":1,"b":2}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `["a\c"]`},
			},
			{
				jsonpath:     `$['a.b']`,
				inputJSON:    `{"a.b":1,"a":{"b":2}}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$["a"]`,
				inputJSON:    `{"a":1}`,
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
				expectedErr: ErrorMemberNotExist{path: `['*']`},
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
				jsonpath:     `$['\'']`,
				inputJSON:    `{"'":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$["\""]`,
				inputJSON:    `{"\"":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$['\\']`,
				inputJSON:    `{"\\":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$["\\"]`,
				inputJSON:    `{"\\":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:     `$[':@."$,*\'\\']`,
				inputJSON:    `{":@.\"$,*'\\": 1}`,
				expectedJSON: `[1]`,
			},
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
				expectedErr:  ErrorTypeUnmatched{expectedType: `object`, foundType: `[]interface {}`, path: `['']`},
			},
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
				jsonpath:    `$['a','b',0]`,
				inputJSON:   `{"b":2,"a":1,"c":3}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `['a','b',0]`},
			},
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
			{
				jsonpath:    `$['c','d']`,
				inputJSON:   `{"a":1,"b":2}`,
				expectedErr: ErrorNoneMatched{path: `['c','d']`},
			},
			{
				jsonpath:     `$['a','d']`,
				inputJSON:    `{"a":1,"b":2}`,
				expectedJSON: `[1]`,
			},
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
			{
				jsonpath:     `$[0]['a','b']`,
				inputJSON:    `[{"a":1,"b":2},{"a":3,"b":4},{"a":5,"b":6}]`,
				expectedJSON: `[1,2]`,
			},
			{
				jsonpath:     `$[0:2]['b','a']`,
				inputJSON:    `[{"a":1,"b":2},{"a":3,"b":4},{"a":5,"b":6}]`,
				expectedJSON: `[2,1,4,3]`,
			},
			{
				jsonpath:     `$['a'].b`,
				inputJSON:    `{"b":2,"a":{"b":1}}`,
				expectedJSON: `[1]`,
			},
		},
	}

	execTestRetrieveTestGroups(t, testGroups)
}

func TestRetrieve_bracketNotation_asterisk(t *testing.T) {
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
		`empty`: []TestCase{
			{
				jsonpath:    `$[*]`,
				inputJSON:   `[]`,
				expectedErr: ErrorNoneMatched{path: `[*]`},
			},
			{
				jsonpath:    `$[*]`,
				inputJSON:   `{}`,
				expectedErr: ErrorNoneMatched{path: `[*]`},
			},
		},
		`apply-to-value-group`: []TestCase{
			{
				jsonpath:     `$[0:2][*]`,
				inputJSON:    `[[1,2],[3,4],[5,6]]`,
				expectedJSON: `[1,2,3,4]`,
			},
		},
		`child-after-asterisk`: []TestCase{
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
	}

	execTestRetrieveTestGroups(t, testGroups)
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

	execTestRetrieveTestGroups(t, testGroups)
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
				expectedErr: ErrorIndexOutOfRange{path: `[3]`},
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
				expectedErr: ErrorIndexOutOfRange{path: `[-4]`},
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
		`empty`: []TestCase{
			{
				jsonpath:    `$[0]`,
				inputJSON:   `[]`,
				expectedErr: ErrorIndexOutOfRange{path: `[0]`},
			},
			{
				jsonpath:    `$[1]`,
				inputJSON:   `[]`,
				expectedErr: ErrorIndexOutOfRange{path: `[1]`},
			},
			{
				jsonpath:    `$[-1]`,
				inputJSON:   `[]`,
				expectedErr: ErrorIndexOutOfRange{path: `[-1]`},
			},
		},
		`big-number`: []TestCase{
			{
				jsonpath:    `$[1000000000000000000]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorIndexOutOfRange{path: `[1000000000000000000]`},
			},
			{
				jsonpath:    `$[-1000000000000000000]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorIndexOutOfRange{path: `[-1000000000000000000]`},
			},
		},
		`not-array`: []TestCase{
			{
				jsonpath:    `$[0]`,
				inputJSON:   `{"a":1,"b":2}`,
				expectedErr: ErrorTypeUnmatched{expectedType: `array`, foundType: `map[string]interface {}`, path: `[0]`},
			},
			{
				jsonpath:    `$[0]`,
				inputJSON:   `"abc"`,
				expectedErr: ErrorTypeUnmatched{expectedType: `array`, foundType: `string`, path: `[0]`},
			},
			{
				jsonpath:    `$[0]`,
				inputJSON:   `123`,
				expectedErr: ErrorTypeUnmatched{expectedType: `array`, foundType: `float64`, path: `[0]`},
			},
			{
				jsonpath:    `$[0]`,
				inputJSON:   `true`,
				expectedErr: ErrorTypeUnmatched{expectedType: `array`, foundType: `bool`, path: `[0]`},
			},
			{
				jsonpath:    `$[0]`,
				inputJSON:   `null`,
				expectedErr: ErrorTypeUnmatched{expectedType: `array`, foundType: `null`, path: `[0]`},
			},
			{
				jsonpath:    `$[0]`,
				inputJSON:   `{}`,
				expectedErr: ErrorTypeUnmatched{expectedType: `array`, foundType: `map[string]interface {}`, path: `[0]`},
			},
		},
	}

	execTestRetrieveTestGroups(t, testGroups)
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
		`asterisk`: []TestCase{
			{
				jsonpath:     `$[*,0]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second","third","first"]`,
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
				expectedErr: ErrorNoneMatched{path: `[3,3]`},
			},
		},
		`array`: []TestCase{
			{
				jsonpath:     `$[0,1]`,
				inputJSON:    `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
				expectedJSON: `[["11","12","13"],["21","22","23"]]`,
			},
		},
	}

	execTestRetrieveTestGroups(t, testGroups)
}

func TestRetrieve_arraySlice_StartToEnd(t *testing.T) {
	testGroups := TestGroup{
		`start-zero`: []TestCase{
			{
				jsonpath:    `$[0:0]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[0:0]`},
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
				expectedErr: ErrorNoneMatched{path: `[1:1]`},
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
				expectedErr: ErrorNoneMatched{path: `[2:2]`},
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
				expectedErr: ErrorNoneMatched{path: `[2:1]`},
			},
			{
				jsonpath:    `$[2:0]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[2:0]`},
			},
		},
		`start-after-last`: []TestCase{
			{
				jsonpath:    `$[3:2]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[3:2]`},
			},
			{
				jsonpath:    `$[3:3]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[3:3]`},
			},
			{
				jsonpath:    `$[3:4]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[3:4]`},
			},
		},
		`start-minus-to-minus-forward`: []TestCase{
			{
				jsonpath:    `$[-1:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[-1:-1]`},
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
				expectedErr: ErrorNoneMatched{path: `[-1:-2]`},
			},
			{
				jsonpath:    `$[-1:-3]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[-1:-3]`},
			},
		},
		`start-minus-to-plus`: []TestCase{
			{
				jsonpath:    `$[-1:2]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[-1:2]`},
			},
			{
				jsonpath:     `$[-1:3]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third"]`,
			},
			{
				jsonpath:    `$[-2:1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[-2:1]`},
			},
			{
				jsonpath:     `$[-2:2]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second"]`,
			},
			{
				jsonpath:    `$[-3:0]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[-3:0]`},
			},
			{
				jsonpath:     `$[-3:1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first"]`,
			},
			{
				jsonpath:    `$[-4:0]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[-4:0]`},
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
				expectedErr: ErrorNoneMatched{path: `[0:-3]`},
			},
			{
				jsonpath:    `$[0:-4]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[0:-4]`},
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
				expectedErr: ErrorNoneMatched{path: `[1:-2]`},
			},
		},
		`start-last-to-minus`: []TestCase{
			{
				jsonpath:    `$[2:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[2:-1]`},
			},
			{
				jsonpath:    `$[2:-2]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[2:-2]`},
			},
		},
		`omitted-start`: []TestCase{
			{
				jsonpath:    `$[:0]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[:0]`},
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
				expectedErr: ErrorNoneMatched{path: `[:-3]`},
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
				expectedErr: ErrorNoneMatched{path: `[3:]`},
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
				expectedErr: ErrorNoneMatched{path: `[1000000000000000000:1]`},
			},
			{
				jsonpath:    `$[1:-1000000000000000000]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[1:-1000000000000000000]`},
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
				expectedErr: ErrorTypeUnmatched{expectedType: `array`, foundType: `map[string]interface {}`, path: `[1:2]`},
			},
			{
				jsonpath:    `$[:]`,
				inputJSON:   `{"first":1,"second":2,"third":3}`,
				expectedErr: ErrorTypeUnmatched{expectedType: `array`, foundType: `map[string]interface {}`, path: `[:]`},
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
	}

	execTestRetrieveTestGroups(t, testGroups)
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
				jsonpath:     `$[0:2:0]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first","second"]`,
			},
		},
		`minus::start-variation`: []TestCase{
			{
				jsonpath:    `$[0:1:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[0:1:-1]`},
			},
			{
				jsonpath:    `$[1:1:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[1:1:-1]`},
			},
			{
				jsonpath:     `$[2:1:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third"]`,
			},
			{
				jsonpath:    `$[3:1:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[3:1:-1]`},
			},
			{
				jsonpath:    `$[4:1:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[4:1:-1]`},
			},
			{
				jsonpath:     `$[5:1:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third"]`,
			},
			{
				jsonpath:     `$[6:1:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third"]`,
			},
		},
		`minus::end-variation::start-0`: []TestCase{
			{
				jsonpath:     `$[0:-2:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first"]`,
			},
			{
				jsonpath:     `$[0:-1:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["first"]`,
			},
			{
				jsonpath:    `$[0:0:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[0:0:-1]`},
			},
			{
				jsonpath:    `$[0:1:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[0:1:-1]`},
			},
			{
				jsonpath:    `$[0:2:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[0:2:-1]`},
			},
			{
				jsonpath:    `$[0:3:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[0:3:-1]`},
			},
			{
				jsonpath:    `$[0:4:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[0:4:-1]`},
			},
		},
		`minus::end-variation::start-1`: []TestCase{
			{
				jsonpath:     `$[1:-2:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second","first"]`,
			},
			{
				jsonpath:     `$[1:-1:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second","first"]`,
			},
			{
				jsonpath:     `$[1:0:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second"]`,
			},
			{
				jsonpath:    `$[1:1:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[1:1:-1]`},
			},
			{
				jsonpath:     `$[1:2:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second","first"]`,
				expectedErr:  ErrorNoneMatched{path: `[1:2:-1]`},
			},
			{
				jsonpath:     `$[1:3:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["second"]`,
			},
			{
				jsonpath:    `$[1:4:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[1:4:-1]`},
			},
			{
				jsonpath:    `$[1:5:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[1:5:-1]`},
			},
		},
		`minus::end-variation::start-2`: []TestCase{
			{
				jsonpath:     `$[2:-2:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third","second","first"]`,
			},
			{
				jsonpath:     `$[2:-1:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third","second","first"]`,
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
				expectedErr: ErrorNoneMatched{path: `[2:2:-1]`},
			},
			{
				jsonpath:     `$[2:3:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third","second"]`,
			},
			{
				jsonpath:     `$[2:4:-1]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third"]`,
			},
			{
				jsonpath:    `$[2:5:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[2:5:-1]`},
			},
			{
				jsonpath:    `$[2:6:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[2:6:-1]`},
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
				jsonpath:     `$[2:-1:-2]`,
				inputJSON:    `["first","second","third"]`,
				expectedJSON: `["third","first"]`,
			},
			{
				jsonpath:    `$[-1:0:-1]`,
				inputJSON:   `["first","second","third"]`,
				expectedErr: ErrorNoneMatched{path: `[-1:0:-1]`},
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
				expectedErr: ErrorTypeUnmatched{expectedType: `array`, foundType: `map[string]interface {}`, path: `[2:1:-1]`},
			},
			{
				jsonpath:    `$[::-1]`,
				inputJSON:   `{"first":1,"second":2,"third":3}`,
				expectedErr: ErrorTypeUnmatched{expectedType: `array`, foundType: `map[string]interface {}`, path: `[::-1]`},
			},
			{
				jsonpath:    `$[2:1:-1]`,
				inputJSON:   `"value"`,
				expectedErr: ErrorTypeUnmatched{expectedType: `array`, foundType: `string`, path: `[2:1:-1]`},
			},
			{
				jsonpath:    `$[2:1:-1]`,
				inputJSON:   `1`,
				expectedErr: ErrorTypeUnmatched{expectedType: `array`, foundType: `float64`, path: `[2:1:-1]`},
			},
			{
				jsonpath:    `$[2:1:-1]`,
				inputJSON:   `true`,
				expectedErr: ErrorTypeUnmatched{expectedType: `array`, foundType: `bool`, path: `[2:1:-1]`},
			},
			{
				jsonpath:    `$[2:1:-1]`,
				inputJSON:   `null`,
				expectedErr: ErrorTypeUnmatched{expectedType: `array`, foundType: `null`, path: `[2:1:-1]`},
			},
		},
	}

	execTestRetrieveTestGroups(t, testGroups)
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
				expectedErr: ErrorNoneMatched{path: `[?(!@)]`},
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
				expectedErr: ErrorNoneMatched{path: `[?(@.c)]`},
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
				jsonpath:     `$[?(@)]`,
				inputJSON:    `{"a":1}`,
				expectedJSON: `[1]`,
			},
			{
				jsonpath:    `$[?(!@)]`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorNoneMatched{path: `[?(!@)]`},
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
		`asterisk-identifier`: []TestCase{
			{
				jsonpath:    `$.*[?(@.a)]`,
				inputJSON:   `[{"a":1},{"b":2}]`,
				expectedErr: ErrorNoneMatched{path: `.*[?(@.a)]`},
			},
		},
		`root`: []TestCase{
			{
				jsonpath:     `$[?($[0].a)]`,
				inputJSON:    `[{"a":1},{"b":2}]`,
				expectedJSON: `[{"a":1},{"b":2}]`,
			},
			{
				jsonpath:    `$[?(!$[0].a)]`,
				inputJSON:   `[{"a":1},{"b":2}]`,
				expectedErr: ErrorNoneMatched{path: `[?(!$[0].a)]`},
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
				expectedErr: ErrorNoneMatched{path: `[?(@['c','d'])]`},
			},
		},
		`current-asterisk`: []TestCase{
			{
				jsonpath:     `$[?(@.*)]`,
				inputJSON:    `[{"a":1},{"b":2}]`,
				expectedJSON: `[{"a":1},{"b":2}]`,
			},
			{
				jsonpath:    `$[?(@.*)]`,
				inputJSON:   `[1,2]`,
				expectedErr: ErrorNoneMatched{path: `[?(@.*)]`},
			},
		},
		`asterisk-qualifier`: []TestCase{
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
				expectedErr: ErrorNoneMatched{path: `[?(@[*])]`},
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
	}

	execTestRetrieveTestGroups(t, testGroups)
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
				expectedErr: ErrorNoneMatched{path: `[?(@.a=='ab')]`},
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
				expectedErr: ErrorNoneMatched{path: `[?(@.a!='ab')]`},
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
				expectedErr: ErrorNoneMatched{path: `[?(1 > @.a)]`},
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
				expectedErr: ErrorNoneMatched{path: `[?(1.00001 >= @.a)]`},
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
				expectedErr: ErrorNoneMatched{path: `[?(1 < @.a)]`},
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
				expectedErr: ErrorNoneMatched{path: `[?(1.000001 <= @.a)]`},
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
				jsonpath:     `$[?(@.a=='a\b')]`,
				inputJSON:    `[{"a":"ab"}]`,
				expectedJSON: `[{"a":"ab"}]`,
			},
			{
				jsonpath:     `$[?(@.a=="a\b")]`,
				inputJSON:    `[{"a":"ab"}]`,
				expectedJSON: `[{"a":"ab"}]`,
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
				expectedErr: ErrorNoneMatched{path: `[?(@.a==5)]`},
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
				expectedErr: ErrorNoneMatched{path: `[?(@.a==1)]`},
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
				inputJSON:    `[{"a":null},{"a":false},{"a":true},{"a":0},{"a":1},{"a":"false"}]`,
				expectedJSON: `[{"a":false}]`,
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
				inputJSON:    `[{"a":null},{"a":false},{"a":true},{"a":0},{"a":1},{"a":"false"}]`,
				expectedJSON: `[{"a":true}]`,
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
				inputJSON:    `[{"a":null},{"a":false},{"a":true},{"a":0},{"a":1},{"a":"false"}]`,
				expectedJSON: `[{"a":null}]`,
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
				jsonpath:     `$[?(@.a+10==20)]`,
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
				jsonpath:     `$[?(@.a*2==11)]`,
				inputJSON:    `[{"a":6},{"a":5},{"a":5.5},{"a":-5},{"a*2":10.999},{"a*2":11.0},{"a*2":11.1},{"a*2":5},{"a*2":"11"}]`,
				expectedJSON: `[{"a*2":11}]`,
			},
			{
				jsonpath:     `$[?(@.a/10==5)]`,
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
				inputJSON:    `[{"a":0},{"a":1},{"a":2,"b":1}]`,
				expectedJSON: `[{"a":1}]`,
			},
			{
				jsonpath:     `$[?($[2].b == @.a)]`,
				inputJSON:    `[{"a":0},{"a":1},{"a":2,"b":1}]`,
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
				expectedErr: ErrorNoneMatched{path: `[?(@.a == $.b)]`},
			},
			{
				jsonpath:    `$[?($.b == @.a)]`,
				inputJSON:   `[{"a":1},{"a":2}]`,
				expectedErr: ErrorNoneMatched{path: `[?($.b == @.a)]`},
			},
			{
				jsonpath:    `$[?(@.b == $[0].a)]`,
				inputJSON:   `[{"a":1},{"a":2}]`,
				expectedErr: ErrorNoneMatched{path: `[?(@.b == $[0].a)]`},
			},
			{
				jsonpath:    `$[?($[0].a == @.b)]`,
				inputJSON:   `[{"a":1},{"a":2}]`,
				expectedErr: ErrorNoneMatched{path: `[?($[0].a == @.b)]`},
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
				inputJSON:    `{"a":2,"more":[{"a":2},{"b":{"a":2}},{"a":{"a":2}},[{"a":2}]]}`,
				expectedJSON: `[{"a":2},{"a":2},{"a":2},{"a":2}]`,
			},
			{
				jsonpath:     `$..*[?(@.id>2)]`,
				inputJSON:    `[{"complexity":{"one":[{"name":"first","id":1},{"name":"next","id":2},{"name":"another","id":3},{"name":"more","id":4}],"more":{"name":"next to last","id":5}}},{"name":"last","id":6}]`,
				expectedJSON: `[{"id":5,"name":"next to last"},{"id":3,"name":"another"},{"id":4,"name":"more"}]`,
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
				expectedErr: ErrorNoneMatched{path: `[?(10==20)]`},
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
		`value-group-jsonpath::asterisk-qualifier`: []TestCase{
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
		`value-group-jsonpath::asterisk-dot-child-identifier`: []TestCase{
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
				expectedErr: ErrorNoneMatched{path: `[?(@[1][0]>1)][?(@[1][0]>1)][?(@[1]>1)]`},
			},
		},
	}

	execTestRetrieveTestGroups(t, testGroups)
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

	execTestRetrieveTestGroups(t, testGroups)
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
		`value-group-jsonpath::asterisk-qualifier`: []TestCase{
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
		`value-group-jsonpath::asterisk-dot-child-identifier`: []TestCase{
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

	execTestRetrieveTestGroups(t, testGroups)
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
				jsonpath:     `$[?(@.a || @.b != 2)]`,
				inputJSON:    `[{"a":"a"},{"b":2},{"b":3}]`,
				expectedJSON: `[{"a":"a"},{"b":3}]`,
			},
			{
				jsonpath:     `$[?(@.b != 2 || @.a)]`,
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

	execTestRetrieveTestGroups(t, testGroups)
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
				jsonpath:     "\t" + `$.a` + "\t",
				inputJSON:    `{"a":123}`,
				expectedJSON: `[123]`,
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
		},
	}

	execTestRetrieveTestGroups(t, testGroups)
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
				jsonpath:    `@`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 0, reason: `the use of '@' at the beginning is prohibited`, near: `@`},
			},
			{
				jsonpath:    `$$`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `$`},
			},
		},
		`dot-child-identifier`: []TestCase{
			{
				jsonpath:    `.c`,
				inputJSON:   `{"a":"b","c":{"d":"e"}}`,
				expectedErr: ErrorInvalidSyntax{position: 0, reason: `unrecognized input`, near: `.c`},
			},
			{
				jsonpath:    `$a`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `a`},
			},
			{
				jsonpath:    `$.`,
				inputJSON:   `{"a":1}`,
				expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `.`},
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
			{
				jsonpath:    `$.func(`,
				inputJSON:   `{}`,
				expectedErr: ErrorInvalidSyntax{position: 6, reason: `unrecognized input`, near: `(`},
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
				expectedErr:  ErrorInvalidSyntax{position: 4, reason: `the omission of '$' allowed only at the beginning`, near: `false)]`},
			},
			{
				jsonpath:     `$[?(true)]`,
				inputJSON:    `[0,1,false,true,null,{},[]]`,
				expectedJSON: `[]`,
				expectedErr:  ErrorInvalidSyntax{position: 4, reason: `the omission of '$' allowed only at the beginning`, near: `true)]`},
			},
			{
				jsonpath:     `$[?(null)]`,
				inputJSON:    `[0,1,false,true,null,{},[]]`,
				expectedJSON: `[]`,
				expectedErr:  ErrorInvalidSyntax{position: 4, reason: `the omission of '$' allowed only at the beginning`, near: `null)]`},
			},
			{
				jsonpath:    `$[?(@.a==["b"])]`,
				inputJSON:   `[{"a":["b"]}]`,
				expectedErr: ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `["b"])]`},
			},
			{
				jsonpath:    `$[?(@[0:1]==[1])]`,
				inputJSON:   `[[1,2,3],[1],[2,3],1,2]`,
				expectedErr: ErrorInvalidSyntax{position: 12, reason: `the omission of '$' allowed only at the beginning`, near: `[1])]`},
			},
			{
				jsonpath:    `$[?(@.*==[1,2])]`,
				inputJSON:   `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
				expectedErr: ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `[1,2])]`},
			},
			{
				jsonpath:    `$[?(@.*==['1','2'])]`,
				inputJSON:   `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
				expectedErr: ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `['1','2'])]`},
			},
			{
				jsonpath:    `$[?(@=={"k":"v"})]`,
				inputJSON:   `{}`,
				expectedErr: ErrorInvalidSyntax{position: 7, reason: `the omission of '$' allowed only at the beginning`, near: `{"k":"v"})]`},
			},
			{
				jsonpath:     `$[?(@.a==fAlse)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `fAlse)]`},
			},
			{
				jsonpath:     `$[?(@.a==faLse)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `faLse)]`},
			},
			{
				jsonpath:     `$[?(@.a==falSe)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `falSe)]`},
			},
			{
				jsonpath:     `$[?(@.a==falsE)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `falsE)]`},
			},
			{
				jsonpath:     `$[?(@.a==FaLse)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `FaLse)]`},
			},
			{
				jsonpath:     `$[?(@.a==FalSe)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `FalSe)]`},
			},
			{
				jsonpath:     `$[?(@.a==FalsE)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `FalsE)]`},
			},
			{
				jsonpath:     `$[?(@.a==FaLSE)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `FaLSE)]`},
			},
			{
				jsonpath:     `$[?(@.a==FAlSE)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `FAlSE)]`},
			},
			{
				jsonpath:     `$[?(@.a==FALsE)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `FALsE)]`},
			},
			{
				jsonpath:     `$[?(@.a==FALSe)]`,
				inputJSON:    `[{"a":false}]`,
				expectedJSON: `[{"a":false}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `FALSe)]`},
			},
			{
				jsonpath:     `$[?(@.a==tRue)]`,
				inputJSON:    `[{"a":true}]`,
				expectedJSON: `[{"a":true}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `tRue)]`},
			},
			{
				jsonpath:     `$[?(@.a==trUe)]`,
				inputJSON:    `[{"a":true}]`,
				expectedJSON: `[{"a":true}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `trUe)]`},
			},
			{
				jsonpath:     `$[?(@.a==truE)]`,
				inputJSON:    `[{"a":true}]`,
				expectedJSON: `[{"a":true}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `truE)]`},
			},
			{
				jsonpath:     `$[?(@.a==TrUe)]`,
				inputJSON:    `[{"a":true}]`,
				expectedJSON: `[{"a":true}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `TrUe)]`},
			},
			{
				jsonpath:     `$[?(@.a==TruE)]`,
				inputJSON:    `[{"a":true}]`,
				expectedJSON: `[{"a":true}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `TruE)]`},
			},
			{
				jsonpath:     `$[?(@.a==TrUE)]`,
				inputJSON:    `[{"a":true}]`,
				expectedJSON: `[{"a":true}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `TrUE)]`},
			},
			{
				jsonpath:     `$[?(@.a==TRuE)]`,
				inputJSON:    `[{"a":true}]`,
				expectedJSON: `[{"a":true}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `TRuE)]`},
			},
			{
				jsonpath:     `$[?(@.a==TRUe)]`,
				inputJSON:    `[{"a":true}]`,
				expectedJSON: `[{"a":true}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `TRUe)]`},
			},
			{
				jsonpath:     `$[?(@.a==nUll)]`,
				inputJSON:    `[{"a":null}]`,
				expectedJSON: `[{"a":null}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `nUll)]`},
			},
			{
				jsonpath:     `$[?(@.a==nuLl)]`,
				inputJSON:    `[{"a":null}]`,
				expectedJSON: `[{"a":null}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `nuLl)]`},
			},
			{
				jsonpath:     `$[?(@.a==nulL)]`,
				inputJSON:    `[{"a":null}]`,
				expectedJSON: `[{"a":null}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `nulL)]`},
			},
			{
				jsonpath:     `$[?(@.a==NuLl)]`,
				inputJSON:    `[{"a":null}]`,
				expectedJSON: `[{"a":null}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `NuLl)]`},
			},
			{
				jsonpath:     `$[?(@.a==NulL)]`,
				inputJSON:    `[{"a":null}]`,
				expectedJSON: `[{"a":null}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `NulL)]`},
			},
			{
				jsonpath:     `$[?(@.a==NuLL)]`,
				inputJSON:    `[{"a":null}]`,
				expectedJSON: `[{"a":null}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `NuLL)]`},
			},
			{
				jsonpath:     `$[?(@.a==NUlL)]`,
				inputJSON:    `[{"a":null}]`,
				expectedJSON: `[{"a":null}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `NUlL)]`},
			},
			{
				jsonpath:     `$[?(@.a==NULl)]`,
				inputJSON:    `[{"a":null}]`,
				expectedJSON: `[{"a":null}]`,
				expectedErr:  ErrorInvalidSyntax{position: 9, reason: `the omission of '$' allowed only at the beginning`, near: `NULl)]`},
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
				expectedErr: ErrorInvalidSyntax{position: 13, reason: `the omission of '$' allowed only at the beginning`, near: `false)]`},
			},
			{
				jsonpath:    `$[?(@.a>1 && true)]`,
				inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
				expectedErr: ErrorInvalidSyntax{position: 13, reason: `the omission of '$' allowed only at the beginning`, near: `true)]`},
			},
			{
				jsonpath:    `$[?(@.a>1 || false)]`,
				inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
				expectedErr: ErrorInvalidSyntax{position: 13, reason: `the omission of '$' allowed only at the beginning`, near: `false)]`},
			},
			{
				jsonpath:    `$[?(@.a>1 || true)]`,
				inputJSON:   `[{"a":1},{"a":2},{"a":3}]`,
				expectedErr: ErrorInvalidSyntax{position: 13, reason: `the omission of '$' allowed only at the beginning`, near: `true)]`},
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
				expectedErr: ErrorInvalidSyntax{position: 4, reason: `the omission of '$' allowed only at the beginning`, near: `a=~/123/)]`},
			},
		},
	}

	execTestRetrieveTestGroups(t, testGroups)
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

	execTestRetrieveTestGroups(t, testGroups)
}

func TestRetrieve_notSupported(t *testing.T) {
	testGroups := TestGroup{
		`Not supported`: []TestCase{
			{
				jsonpath:    `$[(command)]`,
				inputJSON:   `{}`,
				expectedErr: ErrorNotSupported{feature: `script`, path: `[(command)]`},
			},
		},
	}

	execTestRetrieveTestGroups(t, testGroups)
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
				expectedErr:   ErrorNoneMatched{path: `[?(@.a > 123.46)].a`},
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

	execTestRetrieveTestGroups(t, testGroups)
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

				expectedErr: ErrorFunctionFailed{function: `.errFilter()`, err: fmt.Errorf(`filter error`)},
			},
			{
				jsonpath:  `$.*.errFilter()`,
				inputJSON: `[122.345,123.45,123.456]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter`: errFilterFunc,
				},

				expectedErr: ErrorNoneMatched{path: `.*.errFilter()`},
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
				expectedErr: ErrorFunctionFailed{function: `.errFilter()`, err: fmt.Errorf(`filter error`)},
			},
			{
				jsonpath:  `$.*.twice().errFilter()`,
				inputJSON: `[122.345,123.45,123.456]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter`: errFilterFunc,
					`twice`:     twiceFunc,
				},

				expectedErr: ErrorNoneMatched{path: `.*.twice().errFilter()`},
			}, {
				jsonpath:  `$.errFilter().twice()`,
				inputJSON: `[122.345,123.45,123.456]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter`: errFilterFunc,
					`twice`:     twiceFunc,
				},

				expectedErr: ErrorFunctionFailed{function: `.errFilter()`, err: fmt.Errorf(`filter error`)},
			},
			{
				jsonpath:  `$.*.errFilter().twice()`,
				inputJSON: `[122.345,123.45,123.456]`,
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter`: errFilterFunc,
					`twice`:     twiceFunc,
				},

				expectedErr: ErrorNoneMatched{path: `.*.errFilter().twice()`},
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
				expectedErr: ErrorFunctionFailed{function: `.errFilter()`, err: fmt.Errorf(`filter error`)},
			},
		},
		`aggregate-error`: []TestCase{
			{
				jsonpath:  `$.*.errAggregate()`,
				inputJSON: `[122.345,123.45,123.456]`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`errAggregate`: errAggregateFunc,
				},
				expectedErr: ErrorFunctionFailed{function: `.errAggregate()`, err: fmt.Errorf(`aggregate error`)},
			},
			{
				jsonpath:  `$.*.max().errAggregate()`,
				inputJSON: `[122.345,123.45,123.456]`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`errAggregate`: errAggregateFunc,
					`max`:          maxFunc,
				},
				expectedErr: ErrorFunctionFailed{function: `.errAggregate()`, err: fmt.Errorf(`aggregate error`)},
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
				expectedErr: ErrorFunctionFailed{function: `.errAggregate()`, err: fmt.Errorf(`aggregate error`)},
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
				expectedErr: ErrorFunctionFailed{function: `.errAggregate()`, err: fmt.Errorf(`aggregate error`)},
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
				expectedErr: ErrorFunctionFailed{function: `.errAggregate()`, err: fmt.Errorf(`aggregate error`)},
			},
			{
				jsonpath:  `$.a.max()`,
				inputJSON: `{}`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxFunc,
				},
				expectedErr: ErrorMemberNotExist{path: `.a`},
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

	execTestRetrieveTestGroups(t, testGroups)
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
		`convert-srcJSON`: []TestCase{
			{
				jsonpath:     `$[1]`,
				inputJSON:    `[1,2,3]`,
				accessorMode: true,
				resultValidator: func(src interface{}, actualObject []interface{}) error {
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
						return fmt.Errorf(`Delx2 -> Get : expect<%d> != actual<%d>`, 5, accessor.Get())
					}
					accessor.Set(6) // srcArray:[1] , accessor:[1,6,3]
					if len(srcArray) != 1 {
						return fmt.Errorf(`Delx2 -> Set -> Len : expect<%d> != actual<%d>`, 1, len(srcArray))
					}
					srcArray = append(srcArray, 7) // srcArray:[1,7] , accessor:[1,7,3]
					if len(srcArray) != 2 || accessor.Get() != 7 {
						return fmt.Errorf(`Delx2 -> Add -> Get : expect<%d> != actual<%d>`, 7, accessor.Get())
					}
					srcArray = append(srcArray, 8) // srcArray:[1,7,8]    , accessor:[1,7,8]
					srcArray = append(srcArray, 9) // srcArray:[1,7,8,9]  , accessor:[1,7,8,9]
					srcArray[1] = 10               // srcArray:[1,10,8,9] , accessor:[1,10,8,9]
					if len(srcArray) != 4 || accessor.Get() != 10 {
						return fmt.Errorf(`Delx2 -> Addx3 -> Update -> Get : expect<%d> != actual<%d>`, 10, accessor.Get())
					}
					return nil
				},
			},
			{
				jsonpath:     `$.a`,
				inputJSON:    `{"a":1}`,
				accessorMode: true,
				resultValidator: func(src interface{}, actualObject []interface{}) error {
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
				},
			},
		},
	}

	execTestRetrieveTestGroups(t, testGroups)
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
		t.Errorf("%w", err)
		return
	}
	var src2 interface{}
	if err := json.Unmarshal([]byte(srcJSON2), &src2); err != nil {
		t.Errorf("%w", err)
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
		t.Errorf("%w", err)
		return
	}
	actualOutputJSON2, err := json.Marshal(actualObject2)
	if err != nil {
		t.Errorf("%w", err)
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
		t.Errorf("%w", err)
		return
	}
	var src2 interface{}
	if err := json.Unmarshal([]byte(srcJSON2), &src2); err != nil {
		t.Errorf("%w", err)
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
		t.Errorf("%w", err)
		return
	}
	actualOutputJSON2, err := json.Marshal(actualObject2)
	if err != nil {
		t.Errorf("%w", err)
		return
	}

	if string(actualOutputJSON1) != string(expectedOutput1) || string(actualOutputJSON2) != string(expectedOutput2) {
		t.Errorf("actualOutputJSON1<%s> != expectedOutput1<%s> || actualOutputJSON2<%s> != expectedOutput2<%s>\n",
			string(actualOutputJSON1), string(expectedOutput1), string(actualOutputJSON2), string(expectedOutput2))
		return
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

	err := parseError{p: &parser}
	_ = err.Error()

	parser.Parse(1)

	Pretty(true)(&parser)
	parser.PrintSyntaxTree()
	Size(10)(&parser)

	parser.Init(func(p *pegJSONPathParser) error {
		return fmt.Errorf(`test error`)
	})

	os.Stdout = stdoutBackup
}
