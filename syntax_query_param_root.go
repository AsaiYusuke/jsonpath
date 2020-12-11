package jsonpath

type syntaxQueryParamRoot struct {
	param syntaxRootIdentifier
}

func (e syntaxQueryParamRoot) compute(root interface{}, currentMap map[int]interface{}) map[int]interface{} {
	result := make(map[int]interface{}, 0)
	values := make([]interface{}, 0)
	if err := e.param.retrieve(root, root, &values); err != nil {
		return result
	}
	// e.param.isMultiValue() should always be false.
	for index := range currentMap {
		result[index] = values[0]
	}
	return result
}
