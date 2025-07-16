package jsonpath

import (
	"testing"
)

func TestRetrieve_aggregate_function_max(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*.max()`,
		inputJSON:    `[122.345,123.45,123.456]`,
		expectedJSON: `[123.456]`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`max`: maxFunc,
		},
	}
	runTestCase(t, testCase, "TestRetrieve_aggregate_function_max")
}

func TestRetrieve_filter_aggregate_max_in_condition(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.max())]`,
		inputJSON:    `[122.345,123.45,123.456]`,
		expectedJSON: `[122.345,123.45,123.456]`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`max`: maxFunc,
		},
	}
	runTestCase(t, testCase, "TestRetrieve_filter_aggregate_max_in_condition")
}

func TestRetrieve_filter_aggregate_max_equal_comparison(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.max() == 123.45)]`,
		inputJSON:    `[122.345,123.45,123.456]`,
		expectedJSON: `[123.45]`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`max`: maxFunc,
		},
	}
	runTestCase(t, testCase, "TestRetrieve_filter_aggregate_max_equal_comparison")
}

func TestRetrieve_filter_aggregate_max_not_equal_comparison(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(123.45 != @.max())]`,
		inputJSON:    `[122.345,123.45,123.456]`,
		expectedJSON: `[122.345,123.456]`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`max`: maxFunc,
		},
	}
	runTestCase(t, testCase, "TestRetrieve_filter_aggregate_max_not_equal_comparison")
}

func TestRetrieve_filter_aggregate_max_array_comparison(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.max() != 123.45)]`,
		inputJSON:    `[[122.345,123.45,123.456],[122.345,123.45]]`,
		expectedJSON: `[[122.345,123.45,123.456]]`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`max`: maxFunc,
		},
	}
	runTestCase(t, testCase, "TestRetrieve_filter_aggregate_max_array_comparison")
}

func TestRetrieve_filter_aggregate_max_cross_array_comparison(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[?(@.max() == $[1].max())]`,
		inputJSON:    `[[122.345,123.45,123.456],[122.345,123.45]]`,
		expectedJSON: `[[122.345,123.45]]`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`max`: maxFunc,
		},
	}
	runTestCase(t, testCase, "TestRetrieve_filter_aggregate_max_cross_array_comparison")
}

func TestRetrieve_aggregate_filter_mix_max_twice(t *testing.T) {
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
	runTestCase(t, testCase, "TestRetrieve_aggregate_filter_mix_max_twice")
}

func TestRetrieve_filter_aggregate_mix_twice_max(t *testing.T) {
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
	runTestCase(t, testCase, "TestRetrieve_filter_aggregate_mix_twice_max")
}

func TestRetrieve_filter_error_simple(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.errFilter()`,
		inputJSON: `[122.345,123.45,123.456]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`errFilter`: errFilterFunc,
		},
		expectedErr: createErrorFunctionFailed(`.errFilter()`, `filter error`),
	}
	runTestCase(t, testCase, "TestRetrieve_filter_error_simple")
}

func TestRetrieve_filter_error_wildcard(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.errFilter()`,
		inputJSON: `[122.345,123.45,123.456]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`errFilter`: errFilterFunc,
		},
		expectedErr: createErrorFunctionFailed(`.errFilter()`, `filter error`),
	}
	runTestCase(t, testCase, "TestRetrieve_filter_error_wildcard")
}

func TestRetrieve_filter_error_after_max(t *testing.T) {
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
	runTestCase(t, testCase, "TestRetrieve_filter_error_after_max")
}

func TestRetrieve_filter_error_after_twice(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.twice().errFilter()`,
		inputJSON: `[122.345,123.45,123.456]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`errFilter`: errFilterFunc,
			`twice`:     twiceFunc,
		},
		expectedErr: createErrorFunctionFailed(`.errFilter()`, `filter error`),
	}
	runTestCase(t, testCase, "TestRetrieve_filter_error_after_twice")
}

func TestRetrieve_filter_error_before_twice(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.errFilter().twice()`,
		inputJSON: `[122.345,123.45,123.456]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`errFilter`: errFilterFunc,
			`twice`:     twiceFunc,
		},
		expectedErr: createErrorFunctionFailed(`.errFilter()`, `filter error`),
	}
	runTestCase(t, testCase, "TestRetrieve_filter_error_before_twice")
}

func TestRetrieve_filter_error_complex_chain(t *testing.T) {
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
	runTestCase(t, testCase, "TestRetrieve_filter_error_complex_chain")
}

func TestRetrieve_aggregate_error_simple(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.errAggregate()`,
		inputJSON: `[122.345,123.45,123.456]`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`errAggregate`: errAggregateFunc,
		},
		expectedErr: createErrorFunctionFailed(`.errAggregate()`, `aggregate error`),
	}
	runTestCase(t, testCase, "TestRetrieve_aggregate_error_simple")
}

func TestRetrieve_aggregate_error_after_max(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.max().errAggregate()`,
		inputJSON: `[122.345,123.45,123.456]`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`errAggregate`: errAggregateFunc,
			`max`:          maxFunc,
		},
		expectedErr: createErrorFunctionFailed(`.errAggregate()`, `aggregate error`),
	}
	runTestCase(t, testCase, "TestRetrieve_aggregate_error_after_max")
}

