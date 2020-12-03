package jsonpath

type syntaxBasicCompareQuery struct {
	leftParam  syntaxQuery
	rightParam syntaxQuery
	comparator syntaxComparator
}

func (q syntaxBasicCompareQuery) compute(root interface{}, currentMap map[int]interface{}) map[int]interface{} {
	isLeftLiteral, leftValues := q.getComputeParameters(root, currentMap, q.leftParam)
	isRightLiteral, rightValues := q.getComputeParameters(root, currentMap, q.rightParam)

	if isLeftLiteral && isRightLiteral {
		leftValue := leftValues[0]
		rightValue := rightValues[0]
		if q.comparator.comparator(leftValue, rightValue) {
			return currentMap
		}
		return nil
	}

	if !isLeftLiteral && isRightLiteral {
		rightValue := rightValues[0]
		result := make(map[int]interface{}, len(leftValues))
		for leftIndex, leftValue := range leftValues {
			if q.comparator.comparator(leftValue, rightValue) {
				result[leftIndex] = leftValue
			}
		}
		return result
	}

	if isLeftLiteral && !isRightLiteral {
		leftValue := leftValues[0]
		result := make(map[int]interface{}, len(rightValues))
		for rightIndex, rightValue := range rightValues {
			if q.comparator.comparator(leftValue, rightValue) {
				result[rightIndex] = rightValue
			}
		}
		return result
	}

	return nil
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

	return isLiteral, param.compute(root, currentMap)
}
