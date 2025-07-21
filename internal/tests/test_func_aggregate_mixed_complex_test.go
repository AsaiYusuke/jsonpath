package tests

import (
	"testing"
)

func TestAggregateFunction_MixedComplexCases(t *testing.T) {
	testGroups := map[string][]TestCase{
		"advanced-aggregate-filter": {
			{
				jsonpath:     `$.*.max().twice().max()`,
				inputJSON:    `[122.345,123.45,123.456]`,
				expectedJSON: `[246.912]`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxAggregate,
				},
				filters: map[string]func(interface{}) (interface{}, error){
					`twice`: twiceFilter,
				},
			},
			{
				jsonpath:     `$.*.max().twice().twice()`,
				inputJSON:    `[122.345,123.45,123.456]`,
				expectedJSON: `[493.824]`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxAggregate,
				},
				filters: map[string]func(interface{}) (interface{}, error){
					`twice`: twiceFilter,
				},
			},
			{
				jsonpath:     `$.*.max().twice().twice().max()`,
				inputJSON:    `[122.345,123.45,123.456]`,
				expectedJSON: `[493.824]`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxAggregate,
				},
				filters: map[string]func(interface{}) (interface{}, error){
					`twice`: twiceFilter,
				},
			},
			{
				jsonpath:     `$.*.max().twice()`,
				inputJSON:    `[122.345,123.45,123.456]`,
				expectedJSON: `[246.912]`,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxAggregate,
				},
				filters: map[string]func(interface{}) (interface{}, error){
					`twice`: twiceFilter,
				},
			},
			{
				jsonpath:    `$.*.max().errFilter()`,
				inputJSON:   `[122.345,123.45,123.456]`,
				expectedErr: createErrorFunctionFailed(".errFilter()", "function failed (function=errFilter, error=test error)"),
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxAggregate,
				},
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter`: errorFilter,
				},
			},
			{
				jsonpath:    `$.*.max().errFilter().max()`,
				inputJSON:   `[122.345,123.45,123.456]`,
				expectedErr: createErrorFunctionFailed(".errFilter()", "function failed (function=errFilter, error=test error)"),
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxAggregate,
				},
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter`: errorFilter,
				},
			},
			{
				jsonpath:    `$.*.max().errFilter().twice()`,
				inputJSON:   `[122.345,123.45,123.456]`,
				expectedErr: createErrorFunctionFailed(".errFilter()", "function failed (function=errFilter, error=test error)"),
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxAggregate,
				},
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter`: errorFilter,
					`twice`:     twiceFilter,
				},
			},
			{
				jsonpath:    `$.*.max().errFilter().twice().max()`,
				inputJSON:   `[122.345,123.45,123.456]`,
				expectedErr: createErrorFunctionFailed(".errFilter()", "function failed (function=errFilter, error=test error)"),
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxAggregate,
				},
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter`: errorFilter,
					`twice`:     twiceFilter,
				},
			},
			{
				jsonpath:    `$.*.twice().errFilter()`,
				inputJSON:   `[122.345,123.45,123.456]`,
				expectedErr: createErrorFunctionFailed(".errFilter()", "function failed (function=errFilter, error=test error)"),
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter`: errorFilter,
					`twice`:     twiceFilter,
				},
			},
			{
				jsonpath:    `$.*.max().twice()`,
				inputJSON:   `[122.345,123.45,123.456]`,
				expectedErr: createErrorFunctionFailed(".twice()", "function failed (function=errFilter, error=test error)"),
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxAggregate,
				},
				filters: map[string]func(interface{}) (interface{}, error){
					`twice`: errorFilter,
				},
			},
			{
				jsonpath:    `$.a.max()`,
				inputJSON:   `{}`,
				expectedErr: createErrorMemberNotExist(".a"),
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxAggregate,
				},
			},
			{
				jsonpath:    `$.a.max().twice()`,
				inputJSON:   `{}`,
				expectedErr: createErrorMemberNotExist(".a"),
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxAggregate,
				},
				filters: map[string]func(interface{}) (interface{}, error){
					`twice`: twiceFilter,
				},
			},
			{
				jsonpath:    `$.a.max()`,
				inputJSON:   `{"a": null}`,
				expectedErr: createErrorFunctionFailed(".max()", "function failed (function=max, error=non-numeric value)"),
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxAggregate,
				},
			},
			{
				jsonpath:    `$.a.max()`,
				inputJSON:   `{"a": "string"}`,
				expectedErr: createErrorFunctionFailed(".max()", "function failed (function=max, error=non-numeric value)"),
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxAggregate,
				},
			},
			{
				jsonpath:    `$.*.max()[-1]`,
				inputJSON:   `[122.345,123.45,123.456]`,
				expectedErr: createErrorInvalidSyntax(9, "unrecognized input", "[-1]"),
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxAggregate,
				},
			},
			{
				jsonpath:    `$.*.max().twice()[-1]`,
				inputJSON:   `[122.345,123.45,123.456]`,
				expectedErr: createErrorInvalidSyntax(17, "unrecognized input", "[-1]"),
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`max`: maxAggregate,
				},
				filters: map[string]func(interface{}) (interface{}, error){
					`twice`: twiceFilter,
				},
			},
		},
	}

	runTestGroups(t, testGroups)
}
