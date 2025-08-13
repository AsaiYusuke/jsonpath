package syntax

type syntaxSingleIndexProvider interface {
	getIndex(srcLength int) int
}
