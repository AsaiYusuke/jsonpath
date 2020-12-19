package jsonpath

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

type TestGroup struct {
	name      string
	testCases []TestCase
}

type TestCase struct {
	jsonpath     string
	inputJSON    string
	expectedJSON string
	expectedErr  error
	filters      map[string]func(interface{}) (interface{}, error)
	aggregates   map[string]func([]interface{}) (interface{}, error)
}

func execTestRetrieve(t *testing.T, inputJSON interface{}, testCase TestCase) {
	jsonPath := testCase.jsonpath
	expectedOutputJSON := testCase.expectedJSON
	expectedError := testCase.expectedErr
	actualObject, err := Retrieve(jsonPath, inputJSON)
	if err != nil {
		if reflect.TypeOf(expectedError) == reflect.TypeOf(err) &&
			fmt.Sprintf(`%s`, expectedError) == fmt.Sprintf(`%s`, err) {
			return
		}
		t.Errorf("expected error<%s> != actual error<%s>\n",
			expectedError, err)
		return
	}
	if expectedError != nil {
		t.Errorf("expected error<%w> != actual error<none>\n", expectedError)
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
}

func TestRetrieve(t *testing.T) {
	testGroups := []TestGroup{
		{
			name: `dot-notation`,
			testCases: []TestCase{
				{
					jsonpath:     `$`,
					inputJSON:    `{"a":"b","c":{"d":"e"}}`,
					expectedJSON: `[{"a":"b","c":{"d":"e"}}]`,
				},
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
					jsonpath:     `a`,
					inputJSON:    `{"a":"b","c":{"d":"e"}}`,
					expectedJSON: `["b"]`,
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
					jsonpath:     `$.length`,
					inputJSON:    `["length",1,2]`,
					expectedJSON: ``,
					expectedErr:  ErrorTypeUnmatched{`object`, `[]interface {}`, `.length`},
				},
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
					jsonpath:     `$`,
					inputJSON:    `{"$":1}`,
					expectedJSON: `[{"$":1}]`,
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
					jsonpath:     `$.\` + "\t",
					inputJSON:    `{"":123}`,
					expectedJSON: ``,
					expectedErr:  ErrorMemberNotExist{`.\` + "\t"},
				},
				{
					jsonpath:     `$.\` + "\r",
					inputJSON:    `{"":123}`,
					expectedJSON: ``,
					expectedErr:  ErrorMemberNotExist{`.\` + "\r"},
				},
				{
					jsonpath:     `$.\` + "\n",
					inputJSON:    `{"":123}`,
					expectedJSON: ``,
					expectedErr:  ErrorMemberNotExist{`.\` + "\n"},
				},
				{
					jsonpath:     `$.\a`,
					inputJSON:    `{"a":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `.\a`},
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
					jsonpath:     `$.a\` + "\t" + `b`,
					inputJSON:    `{"ab":123}`,
					expectedJSON: ``,
					expectedErr:  ErrorMemberNotExist{`.a\` + "\t" + `b`},
				},
				{
					jsonpath:     `$.a\` + "\r" + `b`,
					inputJSON:    `{"ab":123}`,
					expectedJSON: ``,
					expectedErr:  ErrorMemberNotExist{`.a\` + "\r" + `b`},
				},
				{
					jsonpath:     `$.a\` + "\n" + `b`,
					inputJSON:    `{"ab":123}`,
					expectedJSON: ``,
					expectedErr:  ErrorMemberNotExist{`.a\` + "\n" + `b`},
				},
				{
					jsonpath:     `$.a\a`,
					inputJSON:    `{"aa":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{3, `unrecognized input`, `\a`},
				},
				{
					jsonpath:     `$.\`,
					inputJSON:    `{"\\":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `.\`},
				},
				{
					jsonpath:     `$.(`,
					inputJSON:    `{"(":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `.(`},
				},
				{
					jsonpath:     `$.)`,
					inputJSON:    `{")":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `.)`},
				},
				{
					jsonpath:     `$.=`,
					inputJSON:    `{"=":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `.=`},
				},
				{
					jsonpath:     `$.!`,
					inputJSON:    `{"!":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `.!`},
				},
				{
					jsonpath:     `$.>`,
					inputJSON:    `{">":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `.>`},
				},
				{
					jsonpath:     `$.<`,
					inputJSON:    `{"<":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `.<`},
				},
				{
					jsonpath:     `$. `,
					inputJSON:    `{" ":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `. `},
				},
				{
					jsonpath:     `$.` + "\t",
					inputJSON:    `{"":123}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `.` + "\t"},
				},
				{
					jsonpath:     `$.` + "\r",
					inputJSON:    `{"":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `.` + "\r"},
				},
				{
					jsonpath:     `$.` + "\n",
					inputJSON:    `{"":123}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `.` + "\n"},
				},
				{
					jsonpath:     `$.a\b`,
					inputJSON:    `{"a\\b":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{3, `unrecognized input`, `\b`},
				},
				{
					jsonpath:     `$.a(b`,
					inputJSON:    `{"(":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{3, `unrecognized input`, `(b`},
				},
				{
					jsonpath:     `$.a)b`,
					inputJSON:    `{")":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{3, `unrecognized input`, `)b`},
				},
				{
					jsonpath:     `$.a=b`,
					inputJSON:    `{"=":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{3, `unrecognized input`, `=b`},
				},
				{
					jsonpath:     `$.a!b`,
					inputJSON:    `{"!":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{3, `unrecognized input`, `!b`},
				},
				{
					jsonpath:     `$.a>b`,
					inputJSON:    `{">":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{3, `unrecognized input`, `>b`},
				},
				{
					jsonpath:     `$.a<b`,
					inputJSON:    `{"<":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{3, `unrecognized input`, `<b`},
				},
				{
					jsonpath:     `$.a b`,
					inputJSON:    `{" ":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `unrecognized input`, `b`},
				},
				{
					jsonpath:     `$.a` + "\t" + `b`,
					inputJSON:    `{"":123}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `unrecognized input`, `b`},
				},
				{
					jsonpath:     `$.a` + "\r" + `b`,
					inputJSON:    `{"":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{3, `unrecognized input`, "\r" + `b`},
				},
				{
					jsonpath:     `$.a` + "\n" + `b`,
					inputJSON:    `{"":123}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{3, `unrecognized input`, "\n" + `b`},
				},
				{
					jsonpath:     `$.ﾃｽﾄソポァゼゾタダＡボマミ①`,
					inputJSON:    `{"ﾃｽﾄソポァゼゾタダＡボマミ①":1}`,
					expectedJSON: `[1]`,
				},
				{
					jsonpath:     `$.d`,
					inputJSON:    `{"a":"b","c":{"d":"e"}}`,
					expectedJSON: ``,
					expectedErr:  ErrorMemberNotExist{`.d`},
				},
				{
					jsonpath:     `$.2`,
					inputJSON:    `{"a":1,"2":2,"3":{"2":1}}`,
					expectedJSON: `[2]`,
				},
				{
					jsonpath:     `$.2`,
					inputJSON:    `["a","b",{"2":1}]`,
					expectedJSON: ``,
					expectedErr:  ErrorTypeUnmatched{`object`, `[]interface {}`, `.2`},
				},
				{
					jsonpath:     `$.-1`,
					inputJSON:    `["a","b",{"2":1}]`,
					expectedJSON: ``,
					expectedErr:  ErrorTypeUnmatched{`object`, `[]interface {}`, `.-1`},
				},
				{
					jsonpath:     `$.a.d`,
					inputJSON:    `{"a":"b","c":{"d":"e"}}`,
					expectedJSON: ``,
					expectedErr:  ErrorTypeUnmatched{`object/array`, `string`, `.d`},
				},
				{
					jsonpath:     `$.a.d`,
					inputJSON:    `{"a":123}`,
					expectedJSON: ``,
					expectedErr:  ErrorTypeUnmatched{`object/array`, `float64`, `.d`},
				},
				{
					jsonpath:     `$.a.d`,
					inputJSON:    `{"a":true}`,
					expectedJSON: ``,
					expectedErr:  ErrorTypeUnmatched{`object/array`, `bool`, `.d`},
				},
				{
					jsonpath:     `$.a.d`,
					inputJSON:    `{"a":null}`,
					expectedJSON: ``,
					expectedErr:  ErrorTypeUnmatched{`object/array`, `null`, `.d`},
				},
				{
					jsonpath:     `$.a`,
					inputJSON:    `[1,2]`,
					expectedJSON: ``,
					expectedErr:  ErrorTypeUnmatched{`object`, `[]interface {}`, `.a`},
				},
				{
					jsonpath:     `$.a`,
					inputJSON:    `[{"a":1}]`,
					expectedJSON: ``,
					expectedErr:  ErrorTypeUnmatched{`object`, `[]interface {}`, `.a`},
				},
			},
		},
		{
			name: `dot-notation-recursive-descent`,
			testCases: []TestCase{
				{
					jsonpath:     `$.a..b`,
					inputJSON:    `{"a":{"b":1,"c":{"b":2},"d":["b",{"a":3,"b":4}]},"b":5}`,
					expectedJSON: `[1,2,4]`,
				},
				{
					jsonpath:     `$..a`,
					inputJSON:    `{"a":"b","c":{"a":"d"},"e":["a",{"a":{"a":"h"}}]}`,
					expectedJSON: `["b","d",{"a":"h"},"h"]`,
				},
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
					jsonpath:     `$..x`,
					inputJSON:    `{"a":"b","c":{"a":"d"},"e":["f",{"g":{"a":"h"}}]}`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`..x`},
				},
				{
					jsonpath:     `$..a.x`,
					inputJSON:    `{"a":"b","c":{"a":"d"},"e":["f",{"g":{"a":"h"}}]}`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`..a.x`},
				},
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
				{
					jsonpath:     `$..[?(@.a)]`,
					inputJSON:    `{"a":1,"b":[{"a":2},{"b":{"a":3}},{"a":{"a":4}}]}`,
					expectedJSON: `[{"a":2},{"a":{"a":4}},{"a":3},{"a":4}]`,
				},
				{
					jsonpath:     `$..['a','b']`,
					inputJSON:    `[{"a":1,"b":2,"c":{"a":3}},{"a":4},{"b":5},{"a":6,"b":7},{"d":{"b":8}}]`,
					expectedJSON: `[1,2,3,4,5,6,7,8]`,
				},
			},
		},
		{
			name: `dot-notation-asterisk`,
			testCases: []TestCase{
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
					jsonpath:     `$.*.*`,
					inputJSON:    `[[1,2,3],[4,5,6]]`,
					expectedJSON: `[1,2,3,4,5,6]`,
				},
				{
					jsonpath:     `$.*.a.*`,
					inputJSON:    `[{"a":[1]}]`,
					expectedJSON: `[1]`,
				},
				{
					jsonpath:     `$..[*]`,
					inputJSON:    `{"a":1,"b":[2,3],"c":{"d":"e","f":[4,5]}}`,
					expectedJSON: `[1,[2,3],{"d":"e","f":[4,5]},2,3,"e",[4,5],4,5]`,
				},
				{
					jsonpath:     `$.*`,
					inputJSON:    `{}`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`.*`},
				},
				{
					jsonpath:     `$.*`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`.*`},
				},
				{
					jsonpath:     `$..*`,
					inputJSON:    `"a"`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`..*`},
				},
				{
					jsonpath:     `$..*`,
					inputJSON:    `true`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`..*`},
				},
				{
					jsonpath:     `$..*`,
					inputJSON:    `1`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`..*`},
				},
				{
					jsonpath:     `$.*['a','b']`,
					inputJSON:    `[{"a":1,"b":2,"c":3},{"a":4,"b":5,"d":6}]`,
					expectedJSON: `[1,2,4,5]`,
				},
			},
		},
		{
			name: `bracket-notation`,
			testCases: []TestCase{
				{
					jsonpath:     `$['a']`,
					inputJSON:    `{"a":"b","c":{"d":"e"}}`,
					expectedJSON: `["b"]`,
				},
				{
					jsonpath:     `$['d']`,
					inputJSON:    `{"a":"b","c":{"d":"e"}}`,
					expectedJSON: ``,
					expectedErr:  ErrorMemberNotExist{`['d']`},
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
					jsonpath:     `$['a\c']`,
					inputJSON:    `{"ac":1,"b":2}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `['a\c']`},
				},
				{
					jsonpath:     `$["a\c"]`,
					inputJSON:    `{"ac":1,"b":2}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `["a\c"]`},
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
					jsonpath:     `$['*']`,
					inputJSON:    `{"a":1,"b":2}`,
					expectedJSON: ``,
					expectedErr:  ErrorMemberNotExist{`['*']`},
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
					expectedErr:  ErrorTypeUnmatched{`object`, `[]interface {}`, `['']`},
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
					jsonpath:     `$['a','b',0]`,
					inputJSON:    `{"b":2,"a":1,"c":3}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `['a','b',0]`},
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
					jsonpath:     `$['c','d']`,
					inputJSON:    `{"a":1,"b":2}`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`['c','d']`},
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
		},
		{
			name: `bracket-notation-asterisk`,
			testCases: []TestCase{
				{
					jsonpath:     `$[*]`,
					inputJSON:    `["a",123,true,{"b":"c"},[0,1],null]`,
					expectedJSON: `["a",123,true,{"b":"c"},[0,1],null]`,
				},
				{
					jsonpath:     `$[*]`,
					inputJSON:    `{"a":[1],"b":[2,3]}`,
					expectedJSON: `[[1],[2,3]]`,
				},
				{
					jsonpath:     `$[*]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[*]`},
				},
				{
					jsonpath:     `$[*]`,
					inputJSON:    `{}`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[*]`},
				},
				{
					jsonpath:     `$[0:2][*]`,
					inputJSON:    `[[1,2],[3,4],[5,6]]`,
					expectedJSON: `[1,2,3,4]`,
				},
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
		},
		{
			name: `Value type`,
			testCases: []TestCase{
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
		},
		{
			name: `Array-index`,
			testCases: []TestCase{
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
					jsonpath:     `$[3]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorIndexOutOfRange{`[3]`},
				},
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
					jsonpath:     `$[-4]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorIndexOutOfRange{`[-4]`},
				},
				{
					jsonpath:     `$[0][1]`,
					inputJSON:    `[["a","b"],["c"],["d"]]`,
					expectedJSON: `["b"]`,
				},
				{
					jsonpath:     `$[0]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorIndexOutOfRange{`[0]`},
				},
				{
					jsonpath:     `$[1]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorIndexOutOfRange{`[1]`},
				},
				{
					jsonpath:     `$[-1]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorIndexOutOfRange{`[-1]`},
				},
				{
					jsonpath:     `$[1000000000000000000]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorIndexOutOfRange{`[1000000000000000000]`},
				},
				{
					jsonpath:     `$[0]`,
					inputJSON:    `{"a":1,"b":2}`,
					expectedJSON: ``,
					expectedErr:  ErrorTypeUnmatched{`array`, `map[string]interface {}`, `[0]`},
				},
				{
					jsonpath:     `$[0]`,
					inputJSON:    `"abc"`,
					expectedJSON: ``,
					expectedErr:  ErrorTypeUnmatched{`array`, `string`, `[0]`},
				},
				{
					jsonpath:     `$[0]`,
					inputJSON:    `123`,
					expectedJSON: ``,
					expectedErr:  ErrorTypeUnmatched{`array`, `float64`, `[0]`},
				},
				{
					jsonpath:     `$[0]`,
					inputJSON:    `true`,
					expectedJSON: ``,
					expectedErr:  ErrorTypeUnmatched{`array`, `bool`, `[0]`},
				},
				{
					jsonpath:     `$[0]`,
					inputJSON:    `null`,
					expectedJSON: ``,
					expectedErr:  ErrorTypeUnmatched{`array`, `null`, `[0]`},
				},
				{
					jsonpath:     `$[0]`,
					inputJSON:    `{}`,
					expectedJSON: ``,
					expectedErr:  ErrorTypeUnmatched{`array`, `map[string]interface {}`, `[0]`},
				},
			},
		},
		{
			name: `Array-union`,
			testCases: []TestCase{
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
					jsonpath:     `$[2,0,1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["third","first","second"]`,
				},
				{
					jsonpath:     `$[0,3]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["first"]`,
				},
				{
					jsonpath:     `$[0,-1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["first","third"]`,
				},
				{
					jsonpath:     `$[0,1]`,
					inputJSON:    `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
					expectedJSON: `[["11","12","13"],["21","22","23"]]`,
				},
				{
					jsonpath:     `$[*]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["first","second","third"]`,
				},
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
					jsonpath:     `$[1:2,0]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["second","first"]`,
				},
				{
					jsonpath:     `$[:2,0]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["first","second","first"]`,
				},
			},
		},
		{
			name: `Array-slice-start-to-end`,
			testCases: []TestCase{
				{
					jsonpath:     `$[0:0]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[0:0]`},
				},
				{
					jsonpath:     `$[0:3]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["first","second","third"]`,
				},
				{
					jsonpath:     `$[0:2]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["first","second"]`,
				},
				{
					jsonpath:     `$[1:1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[1:1]`},
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
				{
					jsonpath:     `$[2:1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[2:1]`},
				},
				{
					jsonpath:     `$[3:2]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[3:2]`},
				},
				{
					jsonpath:     `$[3:3]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[3:3]`},
				},
				{
					jsonpath:     `$[3:4]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[3:4]`},
				},
				{
					jsonpath:     `$[-1:-1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[-1:-1]`},
				},
				{
					jsonpath:     `$[-2:-1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["second"]`,
				},
				{
					jsonpath:     `$[-1:-2]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[-1:-2]`},
				},
				{
					jsonpath:     `$[-1:3]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["third"]`,
				},
				{
					jsonpath:     `$[-1:2]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[-1:2]`},
				},
				{
					jsonpath:     `$[-4:3]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["first","second","third"]`,
				},
				{
					jsonpath:     `$[0:-1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["first","second"]`,
				},
				{
					jsonpath:     `$[0:-3]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[0:-3]`},
				},
				{
					jsonpath:     `$[0:-4]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[0:-4]`},
				},
				{
					jsonpath:     `$[1:-2]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[1:-2]`},
				},
				{
					jsonpath:     `$[1:-1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["second"]`,
				},
				{
					jsonpath:     `$[:2]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["first","second"]`,
				},
				{
					jsonpath:     `$[1:]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["second","third"]`,
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
					jsonpath:     `$[-4:]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["first","second","third"]`,
				},
				{
					jsonpath:     `$[:]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["first","second","third"]`,
				},
				{
					jsonpath:     `$[-1000000000000000000:1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["first"]`,
				},
				{
					jsonpath:     `$[1000000000000000000:1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[1000000000000000000:1]`},
				},
				{
					jsonpath:     `$[1:1000000000000000000]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["second","third"]`,
				},
				{
					jsonpath:     `$[1:2]`,
					inputJSON:    `{"first":1,"second":2,"third":3}`,
					expectedJSON: ``,
					expectedErr:  ErrorTypeUnmatched{`array`, `map[string]interface {}`, `[1:2]`},
				},
				{
					jsonpath:     `$[:]`,
					inputJSON:    `{"first":1,"second":2,"third":3}`,
					expectedJSON: ``,
					expectedErr:  ErrorTypeUnmatched{`array`, `map[string]interface {}`, `[:]`},
				},
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
			},
		},
		{
			name: `Array-slice-step`,
			testCases: []TestCase{
				{
					jsonpath:     `$[0:2:1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["first","second"]`,
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
					jsonpath:     `$[0:2:0]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["first","second"]`,
				},
				{
					jsonpath:     `$[0:3:-1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[0:3:-1]`},
				},
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
					jsonpath:     `$[2:-1:-2]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["third","first"]`,
				},
				{
					jsonpath:     `$[3:1:-1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[3:1:-1]`},
				},
				{
					jsonpath:     `$[4:1:-1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[4:1:-1]`},
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
				{
					jsonpath:     `$[2:2:-1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["third","second","first"]`,
				},
				{
					jsonpath:     `$[2:3:-1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["third","second"]`,
				},
				{
					jsonpath:     `$[2:5:-1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[2:5:-1]`},
				},
				{
					jsonpath:     `$[2:6:-1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[2:6:-1]`},
				},
				{
					jsonpath:     `$[2:7:-1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[2:7:-1]`},
				},
				{
					jsonpath:     `$[-1:0:-1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[-1:0:-1]`},
				},
				{
					jsonpath:     `$[2:-1:-1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["third","second","first"]`,
				},
				{
					jsonpath:     `$[0:3:]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["first","second","third"]`,
				},
				{
					jsonpath:     `$[::]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["first","second","third"]`,
				},
				{
					jsonpath:     `$[1::-1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: `["second","first"]`,
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
					jsonpath:     `$[2:1:-1]`,
					inputJSON:    `{"first":1,"second":2,"third":3}`,
					expectedJSON: ``,
					expectedErr:  ErrorTypeUnmatched{`array`, `map[string]interface {}`, `[2:1:-1]`},
				},
				{
					jsonpath:     `$[::-1]`,
					inputJSON:    `{"first":1,"second":2,"third":3}`,
					expectedJSON: ``,
					expectedErr:  ErrorTypeUnmatched{`array`, `map[string]interface {}`, `[::-1]`},
				},
			},
		},
		{
			name: `Filter-exist`,
			testCases: []TestCase{
				{
					jsonpath:     `$[?(@)]`,
					inputJSON:    `["a","b"]`,
					expectedJSON: `["a","b"]`,
				},
				{
					jsonpath:     `$[?(!@)]`,
					inputJSON:    `["a","b"]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[?(!@)]`},
				},
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
					jsonpath:     `$[?(@.c)]`,
					inputJSON:    `[{"a":1},{"b":2}]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[?(@.c)]`},
				},
				{
					jsonpath:     `$[?(!@.c)]`,
					inputJSON:    `[{"a":1},{"b":2}]`,
					expectedJSON: `[{"a":1},{"b":2}]`,
				},
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
					jsonpath:     `$[?(@)]`,
					inputJSON:    `{"a":1}`,
					expectedJSON: `[1]`,
				},
				{
					jsonpath:     `$[?(!@)]`,
					inputJSON:    `{"a":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[?(!@)]`},
				},
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
				{
					jsonpath:     `$.*[?(@.a)]`,
					inputJSON:    `[{"a":1},{"b":2}]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`.*[?(@.a)]`},
				},
				{
					jsonpath:     `$[?($[0].a)]`,
					inputJSON:    `[{"a":1},{"b":2}]`,
					expectedJSON: `[{"a":1},{"b":2}]`,
				},
				{
					jsonpath:     `$[?(!$[0].a)]`,
					inputJSON:    `[{"a":1},{"b":2}]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[?(!$[0].a)]`},
				},
				{
					jsonpath:     `$[?(@['a','b'])]`,
					inputJSON:    `[{"a":1},{"b":2}]`,
					expectedJSON: `[{"a":1},{"b":2}]`,
				},
				{
					jsonpath:     `$[?(@.*)]`,
					inputJSON:    `[{"a":1},{"b":2}]`,
					expectedJSON: `[{"a":1},{"b":2}]`,
				},
				{
					jsonpath:     `$[?(@[0:1])]`,
					inputJSON:    `[[{"a":1}],[]]`,
					expectedJSON: `[[{"a":1}]]`,
				},
				{
					jsonpath:     `$[?(@[*])]`,
					inputJSON:    `[[{"a":1}],[]]`,
					expectedJSON: `[[{"a":1}]]`,
				},
				{
					jsonpath:     `$[?(@[0,1])]`,
					inputJSON:    `[[{"a":1}],[]]`,
					expectedJSON: `[[{"a":1}]]`,
				},
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
		},
		{
			name: `Filter-compare`,
			testCases: []TestCase{
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
					jsonpath:     `$[?(@.a=='ab')]`,
					inputJSON:    `[{"a":"ab"}]`,
					expectedJSON: `[{"a":"ab"}]`,
				},
				{
					jsonpath:     `$[?(@.a!='ab')]`,
					inputJSON:    `[{"a":"ab"}]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[?(@.a!='ab')]`},
				},
				{
					jsonpath:     `$[?(@.a=='a\b')]`,
					inputJSON:    `[{"a":"ab"}]`,
					expectedJSON: `[{"a":"ab"}]`,
				},
				{
					jsonpath:     `$[?(@.a!='a\b')]`,
					inputJSON:    `[{"a":"ab"}]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[?(@.a!='a\b')]`},
				},
				{
					jsonpath:     `$[?(@.a=="ab")]`,
					inputJSON:    `[{"a":"ab"}]`,
					expectedJSON: `[{"a":"ab"}]`,
				},
				{
					jsonpath:     `$[?(@.a!="ab")]`,
					inputJSON:    `[{"a":"ab"}]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[?(@.a!="ab")]`},
				},
				{
					jsonpath:     `$[?(@.a=="a\b")]`,
					inputJSON:    `[{"a":"ab"}]`,
					expectedJSON: `[{"a":"ab"}]`,
				},
				{
					jsonpath:     `$[?(@.a!="a\b")]`,
					inputJSON:    `[{"a":"ab"}]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[?(@.a!="a\b")]`},
				},
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
					jsonpath:     `$[?(@.a == 2)].b`,
					inputJSON:    `[{"a":0},{"a":1},{"a":2,"b":4}]`,
					expectedJSON: `[4]`,
				},
				{
					jsonpath:     `$[?(@.a.b == 1)]`,
					inputJSON:    `[{"a":1},{"a":{"b":1}},{"a":{"a":1}}]`,
					expectedJSON: `[{"a":{"b":1}}]`,
				},
				{
					jsonpath:     `$..*[?(@.id>2)]`,
					inputJSON:    `[{"complexity":{"one":[{"name":"first","id":1},{"name":"next","id":2},{"name":"another","id":3},{"name":"more","id":4}],"more":{"name":"next to last","id":5}}},{"name":"last","id":6}]`,
					expectedJSON: `[{"id":5,"name":"next to last"},{"id":3,"name":"another"},{"id":4,"name":"more"}]`,
				},
				{
					jsonpath:     `$..[?(@.a==2)]`,
					inputJSON:    `{"a":2,"more":[{"a":2},{"b":{"a":2}},{"a":{"a":2}},[{"a":2}]]}`,
					expectedJSON: `[{"a":2},{"a":2},{"a":2},{"a":2}]`,
				},
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
					jsonpath:     `$[?(10==10)]`,
					inputJSON:    `[{"a":10},{"a":20},{"a":30},{"a+10":20}]`,
					expectedJSON: `[{"a":10},{"a":20},{"a":30},{"a+10":20}]`,
				},
				{
					jsonpath:     `$[?(10==20)]`,
					inputJSON:    `[{"a":10},{"a":20},{"a":30},{"a+10":20}]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[?(10==20)]`},
				},
				{
					jsonpath:     `$[?(@.a==@.a)]`,
					inputJSON:    `[{"a":10},{"a":20},{"a":30},{"a+10":20}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `comparison between two current nodes is prohibited`, `@.a==@.a)]`},
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
					jsonpath:  `$[?(@['a<=b']<2.1)]`,
					inputJSON: `[{"a<=b":1.9},{"a":2},{"a":2.1},{"b":3},{"a<=b":"test"}]`,
					// The character '<' is encoded to \u003c using Go's json.Marshal()
					expectedJSON: `[{"a\u003c=b":1.9}]`,
				},
				{
					jsonpath:     `$[?(@[-1]==2)]`,
					inputJSON:    `[[0,1],[0,2],[2],["2"],["a","b"],["b"]]`,
					expectedJSON: `[[0,2],[2]]`,
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
					jsonpath:     `$[?(@[1]=="b")]`,
					inputJSON:    `{"a":["a","b"],"b":["b"]}`,
					expectedJSON: `[["a","b"]]`,
				},
				{
					jsonpath:  `$[?(@.a*2==11)]`,
					inputJSON: `[{"a":6},{"a":5},{"a":5.5},{"a":-5},{"a*2":10.999},{"a*2":11.0},{"a*2":11.1},{"a*2":5},{"a*2":"11"}]`,
					// The number 11.0 is converted to 11 using Go's json.Marshal().
					expectedJSON: `[{"a*2":11}]`,
				},
				{
					jsonpath:     `$[?(@.a/10==5)]`,
					inputJSON:    `[{"a":60},{"a":50},{"a":51},{"a":-50},{"a/10":5},{"a/10":"5"}]`,
					expectedJSON: `[{"a/10":5}]`,
				},
				{
					jsonpath:  `$[?(@.a==5)]`,
					inputJSON: `[{"a":4.9},{"a":5.0},{"a":5.1},{"a":5},{"a":-5},{"a":"5"},{"a":"a"},{"a":true},{"a":null},{"a":{}},{"a":[]},{"b":5},{"a":{"a":5}},{"a":[{"a":5}]}]`,
					// The number 5.0 is converted to 5 using Go's json.Marshal().
					expectedJSON: `[{"a":5},{"a":5}]`,
				},
				{
					jsonpath:  `$[?(@==5)]`,
					inputJSON: `[4.999999,5.00000,5.00001,5,-5,"5","a",null,{},[],{"a":5},[5]]`,
					// The number 5.00000 is converted to 5 using Go's json.Marshal().
					expectedJSON: `[5,5]`,
				},
				{
					jsonpath:     `$[?(@.a==5)]`,
					inputJSON:    `[{"a":4.9},{"a":5.1},{"a":-5},{"a":"5"},{"a":"a"},{"a":true},{"a":null},{"a":{}},{"a":[]},{"b":5},{"a":{"a":5}},{"a":[{"a":5}]}]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[?(@.a==5)]`},
				},
				{
					jsonpath:  `$[?(@.a==1)]`,
					inputJSON: `{"a":{"a":0.999999},"b":{"a":1.0},"c":{"a":1.00001},"d":{"a":1},"e":{"a":-1},"f":{"a":"1"},"g":{"a":[1]}}`,
					// The number 1.0 is converted to 5 using Go's json.Marshal().
					expectedJSON: `[{"a":1},{"a":1}]`,
				},
				{
					jsonpath:     `$[?(@.a==1)]`,
					inputJSON:    `{"a":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[?(@.a==1)]`},
				},
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
				{
					jsonpath:     `$[?(@[0:1]==1)]`,
					inputJSON:    `[[1,2,3],[1],[2,3],1,2]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@[0:1]==1)]`},
				},
				{
					jsonpath:     `$[?(@[0:2]==1)]`,
					inputJSON:    `[[1,2,3],[1],[2,3],1,2]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@[0:2]==1)]`},
				},
				{
					jsonpath:     `$[?(@[*]==1)]`,
					inputJSON:    `[[1,2,3],[1],[2,3],1,2]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@[*]==1)]`},
				},
				{
					jsonpath:     `$[?(@[0,1]==1)]`,
					inputJSON:    `[[1,2,3],[1],[2,3],1,2]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@[0,1]==1)]`},
				},
				{
					jsonpath:     `$[?(@..a==123)]`,
					inputJSON:    `[{"a":"123"},{"a":123}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@..a==123)]`},
				},
				{
					jsonpath:     `$[?(@['a','b']==123)]`,
					inputJSON:    `[{"a":"123"},{"a":123}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@['a','b']==123)]`},
				},
				{
					jsonpath:     `$[?(@.*==2)]`,
					inputJSON:    `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@.*==2)]`},
				},
				{
					jsonpath:  `$[?(@.a==-0.123e2)]`,
					inputJSON: `[{"a":-12.3,"b":1},{"a":-0.123e2,"b":2},{"a":-0.123},{"a":-12},{"a":12.3},{"a":2},{"a":"-0.123e2"}]`,
					// The number -0.123e2 is converted to -12.3 using Go's json.Marshal().
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
				{
					jsonpath:     `$[?(@.a=="value")]`,
					inputJSON:    `[{"a":"value"},{"a":0},{"a":1},{"a":-1},{"a":"val"},{"a":true},{"a":{}},{"a":[]},{"a":["b"]},{"a":{"a":"value"}},{"b":"value"}]`,
					expectedJSON: `[{"a":"value"}]`,
				},
				{
					jsonpath:  `$[?(@.a=="~!@#$%^&*()-_=+[]\\{}|;':\",./<>?")]`,
					inputJSON: `[{"a":"~!@#$%^&*()-_=+[]\\{}|;':\",./<>?"}]`,
					// The character ['&','<','>'] is encoded to [\u0026,\u003c,\u003e] using Go's json.Marshal()
					expectedJSON: `[{"a":"~!@#$%^\u0026*()-_=+[]\\{}|;':\",./\u003c\u003e?"}]`,
				},
				{
					jsonpath:     `$[?(@.a=='value')]`,
					inputJSON:    `[{"a":"value"},{"a":0},{"a":1},{"a":-1},{"a":"val"},{"a":{}},{"a":[]},{"a":["b"]},{"a":{"a":"value"}},{"b":"value"}]`,
					expectedJSON: `[{"a":"value"}]`,
				},
				{
					jsonpath:  `$[?(@.a=='~!@#$%^&*()-_=+[]\\{}|;\':",./<>?')]`,
					inputJSON: `[{"a":"~!@#$%^&*()-_=+[]\\{}|;':\",./<>?"}]`,
					// The character ['&','<','>'] is encoded to [\u0026,\u003c,\u003e] using Go's json.Marshal()
					expectedJSON: `[{"a":"~!@#$%^\u0026*()-_=+[]\\{}|;':\",./\u003c\u003e?"}]`,
				},
				{
					jsonpath:     `$.a[?(@.b==$.c)]`,
					inputJSON:    `{"a":[{"b":123},{"b":123.456},{"b":"123.456"}],"c":123.456}`,
					expectedJSON: `[{"b":123.456}]`,
				},
				{
					jsonpath:     `$[?(@[*]>=2)]`,
					inputJSON:    `[[1,2],[3,4],[5,6]]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@[*]>=2)]`},
				},
				{
					jsonpath:     `$[?(@==$[1])]`,
					inputJSON:    `[[1],[2],[2],[3]]`,
					expectedJSON: `[[2],[2]]`,
				},
				{
					jsonpath:     `$[?(@==$[1])]`,
					inputJSON:    `[{"a":[1]},{"a":[2]},{"a":[2]},{"a":[3]}]`,
					expectedJSON: `[{"a":[2]},{"a":[2]}]`,
				},
				{
					jsonpath:     `$.*[?(@==1)]`,
					inputJSON:    `[{"a":1},{"b":2}]`,
					expectedJSON: `[1]`,
				},
				{
					jsonpath:     `$.*[?(@==1)]`,
					inputJSON:    `[[1],{"b":2}]`,
					expectedJSON: `[1]`,
				},
				{
					jsonpath:     `$.x[?(@[*]>=$.y[*])]`,
					inputJSON:    `{"x":[[1,2],[3,4],[5,6]],"y":[3,4,5]}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{6, `JSONPath that returns a value group is prohibited`, `@[*]>=$.y[*])]`},
				},
				{
					jsonpath:     `$.x[?(@[*]>=$.y.a[0:1])]`,
					inputJSON:    `{"x":[[1,2],[3,4],[5,6]],"y":{"a":[3,4,5]}}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{6, `JSONPath that returns a value group is prohibited`, `@[*]>=$.y.a[0:1])]`},
				},
				{
					jsonpath:     `$[?(@.a == $.b)]`,
					inputJSON:    `[{"a":1},{"a":2}]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[?(@.a == $.b)]`},
				},
				{
					jsonpath:     `$[?($.b == @.a)]`,
					inputJSON:    `[{"a":1},{"a":2}]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[?($.b == @.a)]`},
				},
				{
					jsonpath:     `$[?(@.b == $[0].a)]`,
					inputJSON:    `[{"a":1},{"a":2}]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[?(@.b == $[0].a)]`},
				},
				{
					jsonpath:     `$[?($[0].a == @.b)]`,
					inputJSON:    `[{"a":1},{"a":2}]`,
					expectedJSON: ``,
					expectedErr:  ErrorNoneMatched{`[?($[0].a == @.b)]`},
				},
			},
		},
		{
			name: `Sub-filter`,
			testCases: []TestCase{
				{
					jsonpath:     `$[?(@.a[?(@.b>1)])]`,
					inputJSON:    `[{"a":[{"b":1},{"b":2}]},{"a":[{"b":1}]}]`,
					expectedJSON: `[{"a":[{"b":1},{"b":2}]}]`,
				},
				{
					jsonpath:     `$[?(@.a[?(@.b)] > 1)]`,
					inputJSON:    `[{"a":[{"b":1},{"b":2}]},{"a":[{"b":1}]}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@.a[?(@.b)] > 1)]`},
				},
				{
					jsonpath:     `$[?(@.a[?(@.b)] > 1)]`,
					inputJSON:    `[{"a":[{"b":2}]},{"a":[{"b":1}]}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@.a[?(@.b)] > 1)]`},
				},
				{
					jsonpath:     `$[?(@.a[?(@.b)] > 1)]`,
					inputJSON:    `[{"a":[{"c":2}]},{"a":[{"d":1}]}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@.a[?(@.b)] > 1)]`},
				},
			},
		},
		{
			name: `Regex`,
			testCases: []TestCase{
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
					jsonpath:     `$[?(@.a=~/^\d+[a-d]\/\\$/)]`,
					inputJSON:    `[{"a":"012b/\\"},{"a":"ab/\\"},{"a":"1b\\"},{"a":"1b//"},{"a":"1b/\""}]`,
					expectedJSON: `[{"a":"012b/\\"}]`,
				},
				{
					jsonpath:     `$[?(@.a=~/テスト/)]`,
					inputJSON:    `[{"a":"123テストabc"}]`,
					expectedJSON: `[{"a":"123テストabc"}]`,
				},
				{
					jsonpath:     `$[?(@.a=~/(?i)CASE/)]`,
					inputJSON:    `[{"a":"case"},{"a":"CASE"},{"a":"Case"},{"a":"abc"}]`,
					expectedJSON: `[{"a":"case"},{"a":"CASE"},{"a":"Case"}]`,
				},
				{
					jsonpath:     `$[?($..a=~/123/)]`,
					inputJSON:    `[{"a":"123"},{"a":123}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `$..a=~/123/)]`},
				},
				{
					jsonpath:     `$[?($..a=~/123/)]`,
					inputJSON:    `[{"b":"123"},{"a":"123"}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `$..a=~/123/)]`},
				},
				{
					jsonpath:     `$[?(@['a','b']=~/123/)]`,
					inputJSON:    `[{"b":"123"},{"a":"123"}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@['a','b']=~/123/)]`},
				},
				{
					jsonpath:     `$[?(@.*=~/123/)]`,
					inputJSON:    `[{"b":"123"},{"a":"123"}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@.*=~/123/)]`},
				},
				{
					jsonpath:     `$[?(@[0:1]=~/123/)]`,
					inputJSON:    `[{"b":["123"]},{"a":["123"]}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@[0:1]=~/123/)]`},
				},
				{
					jsonpath:     `$[?(@[*]=~/123/)]`,
					inputJSON:    `[{"b":"123"},{"a":"123"}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@[*]=~/123/)]`},
				},
				{
					jsonpath:     `$[?(@[0,1]=~/123/)]`,
					inputJSON:    `[{"b":["123"]},{"a":[123,"123"]}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@[0,1]=~/123/)]`},
				},
				{
					jsonpath:     `$[?(@.a[?(@.b)]=~/123/)]`,
					inputJSON:    `[{"b":"123"},{"a":"123"}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@.a[?(@.b)]=~/123/)]`},
				},
			},
		},
		{
			name: `Filter-logical-combination`,
			testCases: []TestCase{
				{
					jsonpath:     `$[?(@.a || @.b)]`,
					inputJSON:    `[{"a":1},{"b":2},{"c":3}]`,
					expectedJSON: `[{"a":1},{"b":2}]`,
				},
				{
					jsonpath:     `$[?(@.a && @.b)]`,
					inputJSON:    `[{"a":1},{"b":2},{"a":3,"b":4}]`,
					expectedJSON: `[{"a":3,"b":4}]`,
				},
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
					jsonpath:     `$[?(@.a>1 && @.a<3)]`,
					inputJSON:    `[{"a":1},{"a":1.1},{"a":2.9},{"a":3}]`,
					expectedJSON: `[{"a":1.1},{"a":2.9}]`,
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
					jsonpath:     `$[?(@.a && (@.b || @.c))]`,
					inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
					expectedJSON: `[{"a":2,"b":2},{"a":3,"b":3,"c":3},{"a":5,"c":5}]`,
				},
				{
					jsonpath:     `$[?(@.a && @.b || @.c)]`,
					inputJSON:    `[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
					expectedJSON: `[{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6}]`,
				},
				{
					jsonpath:     `$[?(@.a =~ /a/ && @.b == 2)]`,
					inputJSON:    `[{"a":"a"},{"a":"a","b":2}]`,
					expectedJSON: `[{"a":"a","b":2}]`,
				},
			},
		},
		{
			name: `Space`,
			testCases: []TestCase{
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
					jsonpath:     `$.a` + "\n",
					inputJSON:    `{"a":123}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{3, `unrecognized input`, "\n"},
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
		},
		{
			name: `Invalid syntax`,
			testCases: []TestCase{
				{
					jsonpath:     ``,
					inputJSON:    `{"a":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{0, `unrecognized input`, ``},
				},
				{
					jsonpath:     `@`,
					inputJSON:    `{"a":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{0, `the use of '@' at the beginning is prohibited`, `@`},
				},
				{
					jsonpath:     `$$`,
					inputJSON:    `{"a":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `$`},
				},
				{
					jsonpath:     `$.`,
					inputJSON:    `{"a":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `.`},
				},
				{
					jsonpath:     `$..`,
					inputJSON:    `{"a":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `..`},
				},
				{
					jsonpath:     `$.a..`,
					inputJSON:    `{"a":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{3, `unrecognized input`, `..`},
				},
				{
					jsonpath:     `$..a..`,
					inputJSON:    `{"a":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `unrecognized input`, `..`},
				},
				{
					jsonpath:     `$...a`,
					inputJSON:    `{"a":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `...a`},
				},
				{
					jsonpath:     `$a`,
					inputJSON:    `{"a":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `a`},
				},
				{
					jsonpath:     `$['a]`,
					inputJSON:    `{"a":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `['a]`},
				},
				{
					jsonpath:     `$["a]`,
					inputJSON:    `{"a":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `["a]`},
				},
				{
					jsonpath:     `$.['a']`,
					inputJSON:    `{"a":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `.['a']`},
				},
				{
					jsonpath:     `$.["a"]`,
					inputJSON:    `{"a":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `.["a"]`},
				},
				{
					jsonpath:     `$[0].[1]`,
					inputJSON:    `[["a","b"],["c"],["d"]]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `unrecognized input`, `.[1]`},
				},
				{
					jsonpath:     `$[0].[1,2]`,
					inputJSON:    `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `unrecognized input`, `.[1,2]`},
				},
				{
					jsonpath:     `$[0,1].[1]`,
					inputJSON:    `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{6, `unrecognized input`, `.[1]`},
				},
				{
					jsonpath:     `$[0,1].[1,2]`,
					inputJSON:    `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{6, `unrecognized input`, `.[1,2]`},
				},
				{
					jsonpath:     `$[0:2].[1,2]`,
					inputJSON:    `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{6, `unrecognized input`, `.[1,2]`},
				},
				{
					jsonpath:     `$[0,1].[1:3]`,
					inputJSON:    `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{6, `unrecognized input`, `.[1:3]`},
				},
				{
					jsonpath:     `$.a.b[]`,
					inputJSON:    `{"a":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{5, `unrecognized input`, `[]`},
				},
				{
					jsonpath:     `.c`,
					inputJSON:    `{"a":"b","c":{"d":"e"}}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{0, `unrecognized input`, `.c`},
				},
				{
					jsonpath:     `$()`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `()`},
				},
				{
					jsonpath:     `$(a)`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `(a)`},
				},
				{
					jsonpath:     `$['a'.'b']`,
					inputJSON:    `["a"]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `['a'.'b']`},
				},
				{
					jsonpath:     `$[a.b]`,
					inputJSON:    `[{"a":{"b":1}}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[a.b]`},
				},
				{
					jsonpath:     `$['a'b']`,
					inputJSON:    `["a"]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `['a'b']`},
				},
				{
					jsonpath:     `$['a\\'b']`,
					inputJSON:    `["a"]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `['a\\'b']`},
				},
				{
					jsonpath:     `$['ab\']`,
					inputJSON:    `["a"]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `['ab\']`},
				},
				{
					jsonpath:     `$.[a]`,
					inputJSON:    `{"a":1}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `.[a]`},
				},
				{
					jsonpath:     `$[`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[`},
				},
				{
					jsonpath:     `$[0`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[0`},
				},
				{
					jsonpath:     `$[]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[]`},
				},
				{
					jsonpath:     `$[a]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[a]`},
				},
				{
					jsonpath:     `$[0,]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[0,]`},
				},
				{
					jsonpath:     `$[0,a]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[0,a]`},
				},
				{
					jsonpath:     `$[0,10000000000000000000,]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[0,10000000000000000000,]`},
				},
				{
					jsonpath:     `$[a:1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[a:1]`},
				},
				{
					jsonpath:     `$[0:a]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[0:a]`},
				},
				{
					jsonpath:     `$[0:10000000000000000000:a]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[0:10000000000000000000:a]`},
				},
				{
					jsonpath:     `$[?()]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?()]`},
				},
				{
					jsonpath:     `$[?@a]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?@a]`},
				},
				{
					jsonpath:     `$[?(@.a!!=1)]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a!!=1)]`},
				},
				{
					jsonpath:     `$[?(@.a!=)]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a!=)]`},
				},
				{
					jsonpath:     `$[?(@.a<=)]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a<=)]`},
				},
				{
					jsonpath:     `$[?(@.a<)]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a<)]`},
				},
				{
					jsonpath:     `$[?(@.a>=)]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a>=)]`},
				},
				{
					jsonpath:     `$[?(@.a>)]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a>)]`},
				},
				{
					jsonpath:     `$[?(!=@.a)]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(!=@.a)]`},
				},
				{
					jsonpath:     `$[?(<=@.a)]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(<=@.a)]`},
				},
				{
					jsonpath:     `$[?(<@.a)]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(<@.a)]`},
				},
				{
					jsonpath:     `$[?(>=@.a)]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(>=@.a)]`},
				},
				{
					jsonpath:     `$[?(>@.a)]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(>@.a)]`},
				},
				{
					jsonpath:     `$[?(@.a===1)]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a===1)]`},
				},
				{
					jsonpath:     `$[?(@.a=='abc`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a=='abc`},
				},
				{
					jsonpath:     `$[?(@.a=="abc`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a=="abc`},
				},
				{
					jsonpath:     `$[?(@.a==["b"])]`,
					inputJSON:    `[{"a":["b"]}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `["b"])]`},
				},
				{
					jsonpath:     `$[?(@[0:1]==[1])]`,
					inputJSON:    `[[1,2,3],[1],[2,3],1,2]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{12, `the omission of '$' allowed only at the beginning`, `[1])]`},
				},
				{
					jsonpath:     `$[?(@.*==[1,2])]`,
					inputJSON:    `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `[1,2])]`},
				},
				{
					jsonpath:     `$[?(@.*==['1','2'])]`,
					inputJSON:    `[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `['1','2'])]`},
				},
				{
					jsonpath:     `$[?((@.a<2)==false)]`,
					inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?((@.a<2)==false)]`},
				},
				{
					jsonpath:     `$[?((@.a<2)==true)]`,
					inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?((@.a<2)==true)]`},
				},
				{
					jsonpath:     `$[?((@.a<2)==1)]`,
					inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?((@.a<2)==1)]`},
				},
				{
					jsonpath:     `$[?(false)]`,
					inputJSON:    `[0,1,false,true,null,{},[]]`,
					expectedJSON: `[]`,
					expectedErr:  ErrorInvalidSyntax{4, `the omission of '$' allowed only at the beginning`, `false)]`},
				},
				{
					jsonpath:     `$[?(true)]`,
					inputJSON:    `[0,1,false,true,null,{},[]]`,
					expectedJSON: `[]`,
					expectedErr:  ErrorInvalidSyntax{4, `the omission of '$' allowed only at the beginning`, `true)]`},
				},
				{
					jsonpath:     `$[?(null)]`,
					inputJSON:    `[0,1,false,true,null,{},[]]`,
					expectedJSON: `[]`,
					expectedErr:  ErrorInvalidSyntax{4, `the omission of '$' allowed only at the beginning`, `null)]`},
				},
				{
					jsonpath:     `$[?(@.a>1 && )]`,
					inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a>1 && )]`},
				},
				{
					jsonpath:     `$[?(@.a>1 || )]`,
					inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a>1 || )]`},
				},
				{
					jsonpath:     `$[?( && @.a>1 )]`,
					inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?( && @.a>1 )]`},
				},
				{
					jsonpath:     `$[?( || @.a>1 )]`,
					inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?( || @.a>1 )]`},
				},
				{
					jsonpath:     `$[?(@.a>1 && false)]`,
					inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{13, `the omission of '$' allowed only at the beginning`, `false)]`},
				},
				{
					jsonpath:     `$[?(@.a>1 && true)]`,
					inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{13, `the omission of '$' allowed only at the beginning`, `true)]`},
				},
				{
					jsonpath:     `$[?(@.a>1 || false)]`,
					inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{13, `the omission of '$' allowed only at the beginning`, `false)]`},
				},
				{
					jsonpath:     `$[?(@.a>1 || true)]`,
					inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{13, `the omission of '$' allowed only at the beginning`, `true)]`},
				},
				{
					jsonpath:     `$[?(@.a>1 && ())]`,
					inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a>1 && ())]`},
				},
				{
					jsonpath:     `$[?(((@.a>1)))]`,
					inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
					expectedJSON: `[{"a":2},{"a":3}]`,
				},
				{
					jsonpath:     `$[?((@.a>1 )]`,
					inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?((@.a>1 )]`},
				},
				{
					jsonpath:     `$[?((@.a>1`,
					inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?((@.a>1`},
				},
				{
					jsonpath:     `$[?(!(@.a==2))]`,
					inputJSON:    `[{"a":1.9999},{"a":2},{"a":2.0001},{"a":"2"},{"a":true},{"a":{}},{"a":[]},{"a":["b"]},{"a":{"a":"value"}},{"b":"value"}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(!(@.a==2))]`},
				},
				{
					jsonpath:     `$[?(!(@.a<2))]`,
					inputJSON:    `[{"a":1.9999},{"a":2},{"a":2.0001},{"a":"2"},{"a":true},{"a":{}},{"a":[]},{"a":["b"]},{"a":{"a":"value"}},{"b":"value"}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(!(@.a<2))]`},
				},
				{
					jsonpath:     `$[?(@.a==fAlse)]`,
					inputJSON:    `[{"a":false}]`,
					expectedJSON: `[{"a":false}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `fAlse)]`},
				},
				{
					jsonpath:     `$[?(@.a==faLse)]`,
					inputJSON:    `[{"a":false}]`,
					expectedJSON: `[{"a":false}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `faLse)]`},
				},
				{
					jsonpath:     `$[?(@.a==falSe)]`,
					inputJSON:    `[{"a":false}]`,
					expectedJSON: `[{"a":false}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `falSe)]`},
				},
				{
					jsonpath:     `$[?(@.a==falsE)]`,
					inputJSON:    `[{"a":false}]`,
					expectedJSON: `[{"a":false}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `falsE)]`},
				},
				{
					jsonpath:     `$[?(@.a==FaLse)]`,
					inputJSON:    `[{"a":false}]`,
					expectedJSON: `[{"a":false}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `FaLse)]`},
				},
				{
					jsonpath:     `$[?(@.a==FalSe)]`,
					inputJSON:    `[{"a":false}]`,
					expectedJSON: `[{"a":false}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `FalSe)]`},
				},
				{
					jsonpath:     `$[?(@.a==FalsE)]`,
					inputJSON:    `[{"a":false}]`,
					expectedJSON: `[{"a":false}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `FalsE)]`},
				},
				{
					jsonpath:     `$[?(@.a==FaLSE)]`,
					inputJSON:    `[{"a":false}]`,
					expectedJSON: `[{"a":false}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `FaLSE)]`},
				},
				{
					jsonpath:     `$[?(@.a==FAlSE)]`,
					inputJSON:    `[{"a":false}]`,
					expectedJSON: `[{"a":false}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `FAlSE)]`},
				},
				{
					jsonpath:     `$[?(@.a==FALsE)]`,
					inputJSON:    `[{"a":false}]`,
					expectedJSON: `[{"a":false}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `FALsE)]`},
				},
				{
					jsonpath:     `$[?(@.a==FALSe)]`,
					inputJSON:    `[{"a":false}]`,
					expectedJSON: `[{"a":false}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `FALSe)]`},
				},
				{
					jsonpath:     `$[?(@.a==tRue)]`,
					inputJSON:    `[{"a":true}]`,
					expectedJSON: `[{"a":true}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `tRue)]`},
				},
				{
					jsonpath:     `$[?(@.a==trUe)]`,
					inputJSON:    `[{"a":true}]`,
					expectedJSON: `[{"a":true}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `trUe)]`},
				},
				{
					jsonpath:     `$[?(@.a==truE)]`,
					inputJSON:    `[{"a":true}]`,
					expectedJSON: `[{"a":true}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `truE)]`},
				},
				{
					jsonpath:     `$[?(@.a==TrUe)]`,
					inputJSON:    `[{"a":true}]`,
					expectedJSON: `[{"a":true}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `TrUe)]`},
				},
				{
					jsonpath:     `$[?(@.a==TruE)]`,
					inputJSON:    `[{"a":true}]`,
					expectedJSON: `[{"a":true}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `TruE)]`},
				},
				{
					jsonpath:     `$[?(@.a==TrUE)]`,
					inputJSON:    `[{"a":true}]`,
					expectedJSON: `[{"a":true}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `TrUE)]`},
				},
				{
					jsonpath:     `$[?(@.a==TRuE)]`,
					inputJSON:    `[{"a":true}]`,
					expectedJSON: `[{"a":true}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `TRuE)]`},
				},
				{
					jsonpath:     `$[?(@.a==TRUe)]`,
					inputJSON:    `[{"a":true}]`,
					expectedJSON: `[{"a":true}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `TRUe)]`},
				},
				{
					jsonpath:     `$[?(@.a==nUll)]`,
					inputJSON:    `[{"a":null}]`,
					expectedJSON: `[{"a":null}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `nUll)]`},
				},
				{
					jsonpath:     `$[?(@.a==nuLl)]`,
					inputJSON:    `[{"a":null}]`,
					expectedJSON: `[{"a":null}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `nuLl)]`},
				},
				{
					jsonpath:     `$[?(@.a==nulL)]`,
					inputJSON:    `[{"a":null}]`,
					expectedJSON: `[{"a":null}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `nulL)]`},
				},
				{
					jsonpath:     `$[?(@.a==NuLl)]`,
					inputJSON:    `[{"a":null}]`,
					expectedJSON: `[{"a":null}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `NuLl)]`},
				},
				{
					jsonpath:     `$[?(@.a==NulL)]`,
					inputJSON:    `[{"a":null}]`,
					expectedJSON: `[{"a":null}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `NulL)]`},
				},
				{
					jsonpath:     `$[?(@.a==NuLL)]`,
					inputJSON:    `[{"a":null}]`,
					expectedJSON: `[{"a":null}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `NuLL)]`},
				},
				{
					jsonpath:     `$[?(@.a==NUlL)]`,
					inputJSON:    `[{"a":null}]`,
					expectedJSON: `[{"a":null}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `NUlL)]`},
				},
				{
					jsonpath:     `$[?(@.a==NULl)]`,
					inputJSON:    `[{"a":null}]`,
					expectedJSON: `[{"a":null}]`,
					expectedErr:  ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `NULl)]`},
				},
				{
					jsonpath:     `$[?(@=={"k":"v"})]`,
					inputJSON:    `{}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{7, `the omission of '$' allowed only at the beginning`, `{"k":"v"})]`},
				},
				{
					jsonpath:     `$[?(@.a=~/abc)]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a=~/abc)]`},
				},
				{
					jsonpath:     `$[?(@.a=~///)]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a=~///)]`},
				},
				{
					jsonpath:     `$[?(@.a=~s/a/b/)]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a=~s/a/b/)]`},
				},
				{
					jsonpath:     `$[?(@.a=~@abc@)]`,
					inputJSON:    `[]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a=~@abc@)]`},
				},
				{
					jsonpath:     `$[?(a=~/123/)]`,
					inputJSON:    `[{"a":"123"},{"a":123}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{4, `the omission of '$' allowed only at the beginning`, `a=~/123/)]`},
				},
				{
					jsonpath:     `$[?(@.a=2)]`,
					inputJSON:    `[{"a":2}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a=2)]`},
				},
				{
					jsonpath:     `$[?(@.a<>2)]`,
					inputJSON:    `[{"a":2}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a<>2)]`},
				},
				{
					jsonpath:     `$[?(@.a=<2)]`,
					inputJSON:    `[{"a":2}]`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a=<2)]`},
				},
				{
					jsonpath:     `$[?(@.a),?(@.b)]`,
					inputJSON:    `{}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a),?(@.b)]`},
				},
				{
					jsonpath:     `$[?(@.a & @.b)]`,
					inputJSON:    `{}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a & @.b)]`},
				},
				{
					jsonpath:     `$[?(@.a | @.b)]`,
					inputJSON:    `{}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a | @.b)]`},
				},
				{
					jsonpath:     `$[()]`,
					inputJSON:    `{}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[()]`},
				},
				{
					jsonpath:     `$[(`,
					inputJSON:    `{}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[(`},
				},
				{
					jsonpath:     `$[(]`,
					inputJSON:    `{}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{1, `unrecognized input`, `[(]`},
				},
				{
					jsonpath:     `$.func(`,
					inputJSON:    `{}`,
					expectedJSON: ``,
					expectedErr:  ErrorInvalidSyntax{6, `unrecognized input`, `(`},
				},
			},
		},
		{
			name: `Invalid argument format`,
			testCases: []TestCase{
				{
					jsonpath:     `$[10000000000000000000]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr: ErrorInvalidArgument{
						argument: `10000000000000000000`,
						err:      fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
					},
				},
				{
					jsonpath:     `$[0,10000000000000000000]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr: ErrorInvalidArgument{
						argument: `10000000000000000000`,
						err:      fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
					},
				},
				{
					jsonpath:     `$[10000000000000000000:1]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr: ErrorInvalidArgument{
						argument: `10000000000000000000`,
						err:      fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
					},
				},
				{
					jsonpath:     `$[1:10000000000000000000]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr: ErrorInvalidArgument{
						argument: `10000000000000000000`,
						err:      fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
					},
				},
				{
					jsonpath:     `$[0:3:10000000000000000000]`,
					inputJSON:    `["first","second","third"]`,
					expectedJSON: ``,
					expectedErr: ErrorInvalidArgument{
						argument: `10000000000000000000`,
						err:      fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
					},
				},
				{
					jsonpath:     `$[?(@.a==1e1abc)]`,
					inputJSON:    `{}`,
					expectedJSON: ``,
					expectedErr: ErrorInvalidArgument{
						argument: `1e1abc`,
						err:      fmt.Errorf(`strconv.ParseFloat: parsing "1e1abc": invalid syntax`),
					},
				},
			},
		},
		{
			name: `Not supported`,
			testCases: []TestCase{
				{
					jsonpath:     `$[(command)]`,
					inputJSON:    `{}`,
					expectedJSON: ``,
					expectedErr:  ErrorNotSupported{`script`, `[(command)]`},
				},
			},
		},
	}

	for _, testGroup := range testGroups {
		for _, testCase := range testGroup.testCases {
			jsonPath := testCase.jsonpath
			srcJSON := testCase.inputJSON
			t.Run(
				fmt.Sprintf(`%s <%s> <%s>`, testGroup.name, jsonPath, srcJSON),
				func(t *testing.T) {
					var src interface{}
					if err := json.Unmarshal([]byte(srcJSON), &src); err != nil {
						t.Errorf("%w", err)
						return
					}
					execTestRetrieve(t, src, testCase)
				})
		}
	}
}

func TestRetrieve_jsonNumber(t *testing.T) {
	testGroups := []TestGroup{
		{
			name: `filter`,
			testCases: []TestCase{
				{
					jsonpath:     `$[?(@.a > 123)].a`,
					inputJSON:    `[{"a":123.456}]`,
					expectedJSON: `[123.456]`,
				},
				{
					jsonpath:     `$[?(@.a > 123.46)].a`,
					inputJSON:    `[{"a":123.456}]`,
					expectedJSON: `[]`,
					expectedErr:  ErrorNoneMatched{`[?(@.a > 123.46)].a`},
				},
				{
					jsonpath:     `$[?(@.a > 122)].a`,
					inputJSON:    `[{"a":123}]`,
					expectedJSON: `[123]`,
				},
				{
					jsonpath:     `$[?(123 < @.a)].a`,
					inputJSON:    `[{"a":123.456}]`,
					expectedJSON: `[123.456]`,
				},
				{
					jsonpath:     `$[?(@.a==-0.123e2)]`,
					inputJSON:    `[{"a":-12.3,"b":1},{"a":-0.123e2,"b":2},{"a":-0.123},{"a":-12},{"a":12.3},{"a":2},{"a":"-0.123e2"}]`,
					expectedJSON: `[{"a":-12.3,"b":1},{"a":-0.123e2,"b":2}]`,
				},
				{
					jsonpath:     `$[?(@.a==11)]`,
					inputJSON:    `[{"a":10.999},{"a":11.00},{"a":11.10}]`,
					expectedJSON: `[{"a":11.00}]`,
				},
			},
		},
	}

	for _, testGroup := range testGroups {
		for _, testCase := range testGroup.testCases {
			jsonPath := testCase.jsonpath
			srcJSON := testCase.inputJSON
			t.Run(
				fmt.Sprintf(`%s <%s> <%s>`, testGroup.name, jsonPath, srcJSON),
				func(t *testing.T) {
					var src interface{}
					reader := strings.NewReader(srcJSON)
					decoder := json.NewDecoder(reader)
					decoder.UseNumber()
					if err := decoder.Decode(&src); err != nil {
						t.Errorf("%w", err)
						return
					}
					execTestRetrieve(t, src, testCase)
				})
		}
	}
}

func TestRetrieveConfigFunction(t *testing.T) {
	twiceFunc := func(param interface{}) (interface{}, error) {
		if input, ok := param.(float64); ok {
			return input * 2, nil
		}
		return nil, fmt.Errorf(`type error`)
	}
	quarterFunc := func(param interface{}) (interface{}, error) {
		if input, ok := param.(float64); ok {
			return input / 4, nil
		}
		return nil, fmt.Errorf(`type error`)
	}
	maxFunc := func(param []interface{}) (interface{}, error) {
		var result float64
		for _, value := range param {
			if result < value.(float64) {
				result = value.(float64)
			}
		}
		return result, nil
	}
	minFunc := func(param []interface{}) (interface{}, error) {
		var result float64 = 999
		for _, value := range param {
			if result > value.(float64) {
				result = value.(float64)
			}
		}
		return result, nil
	}
	errAggregateFunc := func(param []interface{}) (interface{}, error) {
		return nil, fmt.Errorf(`aggregate error`)
	}
	errFilterFunc := func(param interface{}) (interface{}, error) {
		return nil, fmt.Errorf(`filter error`)
	}

	testGroups := []TestGroup{
		{
			name: `filter-function`,
			testCases: []TestCase{
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
		},
		{
			name: `aggregate-function`,
			testCases: []TestCase{
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
		},
		{
			name: `aggregate-filter-mix`,
			testCases: []TestCase{
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
		},
		{
			name: `filter-error`,
			testCases: []TestCase{
				{
					jsonpath:     `$.errFilter()`,
					inputJSON:    `[122.345,123.45,123.456]`,
					expectedJSON: ``,
					filters: map[string]func(interface{}) (interface{}, error){
						`errFilter`: errFilterFunc,
					},

					expectedErr: ErrorFunctionFailed{function: `.errFilter()`, err: fmt.Errorf(`filter error`)},
				},
				{
					jsonpath:     `$.*.errFilter()`,
					inputJSON:    `[122.345,123.45,123.456]`,
					expectedJSON: ``,
					filters: map[string]func(interface{}) (interface{}, error){
						`errFilter`: errFilterFunc,
					},

					expectedErr: ErrorNoneMatched{path: `.*.errFilter()`},
				},
				{
					jsonpath:     `$.*.max().errFilter()`,
					inputJSON:    `[122.345,123.45,123.456]`,
					expectedJSON: ``,
					filters: map[string]func(interface{}) (interface{}, error){
						`errFilter`: errFilterFunc,
					},
					aggregates: map[string]func([]interface{}) (interface{}, error){
						`max`: maxFunc,
					},
					expectedErr: ErrorFunctionFailed{function: `.errFilter()`, err: fmt.Errorf(`filter error`)},
				},
				{
					jsonpath:     `$.*.twice().errFilter()`,
					inputJSON:    `[122.345,123.45,123.456]`,
					expectedJSON: ``,
					filters: map[string]func(interface{}) (interface{}, error){
						`errFilter`: errFilterFunc,
						`twice`:     twiceFunc,
					},

					expectedErr: ErrorNoneMatched{path: `.*.twice().errFilter()`},
				}, {
					jsonpath:     `$.errFilter().twice()`,
					inputJSON:    `[122.345,123.45,123.456]`,
					expectedJSON: ``,
					filters: map[string]func(interface{}) (interface{}, error){
						`errFilter`: errFilterFunc,
						`twice`:     twiceFunc,
					},

					expectedErr: ErrorFunctionFailed{function: `.errFilter()`, err: fmt.Errorf(`filter error`)},
				},
				{
					jsonpath:     `$.*.errFilter().twice()`,
					inputJSON:    `[122.345,123.45,123.456]`,
					expectedJSON: ``,
					filters: map[string]func(interface{}) (interface{}, error){
						`errFilter`: errFilterFunc,
						`twice`:     twiceFunc,
					},

					expectedErr: ErrorNoneMatched{path: `.*.errFilter().twice()`},
				},
				{
					jsonpath:     `$.*.max().errFilter().twice()`,
					inputJSON:    `[122.345,123.45,123.456]`,
					expectedJSON: ``,
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
		},
		{
			name: `aggregate-error`,
			testCases: []TestCase{
				{
					jsonpath:     `$.*.errAggregate()`,
					inputJSON:    `[122.345,123.45,123.456]`,
					expectedJSON: ``,
					aggregates: map[string]func([]interface{}) (interface{}, error){
						`errAggregate`: errAggregateFunc,
					},
					expectedErr: ErrorFunctionFailed{function: `.errAggregate()`, err: fmt.Errorf(`aggregate error`)},
				},
				{
					jsonpath:     `$.*.max().errAggregate()`,
					inputJSON:    `[122.345,123.45,123.456]`,
					expectedJSON: ``,
					aggregates: map[string]func([]interface{}) (interface{}, error){
						`errAggregate`: errAggregateFunc,
						`max`:          maxFunc,
					},
					expectedErr: ErrorFunctionFailed{function: `.errAggregate()`, err: fmt.Errorf(`aggregate error`)},
				},
				{
					jsonpath:     `$.*.twice().errAggregate()`,
					inputJSON:    `[122.345,123.45,123.456]`,
					expectedJSON: ``,
					filters: map[string]func(interface{}) (interface{}, error){
						`twice`: twiceFunc,
					},
					aggregates: map[string]func([]interface{}) (interface{}, error){
						`errAggregate`: errAggregateFunc,
					},
					expectedErr: ErrorFunctionFailed{function: `.errAggregate()`, err: fmt.Errorf(`aggregate error`)},
				},
				{
					jsonpath:     `$.*.errAggregate().twice()`,
					inputJSON:    `[122.345,123.45,123.456]`,
					expectedJSON: ``,
					filters: map[string]func(interface{}) (interface{}, error){
						`twice`: twiceFunc,
					},
					aggregates: map[string]func([]interface{}) (interface{}, error){
						`errAggregate`: errAggregateFunc,
					},
					expectedErr: ErrorFunctionFailed{function: `.errAggregate()`, err: fmt.Errorf(`aggregate error`)},
				},
				{
					jsonpath:     `$.*.max().errAggregate().twice()`,
					inputJSON:    `[122.345,123.45,123.456]`,
					expectedJSON: ``,
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
					jsonpath:     `$.a.max()`,
					inputJSON:    `{}`,
					expectedJSON: ``,
					aggregates: map[string]func([]interface{}) (interface{}, error){
						`max`: maxFunc,
					},
					expectedErr: ErrorMemberNotExist{path: `.a`},
				},
			},
		},
		{
			name: `function-syntax-check`,
			testCases: []TestCase{
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
					jsonpath:     `$.*.unknown()`,
					inputJSON:    `[123.456,256]`,
					expectedJSON: ``,
					expectedErr:  ErrorFunctionNotFound{function: `.unknown()`},
				},
			},
		},
	}

	for _, testGroup := range testGroups {
		for _, testCase := range testGroup.testCases {
			jsonPath := testCase.jsonpath
			srcJSON := testCase.inputJSON
			expectedJSON := testCase.expectedJSON
			filterFunctions := testCase.filters
			aggregateFunctions := testCase.aggregates
			expectedError := testCase.expectedErr
			t.Run(
				fmt.Sprintf(`%s <%s> <%s>`, testGroup.name, jsonPath, srcJSON),
				func(t *testing.T) {
					var src interface{}
					if err := json.Unmarshal([]byte(srcJSON), &src); err != nil {
						t.Errorf("%w", err)
						return
					}
					config := Config{}
					for id, function := range filterFunctions {
						config.SetFilterFunction(id, function)
					}
					for id, function := range aggregateFunctions {
						config.SetAggregateFunction(id, function)
					}
					actualObject, err := Retrieve(jsonPath, src, config)
					if err != nil {
						if reflect.TypeOf(expectedError) == reflect.TypeOf(err) &&
							fmt.Sprintf(`%s`, expectedError) == fmt.Sprintf(`%s`, err) {
							return
						}
						t.Errorf("expected error<%s> != actual error<%s>\n",
							expectedError, err)
						return
					}
					if expectedError != nil {
						t.Errorf("expected error<%w> != actual error<none>\n", expectedError)
						return
					}
					actualOutputJSON, err := json.Marshal(actualObject)
					if err != nil {
						t.Errorf("%w", err)
						return
					}
					if string(expectedJSON) != string(actualOutputJSON) {
						t.Errorf("expectedJSON<%s> == actualOutputJSON<%s>\n",
							string(expectedJSON), string(actualOutputJSON))
						return
					}
				})
		}
	}
}

func TestParserFuncExecTwice(t *testing.T) {
	jsonpath := `$.a`
	srcJSON1 := `{"a":1}`
	srcJSON2 := `{"a":2}`

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

	if string(actualOutputJSON1) == string(actualOutputJSON2) {
		t.Errorf("actualOutputJSON1<%s> == expectedOutputJSON2<%s>\n",
			string(actualOutputJSON1), string(actualOutputJSON2))
		return
	}
}

func TestParserExecuteFunctions(t *testing.T) {
	stdoutBackup := os.Stdout
	os.Stdout = nil

	parser := pegJSONPathParser{Buffer: `$`}
	parser.Init()
	parser.Parse()
	parser.Execute()

	parser.AST().isZero()
	parser.Print()
	parser.PreOrder()
	parser.PrintSyntax()
	parser.PrintSyntaxTree()
	parser.Error()
	parser.Expand(10)
	parser.Highlighter()

	err := parseError{p: &parser}
	_ = err.Error()

	os.Stdout = stdoutBackup
}
