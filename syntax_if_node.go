package jsonpath

type syntaxNode interface {
	retrieve(root, current interface{}, result *[]interface{}) error
	setText(text string)
	setMultiValue()
	isMultiValue() bool
	getConnectedText() string
	setNext(next syntaxNode)
	getNext() syntaxNode
	retrieveNext(root interface{}, result *[]interface{}, getter func() interface{}, setter func(interface{})) error
	setAccessorMode(mode bool)
}
