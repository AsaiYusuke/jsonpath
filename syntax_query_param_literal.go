package jsonpath

type syntaxQueryParamLiteral struct {
	literal map[int]interface{}
}

func (l *syntaxQueryParamLiteral) compute(
	_ interface{}, _ map[int]interface{}) map[int]interface{} {

	return l.literal
}
