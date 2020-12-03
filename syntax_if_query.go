package jsonpath

type syntaxQuery interface {
	compute(root interface{}, currentMap map[int]interface{}) map[int]interface{}
}
