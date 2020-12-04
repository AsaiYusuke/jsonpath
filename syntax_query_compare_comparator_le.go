package jsonpath

type syntaxCompareLE struct {
	*syntaxBasicNumericComparator
}

func (c syntaxCompareLE) comparator(left, right interface{}) bool {
	return left.(float64) >= right.(float64)
}
