package syntax

import "encoding/json"

type syntaxCompareGT struct {
}

func (c *syntaxCompareGT) comparator(left []any, right any) bool {
	rightFloat, rightIsFloat := right.(float64)
	rightNumber, rightIsNumber := right.(json.Number)
	rightString, rightIsString := right.(string)
	if !rightIsFloat && !rightIsNumber && !rightIsString {
		return false
	}
	if rightIsNumber {
		rightFloat, _ = rightNumber.Float64()
	}

	var hasValue bool
	for leftIndex := range left {
		if left[leftIndex] == emptyEntity {
			continue
		}
		switch leftValue := left[leftIndex].(type) {
		case float64:
			if (rightIsFloat || rightIsNumber) && leftValue > rightFloat {
				hasValue = true
			} else {
				left[leftIndex] = emptyEntity
			}
		case json.Number:
			leftFloat, _ := leftValue.Float64()
			if (rightIsFloat || rightIsNumber) && leftFloat > rightFloat {
				hasValue = true
			} else {
				left[leftIndex] = emptyEntity
			}
		case string:
			if rightIsString && leftValue > rightString {
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
