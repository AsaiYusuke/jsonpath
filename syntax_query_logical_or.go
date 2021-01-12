package jsonpath

type syntaxLogicalOr struct {
	leftQuery  syntaxQuery
	rightQuery syntaxQuery
}

func (l *syntaxLogicalOr) compute(
	root interface{}, currentList []interface{}, container *bufferContainer) []interface{} {

	leftComputedList := l.leftQuery.compute(root, currentList, container)
	rightComputedList := l.rightQuery.compute(root, currentList, container)
	for index := range rightComputedList {
		if _, ok := leftComputedList[index].(struct{}); ok {
			leftComputedList[index] = rightComputedList[index]
		}
	}

	return leftComputedList
}
