package jsonpath

type syntaxLogicalNot struct {
	param syntaxQuery
}

func (l *syntaxLogicalNot) compute(currentMap map[int]interface{}) map[int]interface{} {
	computedMap := l.param.compute(currentMap)
	resultMap := currentMap
	for index := range computedMap {
		delete(resultMap, index)
	}

	return resultMap
}
