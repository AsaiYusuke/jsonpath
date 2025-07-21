package tests

import (
	"testing"
)

func TestAggregateFunction_MaxFunction(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*.max()`,
		inputJSON:    `[122.345,123.45,123.456]`,
		expectedJSON: `[123.456]`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`max`: maxFunc,
		},
	}
	runTestCase(t, testCase, "TestAggregateFunction_MaxFunction")
}

func TestAggregateFunction_FilterMaxInCondition(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.max())]`,
		inputJSON:    `[122.345,123.45,123.456]`,
		expectedJSON: `[122.345,123.45,123.456]`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`max`: maxFunc,
		},
	}
	runTestCase(t, testCase, "TestAggregateFunction_FilterMaxInCondition")
}

func TestAggregateFunction_FilterMaxEqualComparison(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.max() == 123.45)]`,
		inputJSON:    `[122.345,123.45,123.456]`,
		expectedJSON: `[123.45]`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`max`: maxFunc,
		},
	}
	runTestCase(t, testCase, "TestAggregateFunction_FilterMaxEqualComparison")
}

func TestAggregateFunction_FilterMaxNotEqualComparison(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(123.45 != @.max())]`,
		inputJSON:    `[122.345,123.45,123.456]`,
		expectedJSON: `[122.345,123.456]`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`max`: maxFunc,
		},
	}
	runTestCase(t, testCase, "TestAggregateFunction_FilterMaxNotEqualComparison")
}

func TestAggregateFunction_FilterMaxArrayComparison(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.max() != 123.45)]`,
		inputJSON:    `[[122.345,123.45,123.456],[122.345,123.45]]`,
		expectedJSON: `[[122.345,123.45,123.456]]`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`max`: maxFunc,
		},
	}
	runTestCase(t, testCase, "TestAggregateFunction_FilterMaxArrayComparison")
}

func TestAggregateFunction_FilterMaxCrossArrayComparison(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.max() == $[1].max())]`,
		inputJSON:    `[[122.345,123.45,123.456],[122.345,123.45]]`,
		expectedJSON: `[[122.345,123.45]]`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`max`: maxFunc,
		},
	}
	runTestCase(t, testCase, "TestAggregateFunction_FilterMaxCrossArrayComparison")
}

func TestAggregateFunction_FilterMixMaxTwice(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*.max().twice()`,
		inputJSON:    `[122.345,123.45,123.456]`,
		expectedJSON: `[246.912]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`twice`: twiceFunc,
		},
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`max`: maxFunc,
		},
	}
	runTestCase(t, testCase, "TestAggregateFunction_FilterMixMaxTwice")
}

func TestAggregateFunction_FilterMixTwiceMax(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*.twice().max()`,
		inputJSON:    `[122.345,123.45,123.456]`,
		expectedJSON: `[246.912]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`twice`: twiceFunc,
		},
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`max`: maxFunc,
		},
	}
	runTestCase(t, testCase, "TestAggregateFunction_FilterMixTwiceMax")
}

func TestAggregateFunction_FilterErrorSimple(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.errFilter()`,
		inputJSON: `[122.345,123.45,123.456]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`errFilter`: errFilterFunc,
		},
		expectedErr: createErrorFunctionFailed(`.errFilter()`, `filter error`),
	}
	runTestCase(t, testCase, "TestAggregateFunction_FilterErrorSimple")
}

func TestAggregateFunction_FilterErrorWildcard(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.errFilter()`,
		inputJSON: `[122.345,123.45,123.456]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`errFilter`: errFilterFunc,
		},
		expectedErr: createErrorFunctionFailed(`.errFilter()`, `filter error`),
	}
	runTestCase(t, testCase, "TestAggregateFunction_FilterErrorWildcard")
}

func TestAggregateFunction_FilterErrorAfterMax(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.max().errFilter()`,
		inputJSON: `[122.345,123.45,123.456]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`errFilter`: errFilterFunc,
		},
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`max`: maxFunc,
		},
		expectedErr: createErrorFunctionFailed(`.errFilter()`, `filter error`),
	}
	runTestCase(t, testCase, "TestAggregateFunction_FilterErrorAfterMax")
}

