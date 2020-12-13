package jsonpath

type syntaxNode interface {
	retrieve(root, current interface{}) error

	setText(text string)
	setMultiValue()
	isMultiValue() bool
	getConnectedText() string
	setNext(next syntaxNode)
	getNext() syntaxNode
	retrieveNext(root, current interface{}) error
	setResultPtr(resultPtr **[]interface{})
}
