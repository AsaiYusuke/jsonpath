package jsonpath

type syntaxQueryParamLiteral struct {
	literal interface{}
}

func (l syntaxQueryParamLiteral) compute(root interface{}, currentMap map[int]interface{}) map[int]interface{} {
	result := make(map[int]interface{}, 0)
	result[0] = l.literal
	return result
}
