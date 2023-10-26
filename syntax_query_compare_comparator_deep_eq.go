package jsonpath

import "reflect"

type syntaxCompareDeepEQ struct {
	*syntaxBasicAnyValueTypeValidator
}

func (c *syntaxCompareDeepEQ) comparator(left []interface{}, right interface{}) bool {
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
