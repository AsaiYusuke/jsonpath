package jsonpath

type syntaxQueryBooleanLiteral struct {
	value bool
}

func (q *syntaxQueryBooleanLiteral) compute(root interface{}, currentList []interface{}) []interface{} {
	if q.value {
		return currentList
	}
	return emptyList
}
