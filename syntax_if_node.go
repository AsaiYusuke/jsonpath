package jsonpath

type syntaxNode interface {
	retrieve(root, current interface{}, result *[]interface{}) error
	setText(text string)
	setValueGroup()
	isValueGroup() bool
	getConnectedText() string
	setNext(next syntaxNode)
	getNext() syntaxNode
	setAccessorMode(mode bool)
}
