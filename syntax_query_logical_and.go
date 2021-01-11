package jsonpath

type syntaxLogicalAnd struct {
	leftQuery  syntaxQuery
	rightQuery syntaxQuery
}

func (l *syntaxLogicalAnd) compute(
	root interface{}, currentList []interface{}) []interface{} {

	leftComputedList := l.leftQuery.compute(root, currentList)
	rightComputedList := l.rightQuery.compute(root, currentList)
	for index := range leftComputedList {
		if _, ok := rightComputedList[index].(struct{}); ok {
			leftComputedList[index] = struct{}{}
		}
	}

	return leftComputedList
}
