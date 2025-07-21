package syntax

type syntaxCurrentNodeIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxCurrentNodeIdentifier) retrieve(
	root, current interface{}, container *bufferContainer) errorRuntime {

	return i.retrieveAnyValueNext(root, current, container)
}
