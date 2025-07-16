package jsonpath

import (
	"testing"
)

func TestRetrieve_arrayUnion_index_duplicate_first(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0,0]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first","first"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayUnion_index_duplicate_first")
}

func TestRetrieve_arrayUnion_index_first_second(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0,1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first","second"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayUnion_index_first_second")
}

func TestRetrieve_arrayUnion_index_first_last(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0,-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first","third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayUnion_index_first_last")
}

func TestRetrieve_arrayUnion_index_multiple_order(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[2,0,1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["third","first","second"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayUnion_index_multiple_order")
}

func TestRetrieve_arrayUnion_wildcard_index_first_all(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0,*]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first","first","second","third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayUnion_wildcard_index_first_all")
}

func TestRetrieve_arrayUnion_wildcard_all_index_first(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[*,0]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first","second","third","first"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayUnion_wildcard_all_index_first")
}

func TestRetrieve_arrayUnion_wildcard_slice_all(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[1:2,*]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["second","first","second","third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayUnion_wildcard_slice_all")
}

func TestRetrieve_arrayUnion_wildcard_all_slice(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[*,1:2]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first","second","third","second"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayUnion_wildcard_all_slice")
}

func TestRetrieve_arrayUnion_wildcard_all_all(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[*,*]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first","second","third","first","second","third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayUnion_wildcard_all_all")
}

func TestRetrieve_arrayUnion_slice_index(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[1:2,0]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["second","first"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayUnion_slice_index")
}

func TestRetrieve_arrayUnion_index_nested_array_access(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0,1]`,
		inputJSON:    `[["11","12","13"],["21","22","23"],["31","32","33"]]`,
		expectedJSON: `[["11","12","13"],["21","22","23"]]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayUnion_index_nested_array_access")
}

func TestRetrieve_arrayUnion_wildcard_duplicates(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[*,0,*]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first","second","third","first","first","second","third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayUnion_wildcard_duplicates")
}

func TestRetrieve_arrayUnion_slice_combinations(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[1:2,0:2]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["second","first","second"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayUnion_slice_combinations")
}

func TestRetrieve_arrayUnion_index_wildcard_mixed(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0,*,0]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first","first","second","third","first"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arrayUnion_index_wildcard_mixed")
}
