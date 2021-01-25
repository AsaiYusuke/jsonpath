package jsonpath

type syntaxRootIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxRootIdentifier) retrieve(
	root, _ interface{}, container *bufferContainer) errorRuntime {

	return i.retrieveAnyValueNext(root, root, container)
}
