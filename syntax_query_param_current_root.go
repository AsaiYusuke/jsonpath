package jsonpath

type syntaxQueryParamCurrentRoot struct {
	param syntaxNode
}

func (e *syntaxQueryParamCurrentRoot) isValueGroupParameter() bool {
	return e.param.isValueGroup()
}

func (e *syntaxQueryParamCurrentRoot) compute(
	root interface{}, currentList []interface{}, container *bufferContainer) []interface{} {

	result := make([]interface{}, len(currentList))

	var hasValue bool
	for index := range currentList {
		values := bufferContainer{}

		if err := e.param.retrieve(root, currentList[index], &values); err != nil {
			result[index] = struct{}{}
			continue
		}
		hasValue = true
		if e.param.isValueGroup() {
			result[index] = values.result
		} else {
			result[index] = values.result[0]
		}
	}
	if hasValue {
		return result
	}
	return []interface{}{struct{}{}}
}
