package syntax

type syntaxQueryParamLiteral struct {
	literal []any
}

func (l *syntaxQueryParamLiteral) compute(
	_ any, _ []any) []any {

	return l.literal
}
