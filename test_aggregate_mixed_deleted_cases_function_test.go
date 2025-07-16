package jsonpath

import (
	"fmt"
	"testing"
)

func TestFunctionAndAggregateDeletedCases(t *testing.T) {
	tests := []struct {
		name     string
		testCase TestCase
	}{
		{
			name: "unknown function with invalid syntax",
			testCase: TestCase{
				jsonpath:    `$.func()(`,
				inputJSON:   `{}`,
				expectedErr: ErrorFunctionNotFound{function: `.func()`},
			},
		},
		{
			name: "unknown function call",
			testCase: TestCase{
				jsonpath:    `$.func(){}`,
				inputJSON:   `{}`,
				expectedErr: ErrorFunctionNotFound{function: `.func()`},
			},
		},
		{
			name: "aggregate function on empty object - no error",
			testCase: TestCase{
				jsonpath:    `$.a.max()`,
				inputJSON:   `{}`,
				expectedErr: ErrorFunctionNotFound{function: `.max()`},
			},
		},
		{
			name: "aggregate function on missing member",
			testCase: TestCase{
				jsonpath:    `$.a.max()`,
				inputJSON:   `{}`,
				expectedErr: ErrorFunctionNotFound{function: `.max()`},
			},
		},
		{
			name: "errAggregate followed by errFilter",
			testCase: TestCase{
				jsonpath:    `$.*.a.b.c.errAggregate().errFilter()`,
				inputJSON:   `[{"a":{"b":1}},{"a":2}]`,
				expectedErr: ErrorFunctionNotFound{function: `.errAggregate()`},
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter`: func(param interface{}) (interface{}, error) {
						return nil, fmt.Errorf("filter error")
					},
				},
			},
		},
		{
			name: "errFilter followed by errAggregate",
			testCase: TestCase{
				jsonpath:    `$.*.a.b.c.errFilter().errAggregate()`,
				inputJSON:   `[{"a":{"b":1}},{"a":2}]`,
				expectedErr: ErrorFunctionNotFound{function: `.errAggregate()`},
				filters: map[string]func(interface{}) (interface{}, error){
					`errFilter`: func(param interface{}) (interface{}, error) {
						return nil, fmt.Errorf("filter error")
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestCase(t, tt.testCase, tt.name)
		})
	}
}
