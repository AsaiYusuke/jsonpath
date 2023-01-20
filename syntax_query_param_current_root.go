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
	containers := make([]bufferContainer, len(currentList))

	var hasValue bool
	for index := range currentList {
		if err := e.param.retrieve(root, currentList[index], &containers[index]); err != nil {
			result[index] = struct{}{}
			continue
		}
		hasValue = true
		if e.param.isValueGroup() {
			result[index] = containers[index].result
		} else {
			result[index] = containers[index].result[0]
		}
	}
	if hasValue {
		return result
	}
	return []interface{}{struct{}{}}
}
