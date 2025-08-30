package syntax

type syntaxComparator interface {
	compare(left []any, right any) bool
}
