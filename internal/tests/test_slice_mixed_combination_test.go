package tests

import (
	"testing"
)

func TestSliceNestedMemberAccess(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:2].a.b`,
		inputJSON:   `[{"b": 1}]`,
		expectedErr: createErrorMemberNotExist(`.a`),
	}
	runTestCase(t, testCase, "TestSliceNestedMemberAccess")
}

func TestSliceWithRecursiveDescent(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[1:3]..value`,
			inputJSON:    `[{"x":1}, {"data":{"value":10}}, {"nested":{"deep":{"value":20}}}, {"value":30}]`,
			expectedJSON: `[10,20]`,
		},
		{
			jsonpath:     `$[0:2]..id`,
			inputJSON:    `[{"id":1,"children":[{"id":2}]}, {"id":3,"data":{"id":4}}, {"id":5}]`,
			expectedJSON: `[1,2,3,4]`,
		},
		{
			jsonpath:    `$[1:2]..missing`,
			inputJSON:   `[{"a":1}, {"b":2}, {"c":3}]`,
			expectedErr: createErrorMemberNotExist(`missing`),
		},
	}
	runTestCases(t, "TestSliceWithRecursiveDescent", tests)
}

func TestSliceWithFilter(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[1:4][?(@ > 80)]`,
			inputJSON:    `[{"score":70}, {"score":85}, {"score":75}, {"score":90}, {"score":95}]`,
			expectedJSON: `[85,90]`,
		},
		{
			jsonpath:     `$[0:3][?(@ == true)]`,
			inputJSON:    `[{"active":false}, {"active":true}, {"active":true}, {"active":true}]`,
			expectedJSON: `[true,true]`,
		},
		{
			jsonpath:    `$[1:3][?(@.missing)]`,
			inputJSON:   `[{"a":{"missing":1}}, {"b":2}, {"c":3}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.missing)]`),
		},
	}
	runTestCases(t, "TestSliceWithFilter", tests)
}

func TestSliceWithMultipleBracketAccess(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[1:3][0]`,
			inputJSON:    `[["a","b"], ["c","d"], ["e","f"], ["g","h"]]`,
			expectedJSON: `["c","e"]`,
		},
		{
			jsonpath:     `$[0:2]['name']`,
			inputJSON:    `[{"name":"Alice","age":30}, {"name":"Bob","age":25}, {"name":"Charlie","age":35}]`,
			expectedJSON: `["Alice","Bob"]`,
		},
		{
			jsonpath:    `$[1:3][5]`,
			inputJSON:   `[["a"], ["b","c"], ["d"]]`,
			expectedErr: createErrorMemberNotExist(`[5]`),
		},
	}
	runTestCases(t, "TestSliceWithMultipleBracketAccess", tests)
}

func TestSliceWithDotNotationChain(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[1:3].user.name`,
			inputJSON:    `[{"user":{"name":"Alice"}}, {"user":{"name":"Bob"}}, {"user":{"name":"Charlie"}}, {"user":{"name":"David"}}]`,
			expectedJSON: `["Bob","Charlie"]`,
		},
		{
			jsonpath:     `$[0:2].data.values`,
			inputJSON:    `[{"data":{"values":[1,2,3]}}, {"data":{"values":[4,5,6]}}, {"data":{"values":[7,8,9]}}]`,
			expectedJSON: `[[1,2,3],[4,5,6]]`,
		},
		{
			jsonpath:    `$[1:3].missing.property`,
			inputJSON:   `[{"a":1}, {"b":2}, {"c":3}]`,
			expectedErr: createErrorMemberNotExist(`.missing`),
		},
	}
	runTestCases(t, "TestSliceWithDotNotationChain", tests)
}

func TestSliceWithWildcardAccess(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[1:3].*`,
			inputJSON:    `[{"a":1,"b":2}, {"c":3,"d":4}, {"e":5,"f":6}, {"g":7}]`,
			expectedJSON: `[3,4,5,6]`,
		},
		{
			jsonpath:     `$[0:2][*]`,
			inputJSON:    `[[10,20,30], [40,50], [60,70,80,90]]`,
			expectedJSON: `[10,20,30,40,50]`,
		},
		{
			jsonpath:    `$[2:4].*`,
			inputJSON:   `[{"a":1}, {"b":2}]`,
			expectedErr: createErrorMemberNotExist(`[2:4]`),
		},
	}
	runTestCases(t, "TestSliceWithWildcardAccess", tests)
}

func TestSliceComplexCombination(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[1:3]..items[?(@.available == true)].name`,
			inputJSON:    `[{"items":[{"name":"item1","available":false}]}, {"items":[{"name":"item2","available":true}]}, {"deep":{"items":[{"name":"item3","available":true}]}}]`,
			expectedJSON: `["item2","item3"]`,
		},
		{
			jsonpath:     `$[0:2].categories[*].products[?(@.price < 100)].id`,
			inputJSON:    `[{"categories":[{"products":[{"id":"p1","price":50},{"id":"p2","price":150}]}]}, {"categories":[{"products":[{"id":"p3","price":75}]}]}]`,
			expectedJSON: `["p1","p3"]`,
		},
		{
			jsonpath:    `$[1:3]..data[?(@.type == "special")].value`,
			inputJSON:   `[{"data":{"type":"normal","value":1}}, {"data":{"type":"other","value":2}}, {"data":{"type":"basic","value":3}}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.type == "special")]`),
		},
	}
	runTestCases(t, "TestSliceComplexCombination", tests)
}

func TestSliceWithUnion(t *testing.T) {
	tests := []TestCase{
		{
			jsonpath:     `$[1:3]['name','age']`,
			inputJSON:    `[{"name":"Alice","age":30}, {"name":"Bob","age":25}, {"name":"Charlie","age":35}, {"name":"David","age":40}]`,
			expectedJSON: `["Bob",25,"Charlie",35]`,
		},
		{
			jsonpath:     `$[0:2][0,2]`,
			inputJSON:    `[["a","b","c"], ["d","e","f"], ["g","h","i"]]`,
			expectedJSON: `["a","c","d","f"]`,
		},
		{
			jsonpath:    `$[1:2]['missing','other']`,
			inputJSON:   `[{"a":1}, {"b":2}, {"c":3}]`,
			expectedErr: createErrorMemberNotExist(`['missing','other']`),
		},
	}
	runTestCases(t, "TestSliceWithUnion", tests)
}
