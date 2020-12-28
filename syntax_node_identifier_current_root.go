package jsonpath

type syntaxCurrentRootIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxCurrentRootIdentifier) retrieve(
	root, current interface{}, result *[]interface{}) error {

	return i.retrieveNext(
		root, result,
		func() interface{} {
			return current
		},
		nil)
}
