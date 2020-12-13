package jsonpath

type syntaxCompareGT struct {
	*syntaxBasicNumericComparator
}

func (c *syntaxCompareGT) comparator(left, right interface{}) bool {
	return left.(float64) < right.(float64)
}
