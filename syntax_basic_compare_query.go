package jsonpath

type syntaxBasicCompareQuery struct {
	leftParam  *syntaxBasicCompareParameter
	rightParam *syntaxBasicCompareParameter
	comparator syntaxComparator
}

func (q *syntaxBasicCompareQuery) compute(
	root interface{}, currentList []interface{}) []interface{} {

	leftValues := q.leftParam.compute(root, currentList)
	q.comparator.typeCast(leftValues)

	rightValues := q.rightParam.compute(root, currentList)
	q.comparator.typeCast(rightValues)

	var leftPartialFound bool
	var rightPartialFound bool
	for leftIndex := range leftValues {
		if _, ok := leftValues[leftIndex].(struct{}); ok {
			continue
		}
		leftPartialFound = true

		for rightIndex := range rightValues {
			if _, ok := rightValues[rightIndex].(struct{}); ok {
				continue
			}
			rightPartialFound = true

			if q.comparator.comparator(leftValues[leftIndex], rightValues[rightIndex]) {
				if q.leftParam.isLiteral && q.rightParam.isLiteral {
					return leftValues
				}
				continue
			}

			if !q.leftParam.isLiteral {
				leftValues[leftIndex] = struct{}{}
				break
			} else if !q.rightParam.isLiteral {
				rightValues[rightIndex] = struct{}{}
			} else {
				return []interface{}{struct{}{}}
			}
		}

		if !rightPartialFound && !q.leftParam.isLiteral {
			leftValues[leftIndex] = struct{}{}
		}
	}

	if !leftPartialFound && !q.rightParam.isLiteral {
		for rightIndex := range rightValues {
			rightValues[rightIndex] = struct{}{}
		}
	}

	if !q.leftParam.isLiteral {
		return leftValues
	}

	return rightValues
}
