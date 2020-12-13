package jsonpath

type syntaxRootIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxRootIdentifier) retrieve(current interface{}) error {
	return i.retrieveNext(**i.srcJSON)
}
