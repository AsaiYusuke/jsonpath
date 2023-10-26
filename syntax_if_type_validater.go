package jsonpath

type syntaxTypeValidator interface {
	validate(values []interface{}) bool
}
