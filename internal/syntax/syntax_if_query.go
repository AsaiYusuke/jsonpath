package syntax

type syntaxQuery interface {
	compute(root interface{}, currentList []interface{}) []interface{}
}