func TestAggregateFunction_FilterErrorAfterTwice(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.twice().errFilter()`,
		inputJSON: `[122.345,123.45,123.456]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`errFilter`: errFilterFunc,
			`twice`:     twiceFunc,
		},
		expectedErr: createErrorFunctionFailed(`.errFilter()`, `filter error`),
	}
	runTestCase(t, testCase, "TestAggregateFunction_FilterErrorAfterTwice")
}

func TestAggregateFunction_FilterErrorBeforeTwice(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.errFilter().twice()`,
		inputJSON: `[122.345,123.45,123.456]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`errFilter`: errFilterFunc,
			`twice`:     twiceFunc,
		},
		expectedErr: createErrorFunctionFailed(`.errFilter()`, `filter error`),
	}
	runTestCase(t, testCase, "TestAggregateFunction_FilterErrorBeforeTwice")
}

func TestAggregateFunction_FilterErrorComplexChain(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.max().errFilter().twice()`,
		inputJSON: `[122.345,123.45,123.456]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`errFilter`: errFilterFunc,
			`twice`:     twiceFunc,
		},
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`max`: maxFunc,
		},
		expectedErr: createErrorFunctionFailed(`.errFilter()`, `filter error`),
	}
	runTestCase(t, testCase, "TestAggregateFunction_FilterErrorComplexChain")
}

func TestAggregateFunction_ErrorSimple(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.errAggregate()`,
		inputJSON: `[122.345,123.45,123.456]`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`errAggregate`: errAggregateFunc,
		},
		expectedErr: createErrorFunctionFailed(`.errAggregate()`, `aggregate error`),
	}
	runTestCase(t, testCase, "TestAggregateFunction_ErrorSimple")
}

func TestAggregateFunction_ErrorAfterMax(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.max().errAggregate()`,
		inputJSON: `[122.345,123.45,123.456]`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`errAggregate`: errAggregateFunc,
			`max`:          maxFunc,
		},
		expectedErr: createErrorFunctionFailed(`.errAggregate()`, `aggregate error`),
	}
	runTestCase(t, testCase, "TestAggregateFunction_ErrorAfterMax")
}

func TestAggregateFunction_ErrorAfterTwice(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.twice().errAggregate()`,
		inputJSON: `[122.345,123.45,123.456]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`twice`: twiceFunc,
		},
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`errAggregate`: errAggregateFunc,
		},
		expectedErr: createErrorFunctionFailed(`.errAggregate()`, `aggregate error`),
	}
	runTestCase(t, testCase, "TestAggregateFunction_ErrorAfterTwice")
}

func TestAggregateFunction_ErrorBeforeTwice(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.errAggregate().twice()`,
		inputJSON: `[122.345,123.45,123.456]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`twice`: twiceFunc,
		},
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`errAggregate`: errAggregateFunc,
		},
		expectedErr: createErrorFunctionFailed(`.errAggregate()`, `aggregate error`),
	}
	runTestCase(t, testCase, "TestAggregateFunction_ErrorBeforeTwice")
}

func TestAggregateFunction_ErrorComplexChain(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.max().errAggregate().twice()`,
		inputJSON: `[122.345,123.45,123.456]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`twice`: twiceFunc,
		},
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`errAggregate`: errAggregateFunc,
			`max`:          maxFunc,
		},
		expectedErr: createErrorFunctionFailed(`.errAggregate()`, `aggregate error`),
	}
	runTestCase(t, testCase, "TestAggregateFunction_ErrorComplexChain")
}

func TestAggregateFunction_MissingMember(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.a.max()`,
		inputJSON: `{}`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`max`: maxFunc,
		},
		expectedErr: createErrorMemberNotExist(`.a`),
	}
	runTestCase(t, testCase, "TestAggregateFunction_MissingMember")
}

func TestAggregateFunction_MissingMember_with_path(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.x.*.errAggregate()`,
		inputJSON: `{"a":[122.345,123.45,123.456]}`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`errAggregate`: errAggregateFunc,
		},
		expectedErr: createErrorMemberNotExist(`.x`),
	}
	runTestCase(t, testCase, "TestAggregateFunction_MissingMember_with_path")
}

