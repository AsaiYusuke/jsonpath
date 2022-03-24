package jsonpath

type syntaxBasicCompareQuery struct {
	leftParam  *syntaxBasicCompareParameter
	rightParam *syntaxBasicCompareParameter
	comparator syntaxComparator
}

func (q *syntaxBasicCompareQuery) compute(
	root interface{}, currentList []interface{}, container *bufferContainer) []interface{} {

	leftValues := q.leftParam.compute(root, currentList, container)
	leftFound := q.comparator.typeCast(leftValues)

	rightValues := q.rightParam.compute(root, currentList, container)
	rightFound := q.comparator.typeCast(rightValues)

	if leftFound && rightFound {
		if q.leftParam.isLiteral {
			for rightIndex := range rightValues {
				if _, ok := rightValues[rightIndex].(struct{}); ok {
					continue
				}

				if !q.comparator.comparator(leftValues[0], rightValues[rightIndex]) {
					rightValues[rightIndex] = struct{}{}
				}
			}
			return rightValues
		}

		if q.rightParam.isLiteral {
			for leftIndex := range leftValues {
				if _, ok := leftValues[leftIndex].(struct{}); ok {
					continue
				}

				if !q.comparator.comparator(leftValues[leftIndex], rightValues[0]) {
					leftValues[leftIndex] = struct{}{}
				}
			}
			return leftValues
		}

		for leftIndex := range leftValues {
			if _, ok := leftValues[leftIndex].(struct{}); ok {
				continue
			}

			if !q.comparator.comparator(leftValues[leftIndex], rightValues[leftIndex]) {
				leftValues[leftIndex] = struct{}{}
			}
		}
		return leftValues
	}

	if !q.leftParam.isLiteral {
		if !rightFound {
			q.setBlankValues(leftValues)
		}
		return leftValues
	}

	if !leftFound {
		q.setBlankValues(rightValues)
	}
	return rightValues
}

func (q *syntaxBasicCompareQuery) setBlankValues(values []interface{}) {
	for index := range values {
		values[index] = struct{}{}
	}
}
