package syntax

type syntaxWildcardSubscript struct {
	*syntaxBasicSubscript
}

func (*syntaxWildcardSubscript) getIndexes(srcLength int) []int {
	result := make([]int, srcLength)

	for index := range srcLength {
		result[index] = index
	}
	return result
}
