package jsonpath

type syntaxNode interface {
	retrieve(root, current interface{}, result *[]interface{}) error

	setText(text string)
	setMultiValue()
	isMultiValue() bool
	getConnectedText() string
	setNext(next *syntaxNode)
	getNext() *syntaxNode
	retrieveNext(root, current interface{}, result *[]interface{}) error
}
