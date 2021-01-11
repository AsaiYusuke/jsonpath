package jsonpath

type syntaxLogicalOr struct {
	leftQuery  syntaxQuery
	rightQuery syntaxQuery
}

func (l *syntaxLogicalOr) compute(
	root interface{}, currentList []interface{}) []interface{} {

	leftComputedList := l.leftQuery.compute(root, currentList)
	rightComputedList := l.rightQuery.compute(root, currentList)
	for index := range rightComputedList {
		if _, ok := leftComputedList[index].(struct{}); ok {
			leftComputedList[index] = rightComputedList[index]
		}
	}

	return leftComputedList
}
