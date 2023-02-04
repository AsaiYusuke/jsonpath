package jsonpath

type syntaxLogicalAnd struct {
	leftQuery  syntaxQuery
	rightQuery syntaxQuery
}

func (l *syntaxLogicalAnd) compute(
	root interface{}, currentList []interface{}) []interface{} {

	leftComputedList := l.leftQuery.compute(root, currentList)
	if len(leftComputedList) == 1 {
		if leftComputedList[0] == struct{}{} {
			return leftComputedList
		}
		return l.rightQuery.compute(root, currentList)
	}

	rightComputedList := l.rightQuery.compute(root, currentList)
	if len(rightComputedList) == 1 {
		if rightComputedList[0] == struct{}{} {
			return rightComputedList
		}
		return leftComputedList
	}

	var hasValue bool
	for index := range rightComputedList {
		if rightComputedList[index] == struct{}{} {
			leftComputedList[index] = struct{}{}
			continue
		}
		if leftComputedList[index] != struct{}{} {
			hasValue = true
		}
	}
	if hasValue {
		return leftComputedList
	}
	return emptyList
}
