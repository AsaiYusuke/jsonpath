package jsonpath

type syntaxLogicalOr struct {
	leftQuery  syntaxQuery
	rightQuery syntaxQuery
}

func (l *syntaxLogicalOr) compute(
	root interface{}, currentList []interface{}) []interface{} {

	leftComputedList := l.leftQuery.compute(root, currentList)
	if len(leftComputedList) == 1 {
		if _, ok := leftComputedList[0].(struct{}); ok {
			return l.rightQuery.compute(root, currentList)
		}
		return leftComputedList
	}

	rightComputedList := l.rightQuery.compute(root, currentList)
	if len(rightComputedList) == 1 {
		if _, ok := rightComputedList[0].(struct{}); ok {
			return leftComputedList
		}
		return rightComputedList
	}

	for index := range rightComputedList {
		if _, ok := rightComputedList[index].(struct{}); !ok {
			leftComputedList[index] = rightComputedList[index]
		}
	}
	return leftComputedList
}
