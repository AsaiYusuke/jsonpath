package jsonpath

type syntaxComparator interface {
	comparator(left []interface{}, right interface{}) bool
	validate(values []interface{}) bool
}
