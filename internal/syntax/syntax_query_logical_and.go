package syntax

type syntaxLogicalAnd struct {
	leftQuery  syntaxQuery
	rightQuery syntaxQuery
}

func (l *syntaxLogicalAnd) compute(
	root any, currentList []any, scratch *syntaxQueryScratch) []any {

	leftComputedList := l.leftQuery.compute(root, currentList, scratch)
	if len(leftComputedList) == 1 {
		if leftComputedList[0] == emptyEntity {
			return leftComputedList
		}
		return l.rightQuery.compute(root, currentList, scratch)
	}

	rightComputedList := l.rightQuery.compute(root, currentList, scratch)
	if len(rightComputedList) == 1 {
		if rightComputedList[0] == emptyEntity {
			return rightComputedList
		}
		return leftComputedList
	}

	var hasValue bool
	for index := range rightComputedList {
		if rightComputedList[index] == emptyEntity {
			leftComputedList[index] = emptyEntity
			continue
		}
		if leftComputedList[index] != emptyEntity {
			hasValue = true
		}
	}
	if hasValue {
		return leftComputedList
	}
	return emptyList
}
