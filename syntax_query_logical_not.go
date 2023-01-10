package jsonpath

type syntaxLogicalNot struct {
	query syntaxQuery
}

func (l *syntaxLogicalNot) compute(
	root interface{}, currentList []interface{}, container *bufferContainer) []interface{} {

	computedList := l.query.compute(root, currentList, container)
	var hasValue bool
	for index := range computedList {
		if _, ok := computedList[index].(struct{}); ok {
			computedList[index] = true
			hasValue = true
		} else {
			computedList[index] = struct{}{}
		}
	}
	if hasValue {
		return computedList
	}
	return []interface{}{struct{}{}}
}
