package jsonpath

type syntaxLogicalAnd struct {
	leftParam  syntaxQuery
	rightParam syntaxQuery
}

func (l *syntaxLogicalAnd) compute(currentMap map[int]interface{}) map[int]interface{} {
	leftComputedMap := l.leftParam.compute(currentMap)
	rightComputedMap := l.rightParam.compute(currentMap)
	for index := range leftComputedMap {
		if _, ok := rightComputedMap[index]; !ok {
			delete(leftComputedMap, index)
		}
	}
	return leftComputedMap
}
