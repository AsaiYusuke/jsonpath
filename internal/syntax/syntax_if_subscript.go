package syntax

type syntaxSubscript interface {
	getIndexes(srcLength int) []int
	isValueGroup() bool
}
