package jsonpath

type syntaxCompareGE struct {
}

func (c syntaxCompareGE) comparator(left, right interface{}) bool {
	leftValue, leftOk := left.(float64)
	rightValue, rightOk := right.(float64)
	return leftOk && rightOk && leftValue <= rightValue
}
