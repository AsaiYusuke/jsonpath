package syntax

import "encoding/json"

type syntaxCompareGE struct {
}

func (c *syntaxCompareGE) comparator(left []any, right any) bool {
	rightFloatValue, rightIsFloat := right.(float64)
	rightNumberValue, rightIsNumber := right.(json.Number)
	rightStringValue, rightIsString := right.(string)
	if !rightIsFloat && !rightIsNumber && !rightIsString {
		return false
	}
	if rightIsNumber {
		rightFloatValue, _ = rightNumberValue.Float64()
	}

	var hasValue bool
	for leftIndex := range left {
		if left[leftIndex] == emptyEntity {
			continue
		}
		switch leftValue := left[leftIndex].(type) {
		case float64:
			if (rightIsFloat || rightIsNumber) && leftValue >= rightFloatValue {
				hasValue = true
			} else {
				left[leftIndex] = emptyEntity
			}
		case json.Number:
			leftFloat, _ := leftValue.Float64()
			if (rightIsFloat || rightIsNumber) && leftFloat >= rightFloatValue {
				hasValue = true
			} else {
				left[leftIndex] = emptyEntity
			}
		case string:
			if rightIsString && leftValue >= rightStringValue {
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
