package jsonpath

type syntaxCurrentRootNodeFilter struct {
	param syntaxCurrentRootIdentifier
}

func (e syntaxCurrentRootNodeFilter) compute(root interface{}, currentMap map[int]interface{}) map[int]interface{} {
	result := make(map[int]interface{}, len(currentMap))
	for index, srcNode := range currentMap {
		values := resultContainer{}
		if err := e.param.retrieve(root, srcNode, &values); err != nil {
			continue
		}
		_result := values.getResult()
		if e.param.isMultiValue() {
			result[index] = _result
		} else {
			result[index] = _result[0]
		}
	}
	return result
}
