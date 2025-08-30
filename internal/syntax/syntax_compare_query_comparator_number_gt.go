package syntax

import "encoding/json"

type syntaxCompareNumberGT struct {
}

func (c *syntaxCompareNumberGT) compare(left []any, right any) bool {
	rightFloatValue, _ := right.(float64)

	var hasValue bool
	for leftIndex := range left {
		if left[leftIndex] == emptyEntity {
			continue
		}
		switch leftValue := left[leftIndex].(type) {
		case float64:
			if leftValue > rightFloatValue {
				hasValue = true
			} else {
				left[leftIndex] = emptyEntity
			}
		case json.Number:
			leftFloat, _ := leftValue.Float64()
			if leftFloat > rightFloatValue {
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
