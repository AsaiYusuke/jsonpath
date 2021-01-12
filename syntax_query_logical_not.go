package jsonpath

type syntaxLogicalNot struct {
	query syntaxQuery
}

func (l *syntaxLogicalNot) compute(
	root interface{}, currentList []interface{}, container *bufferContainer) []interface{} {

	computedList := l.query.compute(root, currentList, container)
	for index := range computedList {
		if _, ok := computedList[index].(struct{}); ok {
			computedList[index] = true
		} else {
			computedList[index] = struct{}{}
		}
	}

	return computedList
}
