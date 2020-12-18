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
	testCases [][]interface{}
}

func execTestRetrieve(t *testing.T, src interface{}, testCase []interface{}) {
	jsonPath := testCase[0].(string)
	expectedOutputJSON := testCase[2].(string)
	var expectedError error
	if len(testCase) > 3 {
		expectedError = testCase[3].(error)
	}
	actualObject, err := Retrieve(jsonPath, src)
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
			`dot-notation`,
			[][]interface{}{
				{
					`$`,
					`{"a":"b","c":{"d":"e"}}`,
					`[{"a":"b","c":{"d":"e"}}]`,
				},
				{
					`$.a`,
					`{"a":"b","c":{"d":"e"}}`,
					`["b"]`,
				},
				{
					`$.c`,
					`{"a":"b","c":{"d":"e"}}`,
					`[{"d":"e"}]`,
				},
				{
					`a`,
					`{"a":"b","c":{"d":"e"}}`,
					`["b"]`,
				},
				{
					`$[0].a`,
					`[{"a":"b","c":{"d":"e"}},{"a":"y"}]`,
					`["b"]`,
				},
				{
					`[0].a`,
					`[{"a":"b","c":{"d":"e"}},{"a":"y"}]`,
					`["b"]`,
				},
				{
					`$[2,0].a`,
					`[{"a":"b","c":{"a":"d"}},{"a":"e"},{"a":"a"}]`,
					`["a","b"]`,
				},
				{
					`$[0:2].a`,
					`[{"a":"b","c":{"d":"e"}},{"a":"a"},{"a":"c"}]`,
					`["b","a"]`,
				},
				{
					`$.a.a2`,
					`{"a":{"a1":"1","a2":"2"},"b":{"b1":"3"}}`,
					`["2"]`,
				},
				{
					`$.null`,
					`{"null":1}`,
					`[1]`,
				},
				{
					`$.true`,
					`{"true":1}`,
					`[1]`,
				},
				{
					`$.false`,
					`{"false":1}`,
					`[1]`,
				},
				{
					`$.in`,
					`{"in":1}`,
					`[1]`,
				},
				{
					`$.length`,
					`{"length":1}`,
					`[1]`,
				},
				{
					`$.length`,
					`["length",1,2]`,
					``,
					ErrorTypeUnmatched{`object`, `[]interface {}`, `.length`},
				},
				{
					`$.a-b`,
					`{"a-b":1}`,
					`[1]`,
				},
				{
					`$.a:b`,
					`{"a:b":1}`,
					`[1]`,
				},
				{
					`$.$`,
					`{"$":1}`,
					`[1]`,
				},
				{
					`$`,
					`{"$":1}`,
					`[{"$":1}]`,
				},
				{
					`$.@`,
					`{"@":1}`,
					`[1]`,
				},
				{
					`$.'a'`,
					`{"'a'":1}`,
					`[1]`,
				},
				{
					`$."a"`,
					`{"\"a\"":1}`,
					`[1]`,
				},
				{
					`$.'a.b'`,
					`{"'a.b'":1,"a":{"b":2},"'a'":{"'b'":3},"'a":{"b'":4}}`,
					`[4]`,
				},
				{
					`$.'a\.b'`,
					`{"'a.b'":1,"a":{"b":2},"'a'":{"'b'":3},"'a":{"b'":4}}`,
					`[1]`,
				},
				{
					`$.\\`,
					`{"\\":1}`,
					`[1]`,
				},
				{
					`$.\.`,
					`{".":1}`,
					`[1]`,
				},
				{
					`$.\[`,
					`{"[":1}`,
					`[1]`,
				},
				{
					`$.\(`,
					`{"(":1}`,
					`[1]`,
				},
				{
					`$.\)`,
					`{")":1}`,
					`[1]`,
				},
				{
					`$.\=`,
					`{"=":1}`,
					`[1]`,
				},
				{
					`$.\!`,
					`{"!":1}`,
					`[1]`,
				},
				{
					`$.\>`,
					`{">":1}`,
					`[1]`,
				},
				{
					`$.\<`,
					`{"<":1}`,
					`[1]`,
				},
				{
					`$.\ `,
					`{" ":1}`,
					`[1]`,
				},
				{
					`$.\` + "\t",
					`{"":123}`,
					``,
					ErrorMemberNotExist{`.\` + "\t"},
				},
				{
					`$.\` + "\r",
					`{"":123}`,
					``,
					ErrorMemberNotExist{`.\` + "\r"},
				},
				{
					`$.\` + "\n",
					`{"":123}`,
					``,
					ErrorMemberNotExist{`.\` + "\n"},
				},
				{
					`$.\a`,
					`{"a":1}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `.\a`},
				},
				{
					`$.a\\b`,
					`{"a\\b":1}`,
					`[1]`,
				},
				{
					`$.a\.b`,
					`{"a.b":1}`,
					`[1]`,
				},
				{
					`$.a\[b`,
					`{"a[b":1}`,
					`[1]`,
				},
				{
					`$.a\(b`,
					`{"a(b":1}`,
					`[1]`,
				},
				{
					`$.a\)b`,
					`{"a)b":1}`,
					`[1]`,
				},
				{
					`$.a\=b`,
					`{"a=b":1}`,
					`[1]`,
				},
				{
					`$.a\!b`,
					`{"a!b":1}`,
					`[1]`,
				},
				{
					`$.a\>b`,
					`{"a>b":1}`,
					`[1]`,
				},
				{
					`$.a\<b`,
					`{"a<b":1}`,
					`[1]`,
				},
				{
					`$.a\ b`,
					`{"a b":1}`,
					`[1]`,
				},
				{
					`$.a\` + "\t" + `b`,
					`{"ab":123}`,
					``,
					ErrorMemberNotExist{`.a\` + "\t" + `b`},
				},
				{
					`$.a\` + "\r" + `b`,
					`{"ab":123}`,
					``,
					ErrorMemberNotExist{`.a\` + "\r" + `b`},
				},
				{
					`$.a\` + "\n" + `b`,
					`{"ab":123}`,
					``,
					ErrorMemberNotExist{`.a\` + "\n" + `b`},
				},
				{
					`$.a\a`,
					`{"aa":1}`,
					``,
					ErrorInvalidSyntax{3, `unrecognized input`, `\a`},
				},
				{
					`$.\`,
					`{"\\":1}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `.\`},
				},
				{
					`$.(`,
					`{"(":1}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `.(`},
				},
				{
					`$.)`,
					`{")":1}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `.)`},
				},
				{
					`$.=`,
					`{"=":1}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `.=`},
				},
				{
					`$.!`,
					`{"!":1}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `.!`},
				},
				{
					`$.>`,
					`{">":1}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `.>`},
				},
				{
					`$.<`,
					`{"<":1}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `.<`},
				},
				{
					`$. `,
					`{" ":1}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `. `},
				},
				{
					`$.` + "\t",
					`{"":123}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `.` + "\t"},
				},
				{
					`$.` + "\r",
					`{"":1}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `.` + "\r"},
				},
				{
					`$.` + "\n",
					`{"":123}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `.` + "\n"},
				},
				{
					`$.a\b`,
					`{"a\\b":1}`,
					``,
					ErrorInvalidSyntax{3, `unrecognized input`, `\b`},
				},
				{
					`$.a(b`,
					`{"(":1}`,
					``,
					ErrorInvalidSyntax{3, `unrecognized input`, `(b`},
				},
				{
					`$.a)b`,
					`{")":1}`,
					``,
					ErrorInvalidSyntax{3, `unrecognized input`, `)b`},
				},
				{
					`$.a=b`,
					`{"=":1}`,
					``,
					ErrorInvalidSyntax{3, `unrecognized input`, `=b`},
				},
				{
					`$.a!b`,
					`{"!":1}`,
					``,
					ErrorInvalidSyntax{3, `unrecognized input`, `!b`},
				},
				{
					`$.a>b`,
					`{">":1}`,
					``,
					ErrorInvalidSyntax{3, `unrecognized input`, `>b`},
				},
				{
					`$.a<b`,
					`{"<":1}`,
					``,
					ErrorInvalidSyntax{3, `unrecognized input`, `<b`},
				},
				{
					`$.a b`,
					`{" ":1}`,
					``,
					ErrorInvalidSyntax{4, `unrecognized input`, `b`},
				},
				{
					`$.a` + "\t" + `b`,
					`{"":123}`,
					``,
					ErrorInvalidSyntax{4, `unrecognized input`, `b`},
				},
				{
					`$.a` + "\r" + `b`,
					`{"":1}`,
					``,
					ErrorInvalidSyntax{3, `unrecognized input`, "\r" + `b`},
				},
				{
					`$.a` + "\n" + `b`,
					`{"":123}`,
					``,
					ErrorInvalidSyntax{3, `unrecognized input`, "\n" + `b`},
				},
				{
					`$.ﾃｽﾄソポァゼゾタダＡボマミ①`,
					`{"ﾃｽﾄソポァゼゾタダＡボマミ①":1}`,
					`[1]`,
				},
				{
					`$.d`,
					`{"a":"b","c":{"d":"e"}}`,
					``,
					ErrorMemberNotExist{`.d`},
				},
				{
					`$.2`,
					`{"a":1,"2":2,"3":{"2":1}}`,
					`[2]`,
				},
				{
					`$.2`,
					`["a","b",{"2":1}]`,
					``,
					ErrorTypeUnmatched{`object`, `[]interface {}`, `.2`},
				},
				{
					`$.-1`,
					`["a","b",{"2":1}]`,
					``,
					ErrorTypeUnmatched{`object`, `[]interface {}`, `.-1`},
				},
				{
					`$.a.d`,
					`{"a":"b","c":{"d":"e"}}`,
					``,
					ErrorTypeUnmatched{`object/array`, `string`, `.d`},
				},
				{
					`$.a.d`,
					`{"a":123}`,
					``,
					ErrorTypeUnmatched{`object/array`, `float64`, `.d`},
				},
				{
					`$.a.d`,
					`{"a":true}`,
					``,
					ErrorTypeUnmatched{`object/array`, `bool`, `.d`},
				},
				{
					`$.a.d`,
					`{"a":null}`,
					``,
					ErrorTypeUnmatched{`object/array`, `null`, `.d`},
				},
				{
					`$.a`,
					`[1,2]`,
					``,
					ErrorTypeUnmatched{`object`, `[]interface {}`, `.a`},
				},
				{
					`$.a`,
					`[{"a":1}]`,
					``,
					ErrorTypeUnmatched{`object`, `[]interface {}`, `.a`},
				},
			},
		},
		{
			`dot-notation-recursive-descent`,
			[][]interface{}{
				{
					`$.a..b`,
					`{"a":{"b":1,"c":{"b":2},"d":["b",{"a":3,"b":4}]},"b":5}`,
					`[1,2,4]`,
				},
				{
					`$..a`,
					`{"a":"b","c":{"a":"d"},"e":["a",{"a":{"a":"h"}}]}`,
					`["b","d",{"a":"h"},"h"]`,
				},
				{
					`$..[1]`,
					`[{"a":["b",{"c":{"a":"d"}}],"e":["f",{"g":{"a":"h"}}]},0]`,
					`[0,{"c":{"a":"d"}},{"g":{"a":"h"}}]`,
				},
				{
					`$..[1].a`,
					`[{"a":["b",{"a":{"a":"d"}}],"e":["f",{"g":{"a":"h"}}]},0]`,
					`[{"a":"d"}]`,
				},
				{
					`$..x`,
					`{"a":"b","c":{"a":"d"},"e":["f",{"g":{"a":"h"}}]}`,
					``,
					ErrorNoneMatched{`..x`},
				},
				{
					`$..a.x`,
					`{"a":"b","c":{"a":"d"},"e":["f",{"g":{"a":"h"}}]}`,
					``,
					ErrorNoneMatched{`..a.x`},
				},
				{
					`$..'a'`,
					`{"'a'":1,"b":{"'a'":2},"c":["'a'",{"d":{"'a'":{"'a'":3}}}]}`,
					`[1,2,{"'a'":3},3]`,
				},
				{
					`$.."a"`,
					`{"\"a\"":1,"b":{"\"a\"":2},"c":["\"a\"",{"d":{"\"a\"":{"\"a\"":3}}}]}`,
					`[1,2,{"\"a\"":3},3]`,
				},
				{
					`$..[?(@.a)]`,
					`{"a":1,"b":[{"a":2},{"b":{"a":3}},{"a":{"a":4}}]}`,
					`[{"a":2},{"a":{"a":4}},{"a":3},{"a":4}]`,
				},
				{
					`$..['a','b']`,
					`[{"a":1,"b":2,"c":{"a":3}},{"a":4},{"b":5},{"a":6,"b":7},{"d":{"b":8}}]`,
					`[1,2,3,4,5,6,7,8]`,
				},
			},
		},
		{
			`dot-notation-asterisk`,
			[][]interface{}{
				{
					`$.*`,
					`[[1],[2,3],123,"a",{"b":"c"},[0,1],null]`,
					`[[1],[2,3],123,"a",{"b":"c"},[0,1],null]`,
				},
				{
					`$.*[1]`,
					`[[1],[2,3],[4,[5,6,7]]]`,
					`[3,[5,6,7]]`,
				},
				{
					`$.*.a`,
					`[{"a":1},{"a":[2,3]}]`,
					`[1,[2,3]]`,
				},
				{
					`$..*`,
					`[{"a":1},{"a":[2,3]},null,true]`,
					`[{"a":1},{"a":[2,3]},null,true,1,[2,3],2,3]`,
				},
				{
					`$.*`,
					`{"a":[1],"b":[2,3],"c":{"d":4}}`,
					`[[1],[2,3],{"d":4}]`,
				},
				{
					`$..*`,
					`{"a":1,"b":[2,3],"c":{"d":4,"e":[5,6]}}`,
					`[1,[2,3],{"d":4,"e":[5,6]},2,3,4,[5,6],5,6]`,
				},
				{
					`$.*.*`,
					`[[1,2,3],[4,5,6]]`,
					`[1,2,3,4,5,6]`,
				},
				{
					`$.*.a.*`,
					`[{"a":[1]}]`,
					`[1]`,
				},
				{
					`$..[*]`,
					`{"a":1,"b":[2,3],"c":{"d":"e","f":[4,5]}}`,
					`[1,[2,3],{"d":"e","f":[4,5]},2,3,"e",[4,5],4,5]`,
				},
				{
					`$.*`,
					`{}`,
					``,
					ErrorNoneMatched{`.*`},
				},
				{
					`$.*`,
					`[]`,
					``,
					ErrorNoneMatched{`.*`},
				},
				{
					`$..*`,
					`"a"`,
					``,
					ErrorNoneMatched{`..*`},
				},
				{
					`$..*`,
					`true`,
					``,
					ErrorNoneMatched{`..*`},
				},
				{
					`$..*`,
					`1`,
					``,
					ErrorNoneMatched{`..*`},
				},
				{
					`$.*['a','b']`,
					`[{"a":1,"b":2,"c":3},{"a":4,"b":5,"d":6}]`,
					`[1,2,4,5]`,
				},
			},
		},
		{
			`bracket-notation`,
			[][]interface{}{
				{
					`$['a']`,
					`{"a":"b","c":{"d":"e"}}`,
					`["b"]`,
				},
				{
					`$['d']`,
					`{"a":"b","c":{"d":"e"}}`,
					``,
					ErrorMemberNotExist{`['d']`},
				},
				{
					`$[0]['a']`,
					`[{"a":"b","c":{"d":"e"}},{"x":"y"}]`,
					`["b"]`,
				},
				{
					`$['a'][0]['b']`,
					`{"a":[{"b":"x"},"y"],"c":{"d":"e"}}`,
					`["x"]`,
				},
				{
					`$[0:2]['b']`,
					`[{"a":1},{"b":3},{"b":2,"c":4}]`,
					`[3]`,
				},
				{
					`$[:]['b']`,
					`[{"a":1},{"b":3},{"b":2,"c":4}]`,
					`[3,2]`,
				},
				{
					`$['a']['a2']`,
					`{"a":{"a1":"1","a2":"2"},"b":{"b1":"3"}}`,
					`["2"]`,
				},
				{
					`$['0']`,
					`{"0":1,"a":2}`,
					`[1]`,
				},
				{
					`$['a\'b']`,
					`{"a'b":1,"b":2}`,
					`[1]`,
				},
				{
					`$['ab\'c']`,
					`{"ab'c":1,"b":2}`,
					`[1]`,
				},
				{
					`$['a\c']`,
					`{"ac":1,"b":2}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `['a\c']`},
				},
				{
					`$["a\c"]`,
					`{"ac":1,"b":2}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `["a\c"]`},
				},
				{
					`$['a.b']`,
					`{"a.b":1,"a":{"b":2}}`,
					`[1]`,
				},
				{
					`$["a"]`,
					`{"a":1}`,
					`[1]`,
				},
				{
					`$[':']`,
					`{":":1,"b":2}`,
					`[1]`,
				},
				{
					`$['[']`,
					`{"[":1,"]":2}`,
					`[1]`,
				},
				{
					`$[']']`,
					`{"[":1,"]":2}`,
					`[2]`,
				},
				{
					`$['$']`,
					`{"$":2}`,
					`[2]`,
				},
				{
					`$['@']`,
					`{"@":2}`,
					`[2]`,
				},
				{
					`$['*']`,
					`{"*":2}`,
					`[2]`,
				},
				{
					`$['*']`,
					`{"a":1,"b":2}`,
					``,
					ErrorMemberNotExist{`['*']`},
				},
				{
					`$['.']`,
					`{".":1}`,
					`[1]`,
				},
				{
					`$[',']`,
					`{",":1}`,
					`[1]`,
				},
				{
					`$['.*']`,
					`{".*":1}`,
					`[1]`,
				},
				{
					`$['"']`,
					`{"\"":1}`,
					`[1]`,
				},
				{
					`$["'"]`,
					`{"'":1}`,
					`[1]`,
				},
				{
					`$['\'']`,
					`{"'":1}`,
					`[1]`,
				},
				{
					`$["\""]`,
					`{"\"":1}`,
					`[1]`,
				},
				{
					`$['\\']`,
					`{"\\":1}`,
					`[1]`,
				},
				{
					`$["\\"]`,
					`{"\\":1}`,
					`[1]`,
				},
				{
					`$[':@."$,*\'\\']`,
					`{":@.\"$,*'\\": 1}`,
					`[1]`,
				},
				{
					`$['']`,
					`{"":1, "''":2}`,
					`[1]`,
				},
				{
					`$[""]`,
					`{"":1, "''":2,"\"\"":3}`,
					`[1]`,
				},
				{
					`$[''][0]`,
					`[1,2,3]`,
					`[1]`,
					ErrorTypeUnmatched{`object`, `[]interface {}`, `['']`},
				},
				{
					`$['a','b']`,
					`{"a":1, "b":2}`,
					`[1,2]`,
				},
				{
					`$['b','a']`,
					`{"a":1, "b":2}`,
					`[2,1]`,
				},
				{
					`$['b','a']`,
					`{"b":2,"a":1}`,
					`[2,1]`,
				},
				{
					`$['a','b']`,
					`{"b":2,"a":1}`,
					`[1,2]`,
				},
				{
					`$['a','b',0]`,
					`{"b":2,"a":1,"c":3}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `['a','b',0]`},
				},
				{
					`$['a','b'].a`,
					`{"a":{"a":1}, "b":{"c":2}}`,
					`[1]`,
				},
				{
					`$['a','b']['a']`,
					`{"a":{"a":1}, "b":{"c":2}}`,
					`[1]`,
				},
				{
					`$['c','d']`,
					`{"a":1,"b":2}`,
					``,
					ErrorNoneMatched{`['c','d']`},
				},
				{
					`$['a','d']`,
					`{"a":1,"b":2}`,
					`[1]`,
				},
				{
					`$['a','a']`,
					`{"b":2,"a":1}`,
					`[1,1]`,
				},
				{
					`$['a','a','b','b']`,
					`{"b":2,"a":1}`,
					`[1,1,2,2]`,
				},
				{
					`$[0]['a','b']`,
					`[{"a":1,"b":2},{"a":3,"b":4},{"a":5,"b":6}]`,
					`[1,2]`,
				},
				{
					`$[0:2]['b','a']`,
					`[{"a":1,"b":2},{"a":3,"b":4},{"a":5,"b":6}]`,
					`[2,1,4,3]`,
				},
				{
					`$['a'].b`,
					`{"b":2,"a":{"b":1}}`,
					`[1]`,
				},
			},
		},
		{
			`bracket-notation-asterisk`,
			[][]interface{}{
				{
					`$[*]`,
					`["a",123,true,{"b":"c"},[0,1],null]`,
					`["a",123,true,{"b":"c"},[0,1],null]`,
				},
				{
					`$[*]`,
					`{"a":[1],"b":[2,3]}`,
					`[[1],[2,3]]`,
				},
				{
					`$[*]`,
					`[]`,
					``,
					ErrorNoneMatched{`[*]`},
				},
				{
					`$[*]`,
					`{}`,
					``,
					ErrorNoneMatched{`[*]`},
				},
				{
					`$[0:2][*]`,
					`[[1,2],[3,4],[5,6]]`,
					`[1,2,3,4]`,
				},
				{
					`$[*].a`,
					`[{"a":1},{"b":2}]`,
					`[1]`,
				},
				{
					`$[*].a`,
					`[{"a":1},{"a":1}]`,
					`[1,1]`,
				},
				{
					`$[*].a`,
					`[{"a":[1,[2]]},{"a":2}]`,
					`[[1,[2]],2]`,
				},
				{
					`$[*].a[*]`,
					`[{"a":[1,[2]]},{"a":2}]`,
					`[1,[2]]`,
				},
			},
		},
		{
			`Value type`,
			[][]interface{}{
				{
					`$.a`,
					`{"a":"string"}`,
					`["string"]`,
				},
				{
					`$.a`,
					`{"a":123}`,
					`[123]`,
				},
				{
					`$.a`,
					`{"a":-123.456}`,
					`[-123.456]`,
				},
				{
					`$.a`,
					`{"a":true}`,
					`[true]`,
				},
				{
					`$.a`,
					`{"a":false}`,
					`[false]`,
				},
				{
					`$.a`,
					`{"a":null}`,
					`[null]`,
				},
				{
					`$.a`,
					`{"a":{"b":"c"}}`,
					`[{"b":"c"}]`,
				},
				{
					`$.a`,
					`{"a":[1,3,5]}`,
					`[[1,3,5]]`,
				},
				{
					`$.a`,
					`{"a":{}}`,
					`[{}]`,
				},
				{
					`$.a`,
					`{"a":[]}`,
					`[[]]`,
				},
				{
					`$`,
					`"a"`,
					`["a"]`,
				},
				{
					`$`,
					`2`,
					`[2]`,
				},
				{
					`$`,
					`false`,
					`[false]`,
				},
				{
					`$`,
					`true`,
					`[true]`,
				},
				{
					`$`,
					`null`,
					`[null]`,
				},
				{
					`$`,
					`{}`,
					`[{}]`,
				},
				{
					`$`,
					`[]`,
					`[[]]`,
				},
			},
		},
		{
			`Array-index`,
			[][]interface{}{
				{
					`$[0]`,
					`["first","second","third"]`,
					`["first"]`,
				},
				{
					`$[1]`,
					`["first","second","third"]`,
					`["second"]`,
				},
				{
					`$[3]`,
					`["first","second","third"]`,
					``,
					ErrorIndexOutOfRange{`[3]`},
				},
				{
					`$[+1]`,
					`["first","second","third"]`,
					`["second"]`,
				},
				{
					`$[01]`,
					`["first","second","third"]`,
					`["second"]`,
				},
				{
					`$[-1]`,
					`["first","second","third"]`,
					`["third"]`,
				},
				{
					`$[-2]`,
					`["first","second","third"]`,
					`["second"]`,
				},
				{
					`$[-3]`,
					`["first","second","third"]`,
					`["first"]`,
				},
				{
					`$[-4]`,
					`["first","second","third"]`,
					``,
					ErrorIndexOutOfRange{`[-4]`},
				},
				{
					`$[0][1]`,
					`[["a","b"],["c"],["d"]]`,
					`["b"]`,
				},
				{
					`$[0]`,
					`[]`,
					``,
					ErrorIndexOutOfRange{`[0]`},
				},
				{
					`$[1]`,
					`[]`,
					``,
					ErrorIndexOutOfRange{`[1]`},
				},
				{
					`$[-1]`,
					`[]`,
					``,
					ErrorIndexOutOfRange{`[-1]`},
				},
				{
					`$[1000000000000000000]`,
					`["first","second","third"]`,
					``,
					ErrorIndexOutOfRange{`[1000000000000000000]`},
				},
				{
					`$[0]`,
					`{"a":1,"b":2}`,
					``,
					ErrorTypeUnmatched{`array`, `map[string]interface {}`, `[0]`},
				},
				{
					`$[0]`,
					`"abc"`,
					``,
					ErrorTypeUnmatched{`array`, `string`, `[0]`},
				},
				{
					`$[0]`,
					`123`,
					``,
					ErrorTypeUnmatched{`array`, `float64`, `[0]`},
				},
				{
					`$[0]`,
					`true`,
					``,
					ErrorTypeUnmatched{`array`, `bool`, `[0]`},
				},
				{
					`$[0]`,
					`null`,
					``,
					ErrorTypeUnmatched{`array`, `null`, `[0]`},
				},
				{
					`$[0]`,
					`{}`,
					``,
					ErrorTypeUnmatched{`array`, `map[string]interface {}`, `[0]`},
				},
			},
		},
		{
			`Array-union`,
			[][]interface{}{
				{
					`$[0,0]`,
					`["first","second","third"]`,
					`["first","first"]`,
				},
				{
					`$[0,1]`,
					`["first","second","third"]`,
					`["first","second"]`,
				},
				{
					`$[2,0,1]`,
					`["first","second","third"]`,
					`["third","first","second"]`,
				},
				{
					`$[0,3]`,
					`["first","second","third"]`,
					`["first"]`,
				},
				{
					`$[0,-1]`,
					`["first","second","third"]`,
					`["first","third"]`,
				},
				{
					`$[0,1]`,
					`[["11","12","13"],["21","22","23"],["31","32","33"]]`,
					`[["11","12","13"],["21","22","23"]]`,
				},
				{
					`$[*]`,
					`["first","second","third"]`,
					`["first","second","third"]`,
				},
				{
					`$[*,0]`,
					`["first","second","third"]`,
					`["first","second","third","first"]`,
				},
				{
					`$[*,1:2]`,
					`["first","second","third"]`,
					`["first","second","third","second"]`,
				},
				{
					`$[1:2,0]`,
					`["first","second","third"]`,
					`["second","first"]`,
				},
				{
					`$[:2,0]`,
					`["first","second","third"]`,
					`["first","second","first"]`,
				},
			},
		},
		{
			`Array-slice-start-to-end`,
			[][]interface{}{
				{
					`$[0:0]`,
					`["first","second","third"]`,
					``,
					ErrorNoneMatched{`[0:0]`},
				},
				{
					`$[0:3]`,
					`["first","second","third"]`,
					`["first","second","third"]`,
				},
				{
					`$[0:2]`,
					`["first","second","third"]`,
					`["first","second"]`,
				},
				{
					`$[1:1]`,
					`["first","second","third"]`,
					``,
					ErrorNoneMatched{`[1:1]`},
				},
				{
					`$[1:2]`,
					`["first","second","third"]`,
					`["second"]`,
				},
				{
					`$[1:3]`,
					`["first","second","third"]`,
					`["second","third"]`,
				},
				{
					`$[2:1]`,
					`["first","second","third"]`,
					``,
					ErrorNoneMatched{`[2:1]`},
				},
				{
					`$[3:2]`,
					`["first","second","third"]`,
					``,
					ErrorNoneMatched{`[3:2]`},
				},
				{
					`$[3:3]`,
					`["first","second","third"]`,
					``,
					ErrorNoneMatched{`[3:3]`},
				},
				{
					`$[3:4]`,
					`["first","second","third"]`,
					``,
					ErrorNoneMatched{`[3:4]`},
				},
				{
					`$[-1:-1]`,
					`["first","second","third"]`,
					``,
					ErrorNoneMatched{`[-1:-1]`},
				},
				{
					`$[-2:-1]`,
					`["first","second","third"]`,
					`["second"]`,
				},
				{
					`$[-1:-2]`,
					`["first","second","third"]`,
					``,
					ErrorNoneMatched{`[-1:-2]`},
				},
				{
					`$[-1:3]`,
					`["first","second","third"]`,
					`["third"]`,
				},
				{
					`$[-1:2]`,
					`["first","second","third"]`,
					``,
					ErrorNoneMatched{`[-1:2]`},
				},
				{
					`$[-4:3]`,
					`["first","second","third"]`,
					`["first","second","third"]`,
				},
				{
					`$[0:-1]`,
					`["first","second","third"]`,
					`["first","second"]`,
				},
				{
					`$[0:-3]`,
					`["first","second","third"]`,
					``,
					ErrorNoneMatched{`[0:-3]`},
				},
				{
					`$[0:-4]`,
					`["first","second","third"]`,
					``,
					ErrorNoneMatched{`[0:-4]`},
				},
				{
					`$[1:-2]`,
					`["first","second","third"]`,
					``,
					ErrorNoneMatched{`[1:-2]`},
				},
				{
					`$[1:-1]`,
					`["first","second","third"]`,
					`["second"]`,
				},
				{
					`$[:2]`,
					`["first","second","third"]`,
					`["first","second"]`,
				},
				{
					`$[1:]`,
					`["first","second","third"]`,
					`["second","third"]`,
				},
				{
					`$[-1:]`,
					`["first","second","third"]`,
					`["third"]`,
				},
				{
					`$[-2:]`,
					`["first","second","third"]`,
					`["second","third"]`,
				},
				{
					`$[-4:]`,
					`["first","second","third"]`,
					`["first","second","third"]`,
				},
				{
					`$[:]`,
					`["first","second","third"]`,
					`["first","second","third"]`,
				},
				{
					`$[-1000000000000000000:1]`,
					`["first","second","third"]`,
					`["first"]`,
				},
				{
					`$[1000000000000000000:1]`,
					`["first","second","third"]`,
					``,
					ErrorNoneMatched{`[1000000000000000000:1]`},
				},
				{
					`$[1:1000000000000000000]`,
					`["first","second","third"]`,
					`["second","third"]`,
				},
				{
					`$[1:2]`,
					`{"first":1,"second":2,"third":3}`,
					``,
					ErrorTypeUnmatched{`array`, `map[string]interface {}`, `[1:2]`},
				},
				{
					`$[:]`,
					`{"first":1,"second":2,"third":3}`,
					``,
					ErrorTypeUnmatched{`array`, `map[string]interface {}`, `[:]`},
				},
				{
					`$[+0:+1]`,
					`["first","second","third"]`,
					`["first"]`,
				},
				{
					`$[01:02]`,
					`["first","second","third"]`,
					`["second"]`,
				},
			},
		},
		{
			`Array-slice-step`,
			[][]interface{}{
				{
					`$[0:2:1]`,
					`["first","second","third"]`,
					`["first","second"]`,
				},
				{
					`$[0:3:2]`,
					`["first","second","third"]`,
					`["first","third"]`,
				},
				{
					`$[0:3:3]`,
					`["first","second","third"]`,
					`["first"]`,
				},
				{
					`$[0:2:2]`,
					`["first","second","third"]`,
					`["first"]`,
				},
				{
					`$[0:2:0]`,
					`["first","second","third"]`,
					`["first","second"]`,
				},
				{
					`$[0:3:-1]`,
					`["first","second","third"]`,
					``,
					ErrorNoneMatched{`[0:3:-1]`},
				},
				{
					`$[2:0:-1]`,
					`["first","second","third"]`,
					`["third","second"]`,
				},
				{
					`$[2:0:-2]`,
					`["first","second","third"]`,
					`["third"]`,
				},
				{
					`$[2:-1:-2]`,
					`["first","second","third"]`,
					`["third","first"]`,
				},
				{
					`$[3:1:-1]`,
					`["first","second","third"]`,
					``,
					ErrorNoneMatched{`[3:1:-1]`},
				},
				{
					`$[4:1:-1]`,
					`["first","second","third"]`,
					``,
					ErrorNoneMatched{`[4:1:-1]`},
				},
				{
					`$[5:1:-1]`,
					`["first","second","third"]`,
					`["third"]`,
				},
				{
					`$[6:1:-1]`,
					`["first","second","third"]`,
					`["third"]`,
				},
				{
					`$[2:2:-1]`,
					`["first","second","third"]`,
					`["third","second","first"]`,
				},
				{
					`$[2:3:-1]`,
					`["first","second","third"]`,
					`["third","second"]`,
				},
				{
					`$[2:5:-1]`,
					`["first","second","third"]`,
					``,
					ErrorNoneMatched{`[2:5:-1]`},
				},
				{
					`$[2:6:-1]`,
					`["first","second","third"]`,
					``,
					ErrorNoneMatched{`[2:6:-1]`},
				},
				{
					`$[2:7:-1]`,
					`["first","second","third"]`,
					``,
					ErrorNoneMatched{`[2:7:-1]`},
				},
				{
					`$[-1:0:-1]`,
					`["first","second","third"]`,
					``,
					ErrorNoneMatched{`[-1:0:-1]`},
				},
				{
					`$[2:-1:-1]`,
					`["first","second","third"]`,
					`["third","second","first"]`,
				},
				{
					`$[0:3:]`,
					`["first","second","third"]`,
					`["first","second","third"]`,
				},
				{
					`$[::]`,
					`["first","second","third"]`,
					`["first","second","third"]`,
				},
				{
					`$[1::-1]`,
					`["first","second","third"]`,
					`["second","first"]`,
				},
				{
					`$[:1:-1]`,
					`["first","second","third"]`,
					`["third"]`,
				},
				{
					`$[::2]`,
					`["first","second","third"]`,
					`["first","third"]`,
				},
				{
					`$[::-1]`,
					`["first","second","third"]`,
					`["third","second","first"]`,
				},
				{
					`$[1:1000000000000000000:1]`,
					`["first","second","third"]`,
					`["second","third"]`,
				},
				{
					`$[1:-1000000000000000000:-1]`,
					`["first","second","third"]`,
					`["second","first"]`,
				},
				{
					`$[-1000000000000000000:3:1]`,
					`["first","second","third"]`,
					`["first","second","third"]`,
				},
				{
					`$[1000000000000000000:0:-1]`,
					`["first","second","third"]`,
					`["third","second"]`,
				},
				{
					`$[0:3:+1]`,
					`["first","second","third"]`,
					`["first","second","third"]`,
				},
				{
					`$[0:3:01]`,
					`["first","second","third"]`,
					`["first","second","third"]`,
				},
				{
					`$[2:1:-1]`,
					`{"first":1,"second":2,"third":3}`,
					``,
					ErrorTypeUnmatched{`array`, `map[string]interface {}`, `[2:1:-1]`},
				},
				{
					`$[::-1]`,
					`{"first":1,"second":2,"third":3}`,
					``,
					ErrorTypeUnmatched{`array`, `map[string]interface {}`, `[::-1]`},
				},
			},
		},
		{
			`Filter-exist`,
			[][]interface{}{
				{
					`$[?(@)]`,
					`["a","b"]`,
					`["a","b"]`,
				},
				{
					`$[?(!@)]`,
					`["a","b"]`,
					``,
					ErrorNoneMatched{`[?(!@)]`},
				},
				{
					`$[?(@.a)]`,
					`[{"b":2},{"a":1},{"a":"value"},{"a":""},{"a":true},{"a":false},{"a":null},{"a":{}},{"a":[]}]`,
					`[{"a":1},{"a":"value"},{"a":""},{"a":true},{"a":false},{"a":null},{"a":{}},{"a":[]}]`,
				},
				{
					`$[?(!@.a)]`,
					`[{"b":2},{"a":1},{"a":"value"},{"a":""},{"a":true},{"a":false},{"a":null},{"a":{}},{"a":[]}]`,
					`[{"b":2}]`,
				},
				{
					`$[?(@.c)]`,
					`[{"a":1},{"b":2}]`,
					``,
					ErrorNoneMatched{`[?(@.c)]`},
				},
				{
					`$[?(!@.c)]`,
					`[{"a":1},{"b":2}]`,
					`[{"a":1},{"b":2}]`,
				},
				{
					`$[?(@[1])]`,
					`[[{"a":1}],[{"b":2},{"c":3}],[],{"d":4}]`,
					`[[{"b":2},{"c":3}]]`,
				},
				{
					`$[?(!@[1])]`,
					`[[{"a":1}],[{"b":2},{"c":3}],[],{"d":4}]`,
					`[[{"a":1}],[],{"d":4}]`,
				},
				{
					`$[?(@[1:3])]`,
					`[[{"a":1}],[{"b":2},{"c":3}],[],{"d":4}]`,
					`[[{"b":2},{"c":3}]]`,
				},
				{
					`$[?(!@[1:3])]`,
					`[[{"a":1}],[{"b":2},{"c":3}],[],{"d":4}]`,
					`[[{"a":1}],[],{"d":4}]`,
				},
				{
					`$[?(@[1:3])]`,
					`[[{"a":1}],[{"b":2},{"c":3},{"e":5}],[],{"d":4}]`,
					`[[{"b":2},{"c":3},{"e":5}]]`,
				},
				{
					`$[?(!@[1:3])]`,
					`[[{"a":1}],[{"b":2},{"c":3},{"e":5}],[],{"d":4}]`,
					`[[{"a":1}],[],{"d":4}]`,
				},
				{
					`$[?(@)]`,
					`{"a":1}`,
					`[1]`,
				},
				{
					`$[?(!@)]`,
					`{"a":1}`,
					``,
					ErrorNoneMatched{`[?(!@)]`},
				},
				{
					`$[?(@.a1)]`,
					`{"a":{"a1":1},"b":{"b1":2}}`,
					`[{"a1":1}]`,
				},
				{
					`$[?(!@.a1)]`,
					`{"a":{"a1":1},"b":{"b1":2}}`,
					`[{"b1":2}]`,
				},
				{
					`$[?(@..a)]`,
					`[{"a":1},{"b":2},{"c":{"a":3}},{"a":{"a":4}}]`,
					`[{"a":1},{"c":{"a":3}},{"a":{"a":4}}]`,
				},
				{
					`$[?(!@..a)]`,
					`[{"a":1},{"b":2},{"c":{"a":3}},{"a":{"a":4}}]`,
					`[{"b":2}]`,
				},
				{
					`$[?(@[1])]`,
					`{"a":["a1"],"b":["b1","b2"],"c":[],"d":4}`,
					`[["b1","b2"]]`,
				},
				{
					`$[?(!@[1])]`,
					`{"a":["a1"],"b":["b1","b2"],"c":[],"d":4}`,
					`[["a1"],[],4]`,
				},
				{
					`$[?(@[1:3])]`,
					`{"a":[],"b":[2],"c":[3,4,5,6],"d":4}`,
					`[[3,4,5,6]]`,
				},
				{
					`$[?(!@[1:3])]`,
					`{"a":[],"b":[2],"c":[3,4,5,6],"d":4}`,
					`[[],[2],4]`,
				},
				{
					`$[?(@[1:3])]`,
					`{"a":[],"b":[2],"c":[3,4],"d":4}`,
					`[[3,4]]`,
				},
				{
					`$[?(!@[1:3])]`,
					`{"a":[],"b":[2],"c":[3,4],"d":4}`,
					`[[],[2],4]`,
				},
				{
					`$.*[?(@.a)]`,
					`[{"a":1},{"b":2}]`,
					``,
					ErrorNoneMatched{`.*[?(@.a)]`},
				},
				{
					`$[?($[0].a)]`,
					`[{"a":1},{"b":2}]`,
					`[{"a":1},{"b":2}]`,
				},
				{
					`$[?(!$[0].a)]`,
					`[{"a":1},{"b":2}]`,
					``,
					ErrorNoneMatched{`[?(!$[0].a)]`},
				},
				{
					`$[?(@['a','b'])]`,
					`[{"a":1},{"b":2}]`,
					`[{"a":1},{"b":2}]`,
				},
				{
					`$[?(@.*)]`,
					`[{"a":1},{"b":2}]`,
					`[{"a":1},{"b":2}]`,
				},
				{
					`$[?(@[0:1])]`,
					`[[{"a":1}],[]]`,
					`[[{"a":1}]]`,
				},
				{
					`$[?(@[*])]`,
					`[[{"a":1}],[]]`,
					`[[{"a":1}]]`,
				},
				{
					`$[?(@[0,1])]`,
					`[[{"a":1}],[]]`,
					`[[{"a":1}]]`,
				},
				{
					`$[?(@.a[?(@.b)])]`,
					`[{"a":[{"b":2},{"c":3}]},{"b":4}]`,
					`[{"a":[{"b":2},{"c":3}]}]`,
				},
				{
					`$[?(@.a[?(@.b > 1)])]`,
					`[{"a":[{"b":1},{"c":3}]},{"a":[{"b":2},{"c":5}]},{"b":4}]`,
					`[{"a":[{"b":2},{"c":5}]}]`,
				},
			},
		},
		{
			`Filter-compare`,
			[][]interface{}{
				{
					`$[?(@.a == 2.1)]`,
					`[{"a":0},{"a":1},{"a":2.0,"b":4},{"a":2.1,"b":5},{"a":2.2,"b":6},{"a":"2.1"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
					`[{"a":2.1,"b":5}]`,
				},
				{
					`$[?(2.1 == @.a)]`,
					`[{"a":0},{"a":1},{"a":2.0,"b":4},{"a":2.1,"b":5},{"a":2.2,"b":6},{"a":"2.1"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
					`[{"a":2.1,"b":5}]`,
				},
				{
					`$[?(@.a != 2)]`,
					`[{"a":0},{"a":1},{"a":2,"b":4},{"a":1.999999},{"a":2.000000000001},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
					`[{"a":0},{"a":1},{"a":1.999999},{"a":2.000000000001},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
				},
				{
					`$[?(2 != @.a)]`,
					`[{"a":0},{"a":1},{"a":2,"b":4},{"a":1.999999},{"a":2.000000000001},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
					`[{"a":0},{"a":1},{"a":1.999999},{"a":2.000000000001},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
				},
				{
					`$[?(@.a < 1)]`,
					`[{"a":-9999999},{"a":0.999999},{"a":1.0000000},{"a":1.0000001},{"a":2},{"a":"0.9"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
					`[{"a":-9999999},{"a":0.999999}]`,
				},
				{
					`$[?(1 > @.a)]`,
					`[{"a":-9999999},{"a":0.999999},{"a":1.0000000},{"a":1.0000001},{"a":2},{"a":"0.9"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
					`[{"a":-9999999},{"a":0.999999}]`,
				},
				{
					`$[?(@.a <= 1.00001)]`,
					`[{"a":0},{"a":1},{"a":1.00001},{"a":1.00002},{"a":2,"b":4},{"a":"0.9"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
					`[{"a":0},{"a":1},{"a":1.00001}]`,
				},
				{
					`$[?(1.00001 >= @.a)]`,
					`[{"a":0},{"a":1},{"a":1.00001},{"a":1.00002},{"a":2,"b":4},{"a":"0.9"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
					`[{"a":0},{"a":1},{"a":1.00001}]`,
				},
				{
					`$[?(@.a > 1)]`,
					`[{"a":0},{"a":0.9999},{"a":1},{"a":1.000001},{"a":2,"b":4},{"a":9999999999},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
					`[{"a":1.000001},{"a":2,"b":4},{"a":9999999999}]`,
				},
				{
					`$[?(1 < @.a)]`,
					`[{"a":0},{"a":0.9999},{"a":1},{"a":1.000001},{"a":2,"b":4},{"a":9999999999},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
					`[{"a":1.000001},{"a":2,"b":4},{"a":9999999999}]`,
				},
				{
					`$[?(@.a >= 1.000001)]`,
					`[{"a":0},{"a":1},{"a":1.000001},{"a":1.0000009},{"a":1.001},{"a":2,"b":4},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
					`[{"a":1.000001},{"a":1.001},{"a":2,"b":4}]`,
				},
				{
					`$[?(1.000001 <= @.a)]`,
					`[{"a":0},{"a":1},{"a":1.000001},{"a":1.0000009},{"a":1.001},{"a":2,"b":4},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
					`[{"a":1.000001},{"a":1.001},{"a":2,"b":4}]`,
				},
				{
					`$[?(@.a=='ab')]`,
					`[{"a":"ab"}]`,
					`[{"a":"ab"}]`,
				},
				{
					`$[?(@.a!='ab')]`,
					`[{"a":"ab"}]`,
					``,
					ErrorNoneMatched{`[?(@.a!='ab')]`},
				},
				{
					`$[?(@.a=='a\b')]`,
					`[{"a":"ab"}]`,
					`[{"a":"ab"}]`,
				},
				{
					`$[?(@.a!='a\b')]`,
					`[{"a":"ab"}]`,
					``,
					ErrorNoneMatched{`[?(@.a!='a\b')]`},
				},
				{
					`$[?(@.a=="ab")]`,
					`[{"a":"ab"}]`,
					`[{"a":"ab"}]`,
				},
				{
					`$[?(@.a!="ab")]`,
					`[{"a":"ab"}]`,
					``,
					ErrorNoneMatched{`[?(@.a!="ab")]`},
				},
				{
					`$[?(@.a=="a\b")]`,
					`[{"a":"ab"}]`,
					`[{"a":"ab"}]`,
				},
				{
					`$[?(@.a!="a\b")]`,
					`[{"a":"ab"}]`,
					``,
					ErrorNoneMatched{`[?(@.a!="a\b")]`},
				},
				{
					`$[?(@.a == $[2].b)]`,
					`[{"a":0},{"a":1},{"a":2,"b":1}]`,
					`[{"a":1}]`,
				},
				{
					`$[?($[2].b == @.a)]`,
					`[{"a":0},{"a":1},{"a":2,"b":1}]`,
					`[{"a":1}]`,
				},
				{
					`$[?(@.a == 2)].b`,
					`[{"a":0},{"a":1},{"a":2,"b":4}]`,
					`[4]`,
				},
				{
					`$[?(@.a.b == 1)]`,
					`[{"a":1},{"a":{"b":1}},{"a":{"a":1}}]`,
					`[{"a":{"b":1}}]`,
				},
				{
					`$..*[?(@.id>2)]`,
					`[{"complexity":{"one":[{"name":"first","id":1},{"name":"next","id":2},{"name":"another","id":3},{"name":"more","id":4}],"more":{"name":"next to last","id":5}}},{"name":"last","id":6}]`,
					`[{"id":5,"name":"next to last"},{"id":3,"name":"another"},{"id":4,"name":"more"}]`,
				},
				{
					`$..[?(@.a==2)]`,
					`{"a":2,"more":[{"a":2},{"b":{"a":2}},{"a":{"a":2}},[{"a":2}]]}`,
					`[{"a":2},{"a":2},{"a":2},{"a":2}]`,
				},
				{
					`$[?(@.a+10==20)]`,
					`[{"a":10},{"a":20},{"a":30},{"a+10":20}]`,
					`[{"a+10":20}]`,
				},
				{
					`$[?(@.a-10==20)]`,
					`[{"a":10},{"a":20},{"a":30},{"a-10":20}]`,
					`[{"a-10":20}]`,
				},
				{
					`$[?(10==10)]`,
					`[{"a":10},{"a":20},{"a":30},{"a+10":20}]`,
					`[{"a":10},{"a":20},{"a":30},{"a+10":20}]`,
				},
				{
					`$[?(10==20)]`,
					`[{"a":10},{"a":20},{"a":30},{"a+10":20}]`,
					``,
					ErrorNoneMatched{`[?(10==20)]`},
				},
				{
					`$[?(@.a==@.a)]`,
					`[{"a":10},{"a":20},{"a":30},{"a+10":20}]`,
					``,
					ErrorInvalidSyntax{4, `comparison between two current nodes is prohibited`, `@.a==@.a)]`},
				},
				{
					`$[?(@['a']<2.1)]`,
					`[{"a":1.9},{"a":2},{"a":2.1},{"a":3},{"a":"test"}]`,
					`[{"a":1.9},{"a":2}]`,
				},
				{
					`$[?(@['$a']<2.1)]`,
					`[{"$a":1.9},{"a":2},{"a":2.1},{"a":3},{"$a":"test"}]`,
					`[{"$a":1.9}]`,
				},
				{
					`$[?(@['@a']<2.1)]`,
					`[{"@a":1.9},{"a":2},{"a":2.1},{"a":3},{"@a":"test"}]`,
					`[{"@a":1.9}]`,
				},
				{
					`$[?(@['a==b']<2.1)]`,
					`[{"a==b":1.9},{"a":2},{"a":2.1},{"b":3},{"a==b":"test"}]`,
					`[{"a==b":1.9}]`,
				},
				{
					`$[?(@['a<=b']<2.1)]`,
					`[{"a<=b":1.9},{"a":2},{"a":2.1},{"b":3},{"a<=b":"test"}]`,
					// The character '<' is encoded to \u003c using Go's json.Marshal()
					`[{"a\u003c=b":1.9}]`,
				},
				{
					`$[?(@[-1]==2)]`,
					`[[0,1],[0,2],[2],["2"],["a","b"],["b"]]`,
					`[[0,2],[2]]`,
				},
				{
					`$[?(@[1]=="b")]`,
					`[[0,1],[0,2],[2],["2"],["a","b"],["b"]]`,
					`[["a","b"]]`,
				},
				{
					`$[?(@[1]=="a\"b")]`,
					`[[0,1],[2],["a","a\"b"],["a\"b"]]`,
					`[["a","a\"b"]]`,
				},
				{
					`$[?(@[1]=='b')]`,
					`[[0,1],[2],["a","b"],["b"]]`,
					`[["a","b"]]`,
				},
				{
					`$[?(@[1]=='a\'b')]`,
					`[[0,1],[2],["a","a'b"],["a'b"]]`,
					`[["a","a'b"]]`,
				},
				{
					`$[?(@[1]=="b")]`,
					`{"a":["a","b"],"b":["b"]}`,
					`[["a","b"]]`,
				},
				{
					`$[?(@.a*2==11)]`,
					`[{"a":6},{"a":5},{"a":5.5},{"a":-5},{"a*2":10.999},{"a*2":11.0},{"a*2":11.1},{"a*2":5},{"a*2":"11"}]`,
					// The number 11.0 is converted to 11 using Go's json.Marshal().
					`[{"a*2":11}]`,
				},
				{
					`$[?(@.a/10==5)]`,
					`[{"a":60},{"a":50},{"a":51},{"a":-50},{"a/10":5},{"a/10":"5"}]`,
					`[{"a/10":5}]`,
				},
				{
					`$[?(@.a==5)]`,
					`[{"a":4.9},{"a":5.0},{"a":5.1},{"a":5},{"a":-5},{"a":"5"},{"a":"a"},{"a":true},{"a":null},{"a":{}},{"a":[]},{"b":5},{"a":{"a":5}},{"a":[{"a":5}]}]`,
					// The number 5.0 is converted to 5 using Go's json.Marshal().
					`[{"a":5},{"a":5}]`,
				},
				{
					`$[?(@==5)]`,
					`[4.999999,5.00000,5.00001,5,-5,"5","a",null,{},[],{"a":5},[5]]`,
					// The number 5.00000 is converted to 5 using Go's json.Marshal().
					`[5,5]`,
				},
				{
					`$[?(@.a==5)]`,
					`[{"a":4.9},{"a":5.1},{"a":-5},{"a":"5"},{"a":"a"},{"a":true},{"a":null},{"a":{}},{"a":[]},{"b":5},{"a":{"a":5}},{"a":[{"a":5}]}]`,
					``,
					ErrorNoneMatched{`[?(@.a==5)]`},
				},
				{
					`$[?(@.a==1)]`,
					`{"a":{"a":0.999999},"b":{"a":1.0},"c":{"a":1.00001},"d":{"a":1},"e":{"a":-1},"f":{"a":"1"},"g":{"a":[1]}}`,
					// The number 1.0 is converted to 5 using Go's json.Marshal().
					`[{"a":1},{"a":1}]`,
				},
				{
					`$[?(@.a==1)]`,
					`{"a":1}`,
					``,
					ErrorNoneMatched{`[?(@.a==1)]`},
				},
				{
					`$[?(@.a==false)]`,
					`[{"a":null},{"a":false},{"a":true},{"a":0},{"a":1},{"a":"false"}]`,
					`[{"a":false}]`,
				},
				{
					`$[?(@.a==FALSE)]`,
					`[{"a":false}]`,
					`[{"a":false}]`,
				},
				{
					`$[?(@.a==False)]`,
					`[{"a":false}]`,
					`[{"a":false}]`,
				},
				{
					`$[?(@.a==true)]`,
					`[{"a":null},{"a":false},{"a":true},{"a":0},{"a":1},{"a":"false"}]`,
					`[{"a":true}]`,
				},
				{
					`$[?(@.a==TRUE)]`,
					`[{"a":true}]`,
					`[{"a":true}]`,
				},
				{
					`$[?(@.a==True)]`,
					`[{"a":true}]`,
					`[{"a":true}]`,
				},
				{
					`$[?(@.a==null)]`,
					`[{"a":null},{"a":false},{"a":true},{"a":0},{"a":1},{"a":"false"}]`,
					`[{"a":null}]`,
				},
				{
					`$[?(@.a==NULL)]`,
					`[{"a":null}]`,
					`[{"a":null}]`,
				},
				{
					`$[?(@.a==Null)]`,
					`[{"a":null}]`,
					`[{"a":null}]`,
				},
				{
					`$[?(@[0:1]==1)]`,
					`[[1,2,3],[1],[2,3],1,2]`,
					``,
					ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@[0:1]==1)]`},
				},
				{
					`$[?(@[0:2]==1)]`,
					`[[1,2,3],[1],[2,3],1,2]`,
					``,
					ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@[0:2]==1)]`},
				},
				{
					`$[?(@[*]==1)]`,
					`[[1,2,3],[1],[2,3],1,2]`,
					``,
					ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@[*]==1)]`},
				},
				{
					`$[?(@[0,1]==1)]`,
					`[[1,2,3],[1],[2,3],1,2]`,
					``,
					ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@[0,1]==1)]`},
				},
				{
					`$[?(@..a==123)]`,
					`[{"a":"123"},{"a":123}]`,
					``,
					ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@..a==123)]`},
				},
				{
					`$[?(@['a','b']==123)]`,
					`[{"a":"123"},{"a":123}]`,
					``,
					ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@['a','b']==123)]`},
				},
				{
					`$[?(@.*==2)]`,
					`[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
					``,
					ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@.*==2)]`},
				},
				{
					`$[?(@.a==-0.123e2)]`,
					`[{"a":-12.3,"b":1},{"a":-0.123e2,"b":2},{"a":-0.123},{"a":-12},{"a":12.3},{"a":2},{"a":"-0.123e2"}]`,
					// The number -0.123e2 is converted to -12.3 using Go's json.Marshal().
					`[{"a":-12.3,"b":1},{"a":-12.3,"b":2}]`,
				},
				{
					`$[?(@.a==-0.123E2)]`,
					`[{"a":-12.3}]`,
					`[{"a":-12.3}]`,
				},
				{
					`$[?(@.a==+0.123e+2)]`,
					`[{"a":-12.3},{"a":12.3}]`,
					`[{"a":12.3}]`,
				},
				{
					`$[?(@.a==-1.23e-1)]`,
					`[{"a":-12.3},{"a":-1.23},{"a":-0.123}]`,
					`[{"a":-0.123}]`,
				},
				{
					`$[?(@.a==010)]`,
					`[{"a":10},{"a":0},{"a":"010"},{"a":"10"}]`,
					`[{"a":10}]`,
				},
				{
					`$[?(@.a=="value")]`,
					`[{"a":"value"},{"a":0},{"a":1},{"a":-1},{"a":"val"},{"a":true},{"a":{}},{"a":[]},{"a":["b"]},{"a":{"a":"value"}},{"b":"value"}]`,
					`[{"a":"value"}]`,
				},
				{
					`$[?(@.a=="~!@#$%^&*()-_=+[]\\{}|;':\",./<>?")]`,
					`[{"a":"~!@#$%^&*()-_=+[]\\{}|;':\",./<>?"}]`,
					// The character ['&','<','>'] is encoded to [\u0026,\u003c,\u003e] using Go's json.Marshal()
					`[{"a":"~!@#$%^\u0026*()-_=+[]\\{}|;':\",./\u003c\u003e?"}]`,
				},
				{
					`$[?(@.a=='value')]`,
					`[{"a":"value"},{"a":0},{"a":1},{"a":-1},{"a":"val"},{"a":{}},{"a":[]},{"a":["b"]},{"a":{"a":"value"}},{"b":"value"}]`,
					`[{"a":"value"}]`,
				},
				{
					`$[?(@.a=='~!@#$%^&*()-_=+[]\\{}|;\':",./<>?')]`,
					`[{"a":"~!@#$%^&*()-_=+[]\\{}|;':\",./<>?"}]`,
					// The character ['&','<','>'] is encoded to [\u0026,\u003c,\u003e] using Go's json.Marshal()
					`[{"a":"~!@#$%^\u0026*()-_=+[]\\{}|;':\",./\u003c\u003e?"}]`,
				},
				{
					`$.a[?(@.b==$.c)]`,
					`{"a":[{"b":123},{"b":123.456},{"b":"123.456"}],"c":123.456}`,
					`[{"b":123.456}]`,
				},
				{
					`$[?(@[*]>=2)]`,
					`[[1,2],[3,4],[5,6]]`,
					``,
					ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@[*]>=2)]`},
				},
				{
					`$[?(@==$[1])]`,
					`[[1],[2],[2],[3]]`,
					`[[2],[2]]`,
				},
				{
					`$[?(@==$[1])]`,
					`[{"a":[1]},{"a":[2]},{"a":[2]},{"a":[3]}]`,
					`[{"a":[2]},{"a":[2]}]`,
				},
				{
					`$.*[?(@==1)]`,
					`[{"a":1},{"b":2}]`,
					`[1]`,
				},
				{
					`$.*[?(@==1)]`,
					`[[1],{"b":2}]`,
					`[1]`,
				},
				{
					`$.x[?(@[*]>=$.y[*])]`,
					`{"x":[[1,2],[3,4],[5,6]],"y":[3,4,5]}`,
					``,
					ErrorInvalidSyntax{6, `JSONPath that returns a value group is prohibited`, `@[*]>=$.y[*])]`},
				},
				{
					`$.x[?(@[*]>=$.y.a[0:1])]`,
					`{"x":[[1,2],[3,4],[5,6]],"y":{"a":[3,4,5]}}`,
					``,
					ErrorInvalidSyntax{6, `JSONPath that returns a value group is prohibited`, `@[*]>=$.y.a[0:1])]`},
				},
				{
					`$[?(@.a == $.b)]`,
					`[{"a":1},{"a":2}]`,
					``,
					ErrorNoneMatched{`[?(@.a == $.b)]`},
				},
				{
					`$[?($.b == @.a)]`,
					`[{"a":1},{"a":2}]`,
					``,
					ErrorNoneMatched{`[?($.b == @.a)]`},
				},
				{
					`$[?(@.b == $[0].a)]`,
					`[{"a":1},{"a":2}]`,
					``,
					ErrorNoneMatched{`[?(@.b == $[0].a)]`},
				},
				{
					`$[?($[0].a == @.b)]`,
					`[{"a":1},{"a":2}]`,
					``,
					ErrorNoneMatched{`[?($[0].a == @.b)]`},
				},
			},
		},
		{
			`Sub-filter`,
			[][]interface{}{
				{
					`$[?(@.a[?(@.b>1)])]`,
					`[{"a":[{"b":1},{"b":2}]},{"a":[{"b":1}]}]`,
					`[{"a":[{"b":1},{"b":2}]}]`,
				},
				{
					`$[?(@.a[?(@.b)] > 1)]`,
					`[{"a":[{"b":1},{"b":2}]},{"a":[{"b":1}]}]`,
					``,
					ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@.a[?(@.b)] > 1)]`},
				},
				{
					`$[?(@.a[?(@.b)] > 1)]`,
					`[{"a":[{"b":2}]},{"a":[{"b":1}]}]`,
					``,
					ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@.a[?(@.b)] > 1)]`},
				},
				{
					`$[?(@.a[?(@.b)] > 1)]`,
					`[{"a":[{"c":2}]},{"a":[{"d":1}]}]`,
					``,
					ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@.a[?(@.b)] > 1)]`},
				},
			},
		},
		{
			`Regex`,
			[][]interface{}{
				{
					`$[?(@.a =~ /ab/)]`,
					`[{"a":"abc"},{"a":1},{"a":"def"}]`,
					`[{"a":"abc"}]`,
				},
				{
					`$[?(@.a =~ /123/)]`,
					`[{"a":123},{"a":"123"},{"a":"12"},{"a":"23"},{"a":"0123"},{"a":"1234"}]`,
					`[{"a":"123"},{"a":"0123"},{"a":"1234"}]`,
				},
				{
					`$[?(@.a=~/^\d+[a-d]\/\\$/)]`,
					`[{"a":"012b/\\"},{"a":"ab/\\"},{"a":"1b\\"},{"a":"1b//"},{"a":"1b/\""}]`,
					`[{"a":"012b/\\"}]`,
				},
				{
					`$[?(@.a=~/テスト/)]`,
					`[{"a":"123テストabc"}]`,
					`[{"a":"123テストabc"}]`,
				},
				{
					`$[?(@.a=~/(?i)CASE/)]`,
					`[{"a":"case"},{"a":"CASE"},{"a":"Case"},{"a":"abc"}]`,
					`[{"a":"case"},{"a":"CASE"},{"a":"Case"}]`,
				},
				{
					`$[?($..a=~/123/)]`,
					`[{"a":"123"},{"a":123}]`,
					``,
					ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `$..a=~/123/)]`},
				},
				{
					`$[?($..a=~/123/)]`,
					`[{"b":"123"},{"a":"123"}]`,
					``,
					ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `$..a=~/123/)]`},
				},
				{
					`$[?(@['a','b']=~/123/)]`,
					`[{"b":"123"},{"a":"123"}]`,
					``,
					ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@['a','b']=~/123/)]`},
				},
				{
					`$[?(@.*=~/123/)]`,
					`[{"b":"123"},{"a":"123"}]`,
					``,
					ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@.*=~/123/)]`},
				},
				{
					`$[?(@[0:1]=~/123/)]`,
					`[{"b":["123"]},{"a":["123"]}]`,
					``,
					ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@[0:1]=~/123/)]`},
				},
				{
					`$[?(@[*]=~/123/)]`,
					`[{"b":"123"},{"a":"123"}]`,
					``,
					ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@[*]=~/123/)]`},
				},
				{
					`$[?(@[0,1]=~/123/)]`,
					`[{"b":["123"]},{"a":[123,"123"]}]`,
					``,
					ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@[0,1]=~/123/)]`},
				},
				{
					`$[?(@.a[?(@.b)]=~/123/)]`,
					`[{"b":"123"},{"a":"123"}]`,
					``,
					ErrorInvalidSyntax{4, `JSONPath that returns a value group is prohibited`, `@.a[?(@.b)]=~/123/)]`},
				},
			},
		},
		{
			`Filter-logical-combination`,
			[][]interface{}{
				{
					`$[?(@.a || @.b)]`,
					`[{"a":1},{"b":2},{"c":3}]`,
					`[{"a":1},{"b":2}]`,
				},
				{
					`$[?(@.a && @.b)]`,
					`[{"a":1},{"b":2},{"a":3,"b":4}]`,
					`[{"a":3,"b":4}]`,
				},
				{
					`$[?(!@.a)]`,
					`[{"a":1},{"b":2},{"a":3,"b":4}]`,
					`[{"b":2}]`,
				},
				{
					`$[?(!@.c)]`,
					`[{"a":1},{"b":2},{"a":3,"b":4}]`,
					`[{"a":1},{"b":2},{"a":3,"b":4}]`,
				},
				{
					`$[?(@.a>1 && @.a<3)]`,
					`[{"a":1},{"a":1.1},{"a":2.9},{"a":3}]`,
					`[{"a":1.1},{"a":2.9}]`,
				},
				{
					`$[?(@.a>2 || @.a<2)]`,
					`[{"a":1},{"a":1.9},{"a":2},{"a":2.1},{"a":3}]`,
					`[{"a":1},{"a":1.9},{"a":2.1},{"a":3}]`,
				},
				{
					`$[?(@.a<2 || @.a>2)]`,
					`[{"a":1},{"a":2},{"a":3}]`,
					`[{"a":1},{"a":3}]`,
				},
				{
					`$[?(@.a && (@.b || @.c))]`,
					`[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
					`[{"a":2,"b":2},{"a":3,"b":3,"c":3},{"a":5,"c":5}]`,
				},
				{
					`$[?(@.a && @.b || @.c)]`,
					`[{"a":1},{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6},{"b":7}]`,
					`[{"a":2,"b":2},{"a":3,"b":3,"c":3},{"b":4,"c":4},{"a":5,"c":5},{"c":6}]`,
				},
				{
					`$[?(@.a =~ /a/ && @.b == 2)]`,
					`[{"a":"a"},{"a":"a","b":2}]`,
					`[{"a":"a","b":2}]`,
				},
			},
		},
		{
			`Space`,
			[][]interface{}{
				{
					` $.a `,
					`{"a":123}`,
					`[123]`,
				},
				{
					"\t" + `$.a` + "\t",
					`{"a":123}`,
					`[123]`,
				},
				{
					`$.a` + "\n",
					`{"a":123}`,
					``,
					ErrorInvalidSyntax{3, `unrecognized input`, "\n"},
				},
				{
					`$[ "a" , "c" ]`,
					`{"a":1,"b":2,"c":3}`,
					`[1,3]`,
				},
				{
					`$[ 0 , 2 : 4 , * ]`,
					`[1,2,3,4,5]`,
					`[1,3,4,1,2,3,4,5]`,
				},
				{
					`$[ ?( @.a == 1 ) ]`,
					`[{"a":1}]`,
					`[{"a":1}]`,
				},
				{
					`$[ ?( @.a != 1 ) ]`,
					`[{"a":2}]`,
					`[{"a":2}]`,
				},
				{
					`$[ ?( @.a <= 1 ) ]`,
					`[{"a":1}]`,
					`[{"a":1}]`,
				},
				{
					`$[ ?( @.a < 1 ) ]`,
					`[{"a":0}]`,
					`[{"a":0}]`,
				},
				{
					`$[ ?( @.a >= 1 ) ]`,
					`[{"a":1}]`,
					`[{"a":1}]`,
				},
				{
					`$[ ?( @.a > 1 ) ]`,
					`[{"a":2}]`,
					`[{"a":2}]`,
				},
				{
					`$[ ?( @.a =~ /a/ ) ]`,
					`[{"a":"abc"}]`,
					`[{"a":"abc"}]`,
				},
				{
					`$[ ?( @.a == 1 && @.b == 2 ) ]`,
					`[{"a":1,"b":2}]`,
					`[{"a":1,"b":2}]`,
				},
				{
					`$[ ?( @.a == 1 || @.b == 2 ) ]`,
					`[{"a":1},{"b":2}]`,
					`[{"a":1},{"b":2}]`,
				},
				{
					`$[ ?( ! @.a ) ]`,
					`[{"a":1},{"b":2}]`,
					`[{"b":2}]`,
				},
			},
		},
		{
			`Invalid syntax`,
			[][]interface{}{
				{
					``,
					`{"a":1}`,
					``,
					ErrorInvalidSyntax{0, `unrecognized input`, ``},
				},
				{
					`@`,
					`{"a":1}`,
					``,
					ErrorInvalidSyntax{0, `the use of '@' at the beginning is prohibited`, `@`},
				},
				{
					`$$`,
					`{"a":1}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `$`},
				},
				{
					`$.`,
					`{"a":1}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `.`},
				},
				{
					`$..`,
					`{"a":1}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `..`},
				},
				{
					`$.a..`,
					`{"a":1}`,
					``,
					ErrorInvalidSyntax{3, `unrecognized input`, `..`},
				},
				{
					`$..a..`,
					`{"a":1}`,
					``,
					ErrorInvalidSyntax{4, `unrecognized input`, `..`},
				},
				{
					`$...a`,
					`{"a":1}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `...a`},
				},
				{
					`$a`,
					`{"a":1}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `a`},
				},
				{
					`$['a]`,
					`{"a":1}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `['a]`},
				},
				{
					`$["a]`,
					`{"a":1}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `["a]`},
				},
				{
					`$.['a']`,
					`{"a":1}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `.['a']`},
				},
				{
					`$.["a"]`,
					`{"a":1}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `.["a"]`},
				},
				{
					`$[0].[1]`,
					`[["a","b"],["c"],["d"]]`,
					``,
					ErrorInvalidSyntax{4, `unrecognized input`, `.[1]`},
				},
				{
					`$[0].[1,2]`,
					`[["11","12","13"],["21","22","23"],["31","32","33"]]`,
					``,
					ErrorInvalidSyntax{4, `unrecognized input`, `.[1,2]`},
				},
				{
					`$[0,1].[1]`,
					`[["11","12","13"],["21","22","23"],["31","32","33"]]`,
					``,
					ErrorInvalidSyntax{6, `unrecognized input`, `.[1]`},
				},
				{
					`$[0,1].[1,2]`,
					`[["11","12","13"],["21","22","23"],["31","32","33"]]`,
					``,
					ErrorInvalidSyntax{6, `unrecognized input`, `.[1,2]`},
				},
				{
					`$[0:2].[1,2]`,
					`[["11","12","13"],["21","22","23"],["31","32","33"]]`,
					``,
					ErrorInvalidSyntax{6, `unrecognized input`, `.[1,2]`},
				},
				{
					`$[0,1].[1:3]`,
					`[["11","12","13"],["21","22","23"],["31","32","33"]]`,
					``,
					ErrorInvalidSyntax{6, `unrecognized input`, `.[1:3]`},
				},
				{
					`$.a.b[]`,
					`{"a":1}`,
					``,
					ErrorInvalidSyntax{5, `unrecognized input`, `[]`},
				},
				{
					`.c`,
					`{"a":"b","c":{"d":"e"}}`,
					``,
					ErrorInvalidSyntax{0, `unrecognized input`, `.c`},
				},
				{
					`$()`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `()`},
				},
				{
					`$(a)`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `(a)`},
				},
				{
					`$['a'.'b']`,
					`["a"]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `['a'.'b']`},
				},
				{
					`$[a.b]`,
					`[{"a":{"b":1}}]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[a.b]`},
				},
				{
					`$['a'b']`,
					`["a"]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `['a'b']`},
				},
				{
					`$['a\\'b']`,
					`["a"]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `['a\\'b']`},
				},
				{
					`$['ab\']`,
					`["a"]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `['ab\']`},
				},
				{
					`$.[a]`,
					`{"a":1}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `.[a]`},
				},
				{
					`$[`,
					`["first","second","third"]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[`},
				},
				{
					`$[0`,
					`["first","second","third"]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[0`},
				},
				{
					`$[]`,
					`["first","second","third"]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[]`},
				},
				{
					`$[a]`,
					`["first","second","third"]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[a]`},
				},
				{
					`$[0,]`,
					`["first","second","third"]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[0,]`},
				},
				{
					`$[0,a]`,
					`["first","second","third"]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[0,a]`},
				},
				{
					`$[0,10000000000000000000,]`,
					`["first","second","third"]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[0,10000000000000000000,]`},
				},
				{
					`$[a:1]`,
					`["first","second","third"]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[a:1]`},
				},
				{
					`$[0:a]`,
					`["first","second","third"]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[0:a]`},
				},
				{
					`$[0:10000000000000000000:a]`,
					`["first","second","third"]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[0:10000000000000000000:a]`},
				},
				{
					`$[?()]`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?()]`},
				},
				{
					`$[?@a]`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?@a]`},
				},
				{
					`$[?(@.a!!=1)]`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a!!=1)]`},
				},
				{
					`$[?(@.a!=)]`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a!=)]`},
				},
				{
					`$[?(@.a<=)]`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a<=)]`},
				},
				{
					`$[?(@.a<)]`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a<)]`},
				},
				{
					`$[?(@.a>=)]`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a>=)]`},
				},
				{
					`$[?(@.a>)]`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a>)]`},
				},
				{
					`$[?(!=@.a)]`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(!=@.a)]`},
				},
				{
					`$[?(<=@.a)]`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(<=@.a)]`},
				},
				{
					`$[?(<@.a)]`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(<@.a)]`},
				},
				{
					`$[?(>=@.a)]`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(>=@.a)]`},
				},
				{
					`$[?(>@.a)]`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(>@.a)]`},
				},
				{
					`$[?(@.a===1)]`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a===1)]`},
				},
				{
					`$[?(@.a=='abc`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a=='abc`},
				},
				{
					`$[?(@.a=="abc`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a=="abc`},
				},
				{
					`$[?(@.a==["b"])]`,
					`[{"a":["b"]}]`,
					``,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `["b"])]`},
				},
				{
					`$[?(@[0:1]==[1])]`,
					`[[1,2,3],[1],[2,3],1,2]`,
					``,
					ErrorInvalidSyntax{12, `the omission of '$' allowed only at the beginning`, `[1])]`},
				},
				{
					`$[?(@.*==[1,2])]`,
					`[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
					``,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `[1,2])]`},
				},
				{
					`$[?(@.*==['1','2'])]`,
					`[[1,2],[2,3],[1],[2],[1,2,3],1,2,3]`,
					``,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `['1','2'])]`},
				},
				{
					`$[?((@.a<2)==false)]`,
					`[{"a":1},{"a":2},{"a":3}]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?((@.a<2)==false)]`},
				},
				{
					`$[?((@.a<2)==true)]`,
					`[{"a":1},{"a":2},{"a":3}]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?((@.a<2)==true)]`},
				},
				{
					`$[?((@.a<2)==1)]`,
					`[{"a":1},{"a":2},{"a":3}]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?((@.a<2)==1)]`},
				},
				{
					`$[?(false)]`,
					`[0,1,false,true,null,{},[]]`,
					`[]`,
					ErrorInvalidSyntax{4, `the omission of '$' allowed only at the beginning`, `false)]`},
				},
				{
					`$[?(true)]`,
					`[0,1,false,true,null,{},[]]`,
					`[]`,
					ErrorInvalidSyntax{4, `the omission of '$' allowed only at the beginning`, `true)]`},
				},
				{
					`$[?(null)]`,
					`[0,1,false,true,null,{},[]]`,
					`[]`,
					ErrorInvalidSyntax{4, `the omission of '$' allowed only at the beginning`, `null)]`},
				},
				{
					`$[?(@.a>1 && )]`,
					`[{"a":1},{"a":2},{"a":3}]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a>1 && )]`},
				},
				{
					`$[?(@.a>1 || )]`,
					`[{"a":1},{"a":2},{"a":3}]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a>1 || )]`},
				},
				{
					`$[?( && @.a>1 )]`,
					`[{"a":1},{"a":2},{"a":3}]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?( && @.a>1 )]`},
				},
				{
					`$[?( || @.a>1 )]`,
					`[{"a":1},{"a":2},{"a":3}]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?( || @.a>1 )]`},
				},
				{
					`$[?(@.a>1 && false)]`,
					`[{"a":1},{"a":2},{"a":3}]`,
					``,
					ErrorInvalidSyntax{13, `the omission of '$' allowed only at the beginning`, `false)]`},
				},
				{
					`$[?(@.a>1 && true)]`,
					`[{"a":1},{"a":2},{"a":3}]`,
					``,
					ErrorInvalidSyntax{13, `the omission of '$' allowed only at the beginning`, `true)]`},
				},
				{
					`$[?(@.a>1 || false)]`,
					`[{"a":1},{"a":2},{"a":3}]`,
					``,
					ErrorInvalidSyntax{13, `the omission of '$' allowed only at the beginning`, `false)]`},
				},
				{
					`$[?(@.a>1 || true)]`,
					`[{"a":1},{"a":2},{"a":3}]`,
					``,
					ErrorInvalidSyntax{13, `the omission of '$' allowed only at the beginning`, `true)]`},
				},
				{
					`$[?(@.a>1 && ())]`,
					`[{"a":1},{"a":2},{"a":3}]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a>1 && ())]`},
				},
				{
					`$[?(((@.a>1)))]`,
					`[{"a":1},{"a":2},{"a":3}]`,
					`[{"a":2},{"a":3}]`,
				},
				{
					`$[?((@.a>1 )]`,
					`[{"a":1},{"a":2},{"a":3}]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?((@.a>1 )]`},
				},
				{
					`$[?((@.a>1`,
					`[{"a":1},{"a":2},{"a":3}]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?((@.a>1`},
				},
				{
					`$[?(!(@.a==2))]`,
					`[{"a":1.9999},{"a":2},{"a":2.0001},{"a":"2"},{"a":true},{"a":{}},{"a":[]},{"a":["b"]},{"a":{"a":"value"}},{"b":"value"}]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(!(@.a==2))]`},
				},
				{
					`$[?(!(@.a<2))]`,
					`[{"a":1.9999},{"a":2},{"a":2.0001},{"a":"2"},{"a":true},{"a":{}},{"a":[]},{"a":["b"]},{"a":{"a":"value"}},{"b":"value"}]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(!(@.a<2))]`},
				},
				{
					`$[?(@.a==fAlse)]`,
					`[{"a":false}]`,
					`[{"a":false}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `fAlse)]`},
				},
				{
					`$[?(@.a==faLse)]`,
					`[{"a":false}]`,
					`[{"a":false}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `faLse)]`},
				},
				{
					`$[?(@.a==falSe)]`,
					`[{"a":false}]`,
					`[{"a":false}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `falSe)]`},
				},
				{
					`$[?(@.a==falsE)]`,
					`[{"a":false}]`,
					`[{"a":false}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `falsE)]`},
				},
				{
					`$[?(@.a==FaLse)]`,
					`[{"a":false}]`,
					`[{"a":false}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `FaLse)]`},
				},
				{
					`$[?(@.a==FalSe)]`,
					`[{"a":false}]`,
					`[{"a":false}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `FalSe)]`},
				},
				{
					`$[?(@.a==FalsE)]`,
					`[{"a":false}]`,
					`[{"a":false}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `FalsE)]`},
				},
				{
					`$[?(@.a==FaLSE)]`,
					`[{"a":false}]`,
					`[{"a":false}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `FaLSE)]`},
				},
				{
					`$[?(@.a==FAlSE)]`,
					`[{"a":false}]`,
					`[{"a":false}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `FAlSE)]`},
				},
				{
					`$[?(@.a==FALsE)]`,
					`[{"a":false}]`,
					`[{"a":false}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `FALsE)]`},
				},
				{
					`$[?(@.a==FALSe)]`,
					`[{"a":false}]`,
					`[{"a":false}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `FALSe)]`},
				},
				{
					`$[?(@.a==tRue)]`,
					`[{"a":true}]`,
					`[{"a":true}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `tRue)]`},
				},
				{
					`$[?(@.a==trUe)]`,
					`[{"a":true}]`,
					`[{"a":true}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `trUe)]`},
				},
				{
					`$[?(@.a==truE)]`,
					`[{"a":true}]`,
					`[{"a":true}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `truE)]`},
				},
				{
					`$[?(@.a==TrUe)]`,
					`[{"a":true}]`,
					`[{"a":true}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `TrUe)]`},
				},
				{
					`$[?(@.a==TruE)]`,
					`[{"a":true}]`,
					`[{"a":true}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `TruE)]`},
				},
				{
					`$[?(@.a==TrUE)]`,
					`[{"a":true}]`,
					`[{"a":true}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `TrUE)]`},
				},
				{
					`$[?(@.a==TRuE)]`,
					`[{"a":true}]`,
					`[{"a":true}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `TRuE)]`},
				},
				{
					`$[?(@.a==TRUe)]`,
					`[{"a":true}]`,
					`[{"a":true}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `TRUe)]`},
				},
				{
					`$[?(@.a==nUll)]`,
					`[{"a":null}]`,
					`[{"a":null}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `nUll)]`},
				},
				{
					`$[?(@.a==nuLl)]`,
					`[{"a":null}]`,
					`[{"a":null}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `nuLl)]`},
				},
				{
					`$[?(@.a==nulL)]`,
					`[{"a":null}]`,
					`[{"a":null}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `nulL)]`},
				},
				{
					`$[?(@.a==NuLl)]`,
					`[{"a":null}]`,
					`[{"a":null}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `NuLl)]`},
				},
				{
					`$[?(@.a==NulL)]`,
					`[{"a":null}]`,
					`[{"a":null}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `NulL)]`},
				},
				{
					`$[?(@.a==NuLL)]`,
					`[{"a":null}]`,
					`[{"a":null}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `NuLL)]`},
				},
				{
					`$[?(@.a==NUlL)]`,
					`[{"a":null}]`,
					`[{"a":null}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `NUlL)]`},
				},
				{
					`$[?(@.a==NULl)]`,
					`[{"a":null}]`,
					`[{"a":null}]`,
					ErrorInvalidSyntax{9, `the omission of '$' allowed only at the beginning`, `NULl)]`},
				},
				{
					`$[?(@=={"k":"v"})]`,
					`{}`,
					``,
					ErrorInvalidSyntax{7, `the omission of '$' allowed only at the beginning`, `{"k":"v"})]`},
				},
				{
					`$[?(@.a=~/abc)]`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a=~/abc)]`},
				},
				{
					`$[?(@.a=~///)]`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a=~///)]`},
				},
				{
					`$[?(@.a=~s/a/b/)]`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a=~s/a/b/)]`},
				},
				{
					`$[?(@.a=~@abc@)]`,
					`[]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a=~@abc@)]`},
				},
				{
					`$[?(a=~/123/)]`,
					`[{"a":"123"},{"a":123}]`,
					``,
					ErrorInvalidSyntax{4, `the omission of '$' allowed only at the beginning`, `a=~/123/)]`},
				},
				{
					`$[?(@.a=2)]`,
					`[{"a":2}]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a=2)]`},
				},
				{
					`$[?(@.a<>2)]`,
					`[{"a":2}]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a<>2)]`},
				},
				{
					`$[?(@.a=<2)]`,
					`[{"a":2}]`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a=<2)]`},
				},
				{
					`$[?(@.a),?(@.b)]`,
					`{}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a),?(@.b)]`},
				},
				{
					`$[?(@.a & @.b)]`,
					`{}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a & @.b)]`},
				},
				{
					`$[?(@.a | @.b)]`,
					`{}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[?(@.a | @.b)]`},
				},
				{
					`$[()]`,
					`{}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[()]`},
				},
				{
					`$[(`,
					`{}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[(`},
				},
				{
					`$[(]`,
					`{}`,
					``,
					ErrorInvalidSyntax{1, `unrecognized input`, `[(]`},
				},
				{
					`$.func(`,
					`{}`,
					``,
					ErrorInvalidSyntax{6, `unrecognized input`, `(`},
				},
			},
		},
		{
			`Invalid argument format`,
			[][]interface{}{
				{
					`$[10000000000000000000]`,
					`["first","second","third"]`,
					``,
					ErrorInvalidArgument{
						`10000000000000000000`,
						fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
					},
				},
				{
					`$[0,10000000000000000000]`,
					`["first","second","third"]`,
					``,
					ErrorInvalidArgument{
						`10000000000000000000`,
						fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
					},
				},
				{
					`$[10000000000000000000:1]`,
					`["first","second","third"]`,
					``,
					ErrorInvalidArgument{
						`10000000000000000000`,
						fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
					},
				},
				{
					`$[1:10000000000000000000]`,
					`["first","second","third"]`,
					``,
					ErrorInvalidArgument{
						`10000000000000000000`,
						fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
					},
				},
				{
					`$[0:3:10000000000000000000]`,
					`["first","second","third"]`,
					``,
					ErrorInvalidArgument{
						`10000000000000000000`,
						fmt.Errorf(`strconv.Atoi: parsing "10000000000000000000": value out of range`),
					},
				},
				{
					`$[?(@.a==1e1abc)]`,
					`{}`,
					``,
					ErrorInvalidArgument{
						`1e1abc`,
						fmt.Errorf(`strconv.ParseFloat: parsing "1e1abc": invalid syntax`),
					},
				},
			},
		},
		{
			`Not supported`,
			[][]interface{}{
				{
					`$[(command)]`,
					`{}`,
					``,
					ErrorNotSupported{`script`, `[(command)]`},
				},
			},
		},
	}

	for _, testGroup := range testGroups {
		for _, testCase := range testGroup.testCases {
			jsonPath := testCase[0].(string)
			srcJSON := testCase[1].(string)
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
			`filter`,
			[][]interface{}{
				{
					`$[?(@.a > 123)].a`,
					`[{"a":123.456}]`,
					`[123.456]`,
				},
				{
					`$[?(@.a > 123.46)].a`,
					`[{"a":123.456}]`,
					`[]`,
					ErrorNoneMatched{`[?(@.a > 123.46)].a`},
				},
				{
					`$[?(@.a > 122)].a`,
					`[{"a":123}]`,
					`[123]`,
				},
				{
					`$[?(123 < @.a)].a`,
					`[{"a":123.456}]`,
					`[123.456]`,
				},
				{
					`$[?(@.a==-0.123e2)]`,
					`[{"a":-12.3,"b":1},{"a":-0.123e2,"b":2},{"a":-0.123},{"a":-12},{"a":12.3},{"a":2},{"a":"-0.123e2"}]`,
					`[{"a":-12.3,"b":1},{"a":-0.123e2,"b":2}]`,
				},
				{
					`$[?(@.a==11)]`,
					`[{"a":10.999},{"a":11.00},{"a":11.10}]`,
					`[{"a":11.00}]`,
				},
			},
		},
	}

	for _, testGroup := range testGroups {
		for _, testCase := range testGroup.testCases {
			jsonPath := testCase[0].(string)
			srcJSON := testCase[1].(string)
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
			`filter-function`,
			[][]interface{}{
				{
					`$.*.twice()`,
					`[123.456,256]`,
					`[246.912,512]`,
					map[string]func(interface{}) (interface{}, error){
						`twice`: twiceFunc,
					},
					map[string]func([]interface{}) (interface{}, error){},
				},
				{
					`$.*.twice().twice()`,
					`[123.456,256]`,
					`[493.824,1024]`,
					map[string]func(interface{}) (interface{}, error){
						`twice`: twiceFunc,
					},
					map[string]func([]interface{}) (interface{}, error){},
				},
				{
					`$.*.twice().quarter()`,
					`[123.456,256]`,
					`[61.728,128]`,
					map[string]func(interface{}) (interface{}, error){
						`twice`:   twiceFunc,
						`quarter`: quarterFunc,
					},
					map[string]func([]interface{}) (interface{}, error){},
				},
				{
					`$.*.quarter().twice()`,
					`[123.456,256]`,
					`[61.728,128]`,
					map[string]func(interface{}) (interface{}, error){
						`twice`:   twiceFunc,
						`quarter`: quarterFunc,
					},
					map[string]func([]interface{}) (interface{}, error){},
				},
				{
					`$[?(@.twice())]`,
					`[123.456,256]`,
					`[123.456,256]`,
					map[string]func(interface{}) (interface{}, error){
						`twice`: twiceFunc,
					},
					map[string]func([]interface{}) (interface{}, error){},
				},
				{
					`$[?(@.twice() == 512)]`,
					`[123.456,256]`,
					`[256]`,
					map[string]func(interface{}) (interface{}, error){
						`twice`: twiceFunc,
					},
					map[string]func([]interface{}) (interface{}, error){},
				},
				{
					`$[?(512 != @.twice())]`,
					`[123.456,256]`,
					`[123.456]`,
					map[string]func(interface{}) (interface{}, error){
						`twice`: twiceFunc,
					},
					map[string]func([]interface{}) (interface{}, error){},
				},
				{
					`$[?(@.twice() == $[0].twice())]`,
					`[123.456,256]`,
					`[123.456]`,
					map[string]func(interface{}) (interface{}, error){
						`twice`: twiceFunc,
					},
					map[string]func([]interface{}) (interface{}, error){},
				},
			},
		},
		{
			`aggregate-function`,
			[][]interface{}{
				{
					`$.*.max()`,
					`[122.345,123.45,123.456]`,
					`[123.456]`,
					map[string]func(interface{}) (interface{}, error){},
					map[string]func([]interface{}) (interface{}, error){
						`max`: maxFunc,
					},
				},
				{
					`$.*.max().max()`,
					`[122.345,123.45,123.456]`,
					`[123.456]`,
					map[string]func(interface{}) (interface{}, error){},
					map[string]func([]interface{}) (interface{}, error){
						`max`: maxFunc,
					},
				},
				{
					`$.*.max().min()`,
					`[122.345,123.45,123.456]`,
					`[123.456]`,
					map[string]func(interface{}) (interface{}, error){},
					map[string]func([]interface{}) (interface{}, error){
						`max`: maxFunc,
						`min`: minFunc,
					},
				},
				{
					`$.*.min().max()`,
					`[122.345,123.45,123.456]`,
					`[122.345]`,
					map[string]func(interface{}) (interface{}, error){},
					map[string]func([]interface{}) (interface{}, error){
						`max`: maxFunc,
						`min`: minFunc,
					},
				},
				{
					`$[?(@.max())]`,
					`[122.345,123.45,123.456]`,
					`[122.345,123.45,123.456]`,
					map[string]func(interface{}) (interface{}, error){},
					map[string]func([]interface{}) (interface{}, error){
						`max`: maxFunc,
					},
				},
				{
					`$[?(@.max() == 123.45)]`,
					`[122.345,123.45,123.456]`,
					`[123.45]`,
					map[string]func(interface{}) (interface{}, error){},
					map[string]func([]interface{}) (interface{}, error){
						`max`: maxFunc,
					},
				},
				{
					`$[?(123.45 != @.max())]`,
					`[122.345,123.45,123.456]`,
					`[122.345,123.456]`,
					map[string]func(interface{}) (interface{}, error){},
					map[string]func([]interface{}) (interface{}, error){
						`max`: maxFunc,
					},
				},
				{
					`$[?(@.max() != 123.45)]`,
					`[[122.345,123.45,123.456],[122.345,123.45]]`,
					`[[122.345,123.45,123.456]]`,
					map[string]func(interface{}) (interface{}, error){},
					map[string]func([]interface{}) (interface{}, error){
						`max`: maxFunc,
					},
				},
				{
					`$[?(@.max() == $[1].max())]`,
					`[[122.345,123.45,123.456],[122.345,123.45]]`,
					`[[122.345,123.45]]`,
					map[string]func(interface{}) (interface{}, error){},
					map[string]func([]interface{}) (interface{}, error){
						`max`: maxFunc,
					},
				},
			},
		},
		{
			`aggregate-filter-mix`,
			[][]interface{}{
				{
					`$.*.max().twice()`,
					`[122.345,123.45,123.456]`,
					`[246.912]`,
					map[string]func(interface{}) (interface{}, error){
						`twice`: twiceFunc,
					},
					map[string]func([]interface{}) (interface{}, error){
						`max`: maxFunc,
					},
				},
				{
					`$.*.twice().max()`,
					`[122.345,123.45,123.456]`,
					`[246.912]`,
					map[string]func(interface{}) (interface{}, error){
						`twice`: twiceFunc,
					},
					map[string]func([]interface{}) (interface{}, error){
						`max`: maxFunc,
					},
				},
			},
		},
		{
			`filter-error`,
			[][]interface{}{
				{
					`$.errFilter()`,
					`[122.345,123.45,123.456]`,
					``,
					map[string]func(interface{}) (interface{}, error){
						`errFilter`: errFilterFunc,
					},
					map[string]func([]interface{}) (interface{}, error){},
					ErrorFunctionFailed{function: `.errFilter()`, err: fmt.Errorf(`filter error`)},
				},
				{
					`$.*.errFilter()`,
					`[122.345,123.45,123.456]`,
					``,
					map[string]func(interface{}) (interface{}, error){
						`errFilter`: errFilterFunc,
					},
					map[string]func([]interface{}) (interface{}, error){},
					ErrorNoneMatched{path: `.*.errFilter()`},
				},
				{
					`$.*.max().errFilter()`,
					`[122.345,123.45,123.456]`,
					``,
					map[string]func(interface{}) (interface{}, error){
						`errFilter`: errFilterFunc,
					},
					map[string]func([]interface{}) (interface{}, error){
						`max`: maxFunc,
					},
					ErrorFunctionFailed{function: `.errFilter()`, err: fmt.Errorf(`filter error`)},
				},
				{
					`$.*.twice().errFilter()`,
					`[122.345,123.45,123.456]`,
					``,
					map[string]func(interface{}) (interface{}, error){
						`errFilter`: errFilterFunc,
						`twice`:     twiceFunc,
					},
					map[string]func([]interface{}) (interface{}, error){},
					ErrorNoneMatched{path: `.*.twice().errFilter()`},
				}, {
					`$.errFilter().twice()`,
					`[122.345,123.45,123.456]`,
					``,
					map[string]func(interface{}) (interface{}, error){
						`errFilter`: errFilterFunc,
						`twice`:     twiceFunc,
					},
					map[string]func([]interface{}) (interface{}, error){},
					ErrorFunctionFailed{function: `.errFilter()`, err: fmt.Errorf(`filter error`)},
				},
				{
					`$.*.errFilter().twice()`,
					`[122.345,123.45,123.456]`,
					``,
					map[string]func(interface{}) (interface{}, error){
						`errFilter`: errFilterFunc,
						`twice`:     twiceFunc,
					},
					map[string]func([]interface{}) (interface{}, error){},
					ErrorNoneMatched{path: `.*.errFilter().twice()`},
				},
				{
					`$.*.max().errFilter().twice()`,
					`[122.345,123.45,123.456]`,
					``,
					map[string]func(interface{}) (interface{}, error){
						`errFilter`: errFilterFunc,
						`twice`:     twiceFunc,
					},
					map[string]func([]interface{}) (interface{}, error){
						`max`: maxFunc,
					},
					ErrorFunctionFailed{function: `.errFilter()`, err: fmt.Errorf(`filter error`)},
				},
			},
		},
		{
			`aggregate-error`,
			[][]interface{}{
				{
					`$.*.errAggregate()`,
					`[122.345,123.45,123.456]`,
					``,
					map[string]func(interface{}) (interface{}, error){},
					map[string]func([]interface{}) (interface{}, error){
						`errAggregate`: errAggregateFunc,
					},
					ErrorFunctionFailed{function: `.errAggregate()`, err: fmt.Errorf(`aggregate error`)},
				},
				{
					`$.*.max().errAggregate()`,
					`[122.345,123.45,123.456]`,
					``,
					map[string]func(interface{}) (interface{}, error){},
					map[string]func([]interface{}) (interface{}, error){
						`errAggregate`: errAggregateFunc,
						`max`:          maxFunc,
					},
					ErrorFunctionFailed{function: `.errAggregate()`, err: fmt.Errorf(`aggregate error`)},
				},
				{
					`$.*.twice().errAggregate()`,
					`[122.345,123.45,123.456]`,
					``,
					map[string]func(interface{}) (interface{}, error){
						`twice`: twiceFunc,
					},
					map[string]func([]interface{}) (interface{}, error){
						`errAggregate`: errAggregateFunc,
					},
					ErrorFunctionFailed{function: `.errAggregate()`, err: fmt.Errorf(`aggregate error`)},
				},
				{
					`$.*.errAggregate().twice()`,
					`[122.345,123.45,123.456]`,
					``,
					map[string]func(interface{}) (interface{}, error){
						`twice`: twiceFunc,
					},
					map[string]func([]interface{}) (interface{}, error){
						`errAggregate`: errAggregateFunc,
					},
					ErrorFunctionFailed{function: `.errAggregate()`, err: fmt.Errorf(`aggregate error`)},
				},
				{
					`$.*.max().errAggregate().twice()`,
					`[122.345,123.45,123.456]`,
					``,
					map[string]func(interface{}) (interface{}, error){
						`twice`: twiceFunc,
					},
					map[string]func([]interface{}) (interface{}, error){
						`errAggregate`: errAggregateFunc,
						`max`:          maxFunc,
					},
					ErrorFunctionFailed{function: `.errAggregate()`, err: fmt.Errorf(`aggregate error`)},
				},
			},
		},
		{
			`function-syntax-check`,
			[][]interface{}{
				{
					`$.*.TWICE()`,
					`[123.456,256]`,
					`[246.912,512]`,
					map[string]func(interface{}) (interface{}, error){
						`TWICE`: twiceFunc,
					},
					map[string]func([]interface{}) (interface{}, error){},
				},
				{
					`$.*.--()`,
					`[123.456,256]`,
					`[246.912,512]`,
					map[string]func(interface{}) (interface{}, error){
						`--`: twiceFunc,
					},
					map[string]func([]interface{}) (interface{}, error){},
				},
				{
					`$.*.__()`,
					`[123.456,256]`,
					`[246.912,512]`,
					map[string]func(interface{}) (interface{}, error){
						`__`: twiceFunc,
					},
					map[string]func([]interface{}) (interface{}, error){},
				},
			},
		},
	}

	for _, testGroup := range testGroups {
		for _, testCase := range testGroup.testCases {
			jsonPath := testCase[0].(string)
			srcJSON := testCase[1].(string)
			expectedJSON := testCase[2].(string)
			filterFunctions := testCase[3].(map[string]func(interface{}) (interface{}, error))
			aggregateFunctions := testCase[4].(map[string]func([]interface{}) (interface{}, error))
			var expectedError error
			if len(testCase) > 5 {
				expectedError = testCase[5].(error)
			}
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
