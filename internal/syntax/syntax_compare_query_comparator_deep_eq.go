package syntax

import "reflect"

type syntaxCompareDeepEQ struct {
}

func (c *syntaxCompareDeepEQ) compare(left []any, right any) bool {
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
