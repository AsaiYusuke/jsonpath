package jsonpath

import (
	"fmt"
	"os"
	"testing"
)

func TestPegParserExecuteFunctions(t *testing.T) {
	stdoutBackup := os.Stdout
	os.Stdout = nil

	parser := pegJSONPathParser{Buffer: `$`}
	parser.Init()
	parser.Parse()
	parser.Execute()

	parser.Print()
	parser.Reset()
	parser.PrintSyntaxTree()
	parser.SprintSyntaxTree()

	err := parseError{p: &parser, max: token32{begin: 0, end: 1}}
	_ = err.Error()

	parser.buffer = []rune{'\n'}
	_ = err.Error()

	parser.Parse(1)
	parser.Parse(3)

	Pretty(true)(&parser)
	parser.PrintSyntaxTree()

	_ = err.Error()

	Size(10)(&parser)

	parser.Init(func(p *pegJSONPathParser) error {
		return fmt.Errorf(`test error`)
	})

	parser.Buffer = ``
	parser.PrintSyntaxTree()

	memoizeFunc := DisableMemoize()
	memoizeFunc(&parser)

	os.Stdout = stdoutBackup
}