func TestRetrieve_aggregate_error_after_twice(t *testing.T) {
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
	runTestCase(t, testCase, "TestRetrieve_aggregate_error_after_twice")
}

func TestRetrieve_aggregate_error_before_twice(t *testing.T) {
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
	runTestCase(t, testCase, "TestRetrieve_aggregate_error_before_twice")
}

func TestRetrieve_aggregate_error_complex_chain(t *testing.T) {
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
	runTestCase(t, testCase, "TestRetrieve_aggregate_error_complex_chain")
}

func TestRetrieve_aggregate_missing_member(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.a.max()`,
		inputJSON: `{}`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`max`: maxFunc,
		},
		expectedErr: createErrorMemberNotExist(`.a`),
	}
	runTestCase(t, testCase, "TestRetrieve_aggregate_missing_member")
}

func TestRetrieve_aggregate_missing_member_with_path(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.x.*.errAggregate()`,
		inputJSON: `{"a":[122.345,123.45,123.456]}`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`errAggregate`: errAggregateFunc,
		},
		expectedErr: createErrorMemberNotExist(`.x`),
	}
	runTestCase(t, testCase, "TestRetrieve_aggregate_missing_member_with_path")
}

func TestRetrieve_filter_error_nested_path(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.a.b.c.errFilter()`,
		inputJSON: `[{"a":{"b":1}},{"a":2}]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`errFilter`: errFilterFunc,
		},
		expectedErr: createErrorTypeUnmatched(`.c`, `object`, `float64`),
	}
	runTestCase(t, testCase, "TestRetrieve_filter_error_nested_path")
}

func TestRetrieve_filter_error_multiple_filters(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.a.b.c.errFilter1().errFilter2()`,
		inputJSON: `[{"a":{"b":1}},{"a":2}]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`errFilter1`: errFilterFunc,
			`errFilter2`: errFilterFunc,
		},
		expectedErr: createErrorTypeUnmatched(`.c`, `object`, `float64`),
	}
	runTestCase(t, testCase, "TestRetrieve_filter_error_multiple_filters")
}

func TestRetrieve_aggregate_filter_error_nested_aggregate(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.a.b.c.errAggregate()`,
		inputJSON: `[{"a":{"b":1}},{"a":2}]`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`errAggregate`: errAggregateFunc,
		},
		expectedErr: createErrorTypeUnmatched(`.c`, `object`, `float64`),
	}
	runTestCase(t, testCase, "TestRetrieve_aggregate_filter_error_nested_aggregate")
}

func TestRetrieve_aggregate_filter_error_multiple_aggregates(t *testing.T) {
	testCase := TestCase{
		jsonpath:  `$.*.a.b.c.errAggregate1().errAggregate2()`,
		inputJSON: `[{"a":{"b":1}},{"a":2}]`,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`errAggregate1`: errAggregateFunc,
			`errAggregate2`: errAggregateFunc,
		},
		expectedErr: createErrorTypeUnmatched(`.c`, `object`, `float64`),
	}
	runTestCase(t, testCase, "TestRetrieve_aggregate_filter_error_multiple_aggregates")
}

func TestRetrieve_aggregate_filter_error_mixed_aggregate_filter(t *testing.T) {
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
	runTestCase(t, testCase, "TestRetrieve_aggregate_filter_error_mixed_aggregate_filter")
}

func TestRetrieve_aggregate_filter_error_mixed_filter_aggregate(t *testing.T) {
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
	runTestCase(t, testCase, "TestRetrieve_aggregate_filter_error_mixed_filter_aggregate")
}

func TestRetrieve_function_syntax_uppercase_name(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*.TWICE()`,
		inputJSON:    `[123.456,256]`,
		expectedJSON: `[246.912,512]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`TWICE`: twiceFunc,
		},
	}
	runTestCase(t, testCase, "TestRetrieve_function_syntax_uppercase_name")
}

func TestRetrieve_function_syntax_numeric_name(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*.123()`,
		inputJSON:    `[123.456,256]`,
		expectedJSON: `[246.912,512]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`123`: twiceFunc,
		},
	}
	runTestCase(t, testCase, "TestRetrieve_function_syntax_numeric_name")
}

func TestRetrieve_function_syntax_dash_name(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*.--()`,
		inputJSON:    `[123.456,256]`,
		expectedJSON: `[246.912,512]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`--`: twiceFunc,
		},
	}
	runTestCase(t, testCase, "TestRetrieve_function_syntax_dash_name")
}

func TestRetrieve_function_syntax_underscore_name(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.*.__()`,
		inputJSON:    `[123.456,256]`,
		expectedJSON: `[246.912,512]`,
		filters: map[string]func(interface{}) (interface{}, error){
			`__`: twiceFunc,
		},
	}
	runTestCase(t, testCase, "TestRetrieve_function_syntax_underscore_name")
}
