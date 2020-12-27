package jsonpath

type syntaxNode interface {
	retrieve(current interface{}) error

	setText(text string)
	setMultiValue()
	isMultiValue() bool
	getConnectedText() string
	setNext(next syntaxNode)
	getNext() syntaxNode
	retrieveNext(getter func() interface{}, setter func(interface{})) error
	setResultPtr(resultPtr **[]interface{})
	setAccessorMode(mode bool)
}
