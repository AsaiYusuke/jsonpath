package syntax

import "github.com/AsaiYusuke/jsonpath/v2/errors"

type syntaxRootNodeIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxRootNodeIdentifier) retrieve(
	root, _ any, results *[]any) errors.ErrorRuntime {
	return i.retrieveAnyValueNext(root, root, results)
}
