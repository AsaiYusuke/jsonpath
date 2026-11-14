package syntax

type syntaxCompareParameter interface {
	compute(root any, currentList []any, scratch *syntaxQueryScratch) []any
}
