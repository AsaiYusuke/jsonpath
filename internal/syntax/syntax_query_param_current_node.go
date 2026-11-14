package syntax

type syntaxQueryParamCurrentNode struct {
	param syntaxNode
}

func (e *syntaxQueryParamCurrentNode) isValueGroupParameter() bool {
	return e.param.isValueGroup()
}

func (e *syntaxQueryParamCurrentNode) compute(
	_ any, currentList []any, scratch *syntaxQueryScratch) []any {

	result := scratch.getSlice(len(currentList))
	copy(result, currentList)
	return result
}
