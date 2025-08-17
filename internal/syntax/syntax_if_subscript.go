package syntax

type syntaxSubscript interface {
	forEachIndex(srcLength int, handleIndex func(index int))
	isValueGroup() bool
}
