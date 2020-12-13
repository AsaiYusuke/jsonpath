package jsonpath

type syntaxRootIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxRootIdentifier) retrieve(root, current interface{}) error {
	return i.retrieveNext(root, root)
}
