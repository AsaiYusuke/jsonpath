package tests

import (
	"testing"
)

func TestSlice_PositiveStepOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[0:3:1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second","third"]`,
		},
		{
			jsonpath:     `$[0:3:2]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","third"]`,
		},
		{
			jsonpath:     `$[0:3:3]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first"]`,
		},
		{
			jsonpath:     `$[0:2:2]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first"]`,
		},
		{
			jsonpath:     `$[0:2:3]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first"]`,
		},
		{
			jsonpath:     `$[0:1:3]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first"]`,
		},
	}

	runTestCases(t, "TestSlice_PositiveStepOperations", testCases)
}

func TestSlice_ZeroStepErrorCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[0:2:0]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[0:2:0]`),
		},
		{
			jsonpath:    `$[2:0:0]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[2:0:0]`),
		},
	}

	runTestCases(t, "TestSlice_ZeroStepErrorCases", testCases)
}

func TestSlice_NegativeStepOperations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[-3:1:-1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[-3:1:-1]`),
		},
		{
			jsonpath:    `$[-2:1:-1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[-2:1:-1]`),
		},
		{
			jsonpath:     `$[-1:1:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["third"]`,
		},
		{
			jsonpath:    `$[0:1:-1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[0:1:-1]`),
		},
		{
			jsonpath:    `$[1:1:-1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[1:1:-1]`),
		},
		{
			jsonpath:     `$[2:1:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["third"]`,
		},
		{
			jsonpath:     `$[3:1:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["third"]`,
		},
		{
			jsonpath:     `$[4:1:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["third"]`,
		},
	}

	runTestCases(t, "TestSlice_NegativeStepOperations", testCases)
}

