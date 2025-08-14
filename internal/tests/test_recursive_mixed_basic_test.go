package tests

import (
	"testing"
)

func TestRecursive_BasicCases(t *testing.T) {
	tests := []TestCase{
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
			jsonpath:    `$.x..a`,
			inputJSON:   `{"a":"b","c":{"a":"d"},"e":["f",{"g":{"a":"h"}}]}`,
			expectedErr: createErrorMemberNotExist(`.x`),
		},
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
	}

	runTestCases(t, "RecursiveDescentDeletedCases", tests)
}
