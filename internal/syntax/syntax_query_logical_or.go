package syntax

type syntaxLogicalOr struct {
	leftQuery  syntaxQuery
	rightQuery syntaxQuery
}

func (l *syntaxLogicalOr) compute(
	root any, currentList []any, scratch *syntaxQueryScratch) []any {

	leftComputedList := l.leftQuery.compute(root, currentList, scratch)
	if len(leftComputedList) == 1 {
		if leftComputedList[0] == emptyEntity {
			return l.rightQuery.compute(root, currentList, scratch)
		}
		return leftComputedList
	}

	rightComputedList := l.rightQuery.compute(root, currentList, scratch)
	if len(rightComputedList) == 1 {
		if rightComputedList[0] == emptyEntity {
			return leftComputedList
		}
		return rightComputedList
	}

	for index := range rightComputedList {
		if rightComputedList[index] != emptyEntity {
			leftComputedList[index] = rightComputedList[index]
		}
	}
	return leftComputedList
}
