package jsonpath

type syntaxCompareGT struct {
	*syntaxBasicNumericComparator
}

func (c *syntaxCompareGT) comparator(left []interface{}, right interface{}) bool {
	var hasValue bool
	for leftIndex := range left {
		if left[leftIndex] == emptyEntity {
			continue
		}
		if left[leftIndex].(float64) < right.(float64) {
			hasValue = true
		} else {
			left[leftIndex] = emptyEntity
		}
	}
	return hasValue
}
