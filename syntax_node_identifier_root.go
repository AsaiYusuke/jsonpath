package jsonpath

type syntaxRootIdentifier struct {
	*syntaxBasicNode

	srcJSON **interface{}
}

func (i *syntaxRootIdentifier) retrieve(_ interface{}) error {
	return i.retrieveNext(**i.srcJSON)
}
