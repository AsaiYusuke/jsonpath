package tests

import (
	"testing"
)

func TestWildcard_DotNotationArrayBasic(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*`,
		inputJSON:    `[[1],[2,3],123,"a",{"b":"c"},[0,1],null]`,
		expectedJSON: `[[1],[2,3],123,"a",{"b":"c"},[0,1],null]`,
	}
	runTestCase(t, testCase, "TestWildcard_DotNotationArrayBasic")
}

func TestWildcard_DotNotationArrayIndexAccess(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*[1]`,
		inputJSON:    `[[1],[2,3],[4,[5,6,7]]]`,
		expectedJSON: `[3,[5,6,7]]`,
	}
	runTestCase(t, testCase, "TestWildcard_DotNotationArrayIndexAccess")
}

func TestWildcard_DotNotationArrayPropertyAccess(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*.a`,
		inputJSON:    `[{"a":1},{"a":[2,3]}]`,
		expectedJSON: `[1,[2,3]]`,
	}
	runTestCase(t, testCase, "TestWildcard_DotNotationArrayPropertyAccess")
}

func TestWildcard_DotNotationArrayRecursive(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$..*`,
		inputJSON:    `[{"a":1},{"a":[2,3]},null,true]`,
		expectedJSON: `[{"a":1},{"a":[2,3]},null,true,1,[2,3],2,3]`,
	}
	runTestCase(t, testCase, "TestWildcard_DotNotationArrayRecursive")
}

func TestWildcard_ArrayMultiProperty(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*['a','b']`,
		inputJSON:    `[{"a":1,"b":2,"c":3},{"a":4,"b":5,"d":6}]`,
		expectedJSON: `[1,2,4,5]`,
	}
	runTestCase(t, testCase, "TestWildcard_ArrayMultiProperty")
}

func TestWildcard_ArrayNoRoot(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `*`,
		inputJSON:    `[[1],[2,3],123,"a",{"b":"c"},[0,1],null]`,
		expectedJSON: `[[1],[2,3],123,"a",{"b":"c"},[0,1],null]`,
	}
	runTestCase(t, testCase, "TestWildcard_ArrayNoRoot")
}

func TestWildcard_ObjectBasic(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*`,
		inputJSON:    `{"a":[1],"b":[2,3],"c":{"d":4}}`,
		expectedJSON: `[[1],[2,3],{"d":4}]`,
	}
	runTestCase(t, testCase, "TestWildcard_ObjectBasic")
}

func TestWildcard_ObjectRecursive(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$..*`,
		inputJSON:    `{"a":1,"b":[2,3],"c":{"d":4,"e":[5,6]}}`,
		expectedJSON: `[1,[2,3],{"d":4,"e":[5,6]},2,3,4,[5,6],5,6]`,
	}
	runTestCase(t, testCase, "TestWildcard_ObjectRecursive")
}

func TestWildcard_ObjectRecursiveBracket(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$..[*]`,
		inputJSON:    `{"a":1,"b":[2,3],"c":{"d":"e","f":[4,5]}}`,
		expectedJSON: `[1,[2,3],{"d":"e","f":[4,5]},2,3,"e",[4,5],4,5]`,
	}
	runTestCase(t, testCase, "TestWildcard_ObjectRecursiveBracket")
}

func TestWildcard_ObjectRecursiveBracketEmptyArray(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$..[*]`,
		inputJSON:   `[]`,
		expectedErr: createErrorMemberNotExist(`[*]`),
	}
	runTestCase(t, testCase, "TestWildcard_ObjectRecursiveBracketEmptyArray")
}

func TestWildcard_ObjectNoRoot(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `*`,
		inputJSON:    `{"a":[1],"b":[2,3],"c":{"d":4}}`,
		expectedJSON: `[[1],[2,3],{"d":4}]`,
	}
	runTestCase(t, testCase, "TestWildcard_ObjectNoRoot")
}

func TestWildcard_TwoWildcardsArray(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*.*`,
		inputJSON:    `[[1,2,3],[4,5,6]]`,
		expectedJSON: `[1,2,3,4,5,6]`,
	}
	runTestCase(t, testCase, "TestWildcard_TwoWildcardsArray")
}

func TestWildcard_TwoWildcardsNested(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*.a.*`,
		inputJSON:    `[{"a":[1]}]`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestWildcard_TwoWildcardsNested")
}

func TestWildcard_EmptyObject(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*`,
		inputJSON:   `{}`,
		expectedErr: createErrorMemberNotExist(`.*`),
	}
	runTestCase(t, testCase, "TestWildcard_EmptyObject")
}

