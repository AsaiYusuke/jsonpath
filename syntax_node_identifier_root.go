package jsonpath

type syntaxRootIdentifier struct {
	*syntaxBasicNode

	srcJSON **interface{}
}

func (i *syntaxRootIdentifier) retrieve(current interface{}) error {
	return i.retrieveNext(**i.srcJSON)
}
