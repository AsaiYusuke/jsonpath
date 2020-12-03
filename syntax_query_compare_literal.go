package jsonpath

type syntaxCompareLiteral struct {
	literal interface{}
}

func (l syntaxCompareLiteral) compute(root interface{}, currentMap map[int]interface{}) map[int]interface{} {
	result := make(map[int]interface{}, 0)
	result[0] = l.literal
	return result
}
