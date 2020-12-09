package jsonpath

type syntaxQueryParamRoot struct {
	param syntaxRootIdentifier
}

func (e syntaxQueryParamRoot) compute(root interface{}, currentMap map[int]interface{}) map[int]interface{} {
	result := make(map[int]interface{}, 0)
	values := resultContainer{}
	if err := e.param.retrieve(root, root, &values); err != nil {
		return result
	}
	var _result interface{}
	// e.param.isMultiValue() should always be false.
	_result = values.getResult()[0]
	for index := range currentMap {
		result[index] = _result
	}
	return result
}
