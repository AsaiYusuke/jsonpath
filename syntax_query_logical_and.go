package jsonpath

type syntaxLogicalAnd struct {
	leftQuery  syntaxQuery
	rightQuery syntaxQuery
}

func (l *syntaxLogicalAnd) compute(
	root interface{}, currentMap map[int]interface{}) map[int]interface{} {

	leftComputedMap := l.leftQuery.compute(root, currentMap)
	rightComputedMap := l.rightQuery.compute(root, currentMap)
	for index := range leftComputedMap {
		if _, ok := rightComputedMap[index]; !ok {
			delete(leftComputedMap, index)
		}
	}

	return leftComputedMap
}
