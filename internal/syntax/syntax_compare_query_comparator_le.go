package syntax

import "encoding/json"

type syntaxCompareLE struct {
}

func (c *syntaxCompareLE) comparator(left []any, right any) bool {
	var rightValue float64
	switch rightType := right.(type) {
	case float64:
		rightValue = rightType
	default:
		return false
	}

	var hasValue bool
	for leftIndex := range left {
		if left[leftIndex] == emptyEntity {
			continue
		}
		switch leftValue := left[leftIndex].(type) {
		case float64:
			if leftValue <= rightValue {
				hasValue = true
			} else {
				left[leftIndex] = emptyEntity
			}
		case json.Number:
			leftFloat, _ := leftValue.Float64()
			if leftFloat <= rightValue {
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
