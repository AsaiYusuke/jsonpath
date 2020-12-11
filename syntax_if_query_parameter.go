package jsonpath

type syntaxQueryParameter interface {
	isMultiValueParameter() bool
	compute(root interface{}, currentMap map[int]interface{}) map[int]interface{}
}
