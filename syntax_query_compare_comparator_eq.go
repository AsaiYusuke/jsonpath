package jsonpath

import "reflect"

type syntaxCompareEQ struct {
	*syntaxBasicAnyValueComparator
}

func (c syntaxCompareEQ) comparator(left, right interface{}) bool {
	return reflect.DeepEqual(left, right)
}
