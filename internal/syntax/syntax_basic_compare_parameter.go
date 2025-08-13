package syntax

type syntaxBasicCompareParameter struct {
	param     syntaxQuery
	isLiteral bool
}

func (p *syntaxBasicCompareParameter) compute(
	root any, currentList []any) []any {

	if _, ok := p.param.(*syntaxQueryParamRootNode); ok {
		currentList = []any{root}
	}

	return p.param.compute(root, currentList)
}
