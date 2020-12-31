package jsonpath

type syntaxLogicalOr struct {
	leftQuery  syntaxQuery
	rightQuery syntaxQuery
}

func (l *syntaxLogicalOr) compute(
	root interface{}, currentMap map[int]interface{}) map[int]interface{} {

	leftComputedMap := l.leftQuery.compute(root, currentMap)
	rightComputedMap := l.rightQuery.compute(root, currentMap)
	for index := range rightComputedMap {
		if _, ok := leftComputedMap[index]; !ok {
			leftComputedMap[index] = struct{}{}
		}
	}

	return leftComputedMap
}
