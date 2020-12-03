package jsonpath

type syntaxLogicalAnd struct {
	leftParam  syntaxQuery
	rightParam syntaxQuery
}

func (l syntaxLogicalAnd) compute(root interface{}, currentMap map[int]interface{}) map[int]interface{} {
	computedMap1 := l.leftParam.compute(root, currentMap)
	computedMap2 := l.rightParam.compute(root, currentMap)
	resultMap := make(map[int]interface{}, len(computedMap1))
	for index := range computedMap1 {
		if _, ok := computedMap2[index]; ok {
			resultMap[index] = 1
		}
	}
	return resultMap
}
