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
	for leftIndex, leftValue := range leftValues {
		for rightIndex, rightValue := range rightValues {
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
	root interface{}, currentMap map[int]interface{},
	param syntaxQuery) (bool, map[int]interface{}) {

	var isLiteral bool

	switch param.(type) {
	case syntaxCompareLiteral:
		isLiteral = true
	case syntaxNodeFilter:
		filter := param.(syntaxNodeFilter)
		if filter.isRoot() {
			isLiteral = true
			currentMap = map[int]interface{}{0: root}
		}
	}

	computedValues := param.compute(root, currentMap)
	result := make(map[int]interface{}, len(computedValues))
	for index, value := range computedValues {
		if cast, ok := q.comparator.typeCast(value); ok {
			result[index] = cast
		}
	}

	return isLiteral, result
}
