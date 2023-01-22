package jsonpath

type syntaxQueryParamLiteral struct {
	literal []interface{}
}

func (l *syntaxQueryParamLiteral) compute(
	_ interface{}, _ []interface{}) []interface{} {

	return l.literal
}
