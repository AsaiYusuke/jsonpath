package jsonpath

type syntaxLogicalOr struct {
	leftParam  syntaxQuery
	rightParam syntaxQuery
}

func (l syntaxLogicalOr) compute(root interface{}, currentMap map[int]interface{}) map[int]interface{} {
	leftComputedMap := l.leftParam.compute(root, currentMap)
	rightComputedMap := l.rightParam.compute(root, currentMap)
	for index := range rightComputedMap {
		if _, ok := leftComputedMap[index]; !ok {
			leftComputedMap[index] = 1
		}
	}
	return leftComputedMap
}
