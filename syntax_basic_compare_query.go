package jsonpath

type syntaxBasicCompareQuery struct {
	leftParam  *syntaxBasicCompareParameter
	rightParam *syntaxBasicCompareParameter
	comparator syntaxComparator
}

func (q *syntaxBasicCompareQuery) compute(
	root interface{}, currentList []interface{}) []interface{} {

	leftValues := q.leftParam.compute(root, currentList)
	leftFound := q.comparator.validate(leftValues)

	rightValues := q.rightParam.compute(root, currentList)
	rightFound := q.comparator.validate(rightValues)

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
