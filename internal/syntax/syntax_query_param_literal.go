package syntax

type syntaxQueryParamLiteral struct {
	literal []any
}

func (l *syntaxQueryParamLiteral) compute(
	_ any, _ []any, _ *syntaxQueryScratch) []any {

	return l.literal
}
