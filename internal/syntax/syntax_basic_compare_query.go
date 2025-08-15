package syntax

type syntaxBasicCompareQuery struct {
	leftParam  syntaxCompareParameter
	rightParam syntaxCompareParameter
	comparator syntaxComparator
}

func (q *syntaxBasicCompareQuery) compute(
	root any, currentList []any) []any {

	leftValues := q.leftParam.compute(root, currentList)
	leftFound := !(len(leftValues) == 1 && leftValues[0] == emptyEntity)
	leftFound = leftFound && q.comparator.validate(leftValues)

	rightValues := q.rightParam.compute(root, currentList)
	rightFound := !(len(rightValues) == 1 && rightValues[0] == emptyEntity)
	rightFound = rightFound && q.comparator.validate(rightValues)

	if leftFound && rightFound {
		// The syntax parser always results in a literal value on the right side as input.
		if q.comparator.comparator(leftValues, rightValues[0]) {
			return leftValues
		}
		return emptyList
	}

	// leftFound == false && rightFound == false
	if leftFound == rightFound {
		if _, ok := q.comparator.(*syntaxCompareDeepEQ); ok {
			return currentList
		}
	}

	return emptyList
}
