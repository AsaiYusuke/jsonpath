package jsonpath

type syntaxCurrentRootIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxCurrentRootIdentifier) retrieve(root, current interface{}) error {
	return i.retrieveNext(root, current)
}
