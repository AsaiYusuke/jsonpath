package jsonpath

type syntaxComparator interface {
	comparator(left, right interface{}) bool
}
