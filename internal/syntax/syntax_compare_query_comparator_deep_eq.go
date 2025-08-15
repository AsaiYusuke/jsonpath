package syntax

import "reflect"

type syntaxCompareDeepEQ struct {
	*syntaxAnyValueTypeValidator
}

func (c *syntaxCompareDeepEQ) comparator(left []any, right any) bool {
	var hasValue bool
	for leftIndex := range left {
		if left[leftIndex] == emptyEntity {
			continue
		}
		if reflect.DeepEqual(left[leftIndex], right) {
			hasValue = true
		} else {
			left[leftIndex] = emptyEntity
		}
	}
	return hasValue
}
