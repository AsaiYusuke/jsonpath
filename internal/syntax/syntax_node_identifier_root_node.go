package syntax

import "github.com/AsaiYusuke/jsonpath/errors"

type syntaxRootNodeIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxRootNodeIdentifier) retrieve(
	root, _ interface{}, container *bufferContainer) errors.ErrorRuntime {

	return i.retrieveAnyValueNext(root, root, container)
}
