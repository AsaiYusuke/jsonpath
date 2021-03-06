package jsonpath

type syntaxWildcardSubscript struct {
	*syntaxBasicSubscript
}

func (*syntaxWildcardSubscript) getIndexes(srcLength int) []int {
	result := make([]int, srcLength)

	for index := 0; index < srcLength; index++ {
		result[index] = index
	}
	return result
}
