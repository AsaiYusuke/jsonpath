package jsonpath

type syntaxRootNodeIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxRootNodeIdentifier) retrieve(
	root, _ interface{}, container *bufferContainer) errorRuntime {

	return i.retrieveAnyValueNext(root, root, container)
}
