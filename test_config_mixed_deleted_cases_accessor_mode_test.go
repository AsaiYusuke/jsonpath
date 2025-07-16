package jsonpath

import (
	"testing"
)

func TestRetrieve_accessor_mode_logical_or_property(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$[?(@.a==11||@.a==33)].a`,
		inputJSON:       `[{"a":11},{"a":22},{"a":33}]`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestRetrieve_accessor_mode_logical_or_property")
}

func TestRetrieve_accessor_mode_self_reference_comparison(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$[?(@==$[1])]`,
		inputJSON:       `[11,22,33]`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestRetrieve_accessor_mode_self_reference_comparison")
}

// TestRetrieve_accessorModeDeleted tests deleted accessor mode test cases

// Get-only validator tests
func TestRetrieve_accessorModeDeleted_get_only_root(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$`,
		inputJSON:       `[1,2,3]`,
		accessorMode:    true,
		resultValidator: getOnlyValidator,
	}
	runTestCase(t, testCase, "TestRetrieve_accessorModeDeleted_get_only_root")
}

func TestRetrieve_accessorModeDeleted_get_only_echo(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$.echo()`,
		inputJSON:    `[122.345,123.45,123.456]`,
		accessorMode: true,
		aggregates: map[string]func([]interface{}) (interface{}, error){
			`echo`: echoAggregateFunc,
		},
		resultValidator: getOnlyValidator,
	}
	runTestCase(t, testCase, "TestRetrieve_accessorModeDeleted_get_only_echo")
}

// Struct changed result validator tests
func TestRetrieve_accessorModeDeleted_slice_struct_changed(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$[1]`,
		inputJSON:       `[1,2,3]`,
		accessorMode:    true,
		resultValidator: sliceStructChangedResultValidator,
	}
	runTestCase(t, testCase, "TestRetrieve_accessorModeDeleted_slice_struct_changed")
}

func TestRetrieve_accessorModeDeleted_map_struct_changed(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$.a`,
		inputJSON:       `{"a":1}`,
		accessorMode:    true,
		resultValidator: mapStructChangedResultValidator,
	}
	runTestCase(t, testCase, "TestRetrieve_accessorModeDeleted_map_struct_changed")
}

// Additional accessor mode tests that were deleted from the original test suite

func TestRetrieve_accessor_mode_basic_property_b_test(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$.b`,
		inputJSON:       `{"a":11,"b":22,"c":33}`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestRetrieve_accessor_mode_basic_property_b_test")
}

func TestRetrieve_accessor_mode_bracket_property_b_test(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$['b']`,
		inputJSON:       `{"a":123,"b":456,"c":789}`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestRetrieve_accessor_mode_bracket_property_b_test")
}

func TestRetrieve_accessor_mode_simple_property_b_test(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `b`,
		inputJSON:       `{"a":11,"b":22,"c":33}`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestRetrieve_accessor_mode_simple_property_b_test")
}

func TestRetrieve_accessor_mode_multi_property_test(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$['a','b','c']`,
		inputJSON:       `{"a":11,"b":22}`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestRetrieve_accessor_mode_multi_property_test")
}

func TestRetrieve_accessor_mode_wildcard_object_test(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$.*`,
		inputJSON:       `{"a":11,"b":22,"c":33}`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestRetrieve_accessor_mode_wildcard_object_test")
}

func TestRetrieve_accessor_mode_wildcard_array_test(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$.*`,
		inputJSON:       `[11,22,33]`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestRetrieve_accessor_mode_wildcard_array_test")
}

func TestRetrieve_accessor_mode_recursive_descent_test(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$..a`,
		inputJSON:       `{"b":{"a":11}, "c":66, "a":77}`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestRetrieve_accessor_mode_recursive_descent_test")
}

func TestRetrieve_accessor_mode_array_index_test(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$[1]`,
		inputJSON:       `[123.456,256,789]`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestRetrieve_accessor_mode_array_index_test")
}

func TestRetrieve_accessor_mode_array_multi_index_test(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$[2,1]`,
		inputJSON:       `[123.456,256,789]`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestRetrieve_accessor_mode_array_multi_index_test")
}

func TestRetrieve_accessor_mode_array_slice_test(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$[0:2]`,
		inputJSON:       `[11,22,33,44]`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestRetrieve_accessor_mode_array_slice_test")
}

func TestRetrieve_accessor_mode_array_wildcard_test(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$[*]`,
		inputJSON:       `[11,22]`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestRetrieve_accessor_mode_array_wildcard_test")
}

func TestRetrieve_accessor_mode_slice_with_property_test(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$[0:2].a`,
		inputJSON:       `[{"a":11},{"a":22},{"a":33}]`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestRetrieve_accessor_mode_slice_with_property_test")
}

func TestRetrieve_accessor_mode_filter_object_test(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$[?(@==11||@==33)]`,
		inputJSON:       `{"a":11,"b":22,"c":33}`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestRetrieve_accessor_mode_filter_object_test")
}

func TestRetrieve_accessor_mode_filter_array_test(t *testing.T) {
	testCase := TestCase{
		jsonpath:        `$[?(@==11||@==33)]`,
		inputJSON:       `[11,22,33]`,
		accessorMode:    true,
		resultValidator: createAccessorModeValidator(),
	}
	runTestCase(t, testCase, "TestRetrieve_accessor_mode_filter_array_test")
}
