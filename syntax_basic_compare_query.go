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

	if leftFound {
		for leftIndex := range leftValues {
			if _, ok := leftValues[leftIndex].(struct{}); ok {
				continue
			}

			if !rightFound {
				if !q.leftParam.isLiteral {
					leftValues[leftIndex] = struct{}{}
				}
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

	} else if !q.rightParam.isLiteral {
		for rightIndex := range rightValues {
			rightValues[rightIndex] = struct{}{}
		}
	}

	if !q.leftParam.isLiteral {
		return leftValues
	}

	return rightValues
}
