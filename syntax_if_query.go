package jsonpath

type syntaxQuery interface {
	compute(root interface{}, currentList []interface{}) []interface{}
}
