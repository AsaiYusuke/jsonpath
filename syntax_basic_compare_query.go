package jsonpath

type syntaxBasicCompareQuery struct {
	leftParam  *syntaxBasicCompareParameter
	rightParam *syntaxBasicCompareParameter
	comparator syntaxComparator
}

func (q *syntaxBasicCompareQuery) compute(
	root interface{}, currentList []interface{}) []interface{} {

	leftValues := q.leftParam.compute(root, currentList)
	leftFound := q.comparator.typeCast(leftValues)

	rightValues := q.rightParam.compute(root, currentList)
	rightFound := q.comparator.typeCast(rightValues)

	if leftFound && rightFound {
		var hasValue bool
		// The syntax parser always results in a literal value on the right side as input.
		for leftIndex := range leftValues {
			if leftValues[leftIndex] == struct{}{} {
				continue
			}
			if q.comparator.comparator(leftValues[leftIndex], rightValues[0]) {
				hasValue = true
			} else {
				leftValues[leftIndex] = struct{}{}
			}
		}
		if hasValue {
			return leftValues
		}
		return emptyList
	}

	// leftFound == false && rightFound == false
	if leftFound == rightFound {
		if _, ok := q.comparator.(*syntaxCompareEQ); ok {
			return currentList
		}
	}

	return emptyList
}
