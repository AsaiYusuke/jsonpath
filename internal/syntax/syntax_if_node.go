package syntax

import "github.com/AsaiYusuke/jsonpath/errors"

type syntaxNode interface {
	retrieve(root, current interface{}, container *bufferContainer) errors.ErrorRuntime
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
