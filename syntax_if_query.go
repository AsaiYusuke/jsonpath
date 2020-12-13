package jsonpath

type syntaxQuery interface {
	compute(currentMap map[int]interface{}) map[int]interface{}
}
