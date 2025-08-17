package syntax

import "github.com/AsaiYusuke/jsonpath/errors"

type syntaxCurrentNodeIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxCurrentNodeIdentifier) retrieve(
	root, current any, container *bufferContainer) errors.ErrorRuntime {

	return i.retrieveAnyValueNext(root, current, container)
}
