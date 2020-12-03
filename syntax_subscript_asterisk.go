package jsonpath

type syntaxAsterisk struct {
	*syntaxBasicSubscript
}

func (syntaxAsterisk) getIndexes(src []interface{}) []int {
	result := make([]int, 0, len(src))
	for i := range src {
		result = append(result, i)
	}
	return result
}
