package syntax

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestPegParserExecuteFunctions(t *testing.T) {
	stdoutBackup := os.Stdout
	os.Stdout = nil

	parser := pegJSONPathParser[uint32]{Buffer: `$`}
	parser.Init()
	parser.Parse()
	parser.Execute()

	parser.Print()
	parser.Reset()
	parser.PrintSyntaxTree()
	parser.SprintSyntaxTree()

	err := parseError[uint32]{p: &parser, maxToken: token[uint32]{begin: 0, end: 1}}
	_ = err.Error()

	parser.buffer = []rune{'\n'}
	_ = err.Error()

	parser.Parse(1)
	parser.Parse(3)

	Pretty[uint32](true)(&parser)
	parser.PrintSyntaxTree()

	_ = err.Error()

	Size[uint32](10)(&parser)

	parser.Init(func(p *pegJSONPathParser[uint32]) error {
		return fmt.Errorf(`test error`)
	})

	parser.Buffer = ``
	parser.PrintSyntaxTree()

	memoizeFunc := DisableMemoize[uint32]()
	memoizeFunc(&parser)
	parser.Parse(1)

	os.Stdout = stdoutBackup
}

func TestPegParserEmptyInputErrorPosition(t *testing.T) {
	parser := pegJSONPathParser[uint32]{Buffer: ""}
	if err := parser.Init(); err != nil {
		t.Fatal(err)
	}

	err := parser.Parse(3)
	if err == nil {
		t.Fatal("expected parsing an empty input as a JSONPath query to fail")
	}
	if got := err.Error(); !strings.Contains(got, "line 1 symbol 1 - line 1 symbol 1") {
		t.Fatalf("unexpected parse error position: %s", got)
	}
}

func TestPegParserRepeatedRuleCalls(t *testing.T) {
	tests := []struct {
		name    string
		rule    pegRule
		buffer  string
		matches bool
	}{
		{name: "empty expression", rule: ruleexpression, matches: true},
		{name: "empty segments", rule: rulesegments, matches: true},
		{name: "missing bracketed selection", rule: rulebracketedSelection},
		{name: "function without dot", rule: rulefunction, buffer: "x"},
		{name: "function without name", rule: rulefunction, buffer: "."},
		{name: "function without opening parenthesis", rule: rulefunction, buffer: ".a"},
		{name: "function without closing parenthesis", rule: rulefunction, buffer: ".a("},
		{name: "member name at end", rule: rulememberNameShorthand},
		{name: "sign at end", rule: rulesignsWithoutHyphenUnderscore},
		{name: "object selector at end", rule: ruleobjectElementSelector},
		{name: "wildcard at end", rule: rulewildcardSelector},
		{name: "array element at end", rule: rulearrayElementSelector},
		{name: "omitted array index", rule: ruleanyIndex, matches: true},
		{name: "slice separator at end", rule: rulesepArraySlice},
		{name: "query at end", rule: rulequery},
		{name: "and query at end", rule: ruleandQuery},
		{name: "basic query at end", rule: rulebasicQuery},
		{name: "query parameter at end", rule: ruleqParam},
		{name: "number or string parameter at end", rule: ruleqNumberOrStringParam},
		{name: "single JSONPath filter at end", rule: rulesingleJsonpathFilter},
		{name: "hex digits at end", rule: rulehexDigits},
		{name: "hex digit at end", rule: rulehexDigit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := pegJSONPathParser[uint32]{Buffer: tt.buffer}
			if err := parser.Init(); err != nil {
				t.Fatal(err)
			}

			for attempt := 0; attempt < 2; attempt++ {
				err := parser.Parse(int(tt.rule))
				if got := err == nil; got != tt.matches {
					t.Fatalf("attempt %d: match = %v, want %v (error: %v)", attempt+1, got, tt.matches, err)
				}
			}
		})
	}
}

func TestPegParserRegexRule(t *testing.T) {
	parser := pegJSONPathParser[uint32]{Buffer: "@.a=~/a/"}
	if err := parser.Init(); err != nil {
		t.Fatal(err)
	}
	if err := parser.Parse(int(rulebasicQuery)); err != nil {
		t.Fatalf("failed to parse regex comparison: %v", err)
	}
	foundRegex := false
	foundComparator := false
	for _, token := range parser.Tokens() {
		if token.pegRule == rulecomparator {
			foundComparator = true
		}
		if token.pegRule == ruleregex {
			foundRegex = true
		}
	}
	if !foundComparator || !foundRegex {
		t.Fatalf("regex comparison tokens missing: comparator=%v regex=%v", foundComparator, foundRegex)
	}
}
