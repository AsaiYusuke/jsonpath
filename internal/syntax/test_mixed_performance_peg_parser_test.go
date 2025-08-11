package syntax

import (
	"fmt"
	"os"
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

	os.Stdout = stdoutBackup
}
