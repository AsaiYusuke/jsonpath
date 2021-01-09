package jsonpath

type syntaxCurrentRootIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxCurrentRootIdentifier) retrieve(
	root, current interface{}, result *[]interface{}) error {

	return i.retrieveAnyValueNext(root, current, result)
}
