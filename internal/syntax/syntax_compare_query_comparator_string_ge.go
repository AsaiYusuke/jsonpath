package syntax

type syntaxCompareStringGE struct {
}

func (c *syntaxCompareStringGE) comparator(left []any, right any) bool {
	rightStringValue, _ := right.(string)

	var hasValue bool
	for leftIndex := range left {
		if left[leftIndex] == emptyEntity {
			continue
		}
		switch leftValue := left[leftIndex].(type) {
		case string:
			if leftValue >= rightStringValue {
				hasValue = true
			} else {
				left[leftIndex] = emptyEntity
			}
		default:
			left[leftIndex] = emptyEntity
		}
	}

	return hasValue
}
