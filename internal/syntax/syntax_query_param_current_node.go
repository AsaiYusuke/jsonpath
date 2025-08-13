package syntax

type syntaxQueryParamCurrentNode struct {
	param syntaxNode
}

func (e *syntaxQueryParamCurrentNode) isValueGroupParameter() bool {
	return e.param.isValueGroup()
}

func (e *syntaxQueryParamCurrentNode) compute(
	root any, currentList []any) []any {

	result := make([]any, len(currentList))

	var hasValue bool

	container := getContainer()
	defer func() {
		putContainer(container)
	}()

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

	if hasValue {
		return result
	}

	return emptyList
}
