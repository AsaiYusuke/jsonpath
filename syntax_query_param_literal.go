package jsonpath

type syntaxQueryParamLiteral struct {
	literal interface{}
}

func (l *syntaxQueryParamLiteral) compute(_ map[int]interface{}) map[int]interface{} {
	return map[int]interface{}{0: l.literal}
}
