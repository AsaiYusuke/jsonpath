package jsonpath

type syntaxQueryParamCurrentRoot struct {
	param     syntaxNode
	srcJSON   **interface{}
	resultPtr *[]interface{}
}

func (e *syntaxQueryParamCurrentRoot) isMultiValueParameter() bool {
	return e.param.isMultiValue()
}

func (e *syntaxQueryParamCurrentRoot) compute(currentMap map[int]interface{}) map[int]interface{} {
	result := make(map[int]interface{}, len(currentMap))
	for index, srcNode := range currentMap {
		values := make([]interface{}, 0)
		e.resultPtr = &values
		if err := e.param.retrieve(srcNode); err != nil {
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
