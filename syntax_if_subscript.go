package jsonpath

type syntaxSubscript interface {
	getIndexes(srcLength int) []int
	isValueGroup() bool
}
