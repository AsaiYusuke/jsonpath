package jsonpath

type syntaxBasicCompareParameter struct {
	param     syntaxQuery
	isLiteral bool
}

func (p *syntaxBasicCompareParameter) get(currentMap map[int]interface{}) map[int]interface{} {

	if param, ok := p.param.(*syntaxQueryParamRoot); ok {
		currentMap = map[int]interface{}{0: **param.srcJSON}
	}
	return p.param.compute(currentMap)
}
