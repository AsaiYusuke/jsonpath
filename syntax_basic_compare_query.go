package jsonpath

type syntaxBasicCompareQuery struct {
	leftParam  *syntaxBasicCompareParameter
	rightParam *syntaxBasicCompareParameter
	comparator syntaxComparator
}

func (q *syntaxBasicCompareQuery) compute(
	root interface{}, currentMap map[int]interface{}) map[int]interface{} {

	leftValues := q.leftParam.compute(root, currentMap)
	q.comparator.typeCast(leftValues)

	rightValues := q.rightParam.compute(root, currentMap)
	q.comparator.typeCast(rightValues)

	for leftIndex := range leftValues {
		for rightIndex := range rightValues {
			if q.comparator.comparator(leftValues[leftIndex], rightValues[rightIndex]) {
				if q.leftParam.isLiteral && q.rightParam.isLiteral {
					return currentMap
				}
			} else {
				if !q.leftParam.isLiteral {
					delete(leftValues, leftIndex)
				} else {
					delete(rightValues, rightIndex)
				}
			}
		}
	}

	if !q.leftParam.isLiteral && len(rightValues) > 0 {
		return leftValues
	}

	if !q.rightParam.isLiteral && len(leftValues) > 0 {
		return rightValues
	}

	return nil
}
