package jsonpath

type syntaxNode interface {
	retrieve(current interface{}) error

	setText(text string)
	setMultiValue()
	isMultiValue() bool
	getConnectedText() string
	setNext(next syntaxNode)
	getNext() syntaxNode
	retrieveNext(current interface{}) error
	setResultPtr(resultPtr **[]interface{})
}
