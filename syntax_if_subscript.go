package jsonpath

type syntaxSubscript interface {
	getIndexes(src []interface{}) []int
	isMultiValue() bool
}
