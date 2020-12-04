package jsonpath

type syntaxCompareGE struct {
	*syntaxBasicNumericComparator
}

func (c syntaxCompareGE) comparator(left, right interface{}) bool {
	return left.(float64) <= right.(float64)
}