func TestSlice_NegativeStepComplexCases(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[0:-2:-1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[0:-2:-1]`),
		},
		{
			jsonpath:    `$[0:-1:-1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[0:-1:-1]`),
		},
		{
			jsonpath:    `$[0:0:-1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[0:0:-1]`),
		},
		{
			jsonpath:    `$[0:2:-1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[0:2:-1]`),
		},
		{
			jsonpath:    `$[0:3:-1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[0:3:-1]`),
		},
		{
			jsonpath:    `$[0:4:-1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[0:4:-1]`),
		},
		{
			jsonpath:     `$[1:-5:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["second","first"]`,
		},
		{
			jsonpath:     `$[1:-4:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["second","first"]`,
		},
		{
			jsonpath:     `$[1:-3:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["second"]`,
		},
		{
			jsonpath:    `$[1:-2:-1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[1:-2:-1]`),
		},
		{
			jsonpath:    `$[1:-1:-1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[1:-1:-1]`),
		},
		{
			jsonpath:     `$[1:0:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["second"]`,
		},
		{
			jsonpath:    `$[1:2:-1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[1:2:-1]`),
		},
		{
			jsonpath:    `$[1:3:-1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[1:3:-1]`),
		},
		{
			jsonpath:     `$[2:-5:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["third","second","first"]`,
		},
		{
			jsonpath:     `$[2:-4:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["third","second","first"]`,
		},
		{
			jsonpath:     `$[2:-3:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["third","second"]`,
		},
		{
			jsonpath:     `$[2:-2:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["third"]`,
		},
		{
			jsonpath:    `$[2:-1:-1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[2:-1:-1]`),
		},
		{
			jsonpath:     `$[2:0:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["third","second"]`,
		},
		{
			jsonpath:    `$[2:2:-1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[2:2:-1]`),
		},
		{
			jsonpath:    `$[2:3:-1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[2:3:-1]`),
		},
		{
			jsonpath:    `$[2:4:-1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[2:4:-1]`),
		},
		{
			jsonpath:    `$[2:5:-1]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[2:5:-1]`),
		},
	}

	runTestCases(t, "TestSlice_NegativeStepComplexCases", testCases)
}

func TestSlice_NegativeStepVariations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[2:0:-2]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["third"]`,
		},
		{
			jsonpath:     `$[2:0:-3]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["third"]`,
		},
		{
			jsonpath:    `$[2:-1:-2]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorMemberNotExist(`[2:-1:-2]`),
		},
		{
			jsonpath:     `$[-1:0:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["third","second"]`,
		},
	}

	runTestCases(t, "TestSlice_NegativeStepVariations", testCases)
}

func TestSlice_OmittedParameters(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[0:3:]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second","third"]`,
		},
		{
			jsonpath:     `$[1::1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["second","third"]`,
		},
		{
			jsonpath:     `$[1::-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["second","first"]`,
		},
		{
			jsonpath:     `$[:1:1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first"]`,
		},
		{
			jsonpath:     `$[:1:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["third"]`,
		},
		{
			jsonpath:     `$[::2]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","third"]`,
		},
		{
			jsonpath:     `$[::-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["third","second","first"]`,
		},
		{
			jsonpath:     `$[::]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second","third"]`,
		},
	}

	runTestCases(t, "TestSlice_OmittedParameters", testCases)
}

func TestSlice_LargeNumbers(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[1:1000000000000000000:1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["second","third"]`,
		},
		{
			jsonpath:     `$[1:-1000000000000000000:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["second","first"]`,
		},
		{
			jsonpath:     `$[-1000000000000000000:3:1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second","third"]`,
		},
		{
			jsonpath:     `$[1000000000000000000:0:-1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["third","second"]`,
		},
		{
			jsonpath:     `$[1:0:-1000000000000000000]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["second"]`,
		},
		{
			jsonpath:     `$[0:1:1000000000000000000]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first"]`,
		},
	}

	runTestCases(t, "TestSlice_LargeNumbers", testCases)
}

func TestSlice_SyntaxVariations(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:     `$[0:3:+1]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second","third"]`,
		},
		{
			jsonpath:     `$[0:3:01]`,
			inputJSON:    `["first","second","third"]`,
			expectedJSON: `["first","second","third"]`,
		},
		{
			jsonpath:    `$[0:3:1.0]`,
			inputJSON:   `["first","second","third"]`,
			expectedErr: createErrorInvalidSyntax(1, `unrecognized input`, `[0:3:1.0]`),
		},
	}

	runTestCases(t, "TestSlice_SyntaxVariations", testCases)
}

func TestSlice_TypeMismatchErrors(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[2:1:-1]`,
			inputJSON:   `{"first":1,"second":2,"third":3}`,
			expectedErr: createErrorTypeUnmatched(`[2:1:-1]`, `array`, `map[string]interface {}`),
		},
		{
			jsonpath:    `$[::-1]`,
			inputJSON:   `{"first":1,"second":2,"third":3}`,
			expectedErr: createErrorTypeUnmatched(`[::-1]`, `array`, `map[string]interface {}`),
		},
		{
			jsonpath:    `$[2:1:-1]`,
			inputJSON:   `"value"`,
			expectedErr: createErrorTypeUnmatched(`[2:1:-1]`, `array`, `string`),
		},
		{
			jsonpath:    `$[2:1:-1]`,
			inputJSON:   `1`,
			expectedErr: createErrorTypeUnmatched(`[2:1:-1]`, `array`, `float64`),
		},
		{
			jsonpath:    `$[2:1:-1]`,
			inputJSON:   `true`,
			expectedErr: createErrorTypeUnmatched(`[2:1:-1]`, `array`, `bool`),
		},
		{
			jsonpath:    `$[2:1:-1]`,
			inputJSON:   `null`,
			expectedErr: createErrorTypeUnmatched(`[2:1:-1]`, `array`, `null`),
		},
	}

	runTestCases(t, "TestSlice_TypeMismatchErrors", testCases)
}

func TestSlice_ChildAccessErrors(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:    `$[-1:-1:-1].a.b`,
			inputJSON:   `[0]`,
			expectedErr: createErrorMemberNotExist(`[-1:-1:-1]`),
		},
		{
			jsonpath:    `$[0:-2:-1].a.b`,
			inputJSON:   `[{"b":1}]`,
			expectedErr: createErrorMemberNotExist(`.a`),
		},
		{
			jsonpath:    `$[1:-3:-1].a.b`,
			inputJSON:   `[{"b":1},{"c":2}]`,
			expectedErr: createErrorMemberNotExist(`.a`),
		},
		{
			jsonpath:    `$[1:-3:-1].a.b.c`,
			inputJSON:   `[{"a":1},{"b":2}]`,
			expectedErr: createErrorTypeUnmatched(`.b`, `object`, `float64`),
		},
	}

	runTestCases(t, "TestSlice_ChildAccessErrors", testCases)
}
