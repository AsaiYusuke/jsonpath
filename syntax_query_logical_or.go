package jsonpath

type syntaxLogicalOr struct {
	leftParam  syntaxQuery
	rightParam syntaxQuery
}

func (l syntaxLogicalOr) compute(root interface{}, currentMap map[int]interface{}) map[int]interface{} {
	computedMap1 := l.leftParam.compute(root, currentMap)
	computedMap2 := l.rightParam.compute(root, currentMap)
	resultMap := computedMap1
	for index := range computedMap2 {
		if _, ok := computedMap1[index]; !ok {
			resultMap[index] = 1
		}
	}
	return resultMap
}
