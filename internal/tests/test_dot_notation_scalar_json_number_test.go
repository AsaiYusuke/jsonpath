package tests

import (
	"encoding/json"
	"strings"
	"testing"
)

var useJSONNumberDecoderFunction = func(srcJSON string, src *any) error {
	reader := strings.NewReader(srcJSON)
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()
	return decoder.Decode(src)
}

func TestDotNotation_JsonNumberFilter(t *testing.T) {
	testCases := []TestCase{
		{
			jsonpath:      `$[?(@.a > 123)].a`,
			inputJSON:     `[{"a":123.456}]`,
			expectedJSON:  `[123.456]`,
			unmarshalFunc: useJSONNumberDecoderFunction,
		},
		{
			jsonpath:      `$[?(@.a > 123.46)].a`,
			inputJSON:     `[{"a":123.456}]`,
			expectedJSON:  `[]`,
			expectedErr:   createErrorMemberNotExist(`[?(@.a > 123.46)]`),
			unmarshalFunc: useJSONNumberDecoderFunction,
		},
		{
			jsonpath:      `$[?(@.a > 122)].a`,
			inputJSON:     `[{"a":123}]`,
			expectedJSON:  `[123]`,
			unmarshalFunc: useJSONNumberDecoderFunction,
		},
		{
			jsonpath:      `$[?(123 < @.a)].a`,
			inputJSON:     `[{"a":123.456}]`,
			expectedJSON:  `[123.456]`,
			unmarshalFunc: useJSONNumberDecoderFunction,
		},
		{
			jsonpath:      `$[?(@.a==-0.123e2)]`,
			inputJSON:     `[{"a":-12.3,"b":1},{"a":-0.123e2,"b":2},{"a":-0.123},{"a":-12},{"a":12.3},{"a":2},{"a":"-0.123e2"}]`,
			expectedJSON:  `[{"a":-12.3,"b":1},{"a":-0.123e2,"b":2}]`,
			unmarshalFunc: useJSONNumberDecoderFunction,
		},
		{
			jsonpath:      `$[?(@.a==11)]`,
			inputJSON:     `[{"a":10.999},{"a":11.00},{"a":11.10}]`,
			expectedJSON:  `[{"a":11.00}]`,
			unmarshalFunc: useJSONNumberDecoderFunction,
		},
	}

	runTestCases(t, "TestDotNotation_JsonNumberFilter", testCases)
}
