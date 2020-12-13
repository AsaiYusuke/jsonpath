package jsonpath

type syntaxLogicalNot struct {
	param syntaxQuery
}

func (l *syntaxLogicalNot) compute(root interface{}, currentMap map[int]interface{}) map[int]interface{} {
	computedMap := l.param.compute(root, currentMap)
	resultMap := currentMap
	for index := range computedMap {
		delete(resultMap, index)
	}

	return resultMap
}
