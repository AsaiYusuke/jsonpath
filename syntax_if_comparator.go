package jsonpath

type syntaxComparator interface {
	comparator(left, right interface{}) bool
	typeCast(values []interface{})
}
