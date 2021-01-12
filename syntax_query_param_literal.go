package jsonpath

type syntaxQueryParamLiteral struct {
	literal []interface{}
}

func (l *syntaxQueryParamLiteral) compute(
	_ interface{}, _ []interface{}, _ *bufferContainer) []interface{} {

	return l.literal
}
