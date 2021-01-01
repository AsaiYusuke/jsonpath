package jsonpath

type syntaxQueryParamCurrentRoot struct {
	param syntaxNode
}

func (e *syntaxQueryParamCurrentRoot) isValueGroupParameter() bool {
	return e.param.isValueGroup()
}

func (e *syntaxQueryParamCurrentRoot) compute(
	root interface{}, currentMap map[int]interface{}) map[int]interface{} {

	result := make(map[int]interface{}, len(currentMap))
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
