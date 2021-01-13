package jsonpath

type syntaxQueryParamRoot struct {
	param syntaxNode
}

func (e *syntaxQueryParamRoot) isValueGroupParameter() bool {
	return e.param.isValueGroup()
}

func (e *syntaxQueryParamRoot) compute(
	root interface{}, currentList []interface{}, container *bufferContainer) []interface{} {

	values := bufferContainer{}

	if err := e.param.retrieve(root, root, &values); err != nil {
		return values.result[:0]
	}

	// e.param.isValueGroup() should always be false.

	return values.result[:1]
}
