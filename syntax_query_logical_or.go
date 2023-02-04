package jsonpath

type syntaxLogicalOr struct {
	leftQuery  syntaxQuery
	rightQuery syntaxQuery
}

func (l *syntaxLogicalOr) compute(
	root interface{}, currentList []interface{}) []interface{} {

	leftComputedList := l.leftQuery.compute(root, currentList)
	if len(leftComputedList) == 1 {
		if leftComputedList[0] == struct{}{} {
			return l.rightQuery.compute(root, currentList)
		}
		return leftComputedList
	}

	rightComputedList := l.rightQuery.compute(root, currentList)
	if len(rightComputedList) == 1 {
		if rightComputedList[0] == struct{}{} {
			return leftComputedList
		}
		return rightComputedList
	}

	for index := range rightComputedList {
		if rightComputedList[index] != struct{}{} {
			leftComputedList[index] = rightComputedList[index]
		}
	}
	return leftComputedList
}
