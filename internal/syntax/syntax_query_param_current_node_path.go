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

	buf := getNodeSlice()

	for index := range currentList {
		*buf = (*buf)[:0]
		if e.param.retrieve(root, currentList[index], buf) != nil {
			result[index] = emptyEntity
			continue
		}
		hasValue = true
		// If e.param.isValueGroup==true,
		// Only the first element is returned because it is an existence check.
		result[index] = (*buf)[0]
	}
	putNodeSlice(buf)

	if hasValue {
		return result
	}

	return emptyList
}
