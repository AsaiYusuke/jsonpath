package syntax

import "encoding/json"

type syntaxCompareDirectEQ struct {
}

func (c *syntaxCompareDirectEQ) comparator(left []any, right any) bool {
	var hasValue bool
	for leftIndex := range left {
		if left[leftIndex] == emptyEntity {
			continue
		}
		switch leftValue := left[leftIndex].(type) {
		case json.Number:
			leftFloat, _ := leftValue.Float64()
			if leftFloat == right {
				hasValue = true
			} else {
				left[leftIndex] = emptyEntity
			}
		default:
			if left[leftIndex] == right {
				hasValue = true
			} else {
				left[leftIndex] = emptyEntity
			}
		}
	}
	return hasValue
}
