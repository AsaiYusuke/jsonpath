package jsonpath

type syntaxCurrentRootIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxCurrentRootIdentifier) retrieve(current interface{}) error {
	return i.retrieveNext(
		func() interface{} {
			return current
		},
		nil)
}
