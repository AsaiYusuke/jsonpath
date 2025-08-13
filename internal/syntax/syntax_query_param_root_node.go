package syntax

type syntaxQueryParamRootNode struct {
	param syntaxNode
}

func (e *syntaxQueryParamRootNode) isValueGroupParameter() bool {
	return e.param.isValueGroup()
}

func (e *syntaxQueryParamRootNode) compute(
	root any, _ []any) []any {

	values := getContainer()
	defer func() {
		putContainer(values)
	}()

	if e.param.retrieve(root, root, values) != nil {
		return emptyList
	}

	if len(values.result) == 1 {
		return []any{values.result[0]}
	}

	return fullList
}
