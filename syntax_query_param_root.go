package jsonpath

type syntaxQueryParamRoot struct {
	param syntaxNode
}

func (e *syntaxQueryParamRoot) isValueGroupParameter() bool {
	return e.param.isValueGroup()
}

func (e *syntaxQueryParamRoot) compute(
	root interface{}, currentList []interface{}) []interface{} {

	values := bufferContainer{}

	if err := e.param.retrieve(root, root, &values); err != nil {
		return []interface{}{struct{}{}}
	}

	if len(values.result) == 1 {
		return values.result
	}

	return []interface{}{true}
}
