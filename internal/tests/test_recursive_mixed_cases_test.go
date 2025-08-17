package tests

import (
	"testing"

	syntax "github.com/AsaiYusuke/jsonpath/v2/internal/syntax"
)

func TestRecursiveBasic_ConditionalRecursive(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$.z[?($..x)]`,
			inputJSON:    `{"x":[], "y":{"x":[]}, "z":[{"b":1},{"b":2},{"b":3}]}`,
			expectedJSON: `[{"b":1},{"b":2},{"b":3}]`,
		},
		{
			jsonpath:     `$..[?(@.a==2)]`,
			inputJSON:    `{"a":2,"x":[{"a":2},{"b":{"a":2}},{"a":{"a":2}},[{"a":2}]]}`,
			expectedJSON: `[{"a":2},{"a":2},{"a":2},{"a":2}]`,
		},
	}

	runTestCases(t, "TestRecursiveBasic_ConditionalRecursive", tests)
}

func TestRecursiveBasic_FilterWithRecursive(t *testing.T) {
	tests := []TestCase{
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
	}

	runTestCases(t, "TestRecursiveBasic_FilterWithRecursive", tests)
}

func TestRecursiveBasic_RecursiveChildSlicesGrow(t *testing.T) {
	tests := []TestCase{
		// Map case 1: number of child elements exceeds the initial capacity 10 (11 elements).
		{
			jsonpath:     `$..*`,
			inputJSON:    `{"1":[],"2":[],"3":[],"4":[],"5":[],"6":[],"7":[],"8":[],"9":[],"10":[],"11":[]}`,
			expectedJSON: `[[],[],[],[],[],[],[],[],[],[],[]]`,
			//              1  2  3  4  5  6  7  8  9  10 11
		},
		// Map case 2: number of child elements exceeds twice the initial capacity 10 (21 elements).
		{
			jsonpath: `$..*`,
			inputJSON: `{
							"01":[],"02":[],"03":[],"04":[],"05":[],"06":[],"07":[],"08":[],"09":[],"10":[],
							"11":[],"12":[],"13":[],"14":[],"15":[],"16":[],"17":[],"18":[],"19":[],"20":[],
							"21":[]
						}`,
			expectedJSON: `[[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[]]`,
			//              1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18 19 20 21
		},
		// Array case 1: number of child elements exceeds the initial capacity 10 (11 elements).
		{
			jsonpath:     `$..*`,
			inputJSON:    `[{},{},{},{},{},{},{},{},{},{},{}]`,
			expectedJSON: `[{},{},{},{},{},{},{},{},{},{},{}]`,
			//              1  2  3  4  5  6  7  8  9  10 11
		},
		// Array case 2: number of child elements exceeds twice the initial capacity 10 (21 elements).
		{
			jsonpath:     `$..*`,
			inputJSON:    `[{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{}]`,
			expectedJSON: `[{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{}]`,
			//              1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18 19 20 21
		},
	}

	// Isolate from other tests that may have enlarged the pool
	runTestCasesSerial(t, "TestRecursiveBasic_RecursiveChildSlicesGrow", tests, func() {
		syntax.ResetNodeSliceSyncPool()
	})
}
