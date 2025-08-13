package syntax

type syntaxQuery interface {
	compute(root any, currentList []any) []any
}
