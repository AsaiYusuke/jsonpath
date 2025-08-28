package syntax

type syntaxQueryParamRootNodePath struct {
	param syntaxNode
}

func (e *syntaxQueryParamRootNodePath) isValueGroupParameter() bool {
	return e.param.isValueGroup()
}

func (e *syntaxQueryParamRootNodePath) compute(
	root any, _ []any) []any {

	buf := getNodeSlice()

	if e.param.retrieve(root, root, buf) != nil {
		putNodeSlice(buf)
		return emptyList
	}

	value := (*buf)[0]
	putNodeSlice(buf)
	return []any{value}
}
