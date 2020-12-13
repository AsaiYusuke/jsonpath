package jsonpath

type syntaxCompareLT struct {
	*syntaxBasicNumericComparator
}

func (c *syntaxCompareLT) comparator(left, right interface{}) bool {
	return left.(float64) > right.(float64)
}
