package jsonpath

type syntaxBasicCompareParameter struct {
	param     syntaxQuery
	isLiteral bool
}

func (p *syntaxBasicCompareParameter) compute(
	root interface{}, currentList []interface{}, container *bufferContainer) []interface{} {

	if _, ok := p.param.(*syntaxQueryParamRoot); ok {
		currentList = []interface{}{root}
	}

	return p.param.compute(root, currentList, container)
}
