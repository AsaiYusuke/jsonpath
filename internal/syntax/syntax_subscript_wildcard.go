package syntax

type syntaxWildcardSubscript struct {
	*syntaxBasicSubscript
}

func (*syntaxWildcardSubscript) forEachIndex(srcLength int, handleIndex func(index int)) {
	for index := range srcLength {
		handleIndex(index)
	}
}
