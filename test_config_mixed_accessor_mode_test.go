package jsonpath

import (
	"fmt"
	"testing"
)

func createAccessorModeValidatorOrig(
	resultIndex int,
	expectedValue1, expectedValue2, expectedValue3 interface{},
	srcGetter func(interface{}) interface{},
	srcSetter func(interface{}, interface{})) func(interface{}, []interface{}) error {
	return func(src interface{}, actualObject []interface{}) error {
		accessor := actualObject[resultIndex].(Accessor)

		getValue := accessor.Get()
		if getValue != expectedValue1 {
			return fmt.Errorf(`Get : expect<%f> != actual<%f>`, expectedValue1, getValue)
		}

		accessor.Set(expectedValue2)

		newSrcValue := srcGetter(src)

		if newSrcValue != expectedValue2 {
			return fmt.Errorf(`Set : expect<%f> != actual<%f>`, expectedValue2, newSrcValue)
		}

		getValue = accessor.Get()
		if getValue != expectedValue2 {
			return fmt.Errorf(`Set -> Get : expect<%f> != actual<%f>`, expectedValue2, getValue)
		}

		srcSetter(src, expectedValue3)

		getValue = accessor.Get()
		if getValue != expectedValue3 {
			return fmt.Errorf(`Src -> Get : expect<%f> != actual<%f>`, expectedValue3, getValue)
		}

		return nil
	}
}

func TestRetrieve_configAccessorMode(t *testing.T) {
	testGroups := TestGroup{
		`getter-setter`: []TestCase{
			{
				jsonpath:     `$.b`,
				inputJSON:    `{"a":11,"b":22,"c":33}`,
				accessorMode: true,
				resultValidator: createAccessorModeValidatorOrig(
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
				resultValidator: createAccessorModeValidatorOrig(
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
				resultValidator: createAccessorModeValidatorOrig(
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
				resultValidator: createAccessorModeValidatorOrig(
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
				inputJSON:    `{"a":11,"b":22}`,
				accessorMode: true,
				resultValidator: createAccessorModeValidatorOrig(
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
				inputJSON:    `[1,2,3]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidatorOrig(
					1, 2.0, 4.0, 5.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[1]
					},
					func(src, value interface{}) {
						src.([]interface{})[1] = value
					}),
			},
			{
				jsonpath:     `$..a`,
				inputJSON:    `{"a":1,"b":{"a":2}}`,
				accessorMode: true,
				resultValidator: createAccessorModeValidatorOrig(
					1, 2.0, 4.0, 5.0,
					func(src interface{}) interface{} {
						return src.(map[string]interface{})[`b`].(map[string]interface{})[`a`]
					},
					func(src, value interface{}) {
						src.(map[string]interface{})[`b`].(map[string]interface{})[`a`] = value
					}),
			},
			{
				jsonpath:     `$[1]`,
				inputJSON:    `[1,2,3]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidatorOrig(
					0, 2.0, 4.0, 5.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[1]
					},
					func(src, value interface{}) {
						src.([]interface{})[1] = value
					}),
			},
			{
				jsonpath:     `$[2,1]`,
				inputJSON:    `[1,2,3]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidatorOrig(
					1, 2.0, 4.0, 5.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[1]
					},
					func(src, value interface{}) {
						src.([]interface{})[1] = value
					}),
			},
			{
				jsonpath:     `$[0:2]`,
				inputJSON:    `[1,2,3]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidatorOrig(
					1, 2.0, 4.0, 5.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[1]
					},
					func(src, value interface{}) {
						src.([]interface{})[1] = value
					}),
			},
			{
				jsonpath:     `$[*]`,
				inputJSON:    `[1,2,3]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidatorOrig(
					1, 2.0, 4.0, 5.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[1]
					},
					func(src, value interface{}) {
						src.([]interface{})[1] = value
					}),
			},
			{
				jsonpath:     `$[0:2].a`,
				inputJSON:    `[{"a":1},{"a":2},{"a":3}]`,
				accessorMode: true,
				resultValidator: createAccessorModeValidatorOrig(
					1, 2.0, 4.0, 5.0,
					func(src interface{}) interface{} {
						return src.([]interface{})[1].(map[string]interface{})[`a`]
					},
					func(src, value interface{}) {
						src.([]interface{})[1].(map[string]interface{})[`a`] = value
					}),
			},
			{
				jsonpath:     `$[?(@==11||@==33)]`,
				inputJSON:    `{"a":11,"b":22,"c":33}`,
				accessorMode: true,
				resultValidator: createAccessorModeValidatorOrig(
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
				resultValidator: createAccessorModeValidatorOrig(
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
				resultValidator: createAccessorModeValidatorOrig(
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
				resultValidator: createAccessorModeValidatorOrig(
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
				resultValidator: getOnlyValidator,
			},
			{
				jsonpath:     `$.echo()`,
				inputJSON:    `[122.345,123.45,123.456]`,
				accessorMode: true,
				aggregates: map[string]func([]interface{}) (interface{}, error){
					`echo`: echoAggregateFunc,
				},
				resultValidator: getOnlyValidator,
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
