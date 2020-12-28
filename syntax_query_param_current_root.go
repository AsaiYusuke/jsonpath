package jsonpath

type syntaxQueryParamCurrentRoot struct {
	param syntaxNode
}

func (e *syntaxQueryParamCurrentRoot) isMultiValueParameter() bool {
	return e.param.isMultiValue()
}

func (e *syntaxQueryParamCurrentRoot) compute(
	root interface{}, currentMap map[int]interface{}) map[int]interface{} {

	result := make(map[int]interface{}, len(currentMap))
	for index, srcNode := range currentMap {
		var values []interface{}
		if err := e.param.retrieve(root, srcNode, &values); err != nil {
			continue
		}
		if e.param.isMultiValue() {
			result[index] = values
		} else {
			result[index] = values[0]
		}
	}

	return result
}
