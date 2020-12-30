package jsonpath

type syntaxBasicCompareQuery struct {
	leftParam  *syntaxBasicCompareParameter
	rightParam *syntaxBasicCompareParameter
	comparator syntaxComparator
}

func (q *syntaxBasicCompareQuery) compute(
	root interface{}, currentMap map[int]interface{}) map[int]interface{} {

	leftValues := q.leftParam.get(root, currentMap)
	for index := range leftValues {
		if cast, ok := q.comparator.typeCast(leftValues[index]); ok {
			leftValues[index] = cast
		} else {
			delete(leftValues, index)
		}
	}

	rightValues := q.rightParam.get(root, currentMap)
	for index := range rightValues {
		if cast, ok := q.comparator.typeCast(rightValues[index]); ok {
			rightValues[index] = cast
		} else {
			delete(rightValues, index)
		}
	}

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
