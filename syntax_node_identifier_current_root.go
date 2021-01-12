package jsonpath

type syntaxCurrentRootIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxCurrentRootIdentifier) retrieve(
	root, current interface{}, container *bufferContainer) error {

	return i.retrieveAnyValueNext(root, current, container)
}
