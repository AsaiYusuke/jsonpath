package jsonpath

type syntaxCurrentRootIdentifier struct {
	*syntaxBasicNode
}

func (i syntaxCurrentRootIdentifier) retrieve(root, current interface{}, result *resultContainer) error {
	return i.retrieveNext(root, current, result)
}
