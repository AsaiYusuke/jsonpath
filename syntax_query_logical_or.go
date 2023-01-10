package jsonpath

type syntaxLogicalOr struct {
	leftQuery  syntaxQuery
	rightQuery syntaxQuery
}

func (l *syntaxLogicalOr) compute(
	root interface{}, currentList []interface{}, container *bufferContainer) []interface{} {

	leftComputedList := l.leftQuery.compute(root, currentList, container)
	if len(leftComputedList) == 1 {
		if _, ok := leftComputedList[0].(struct{}); ok {
			return l.rightQuery.compute(root, currentList, container)
		}
		return leftComputedList
	}

	rightComputedList := l.rightQuery.compute(root, currentList, container)
	if len(rightComputedList) == 1 {
		if _, ok := rightComputedList[0].(struct{}); ok {
			return leftComputedList
		}
		return rightComputedList
	}

	for index := range leftComputedList {
		if _, ok := rightComputedList[index].(struct{}); !ok {
			leftComputedList[index] = rightComputedList[index]
		}
	}
	return leftComputedList
}
