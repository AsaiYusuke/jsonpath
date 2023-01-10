package jsonpath

type syntaxLogicalAnd struct {
	leftQuery  syntaxQuery
	rightQuery syntaxQuery
}

func (l *syntaxLogicalAnd) compute(
	root interface{}, currentList []interface{}, container *bufferContainer) []interface{} {

	leftComputedList := l.leftQuery.compute(root, currentList, container)
	if len(leftComputedList) == 1 {
		if _, ok := leftComputedList[0].(struct{}); ok {
			return leftComputedList
		}
		return l.rightQuery.compute(root, currentList, container)
	}

	rightComputedList := l.rightQuery.compute(root, currentList, container)
	if len(rightComputedList) == 1 {
		if _, ok := rightComputedList[0].(struct{}); ok {
			return rightComputedList
		}
		return leftComputedList
	}

	var hasValue bool
	for index := range rightComputedList {
		if _, ok := rightComputedList[index].(struct{}); ok {
			leftComputedList[index] = struct{}{}
			continue
		}
		if _, ok := leftComputedList[index].(struct{}); !ok {
			hasValue = true
		}
	}
	if hasValue {
		return leftComputedList
	}
	return []interface{}{struct{}{}}
}
