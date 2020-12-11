package jsonpath

type syntaxAsterisk struct {
	*syntaxBasicSubscript
}

func (syntaxAsterisk) getIndexes(src []interface{}) []int {
	index, result := 0, make([]int, len(src))
	for i := range src {
		result[index] = i
		index++
	}
	return result
}
