package jsonpath

type syntaxCurrentRootIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxCurrentRootIdentifier) retrieve(_, _ interface{}, _ *[]interface{}) error {
	return nil
}
