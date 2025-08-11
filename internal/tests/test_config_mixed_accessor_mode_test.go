package tests

import (
	"testing"
)

func TestConfig_AccessorModeOperations(t *testing.T) {
	testGroups := TestGroup{
		`getter-setter`: []TestCase{
			{
				jsonpath:     `$.b`,
				inputJSON:    `{"a":11,"b":22,"c":33}`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					0, 22.0, 33.0, 44.0,
					func(src interface{}) interface{} {
						return src.(map[string]interface{})[`b`]
					},
					func(src, value interface{}) {
						src.(map[string]interface{})[`b`] = value
					}),
			},
			{
				jsonpath:     `$['b']`,
				inputJSON:    `{"a":123,"b":456,"c":789}`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					0, 456.0, 246.0, 369.0,
					func(src interface{}) interface{} {
						return src.(map[string]interface{})[`b`]
					},
					func(src, value interface{}) {
						src.(map[string]interface{})[`b`] = value
					}),
			},
			{
				jsonpath:     `b`,
				inputJSON:    `{"a":11,"b":22,"c":33}`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					0, 22.0, 33.0, 44.0,
					func(src interface{}) interface{} {
						return src.(map[string]interface{})[`b`]
					},
					func(src, value interface{}) {
						src.(map[string]interface{})[`b`] = value
					}),
			},
			{
				jsonpath:     `$['a','b','c']`,
				inputJSON:    `{"a":11,"b":22}`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					1, 22.0, 44.0, 55.0,
					func(src interface{}) interface{} {
						return src.(map[string]interface{})[`b`]
					},
					func(src, value interface{}) {
						src.(map[string]interface{})[`b`] = value
					}),
			},
			{
				jsonpath:     `$.*`,
				inputJSON:    `{"a":11,"b":22,"c":33}`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					1, 22.0, 33.0, 44.0,
					func(src interface{}) interface{} {
						return src.(map[string]interface{})[`b`]
					},
					func(src, value interface{}) {
						src.(map[string]interface{})[`b`] = value
					}),
			},
			{
				jsonpath:     `$.*`,
				inputJSON:    `[11,22,33]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					1, 22.0, 33.0, 44.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[1]
					},
					func(src, value interface{}) {
						src.([]interface{})[1] = value
					}),
			},
			{
				jsonpath:     `$..a`,
				inputJSON:    `{"b":{"a":11}, "c":66, "a":77}`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					1, 11.0, 22.0, 44.0,
					func(src interface{}) interface{} {
						return src.(map[string]interface{})[`b`].(map[string]interface{})[`a`]
					},
					func(src, value interface{}) {
						src.(map[string]interface{})[`b`].(map[string]interface{})[`a`] = value
					}),
			},
			{
				jsonpath:     `$[1]`,
				inputJSON:    `[123.456,256,789]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					0, 256.0, 512.0, 1024.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[1]
					},
					func(src, value interface{}) {
						src.([]interface{})[1] = value
					}),
			},
			{
				jsonpath:     `$[2,1]`,
				inputJSON:    `[123.456,256,789]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					1, 256.0, 512.0, 1024.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[1]
					},
					func(src, value interface{}) {
						src.([]interface{})[1] = value
					}),
			},
			{
				jsonpath:     `$[0:2]`,
				inputJSON:    `[11,22,33,44]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					0, 11.0, 22.0, 44.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[0]
					},
					func(src, value interface{}) {
						src.([]interface{})[0] = value
					}),
			},
			{
				jsonpath:     `$[*]`,
				inputJSON:    `[11,22]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					0, 11.0, 22.0, 44.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[0]
					},
					func(src, value interface{}) {
						src.([]interface{})[0] = value
					}),
			},
			{
				jsonpath:     `$[0:2].a`,
				inputJSON:    `[{"a":11},{"a":22},{"a":33}]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					0, 11.0, 22.0, 44.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[0].(map[string]interface{})[`a`]
					},
					func(src, value interface{}) {
						src.([]interface{})[0].(map[string]interface{})[`a`] = value
					}),
			},
			{
				jsonpath:     `$[?(@==11||@==33)]`,
				inputJSON:    `{"a":11,"b":22,"c":33}`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					1, 33.0, 22.0, 44.0,
					func(src interface{}) interface{} {
						return src.(map[string]interface{})[`c`]
					},
					func(src, value interface{}) {
						src.(map[string]interface{})[`c`] = value
					}),
			},
			{
				jsonpath:     `$[?(@==11||@==33)]`,
				inputJSON:    `[11,22,33]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					1, 33.0, 22.0, 44.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[2]
					},
					func(src, value interface{}) {
						src.([]interface{})[2] = value
					}),
			},
			{
				jsonpath:     `$[?(@.a==11||@.a==33)].a`,
				inputJSON:    `[{"a":11},{"a":22},{"a":33}]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					1, 33.0, 22.0, 44.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[2].(map[string]interface{})[`a`]
					},
					func(src, value interface{}) {
						src.([]interface{})[2].(map[string]interface{})[`a`] = value
					}),
			},
			{
				jsonpath:     `$[?(@==$[1])]`,
				inputJSON:    `[11,22,33]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidator(
					0, 22.0, 33.0, 44.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[1]
					},
					func(src, value interface{}) {
						src.([]interface{})[1] = value
					}),
			},
		},
		`get-only`: []TestCase{
			{
				jsonpath:        `$`,
				inputJSON:       `[1,2,3]`,
				accessorMode:    true,
				resultValidator: createGetOnlyValidator([]interface{}{1.0, 2.0, 3.0}),
			},
			{
				jsonpath:     `$.echo()`,
				inputJSON:    `[122.345,123.45,123.456]`,
				accessorMode: true,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`echo`: echoAggregateFunc,
				},
				resultValidator: createGetOnlyValidator([]interface{}{122.345, 123.45, 123.456}),
			},
		},
		`convert-srcJSON`: []TestCase{
			{
				jsonpath:        `$[1]`,
				inputJSON:       `[1,2,3]`,
				accessorMode:    true,
				resultValidator: sliceStructChangedResultValidator,
			},
			{
				jsonpath:        `$.a`,
				inputJSON:       `{"a":1}`,
				accessorMode:    true,
				resultValidator: mapStructChangedResultValidator,
			},
		},
	}

	runTestGroups(t, testGroups)
}
