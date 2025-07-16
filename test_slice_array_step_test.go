package jsonpath

import (
	"testing"
)

func TestRetrieve_arraySlice_step_variation_step1(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0:3:1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first","second","third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_variation_step1")
}

func TestRetrieve_arraySlice_step_variation_step2(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0:3:2]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first","third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_variation_step2")
}

func TestRetrieve_arraySlice_step_variation_step3(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0:3:3]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_variation_step3")
}

func TestRetrieve_arraySlice_step_variation_end2_step2(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0:2:2]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_variation_end2_step2")
}

func TestRetrieve_arraySlice_step_variation_end2_step3(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0:2:3]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_variation_end2_step3")
}

func TestRetrieve_arraySlice_step_variation_end1_step3(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0:1:3]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_variation_end1_step3")
}

func TestRetrieve_arraySlice_step_zero_forward(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:2:0]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[0:2:0]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_zero_forward")
}

func TestRetrieve_arraySlice_step_zero_backward(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[2:0:0]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[2:0:0]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_zero_backward")
}

func TestRetrieve_arraySlice_step_minus_start_neg3(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[-3:1:-1]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[-3:1:-1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_start_neg3")
}

func TestRetrieve_arraySlice_step_minus_start_neg2(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[-2:1:-1]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[-2:1:-1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_start_neg2")
}

func TestRetrieve_arraySlice_step_minus_start_neg1(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[-1:1:-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_start_neg1")
}

func TestRetrieve_arraySlice_step_minus_start_0(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:1:-1]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[0:1:-1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_start_0")
}

func TestRetrieve_arraySlice_step_minus_start_1(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[1:1:-1]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[1:1:-1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_start_1")
}

func TestRetrieve_arraySlice_step_minus_start_2(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[2:1:-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_start_2")
}

func TestRetrieve_arraySlice_step_minus_start_3(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[3:1:-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_start_3")
}

func TestRetrieve_arraySlice_step_minus_start_4(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[4:1:-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_start_4")
}

func TestRetrieve_arraySlice_step_minus_end_start0_neg2(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:-2:-1]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[0:-2:-1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start0_neg2")
}

func TestRetrieve_arraySlice_step_minus_end_start0_neg1(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:-1:-1]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[0:-1:-1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start0_neg1")
}

func TestRetrieve_arraySlice_step_minus_end_start0_0(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:0:-1]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[0:0:-1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start0_0")
}

func TestRetrieve_arraySlice_step_minus_end_start0_2(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:2:-1]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[0:2:-1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start0_2")
}

func TestRetrieve_arraySlice_step_minus_end_start0_3(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:3:-1]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[0:3:-1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start0_3")
}

func TestRetrieve_arraySlice_step_minus_end_start0_4(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:4:-1]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[0:4:-1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start0_4")
}

func TestRetrieve_arraySlice_step_minus_end_start1_neg5(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[1:-5:-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["second","first"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start1_neg5")
}

func TestRetrieve_arraySlice_step_minus_end_start1_neg4(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[1:-4:-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["second","first"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start1_neg4")
}

func TestRetrieve_arraySlice_step_minus_end_start1_neg3(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[1:-3:-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["second"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start1_neg3")
}

func TestRetrieve_arraySlice_step_minus_end_start1_neg2(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[1:-2:-1]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[1:-2:-1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start1_neg2")
}

func TestRetrieve_arraySlice_step_minus_end_start1_neg1(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[1:-1:-1]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[1:-1:-1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start1_neg1")
}

func TestRetrieve_arraySlice_step_minus_end_start1_0(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[1:0:-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["second"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start1_0")
}

func TestRetrieve_arraySlice_step_minus_end_start1_2(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[1:2:-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["second","first"]`,
		expectedErr:  createErrorMemberNotExist(`[1:2:-1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start1_2")
}

func TestRetrieve_arraySlice_step_minus_end_start1_3(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[1:3:-1]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[1:3:-1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start1_3")
}

func TestRetrieve_arraySlice_step_minus_end_start2_neg5(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[2:-5:-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["third","second","first"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start2_neg5")
}

func TestRetrieve_arraySlice_step_minus_end_start2_neg4(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[2:-4:-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["third","second","first"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start2_neg4")
}

func TestRetrieve_arraySlice_step_minus_end_start2_neg3(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[2:-3:-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["third","second"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start2_neg3")
}

func TestRetrieve_arraySlice_step_minus_end_start2_neg2(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[2:-2:-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start2_neg2")
}

func TestRetrieve_arraySlice_step_minus_end_start2_neg1(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[2:-1:-1]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[2:-1:-1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start2_neg1")
}

func TestRetrieve_arraySlice_step_minus_end_start2_0(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[2:0:-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["third","second"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start2_0")
}

func TestRetrieve_arraySlice_step_minus_end_start2_2(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[2:2:-1]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[2:2:-1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start2_2")
}

func TestRetrieve_arraySlice_step_minus_end_start2_3(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[2:3:-1]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[2:3:-1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start2_3")
}

func TestRetrieve_arraySlice_step_minus_end_start2_4(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[2:4:-1]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[2:4:-1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start2_4")
}

func TestRetrieve_arraySlice_step_minus_end_start2_5(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[2:5:-1]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[2:5:-1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_end_start2_5")
}

func TestRetrieve_arraySlice_step_minus_step_variation_step2(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[2:0:-2]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_step_variation_step2")
}

func TestRetrieve_arraySlice_step_minus_step_variation_step3(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[2:0:-3]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_step_variation_step3")
}

func TestRetrieve_arraySlice_step_minus_start_end_variation_negative_step2(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[2:-1:-2]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: createErrorMemberNotExist(`[2:-1:-2]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_start_end_variation_negative_step2")
}

func TestRetrieve_arraySlice_step_minus_start_end_variation_negative_start(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[-1:0:-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["third","second"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_minus_start_end_variation_negative_start")
}

func TestRetrieve_arraySlice_step_omitted_number_omit_step(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0:3:]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first","second","third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_omitted_number_omit_step")
}

func TestRetrieve_arraySlice_step_omitted_number_omit_end_step1(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[1::1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["second","third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_omitted_number_omit_end_step1")
}

func TestRetrieve_arraySlice_step_omitted_number_omit_end_step_neg1(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[1::-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["second","first"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_omitted_number_omit_end_step_neg1")
}

func TestRetrieve_arraySlice_step_omitted_number_omit_start_step1(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[:1:1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_omitted_number_omit_start_step1")
}

func TestRetrieve_arraySlice_step_omitted_number_omit_start_step_neg1(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[:1:-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_omitted_number_omit_start_step_neg1")
}

func TestRetrieve_arraySlice_step_omitted_number_omit_start_end_step2(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[::2]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first","third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_omitted_number_omit_start_end_step2")
}

func TestRetrieve_arraySlice_step_omitted_number_omit_start_end_step_neg1(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[::-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["third","second","first"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_omitted_number_omit_start_end_step_neg1")
}

func TestRetrieve_arraySlice_step_omitted_number_omit_all(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[::]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first","second","third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_omitted_number_omit_all")
}

func TestRetrieve_arraySlice_step_big_number_big_end(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[1:1000000000000000000:1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["second","third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_big_number_big_end")
}

func TestRetrieve_arraySlice_step_big_number_big_neg_end(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[1:-1000000000000000000:-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["second","first"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_big_number_big_neg_end")
}

func TestRetrieve_arraySlice_step_big_number_big_neg_start(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[-1000000000000000000:3:1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first","second","third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_big_number_big_neg_start")
}

func TestRetrieve_arraySlice_step_big_number_big_start(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[1000000000000000000:0:-1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["third","second"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_big_number_big_start")
}

func TestRetrieve_arraySlice_step_big_number_big_neg_step(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[1:0:-1000000000000000000]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["second"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_big_number_big_neg_step")
}

func TestRetrieve_arraySlice_step_big_number_big_step(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0:1:1000000000000000000]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_big_number_big_step")
}

func TestRetrieve_arraySlice_step_syntax_plus_sign(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0:3:+1]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first","second","third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_syntax_plus_sign")
}

func TestRetrieve_arraySlice_step_syntax_leading_zero(t *testing.T) {
	testCase := TestCase{
		jsonpath:     `$[0:3:01]`,
		inputJSON:    `["first","second","third"]`,
		expectedJSON: `["first","second","third"]`,
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_syntax_leading_zero")
}

func TestRetrieve_arraySlice_step_syntax_decimal(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:3:1.0]`,
		inputJSON:   `["first","second","third"]`,
		expectedErr: ErrorInvalidSyntax{position: 1, reason: `unrecognized input`, near: `[0:3:1.0]`},
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_syntax_decimal")
}

func TestRetrieve_arraySlice_step_not_array_object(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[2:1:-1]`,
		inputJSON:   `{"first":1,"second":2,"third":3}`,
		expectedErr: createErrorTypeUnmatched(`[2:1:-1]`, `array`, `map[string]interface {}`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_not_array_object")
}

func TestRetrieve_arraySlice_step_not_array_object_all(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[::-1]`,
		inputJSON:   `{"first":1,"second":2,"third":3}`,
		expectedErr: createErrorTypeUnmatched(`[::-1]`, `array`, `map[string]interface {}`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_not_array_object_all")
}

func TestRetrieve_arraySlice_step_not_array_string(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[2:1:-1]`,
		inputJSON:   `"value"`,
		expectedErr: createErrorTypeUnmatched(`[2:1:-1]`, `array`, `string`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_not_array_string")
}

func TestRetrieve_arraySlice_step_not_array_number(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[2:1:-1]`,
		inputJSON:   `1`,
		expectedErr: createErrorTypeUnmatched(`[2:1:-1]`, `array`, `float64`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_not_array_number")
}

func TestRetrieve_arraySlice_step_not_array_boolean(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[2:1:-1]`,
		inputJSON:   `true`,
		expectedErr: createErrorTypeUnmatched(`[2:1:-1]`, `array`, `bool`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_not_array_boolean")
}

func TestRetrieve_arraySlice_step_not_array_null(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[2:1:-1]`,
		inputJSON:   `null`,
		expectedErr: createErrorTypeUnmatched(`[2:1:-1]`, `array`, `null`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_not_array_null")
}

func TestRetrieve_arraySlice_step_child_error_empty_result(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[-1:-1:-1].a.b`,
		inputJSON:   `[0]`,
		expectedErr: createErrorMemberNotExist(`[-1:-1:-1]`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_child_error_empty_result")
}

func TestRetrieve_arraySlice_step_child_error_missing_member(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[0:-2:-1].a.b`,
		inputJSON:   `[{"b":1}]`,
		expectedErr: createErrorMemberNotExist(`.a`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_child_error_missing_member")
}

func TestRetrieve_arraySlice_step_child_error_missing_member_multi(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[1:-3:-1].a.b`,
		inputJSON:   `[{"b":1},{"c":2}]`,
		expectedErr: createErrorMemberNotExist(`.a`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_child_error_missing_member_multi")
}

func TestRetrieve_arraySlice_step_child_error_type_mismatch(t *testing.T) {
	testCase := TestCase{
		jsonpath:    `$[1:-3:-1].a.b.c`,
		inputJSON:   `[{"a":1},{"b":2}]`,
		expectedErr: createErrorTypeUnmatched(`.b`, `object`, `float64`),
	}
	runTestCase(t, testCase, "TestRetrieve_arraySlice_step_child_error_type_mismatch")
}
