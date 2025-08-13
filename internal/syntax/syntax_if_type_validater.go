package syntax

type syntaxTypeValidator interface {
	validate(values []any) bool
}
