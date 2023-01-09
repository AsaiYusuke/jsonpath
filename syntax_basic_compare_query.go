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

	if !leftFound && !rightFound {
		if _, ok := q.comparator.(*syntaxCompareEQ); ok {
			return currentList
		}
	}

	return []interface{}{struct{}{}}
}
