package jsonpath

type syntaxLogicalNot struct {
	query syntaxQuery
}

func (l *syntaxLogicalNot) compute(
	root interface{}, currentMap map[int]interface{}) map[int]interface{} {

	computedMap := l.query.compute(root, currentMap)
	resultMap := make(map[int]interface{}, 0)
	for index := range currentMap {
		if _, ok := computedMap[index]; !ok {
			resultMap[index] = currentMap[index]
		}
	}

	return resultMap
}
