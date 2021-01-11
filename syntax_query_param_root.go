package jsonpath

type syntaxQueryParamRoot struct {
	param syntaxNode
}

func (e *syntaxQueryParamRoot) isValueGroupParameter() bool {
	return e.param.isValueGroup()
}

func (e *syntaxQueryParamRoot) compute(
	root interface{}, currentList []interface{}) []interface{} {

	values := make([]interface{}, 0, 1)

	if err := e.param.retrieve(root, root, &values); err != nil {
		return []interface{}{}
	}

	// e.param.isValueGroup() should always be false.

	return values[:1]
}
