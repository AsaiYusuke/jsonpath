package jsonpath

type syntaxRootIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxRootIdentifier) retrieve(_, _ interface{}, _ *[]interface{}) error {
	return nil
}
