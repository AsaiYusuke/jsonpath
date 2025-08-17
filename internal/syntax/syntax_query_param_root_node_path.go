package syntax

type syntaxQueryParamRootNodePath struct {
	param syntaxNode
}

func (e *syntaxQueryParamRootNodePath) isValueGroupParameter() bool {
	return e.param.isValueGroup()
}

func (e *syntaxQueryParamRootNodePath) compute(
	root any, _ []any) []any {

	values := getContainer()
	defer func() {
		putContainer(values)
	}()

	if e.param.retrieve(root, root, values) != nil {
		return emptyList
	}

	return []any{values.result[0]}
}
