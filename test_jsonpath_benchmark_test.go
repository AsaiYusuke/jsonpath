package jsonpath

import (
	"encoding/json"
	"testing"
)

func execRetrieve(jsonPath, srcJSON string, b *testing.B) {
	var src interface{}
	if err := json.Unmarshal([]byte(srcJSON), &src); err != nil {
		b.Errorf(`%w`, err)
		return
	}

	for i := 0; i < b.N; i++ {
		if _, err := Retrieve(jsonPath, src); err != nil {
			b.Error(`%w`, err)
		}

	}
}

func execParserFunc(jsonPath, srcJSON string, b *testing.B) {
	var src interface{}
	if err := json.Unmarshal([]byte(srcJSON), &src); err != nil {
		b.Errorf(`%w`, err)
		return
	}

	parserFunc, err := Parse(jsonPath)
	if err != nil {
		b.Errorf(`%w`, err)
		return
	}

	for i := 0; i < b.N; i++ {
		if _, err := parserFunc(src); err != nil {
			b.Error(`%w`, err)
		}

	}
}

// =====================================================================
// Retrieve

func BenchmarkRetrieve_dotNotation(b *testing.B) {
	jsonPath := `$.a`
	srcJSON := `{"a":123.456}`
	execRetrieve(jsonPath, srcJSON, b)
}

func BenchmarkRetrieve_bracketNotation(b *testing.B) {
	jsonPath := `$['a']`
	srcJSON := `{"a":123.456}`
	execRetrieve(jsonPath, srcJSON, b)
}

func BenchmarkRetrieve_asterisk_identifier_dotNotation(b *testing.B) {
	jsonPath := `$.*`
	srcJSON := `{"a":123.456}`
	execRetrieve(jsonPath, srcJSON, b)
}

func BenchmarkRetrieve_asterisk_identifier_bracketNotation(b *testing.B) {
	jsonPath := `$[*]`
	srcJSON := `{"a":123.456}`
	execRetrieve(jsonPath, srcJSON, b)
}

func BenchmarkRetrieve_multi_identifier(b *testing.B) {
	jsonPath := `$['a','a']`
	srcJSON := `{"a":123.456}`
	execRetrieve(jsonPath, srcJSON, b)
}

func BenchmarkRetrieve_qualifier_index(b *testing.B) {
	jsonPath := `$[0]`
	srcJSON := `[{"a":123.456}]`
	execRetrieve(jsonPath, srcJSON, b)
}

func BenchmarkRetrieve_qualifier_slice(b *testing.B) {
	jsonPath := `$[0:1]`
	srcJSON := `[{"a":123.456}]`
	execRetrieve(jsonPath, srcJSON, b)
}

func BenchmarkRetrieve_qualifier_asterisk(b *testing.B) {
	jsonPath := `$[*]`
	srcJSON := `[{"a":123.456}]`
	execRetrieve(jsonPath, srcJSON, b)
}

func BenchmarkRetrieve_qualifier_union(b *testing.B) {
	jsonPath := `$[0,0]`
	srcJSON := `[{"a":123.456}]`
	execRetrieve(jsonPath, srcJSON, b)
}

func BenchmarkRetrieve_filter_logicalOR(b *testing.B) {
	jsonPath := `$[?(@||@)]`
	srcJSON := `[{"a":1}]`
	execRetrieve(jsonPath, srcJSON, b)
}

func BenchmarkRetrieve_filter_logicalAND(b *testing.B) {
	jsonPath := `$[?(@&&@)]`
	srcJSON := `[{"a":1}]`
	execRetrieve(jsonPath, srcJSON, b)
}

func BenchmarkRetrieve_filter_nodeFilter(b *testing.B) {
	jsonPath := `$[?(@.a)]`
	srcJSON := `[{"a":1}]`
	execRetrieve(jsonPath, srcJSON, b)
}

func BenchmarkRetrieve_filter_logicalNOT(b *testing.B) {
	jsonPath := `$[?(!@.a)]`
	srcJSON := `[{"a":1},{"b":1}]`
	execRetrieve(jsonPath, srcJSON, b)
}

func BenchmarkRetrieve_filter_compareEQ(b *testing.B) {
	jsonPath := `$[?(@.a==1)]`
	srcJSON := `[{"a":1}]`
	execRetrieve(jsonPath, srcJSON, b)
}

func BenchmarkRetrieve_filter_compareNE(b *testing.B) {
	jsonPath := `$[?(@.a!=2)]`
	srcJSON := `[{"a":1}]`
	execRetrieve(jsonPath, srcJSON, b)
}

func BenchmarkRetrieve_filter_compareGE(b *testing.B) {
	jsonPath := `$[?(@.a<=2)]`
	srcJSON := `[{"a":1}]`
	execRetrieve(jsonPath, srcJSON, b)
}

func BenchmarkRetrieve_filter_compareGT(b *testing.B) {
	jsonPath := `$[?(@.a<2)]`
	srcJSON := `[{"a":1}]`
	execRetrieve(jsonPath, srcJSON, b)
}

func BenchmarkRetrieve_filter_compareLE(b *testing.B) {
	jsonPath := `$[?(@.a>=0)]`
	srcJSON := `[{"a":1}]`
	execRetrieve(jsonPath, srcJSON, b)
}

func BenchmarkRetrieve_filter_compareLT(b *testing.B) {
	jsonPath := `$[?(@.a>0)]`
	srcJSON := `[{"a":1}]`
	execRetrieve(jsonPath, srcJSON, b)
}

func BenchmarkRetrieve_filter_regex(b *testing.B) {
	jsonPath := `$[?(@.a =~ /ab/)]`
	srcJSON := `[{"a":"abc"}]`
	execRetrieve(jsonPath, srcJSON, b)
}

// =====================================================================
// ParserFunc

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

func BenchmarkParserFunc_asterisk_identifier_dotNotation(b *testing.B) {
	jsonPath := `$.*`
	srcJSON := `{"a":123.456}`
	execParserFunc(jsonPath, srcJSON, b)
}

func BenchmarkParserFunc_asterisk_identifier_bracketNotation(b *testing.B) {
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

func BenchmarkParserFunc_qualifier_asterisk(b *testing.B) {
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
