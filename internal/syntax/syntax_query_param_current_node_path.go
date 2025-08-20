package syntax

type syntaxQueryParamCurrentNodePath struct {
	param syntaxNode
}

func (e *syntaxQueryParamCurrentNodePath) isValueGroupParameter() bool {
	return e.param.isValueGroup()
}

func (e *syntaxQueryParamCurrentNodePath) compute(
	root any, currentList []any) []any {

	result := make([]any, len(currentList))

	var hasValue bool

	container := getContainer()

	for index := range currentList {
		container.result = container.result[:0]
		if e.param.retrieve(root, currentList[index], container) != nil {
			result[index] = emptyEntity
			continue
		}
		hasValue = true
		// If e.param.isValueGroup==true,
		// Only the first element is returned because it is an existence check.
		result[index] = container.result[0]
	}

	putContainer(container)

	if hasValue {
		return result
	}

	return emptyList
}
