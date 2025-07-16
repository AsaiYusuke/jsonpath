package jsonpath

import (
	"testing"
)

func TestConfig_AccessorModeLogicalOrProperty(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$[?(@.a==11||@.a==33)].a`,
		inputJSON:       `[{"a":11},{"a":22},{"a":33}]`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestConfig_AccessorModeLogicalOrProperty")
}

func TestConfig_AccessorModeSelfReferenceComparison(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$[?(@==$[1])]`,
		inputJSON:       `[11,22,33]`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestConfig_AccessorModeSelfReferenceComparison")
}

func TestConfig_AccessorModeDeletedGetOnlyRoot(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$`,
		inputJSON:       `[1,2,3]`,
		accessorMode:    true,
		resultValidator: getOnlyValidator,
	}
	runTestCase(t, testCase, "TestConfig_AccessorModeDeletedGetOnlyRoot")
}

func TestConfig_AccessorModeDeletedGetOnlyEcho(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.echo()`,
		inputJSON:    `[122.345,123.45,123.456]`,
		accessorMode: true,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`echo`: echoAggregateFunc,
		},
		resultValidator: getOnlyValidator,
	}
	runTestCase(t, testCase, "TestConfig_AccessorModeDeletedGetOnlyEcho")
}

func TestConfig_AccessorModeDeletedSliceStructChanged(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$[1]`,
		inputJSON:       `[1,2,3]`,
		accessorMode:    true,
		resultValidator: sliceStructChangedResultValidator,
	}
	runTestCase(t, testCase, "TestConfig_AccessorModeDeletedSliceStructChanged")
}

func TestConfig_AccessorModeDeletedMapStructChanged(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$.a`,
		inputJSON:       `{"a":1}`,
		accessorMode:    true,
		resultValidator: mapStructChangedResultValidator,
	}
	runTestCase(t, testCase, "TestConfig_AccessorModeDeletedMapStructChanged")
}

func TestConfig_AccessorModeBasicPropertyB(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$.b`,
		inputJSON:       `{"a":11,"b":22,"c":33}`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestConfig_AccessorModeBasicPropertyB")
}

func TestConfig_AccessorModeBracketPropertyB(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$['b']`,
		inputJSON:       `{"a":123,"b":456,"c":789}`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestConfig_AccessorModeBracketPropertyB")
}

func TestConfig_AccessorModeSimplePropertyB(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `b`,
		inputJSON:       `{"a":11,"b":22,"c":33}`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestConfig_AccessorModeSimplePropertyB")
}

func TestConfig_AccessorModeMultiProperty(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$['a','b','c']`,
		inputJSON:       `{"a":11,"b":22}`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestConfig_AccessorModeMultiProperty")
}

func TestConfig_AccessorModeWildcardObject(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$.*`,
		inputJSON:       `{"a":11,"b":22,"c":33}`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestConfig_AccessorModeWildcardObject")
}

func TestConfig_AccessorModeWildcardArray(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$.*`,
		inputJSON:       `[11,22,33]`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestConfig_AccessorModeWildcardArray")
}

func TestConfig_AccessorModeRecursiveDescent(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$..a`,
		inputJSON:       `{"b":{"a":11}, "c":66, "a":77}`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestConfig_AccessorModeRecursiveDescent")
}

func TestConfig_AccessorModeArrayIndex(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$[1]`,
		inputJSON:       `[123.456,256,789]`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestConfig_AccessorModeArrayIndex")
}

func TestConfig_AccessorModeArrayMultiIndex(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$[2,1]`,
		inputJSON:       `[123.456,256,789]`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestConfig_AccessorModeArrayMultiIndex")
}

func TestConfig_AccessorModeArraySlice(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$[0:2]`,
		inputJSON:       `[11,22,33,44]`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestConfig_AccessorModeArraySlice")
}

func TestConfig_AccessorModeArrayWildcard(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$[*]`,
		inputJSON:       `[11,22]`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestConfig_AccessorModeArrayWildcard")
}

func TestConfig_AccessorModeSliceWithProperty(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$[0:2].a`,
		inputJSON:       `[{"a":11},{"a":22},{"a":33}]`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestConfig_AccessorModeSliceWithProperty")
}

func TestConfig_AccessorModeFilterObject(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$[?(@==11||@==33)]`,
		inputJSON:       `{"a":11,"b":22,"c":33}`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestConfig_AccessorModeFilterObject")
}

func TestConfig_AccessorModeFilterArray(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$[?(@==11||@==33)]`,
		inputJSON:       `[11,22,33]`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestConfig_AccessorModeFilterArray")
}
