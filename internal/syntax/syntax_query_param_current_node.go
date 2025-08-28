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
	copy(result, currentList)
	return result
}
