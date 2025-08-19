package syntax

import (
	"github.com/AsaiYusuke/jsonpath/v2/config"
	"github.com/AsaiYusuke/jsonpath/v2/errors"
)

type syntaxBasicNode struct {
	path                 string
	remainingPath        string
	remainingPathLen     int
	valueGroup           bool
	next                 syntaxNode
	accessorMode         bool
	preErrMemberNotExist errors.ErrorMemberNotExist
}

func (i *syntaxBasicNode) setPath(path string) {
	i.path = path
}

func (i *syntaxBasicNode) getPath() string {
	return i.path
}

func (i *syntaxBasicNode) setValueGroup() {
	i.valueGroup = true
}

func (i *syntaxBasicNode) isValueGroup() bool {
	return i.valueGroup
}

func (i *syntaxBasicNode) setRemainingPath(path string) {
	i.remainingPath = path
	i.remainingPathLen = len(path)
}

func (i *syntaxBasicNode) getRemainingPath() string {
	return i.remainingPath
}

func (i *syntaxBasicNode) setNext(next syntaxNode) {
	if i.next != nil {
		i.next.setNext(next)
	} else {
		i.next = next
	}
}

func (i *syntaxBasicNode) getNext() syntaxNode {
	return i.next
}

func (i *syntaxBasicNode) newErrMemberNotExist() errors.ErrorMemberNotExist {
	if i.preErrMemberNotExist.ErrorBasicRuntime == nil {
		i.preErrMemberNotExist = errors.NewErrorMemberNotExist(i.path, i.remainingPathLen)
	}
	return i.preErrMemberNotExist
}

func (i *syntaxBasicNode) retrieveAnyValueNext(
	root any, nextSrc any, container *bufferContainer) errors.ErrorRuntime {

	if i.next != nil {
		return i.next.retrieve(root, nextSrc, container)
	}

	if i.accessorMode {
		container.result = append(container.result, config.Accessor{
			Get: func() any { return nextSrc },
			Set: nil,
		})
	} else {
		container.result = append(container.result, nextSrc)
	}

	return nil
}

func (i *syntaxBasicNode) retrieveMapNext(
	root any, currentMap map[string]any, key string, container *bufferContainer) errors.ErrorRuntime {

	nextNode, ok := currentMap[key]
	if !ok {
		return i.newErrMemberNotExist()
	}

	if i.next != nil {
		return i.next.retrieve(root, nextNode, container)
	}

	if i.accessorMode {
		container.result = append(container.result, config.Accessor{
			Get: func() any { return currentMap[key] },
			Set: func(value any) { currentMap[key] = value },
		})
	} else {
		container.result = append(container.result, nextNode)
	}

	return nil
}

func (i *syntaxBasicNode) retrieveListNext(
	root any, currentList []any, index int, container *bufferContainer) errors.ErrorRuntime {

	if i.next != nil {
		return i.next.retrieve(root, currentList[index], container)
	}

	if i.accessorMode {
		container.result = append(container.result, config.Accessor{
			Get: func() any { return currentList[index] },
			Set: func(value any) { currentList[index] = value },
		})
	} else {
		container.result = append(container.result, currentList[index])
	}

	return nil
}

func (i *syntaxBasicNode) setAccessorMode(mode bool) {
	i.accessorMode = mode
}

func (i *syntaxBasicNode) getMostResolvedError(
	newError errors.ErrorRuntime, currentMostResolvedError errors.ErrorRuntime) errors.ErrorRuntime {

	if currentMostResolvedError == nil {
		return newError
	}

	newPathLength := newError.GetRemainingPathLen()
	currentMostResolvedPathLen := currentMostResolvedError.GetRemainingPathLen()

	if currentMostResolvedPathLen > newPathLength {
		return newError
	}

	if currentMostResolvedPathLen == newPathLength {
		if _, ok := currentMostResolvedError.(errors.ErrorTypeUnmatched); ok {
			return newError
		}
	}

	return currentMostResolvedError
}
