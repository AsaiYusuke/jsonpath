package jsonpath

type syntaxBasicCompareQuery struct {
	leftParam  syntaxQuery
	rightParam syntaxQuery
	comparator syntaxComparator
}

func (q syntaxBasicCompareQuery) compute(root interface{}, currentMap map[int]interface{}) map[int]interface{} {
	isLeftLiteral, leftValues := q.getComputeParameters(root, currentMap, q.leftParam)
	isRightLiteral, rightValues := q.getComputeParameters(root, currentMap, q.rightParam)

	result := make(map[int]interface{}, len(leftValues))

	var leftIndex, rightIndex int
	var leftValue, rightValue interface{}
	for leftIndex, leftValue = range leftValues {
		for rightIndex, rightValue = range rightValues {
			if q.comparator.comparator(leftValue, rightValue) {
				if isLeftLiteral && isRightLiteral {
					return currentMap
				} else if !isLeftLiteral {
					result[leftIndex] = leftValue
				} else {
					result[rightIndex] = rightValue
				}
			}
		}
	}

	return result
}

func (q *syntaxBasicCompareQuery) getComputeParameters(
	root interface{}, currentMap map[int]interface{}, param syntaxQuery) (bool, map[int]interface{}) {

	var isLiteral bool

	switch param.(type) {
	case syntaxQueryParamLiteral:
		isLiteral = true
	case syntaxQueryParamRoot:
		isLiteral = true
		currentMap = map[int]interface{}{0: root}
	}

	computedValues := param.compute(root, currentMap)

	for index, value := range computedValues {
		if cast, ok := q.comparator.typeCast(value); ok {
			computedValues[index] = cast
		} else {
			delete(computedValues, index)
		}
	}

	return isLiteral, computedValues
}
