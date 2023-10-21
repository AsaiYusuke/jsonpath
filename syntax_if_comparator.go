package jsonpath

type syntaxComparator interface {
	comparator(left []interface{}, right interface{}) bool
	typeCast(values []interface{}) bool
}
