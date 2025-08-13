package tests

import (
	"testing"
)

func TestFilter_PropertyAccessFilterThenPropertyChain(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[?(@.a)].b.c`,
		inputJSON:   `[{"a":1},{"b":2}]`,
		expectedErr: createErrorMemberNotExist(`.b`),
	}
	runTestCase(t, testCase, "TestFilter_PropertyAccessFilterThenPropertyChain")
}

func TestFilter_ObjectPropertyNotExist(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@.missing)]`,
			inputJSON:   `[{"a":1},{"b":2}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.missing)]`),
		},
		{
			jsonpath:    `$[?(@.user.missing)]`,
			inputJSON:   `[{"user":{"name":"john"}},{"user":{"age":30}}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.user.missing)]`),
		},
		{
			jsonpath:    `$[?(@.nested.deep.property)]`,
			inputJSON:   `[{"nested":{"shallow":"value"}},{"other":"data"}]`,
			expectedErr: createErrorMemberNotExist(`[?(@.nested.deep.property)]`),
		},
	}
	runTestCases(t, "TestFilter_ObjectPropertyNotExist", testCases)
}

func TestFilter_ObjectFilterThenMissingProperty(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@.items)].count`,
			inputJSON:   `[{"items":[1,2,3]},{"other":"data"}]`,
			expectedErr: createErrorMemberNotExist(`.count`),
		},
		{
			jsonpath:    `$[?(@.user.active)].permissions`,
			inputJSON:   `[{"user":{"active":true}},{"user":{"active":false}}]`,
			expectedErr: createErrorMemberNotExist(`.permissions`),
		},
		{
			jsonpath:    `$[?(@.status == 'active')].metadata.created`,
			inputJSON:   `[{"status":"active"},{"status":"inactive","metadata":{"created":"2023-01-01"}}]`,
			expectedErr: createErrorMemberNotExist(`.metadata`),
		},
	}
	runTestCases(t, "TestFilter_ObjectFilterThenMissingProperty", testCases)
}

func TestFilter_ObjectFilterThenTypeUnmatched(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@.score >= 80)].grade.letter`,
			inputJSON:   `[{"score":85,"grade":"A"},{"score":75,"grade":{"letter":"B"}}]`,
			expectedErr: createErrorTypeUnmatched(`.letter`, `object`, `string`),
		},
	}
	runTestCases(t, "TestFilter_ObjectFilterThenTypeUnmatched", testCases)
}

func TestFilter_ObjectNestedFilterErrors(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@.items[?(@.active)])].count`,
			inputJSON:   `[{"items":[{"active":true}]},{"other":"data"}]`,
			expectedErr: createErrorMemberNotExist(`.count`),
		},
		{
			jsonpath:    `$[?(@.users[?(@.role == 'admin')])].summary`,
			inputJSON:   `[{"users":[{"role":"admin"}]},{"users":[{"role":"user"}]}]`,
			expectedErr: createErrorMemberNotExist(`.summary`),
		},
	}
	runTestCases(t, "TestFilter_ObjectNestedFilterErrors", testCases)
}

func TestFilter_ObjectBracketNotationAfterFilter(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@.user)]['profile']`,
			inputJSON:   `[{"user":"john"},{"other":"data"}]`,
			expectedErr: createErrorMemberNotExist(`['profile']`),
		},
		{
			jsonpath:    `$[?(@.active == true)]['settings']['theme']`,
			inputJSON:   `[{"active":true,"settings":"dark"},{"active":false,"settings":{"theme":"light"}}]`,
			expectedErr: createErrorTypeUnmatched(`['theme']`, `object`, `string`),
		},
	}
	runTestCases(t, "TestFilter_ObjectBracketNotationAfterFilter", testCases)
}

func TestFilter_ObjectPropertyAfterSuccessfulFilter(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@.value > 10)].details`,
			inputJSON:   `[{"value":15},{"value":5,"details":"info"}]`,
			expectedErr: createErrorMemberNotExist(`.details`),
		},
		{
			jsonpath:    `$[?(@.name == 'john')].age.years`,
			inputJSON:   `[{"name":"john","age":30},{"name":"jane","age":{"years":25}}]`,
			expectedErr: createErrorTypeUnmatched(`.years`, `object`, `float64`),
		},
	}
	runTestCases(t, "TestFilter_ObjectPropertyAfterSuccessfulFilter", testCases)
}

func TestFilter_ObjectChainedPropertyAccess(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[?(@.category == 'electronics')].product.details.specs`,
			inputJSON:   `[{"category":"electronics","product":{"details":"summary"}},{"category":"books","product":{"details":{"specs":"tech"}}}]`,
			expectedErr: createErrorTypeUnmatched(`.specs`, `object`, `string`),
		},
	}
	runTestCases(t, "TestFilter_ObjectChainedPropertyAccess", testCases)
}
