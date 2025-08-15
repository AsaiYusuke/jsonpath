package syntax

type syntaxCompareParameter interface {
	compute(root any, currentList []any) []any
}
