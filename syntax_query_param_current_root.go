package jsonpath

type syntaxQueryParamCurrentRoot struct {
	param syntaxNode
}

func (e *syntaxQueryParamCurrentRoot) isValueGroupParameter() bool {
	return e.param.isValueGroup()
}

func (e *syntaxQueryParamCurrentRoot) compute(
	root interface{}, currentList []interface{}) []interface{} {

	result := make([]interface{}, len(currentList))

	for index := range currentList {
		var values []interface{}
		if err := e.param.retrieve(root, currentList[index], &values); err != nil {
			result[index] = struct{}{}
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
