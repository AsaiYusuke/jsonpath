package syntax

import "github.com/AsaiYusuke/jsonpath/v2/errors"

type syntaxNode interface {
	retrieve(root, current any, container *bufferContainer) errors.ErrorRuntime
	setPath(path string)
	getPath() string
	setValueGroup()
	isValueGroup() bool
	setRemainingPath(path string)
	getRemainingPath() string
	setNext(next syntaxNode)
	getNext() syntaxNode
	setAccessorMode(mode bool)
}
