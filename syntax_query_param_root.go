package jsonpath

type syntaxQueryParamRoot struct {
	param syntaxNode
}

func (e *syntaxQueryParamRoot) isValueGroupParameter() bool {
	return e.param.isValueGroup()
}

func (e *syntaxQueryParamRoot) compute(
	root interface{}, _ []interface{}) []interface{} {

	values := getContainer()
	defer func() {
		putContainer(values)
	}()

	if e.param.retrieve(root, root, values) != nil {
		return emptyList
	}

	if len(values.result) == 1 {
		return []interface{}{values.result[0]}
	}

	return fullList
}
