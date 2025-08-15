package syntax

type syntaxCompareDirectEQ struct {
	syntaxTypeValidator
}

func (c *syntaxCompareDirectEQ) comparator(left []any, right any) bool {
	var hasValue bool
	for leftIndex := range left {
		if left[leftIndex] == right {
			hasValue = true
		} else {
			left[leftIndex] = emptyEntity
		}
	}
	return hasValue
}
