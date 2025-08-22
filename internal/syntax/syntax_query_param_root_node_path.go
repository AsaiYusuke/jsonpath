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
	defer func() { putNodeSlice(buf) }()

	if e.param.retrieve(root, root, buf) != nil {
		return emptyList
	}

	return []any{(*buf)[0]}
}