func TestWildcard_EmptyArray(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*`,
		inputJSON:   `[]`,
		expectedErr: createErrorMemberNotExist(`.*`),
	}
	runTestCase(t, testCase, "TestWildcard_EmptyArray")
}

func TestWildcard_RecursiveEmptyObject(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$..*`,
		inputJSON:   `{}`,
		expectedErr: createErrorMemberNotExist(`*`),
	}
	runTestCase(t, testCase, "TestWildcard_RecursiveEmptyObject")
}

func TestWildcard_RecursiveEmptyArray(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$..*`,
		inputJSON:   `[]`,
		expectedErr: createErrorMemberNotExist(`*`),
	}
	runTestCase(t, testCase, "TestWildcard_RecursiveEmptyArray")
}

func TestWildcard_RecursiveSingleProp(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$..*`,
		inputJSON:    `{"a":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestWildcard_RecursiveSingleProp")
}

func TestWildcard_RecursiveMultiProp(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$..*`,
		inputJSON:    `{"b":2,"a":1}`,
		expectedJSON: `[1,2]`,
	}
	runTestCase(t, testCase, "TestWildcard_RecursiveMultiProp")
}

func TestWildcard_RecursiveNestedObject(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$..*`,
		inputJSON:    `{"a":{"b":2}}`,
		expectedJSON: `[{"b":2},2]`,
	}
	runTestCase(t, testCase, "TestWildcard_RecursiveNestedObject")
}

func TestWildcard_RecursiveNestedMultiProp(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$..*`,
		inputJSON:    `{"a":{"c":3,"b":2}}`,
		expectedJSON: `[{"b":2,"c":3},2,3]`,
	}
	runTestCase(t, testCase, "TestWildcard_RecursiveNestedMultiProp")
}

func TestWildcard_RecursiveArraySingle(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$..*`,
		inputJSON:    `[1]`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestWildcard_RecursiveArraySingle")
}

func TestWildcard_RecursiveArrayMulti(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$..*`,
		inputJSON:    `[2,1]`,
		expectedJSON: `[2,1]`,
	}
	runTestCase(t, testCase, "TestWildcard_RecursiveArrayMulti")
}

func TestWildcard_RecursiveArrayWithObject(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$..*`,
		inputJSON:    `[{"a":1}]`,
		expectedJSON: `[{"a":1},1]`,
	}
	runTestCase(t, testCase, "TestWildcard_RecursiveArrayWithObject")
}

func TestWildcard_RecursiveArrayWithMultiPropObject(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$..*`,
		inputJSON:    `[{"b":2,"a":1}]`,
		expectedJSON: `[{"a":1,"b":2},1,2]`,
	}
	runTestCase(t, testCase, "TestWildcard_RecursiveArrayWithMultiPropObject")
}

func TestWildcard_RecursiveNoRoot(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `..*`,
		inputJSON:    `{"a":1}`,
		expectedJSON: `[1]`,
	}
	runTestCase(t, testCase, "TestWildcard_RecursiveNoRoot")
}

func TestWildcard_ChildErrorObjectMissingMember(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*.a.b`,
		inputJSON:   `{"a":{"b":1}}`,
		expectedErr: createErrorMemberNotExist(`.a`),
	}
	runTestCase(t, testCase, "TestWildcard_ChildErrorObjectMissingMember")
}

func TestWildcard_ChildErrorArrayMissingMember(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*.a.b`,
		inputJSON:   `[{"b":1}]`,
		expectedErr: createErrorMemberNotExist(`.a`),
	}
	runTestCase(t, testCase, "TestWildcard_ChildErrorArrayMissingMember")
}

func TestWildcard_ChildErrorObjectTypeMismatch(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*.a.b.c`,
		inputJSON:   `{"a":{"b":1},"b":{"a":2}}`,
		expectedErr: createErrorTypeUnmatched(`.b`, `object`, `float64`),
	}
	runTestCase(t, testCase, "TestWildcard_ChildErrorObjectTypeMismatch")
}

func TestWildcard_ChildErrorArrayTypeMismatch(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*.a.b.c`,
		inputJSON:   `[{"b":1},{"a":2}]`,
		expectedErr: createErrorTypeUnmatched(`.b`, `object`, `float64`),
	}
	runTestCase(t, testCase, "TestWildcard_ChildErrorArrayTypeMismatch")
}

func TestWildcard_ChildErrorObjectMissingNested(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*.a.b.c`,
		inputJSON:   `{"a":{"a":1},"b":{"a":{"c":2}}}`,
		expectedErr: createErrorMemberNotExist(`.b`),
	}
	runTestCase(t, testCase, "TestWildcard_ChildErrorObjectMissingNested")
}

func TestWildcard_ChildErrorArrayMissingNested(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$.*.a.b.c`,
		inputJSON:   `[{"a":1},{"a":{"c":2}}]`,
		expectedErr: createErrorMemberNotExist(`.b`),
	}
	runTestCase(t, testCase, "TestWildcard_ChildErrorArrayMissingNested")
}
