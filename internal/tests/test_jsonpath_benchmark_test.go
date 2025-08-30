package tests

import (
	"encoding/json"
	"testing"

	"github.com/AsaiYusuke/jsonpath/v2"
)

func execParserFunc(jsonPath, srcJSON string, b *testing.B) {
	var src any
	if err := json.Unmarshal([]byte(srcJSON), &src); err != nil {
		b.Errorf(`%s`, err)
		return
	}

	parserFunc, err := jsonpath.Parse(jsonPath)
	if err != nil {
		b.Errorf(`%s`, err)
		return
	}

	buf := make([]any, 0, 256)
	args := []*[]any{&buf}

	for b.Loop() {
		if _, err := parserFunc(src, args...); err != nil {
			b.Errorf(`%s`, err)
		}

	}
}

func BenchmarkParserFunc_dotNotation(b *testing.B) {
	jsonPath := `$.a`
	srcJSON := `{"a":123.456}`
	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_bracketNotation(b *testing.B) {
	jsonPath := `$['a']`
	srcJSON := `{"a":123.456}`
	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_wildcard_identifier_dotNotation(b *testing.B) {
	jsonPath := `$.*`
	srcJSON := `{"a":123.456}`
	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_wildcard_identifier_bracketNotation(b *testing.B) {
	jsonPath := `$[*]`
	srcJSON := `{"a":123.456}`
	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_multi_identifier(b *testing.B) {
	jsonPath := `$['a','a']`
	srcJSON := `{"a":123.456}`
	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_qualifier_index(b *testing.B) {
	jsonPath := `$[0]`
	srcJSON := `[{"a":123.456}]`
	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_qualifier_slice(b *testing.B) {
	jsonPath := `$[0:1]`
	srcJSON := `[{"a":123.456}]`
	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_qualifier_wildcard(b *testing.B) {
	jsonPath := `$[*]`
	srcJSON := `[{"a":123.456}]`
	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_qualifier_union(b *testing.B) {
	jsonPath := `$[0,0]`
	srcJSON := `[{"a":123.456}]`
	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_filter_logicalOR(b *testing.B) {
	jsonPath := `$[?(@||@)]`
	srcJSON := `[{"a":1}]`
	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_filter_logicalAND(b *testing.B) {
	jsonPath := `$[?(@&&@)]`
	srcJSON := `[{"a":1}]`
	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_filter_nodeFilter(b *testing.B) {
	jsonPath := `$[?(@.a)]`
	srcJSON := `[{"a":1}]`
	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_filter_logicalNOT(b *testing.B) {
	jsonPath := `$[?(!@.a)]`
	srcJSON := `[{"a":1},{"b":1}]`
	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_filter_compareEQ(b *testing.B) {
	jsonPath := `$[?(@.a==1)]`
	srcJSON := `[{"a":1}]`
	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_filter_compareNE(b *testing.B) {
	jsonPath := `$[?(@.a!=2)]`
	srcJSON := `[{"a":1}]`
	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_filter_compareGE(b *testing.B) {
	jsonPath := `$[?(@.a<=2)]`
	srcJSON := `[{"a":1}]`
	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_filter_compareGT(b *testing.B) {
	jsonPath := `$[?(@.a<2)]`
	srcJSON := `[{"a":1}]`
	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_filter_compareLE(b *testing.B) {
	jsonPath := `$[?(@.a>=0)]`
	srcJSON := `[{"a":1}]`
	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_filter_compareLT(b *testing.B) {
	jsonPath := `$[?(@.a>0)]`
	srcJSON := `[{"a":1}]`
	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_filter_regex(b *testing.B) {
	jsonPath := `$[?(@.a =~ /ab/)]`
	srcJSON := `[{"a":"abc"}]`
	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_recursive_name(b *testing.B) {
	jsonPath := `$..price`
	srcJSON := `{ "store": {
		"book": [ 
		  { "category": "reference",
			"author": "Nigel Rees",
			"title": "Sayings of the Century",
			"price": 8.95
		  },
		  { "category": "fiction",
			"author": "Evelyn Waugh",
			"title": "Sword of Honour",
			"price": 12.99
		  },
		  { "category": "fiction",
			"author": "Herman Melville",
			"title": "Moby Dick",
			"isbn": "0-553-21311-3",
			"price": 8.99
		  },
		  { "category": "fiction",
			"author": "J. R. R. Tolkien",
			"title": "The Lord of the Rings",
			"isbn": "0-395-19395-8",
			"price": 22.99
		  }
		],
		"bicycle": {
		  "color": "red",
		  "price": 19.95
		}
	  }
	}`

	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_recursive_index(b *testing.B) {
	jsonPath := `$..[0]`
	srcJSON := `{ "store": {
		"bicycle": {
		  "gearRatios": [2.1, 2.5, 3.0, 3.5],
		  "tirePressuresPSI": [60, 65, 70, 75],
		  "serviceMileageHistoryKm": [500, 1200, 1800, 2500],
		  "frameSizeOptionsCm": [48, 50, 52, 54],
		  "wheelDiameterOptionsMm": [622, 630, 635, 584]
		}
	  }
	}`

	execParserFunc(jsonPath, srcJSON, b)
}
