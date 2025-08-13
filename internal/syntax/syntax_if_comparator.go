package syntax

type syntaxComparator interface {
	comparator(left []any, right any) bool
	validate(values []any) bool
}
