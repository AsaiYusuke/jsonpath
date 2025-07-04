package jsonpath

type syntaxBasicCompareParameter struct {
	param     syntaxQuery
	isLiteral bool
}

func (p *syntaxBasicCompareParameter) compute(
	root interface{}, currentList []interface{}) []interface{} {

	if _, ok := p.param.(*syntaxQueryParamRootNode); ok {
		currentList = []interface{}{root}
	}

	return p.param.compute(root, currentList)
}
