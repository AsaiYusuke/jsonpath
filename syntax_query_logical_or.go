package jsonpath

type syntaxLogicalOr struct {
	leftParam  syntaxQuery
	rightParam syntaxQuery
}

func (l *syntaxLogicalOr) compute(currentMap map[int]interface{}) map[int]interface{} {
	leftComputedMap := l.leftParam.compute(currentMap)
	rightComputedMap := l.rightParam.compute(currentMap)
	for index := range rightComputedMap {
		if _, ok := leftComputedMap[index]; !ok {
			leftComputedMap[index] = 1
		}
	}
	return leftComputedMap
}
