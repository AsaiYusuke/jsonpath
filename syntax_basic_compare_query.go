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
		for leftIndex := range leftValues {
			if _, ok := leftValues[leftIndex].(struct{}); ok {
				continue
			}

			for rightIndex := range rightValues {
				if _, ok := rightValues[rightIndex].(struct{}); ok {
					continue
				}

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
					leftValues[0] = struct{}{}
					return leftValues[:1]
				}
			}
		}
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
	for leftIndex := range values {
		values[leftIndex] = struct{}{}
	}
}
