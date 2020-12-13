package jsonpath

type syntaxQueryParameter interface {
	isMultiValueParameter() bool
	compute(currentMap map[int]interface{}) map[int]interface{}
}
