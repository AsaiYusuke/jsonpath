package syntax

type syntaxTypeValidator interface {
	validate(values []interface{}) bool
}
