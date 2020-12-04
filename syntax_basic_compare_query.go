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
		leftValue, leftOk := q.comparator.typeCast(leftValues[0])
		rightValue, rightOk := q.comparator.typeCast(rightValues[0])
		if leftOk && rightOk && q.comparator.comparator(leftValue, rightValue) {
			return currentMap
		}
		return nil
	}

	if !isLeftLiteral && isRightLiteral {
		result := make(map[int]interface{}, len(leftValues))
		rightValue, rightOk := q.comparator.typeCast(rightValues[0])
		if rightOk {
			for leftIndex, leftValue := range leftValues {
				leftValue, leftOk := q.comparator.typeCast(leftValue)
				if leftOk && q.comparator.comparator(leftValue, rightValue) {
					result[leftIndex] = leftValue
				}
			}
		}
		return result
	}

	if isLeftLiteral && !isRightLiteral {
		result := make(map[int]interface{}, len(rightValues))
		leftValue, leftOk := q.comparator.typeCast(leftValues[0])
		if leftOk {
			for rightIndex, rightValue := range rightValues {
				rightValue, rightOk := q.comparator.typeCast(rightValue)
				if rightOk && q.comparator.comparator(leftValue, rightValue) {
					result[rightIndex] = rightValue
				}
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
