package jsonpath

type syntaxLogicalNot struct {
	param syntaxQuery
}

func (l *syntaxLogicalNot) compute(
	root interface{}, currentMap map[int]interface{}) map[int]interface{} {

	computedMap := l.param.compute(root, currentMap)
	resultMap := make(map[int]interface{}, 0)
	for index, value := range currentMap {
		if _, ok := computedMap[index]; !ok {
			resultMap[index] = value
		}
	}

	return resultMap
}