func TestAggregateFunction_FilterErrornested_path(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.a.b.c.errFilter()`,
		inputJSON: `[{"a":{"b":1}},{"a":2}]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`errFilter`: errFilterFunc,
		},
		expectedErr: createErrorTypeUnmatched(`.c`, `object`, `float64`),
	}
	runTestCase(t, testCase, "TestAggregateFunction_FilterErrornested_path")
}

func TestAggregateFunction_FilterErrormultiple_filters(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.a.b.c.errFilter1().errFilter2()`,
		inputJSON: `[{"a":{"b":1}},{"a":2}]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`errFilter1`: errFilterFunc,
			`errFilter2`: errFilterFunc,
		},
		expectedErr: createErrorTypeUnmatched(`.c`, `object`, `float64`),
	}
	runTestCase(t, testCase, "TestAggregateFunction_FilterErrormultiple_filters")
}

func TestAggregateFunction_FilterErrorNestedAggregate(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.a.b.c.errAggregate()`,
		inputJSON: `[{"a":{"b":1}},{"a":2}]`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`errAggregate`: errAggregateFunc,
		},
		expectedErr: createErrorTypeUnmatched(`.c`, `object`, `float64`),
	}
	runTestCase(t, testCase, "TestAggregateFunction_FilterErrorNestedAggregate")
}

func TestAggregateFunction_FilterErrorMultipleAggregates(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.a.b.c.errAggregate1().errAggregate2()`,
		inputJSON: `[{"a":{"b":1}},{"a":2}]`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`errAggregate1`: errAggregateFunc,
			`errAggregate2`: errAggregateFunc,
		},
		expectedErr: createErrorTypeUnmatched(`.c`, `object`, `float64`),
	}
	runTestCase(t, testCase, "TestAggregateFunction_FilterErrorMultipleAggregates")
}

func TestAggregateFunction_FilterErrorMixedAggregateFilter(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.a.b.c.errAggregate().errFilter()`,
		inputJSON: `[{"a":{"b":1}},{"a":2}]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`errFilter`: twiceFunc,
		},
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`errAggregate`: errAggregateFunc,
		},
		expectedErr: createErrorTypeUnmatched(`.c`, `object`, `float64`),
	}
	runTestCase(t, testCase, "TestAggregateFunction_FilterErrorMixedAggregateFilter")
}

func TestAggregateFunction_FilterErrorMixedFilterAggregate(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.a.b.c.errFilter().errAggregate()`,
		inputJSON: `[{"a":{"b":1}},{"a":2}]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`errFilter`: twiceFunc,
		},
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`errAggregate`: errAggregateFunc,
		},
		expectedErr: createErrorTypeUnmatched(`.c`, `object`, `float64`),
	}
	runTestCase(t, testCase, "TestAggregateFunction_FilterErrorMixedFilterAggregate")
}

func TestAggregateFunction_FunctionSyntaxUppercaseName(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*.TWICE()`,
		inputJSON:    `[123.456,256]`,
		expectedJSON: `[246.912,512]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`TWICE`: twiceFunc,
		},
	}
	runTestCase(t, testCase, "TestAggregateFunction_FunctionSyntaxUppercaseName")
}

func TestAggregateFunction_FunctionSyntaxNumericName(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*.123()`,
		inputJSON:    `[123.456,256]`,
		expectedJSON: `[246.912,512]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`123`: twiceFunc,
		},
	}
	runTestCase(t, testCase, "TestAggregateFunction_FunctionSyntaxNumericName")
}

func TestAggregateFunction_FunctionSyntaxDashName(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*.--()`,
		inputJSON:    `[123.456,256]`,
		expectedJSON: `[246.912,512]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`--`: twiceFunc,
		},
	}
	runTestCase(t, testCase, "TestAggregateFunction_FunctionSyntaxDashName")
}

func TestAggregateFunction_FunctionSyntaxUnderscoreName(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*.__()`,
		inputJSON:    `[123.456,256]`,
		expectedJSON: `[246.912,512]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`__`: twiceFunc,
		},
	}
	runTestCase(t, testCase, "TestAggregateFunction_FunctionSyntaxUnderscoreName")
}

func TestAggregateFunction_ChainedOperations(t *testing.T) {
	testCases := []TestCase{
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
	}

	runTestCases(t, "TestAggregateFunction_ChainedOperations", testCases)
}
