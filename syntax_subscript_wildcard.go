package jsonpath

type syntaxWildcardSubscript struct {
	*syntaxBasicSubscript
}

func (*syntaxWildcardSubscript) getIndexes(src []interface{}) []int {
	result := make([]int, len(src))
	for index := range src {
		result[index] = index
	}
	return result
}
