package tests

import (
	"testing"
)

func TestFilterSpecialCharacters(t *testing.T) {
	testCases := []TestCase{
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
			jsonpath:     `$[?(@['a<=b']<2.1)]`,
			inputJSON:    `[{"a<=b":1.9},{"a":2},{"a":2.1},{"b":3},{"a<=b":"test"}]`,
			expectedJSON: `[{"a\u003c=b":1.9}]`,
		},
		{
			jsonpath:     `$[?(@.a=="~!@#$%^&*()-_=+[]\\{}|;':\",./<>?")]`,
			inputJSON:    `[{"a":"~!@#$%^&*()-_=+[]\\{}|;':\",./<>?"}]`,
			expectedJSON: `[{"a":"~!@#$%^\u0026*()-_=+[]\\{}|;':\",./\u003c\u003e?"}]`,
		},
	}

	runTestCases(t, "TestFilterSpecialCharacters", testCases)
}
