package syntax

type syntaxCompareQuery struct {
	leftParam  syntaxCompareParameter
	rightParam syntaxCompareParameter
	comparator syntaxComparator
}

func (q *syntaxCompareQuery) compute(
	root any, currentList []any) []any {

	leftValues := q.leftParam.compute(root, currentList)
	if len(leftValues) == 1 && leftValues[0] == emptyEntity {
		if _, ok := q.comparator.(*syntaxCompareDeepEQ); !ok {
			return emptyList
		}
	}

	rightValues := q.rightParam.compute(root, currentList)

	// The syntax parser always results in a literal value on the right side as input.
	if q.comparator.comparator(leftValues, rightValues[0]) {
		return leftValues
	}

	if len(leftValues) == 1 && leftValues[0] == emptyEntity && rightValues[0] == emptyEntity {
		if _, ok := q.comparator.(*syntaxCompareDeepEQ); ok {
			return currentList
		}
	}

	return emptyList
}
