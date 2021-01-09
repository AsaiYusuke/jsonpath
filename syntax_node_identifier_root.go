package jsonpath

type syntaxRootIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxRootIdentifier) retrieve(
	root, _ interface{}, result *[]interface{}) error {
	return i.retrieveAnyValueNext(root, root, result)
}
