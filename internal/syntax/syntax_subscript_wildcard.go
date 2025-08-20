package syntax

type syntaxWildcardSubscript struct {
	*syntaxBasicSubscript
}

func (*syntaxWildcardSubscript) count(srcLength int) int {
	return srcLength
}

func (*syntaxWildcardSubscript) indexAt(srcLength int, ordinal int) int {
	return ordinal
}
