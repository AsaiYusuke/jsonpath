package jsonpath

type syntaxRootIdentifier struct {
	*syntaxBasicNode
}

func (i syntaxRootIdentifier) retrieve(root, current interface{}, result *[]interface{}) error {
	return i.retrieveNext(root, root, result)
}
