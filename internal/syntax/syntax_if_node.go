package syntax

type syntaxNode interface {
	retrieve(root, current interface{}, container *bufferContainer) errorRuntime
	setText(text string)
	getText() string
	setValueGroup()
	isValueGroup() bool
	setConnectedText(text string)
	getConnectedText() string
	setNext(next syntaxNode)
	getNext() syntaxNode
	setAccessorMode(mode bool)
}
