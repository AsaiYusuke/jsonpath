package jsonpath

type syntaxQueryParamLiteral struct {
	literal interface{}
}

func (l *syntaxQueryParamLiteral) compute(root interface{}, currentMap map[int]interface{}) map[int]interface{} {
	return map[int]interface{}{0: l.literal}
}
