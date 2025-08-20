package syntax

type syntaxSubscript interface {
	isValueGroup() bool
	count(srcLength int) int
	indexAt(srcLength int, ordinal int) int
}
