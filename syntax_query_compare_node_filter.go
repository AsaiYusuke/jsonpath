package jsonpath

type syntaxNodeFilter struct {
	param syntaxNode
}

func (e syntaxNodeFilter) compute(root interface{}, currentMap map[int]interface{}) map[int]interface{} {
	result := make(map[int]interface{}, len(currentMap))
	for index, srcNode := range currentMap {
		values := resultContainer{}
		if err := e.param.retrieve(root, srcNode, &values); err != nil {
			continue
		}
		if !values.hasResult() {
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

func (e syntaxNodeFilter) isRoot() bool {
	_, ok := e.param.(syntaxRootIdentifier)
	return ok
}

func (e syntaxNodeFilter) isCurrentRoot() bool {
	_, ok := e.param.(syntaxCurrentRootIdentifier)
	return ok
}
