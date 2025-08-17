package syntax

type syntaxQueryParamRootNode struct {
	param syntaxNode
}

// Dummy to satisfy syntaxQueryJSONPathParameter; not used in normal paths.
func (e *syntaxQueryParamRootNode) isValueGroupParameter() bool {
	return false
}

func (e *syntaxQueryParamRootNode) compute(
	_ any, _ []any) []any {

	return fullList
}
