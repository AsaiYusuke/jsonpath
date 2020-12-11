package jsonpath

type syntaxBasicCompareParameter struct {
	param     syntaxQuery
	isLiteral bool
}

func (p syntaxBasicCompareParameter) get(
	root interface{}, currentMap map[int]interface{}) map[int]interface{} {

	if _, ok := p.param.(syntaxQueryParamRoot); ok {
		currentMap = map[int]interface{}{0: root}
	}
	return p.param.compute(root, currentMap)
}
