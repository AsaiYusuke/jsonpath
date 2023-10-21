package jsonpath

import "reflect"

type syntaxCompareEQ struct {
	*syntaxBasicAnyValueComparator
}

func (c *syntaxCompareEQ) comparator(left []interface{}, right interface{}) bool {
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
