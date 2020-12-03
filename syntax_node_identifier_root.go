package jsonpath

type syntaxRootIdentifier struct {
	*syntaxBasicNode
}

func (i syntaxRootIdentifier) retrieve(root, current interface{}, result *resultContainer) error {
	return i.retrieveNext(root, root, result)
}
