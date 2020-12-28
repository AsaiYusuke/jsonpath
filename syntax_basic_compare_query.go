package jsonpath

type syntaxBasicCompareQuery struct {
	leftParam  *syntaxBasicCompareParameter
	rightParam *syntaxBasicCompareParameter
	comparator syntaxComparator
}

func (q *syntaxBasicCompareQuery) compute(
	root interface{}, currentMap map[int]interface{}) map[int]interface{} {

	leftValues := q.leftParam.get(root, currentMap)
	for index, value := range leftValues {
		if cast, ok := q.comparator.typeCast(value); ok {
			leftValues[index] = cast
		} else {
			delete(leftValues, index)
		}
	}

	rightValues := q.rightParam.get(root, currentMap)
	for index, value := range rightValues {
		if cast, ok := q.comparator.typeCast(value); ok {
			rightValues[index] = cast
		} else {
			delete(rightValues, index)
		}
	}

	var leftIndex, rightIndex int
	var leftValue, rightValue interface{}
	for leftIndex, leftValue = range leftValues {
		for rightIndex, rightValue = range rightValues {
			if q.comparator.comparator(leftValue, rightValue) {
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
