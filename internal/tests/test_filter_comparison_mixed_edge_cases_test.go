package tests

import (
	"testing"
)

func TestFilter_ComparisonEdgeCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[?(1 > @.a)]`,
			inputJSON:    `[{"a":-9999999},{"a":0.999999},{"a":1.0000000},{"a":1.0000001},{"a":2},{"a":"0.9"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
			expectedJSON: `[{"a":-9999999},{"a":0.999999}]`,
		},
		{
			jsonpath:    `$[?(1.00001 >= @.a)]`,
			inputJSON:   `[{"a":1.00002},{"a":2,"b":4},{"a":"0.9"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
			expectedErr: createErrorMemberNotExist(`[?(1.00001 >= @.a)]`),
		},
		{
			jsonpath:    `$[?(1.000001 <= @.a)]`,
			inputJSON:   `[{"a":0},{"a":1},{"a":1.0000009},{"a":"2"},{"a":{}},{"a":[]},{"a":true},{"a":null},{"b":"c"}]`,
			expectedErr: createErrorMemberNotExist(`[?(1.000001 <= @.a)]`),
		},

		{
			jsonpath:     `$[?(@.a=='~!@#$%^&*()-_=+[]\\{}|;\':",./<>?')]`,
			inputJSON:    `[{"a":"~!@#$%^&*()-_=+[]\\{}|;':\",./<>?"}]`,
			expectedJSON: `[{"a":"~!@#$%^\u0026*()-_=+[]\\{}|;':\",./\u003c\u003e?"}]`,
		},
		{
			jsonpath:     `$[?(@.a=='a\bb')]`,
			inputJSON:    `[{"a":"a\bb"},{"b":1}]`,
			expectedJSON: `[{"a":"a\bb"}]`,
		},

		{
			jsonpath:     `$[?(@.a==010)]`,
			inputJSON:    `[{"a":10},{"a":0},{"a":"010"},{"a":"10"}]`,
			expectedJSON: `[{"a":10}]`,
		},

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
	}

	runTestCases(t, "TestFilter_ComparisonEdgeCases", testCases)
}
