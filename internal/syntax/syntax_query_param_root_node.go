package syntax

type syntaxQueryParamRootNode struct {
	param syntaxNode
}

func (e *syntaxQueryParamRootNode) isValueGroupParameter() bool {
	return e.param.isValueGroup()
}

func (e *syntaxQueryParamRootNode) compute(
	_ any, _ []any) []any {

	return fullList
}
