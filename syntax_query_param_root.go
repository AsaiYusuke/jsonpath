package jsonpath

type syntaxQueryParamRoot struct {
	param     syntaxNode
	srcJSON   **interface{}
	resultPtr *[]interface{}
}

func (e *syntaxQueryParamRoot) isMultiValueParameter() bool {
	return e.param.isMultiValue()
}

func (e *syntaxQueryParamRoot) compute(currentMap map[int]interface{}) map[int]interface{} {
	result := make(map[int]interface{}, len(currentMap))
	values := make([]interface{}, 0, 1)
	e.resultPtr = &values
	if err := e.param.retrieve(**e.srcJSON); err != nil {
		return result
	}
	// e.param.isMultiValue() should always be false.
	for index := range currentMap {
		result[index] = values[0]
	}
	return result
}
