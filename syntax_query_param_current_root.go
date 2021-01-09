package jsonpath

type syntaxQueryParamCurrentRoot struct {
	param syntaxNode
}

func (e *syntaxQueryParamCurrentRoot) isValueGroupParameter() bool {
	return e.param != nil && e.param.isValueGroup()
}

func (e *syntaxQueryParamCurrentRoot) compute(
	root interface{}, currentMap map[int]interface{}) map[int]interface{} {

	result := make(map[int]interface{}, len(currentMap))

	if e.param == nil {
		for index := range currentMap {
			result[index] = currentMap[index]
		}
		return result
	}

	for index := range currentMap {
		var values []interface{}
		if err := e.param.retrieve(root, currentMap[index], &values); err != nil {
			continue
		}
		if e.param.isValueGroup() {
			result[index] = values
		} else {
			result[index] = values[0]
		}
	}

	return result
}
