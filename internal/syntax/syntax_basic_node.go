package syntax

import (
	"reflect"
	"sync"

	"github.com/AsaiYusuke/jsonpath/v2/config"
	"github.com/AsaiYusuke/jsonpath/v2/errors"
)

type syntaxNodeErrState struct {
	basicRuntime errors.ErrorBasicRuntime
}

type syntaxBasicNode struct {
	path             string
	remainingPath    string
	remainingPathLen int
	valueGroup       bool
	next             syntaxNode
	accessorMode     bool
	errState         *syntaxNodeErrState
	onceErrState     sync.Once
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

func (i *syntaxBasicNode) ensureErrState() {
	i.onceErrState.Do(func() {
		i.errState = &syntaxNodeErrState{}
		i.errState.basicRuntime = errors.NewErrorBasicRuntime(i.path, i.remainingPathLen)
	})
}

func (i *syntaxBasicNode) newErrMemberNotExist() errors.ErrorMemberNotExist {
	i.ensureErrState()
	return errors.NewErrorMemberNotExist(&i.errState.basicRuntime)
}

func (i *syntaxBasicNode) newErrTypeUnmatched(expected string, current any) errors.ErrorTypeUnmatched {
	i.ensureErrState()
	if current != nil {
		return errors.NewErrorTypeUnmatched(&i.errState.basicRuntime, expected, reflect.TypeOf(current).String())
	}
	return errors.NewErrorTypeUnmatched(&i.errState.basicRuntime, expected, msgTypeNull)
}

func (i *syntaxBasicNode) retrieveAnyValueNext(
	root any, nextSrc any, results *[]any) errors.ErrorRuntime {

	if i.next != nil {
		return i.next.retrieve(root, nextSrc, results)
	}

	if i.accessorMode {
		*results = append(*results, config.Accessor{
			Get: func() any { return nextSrc },
			Set: nil,
		})
		return nil
	}

	*results = append(*results, nextSrc)
	return nil
}

func (i *syntaxBasicNode) retrieveMapNext(
	root any, currentMap map[string]any, key string, results *[]any) errors.ErrorRuntime {

	nextNode, ok := currentMap[key]
	if !ok {
		return i.newErrMemberNotExist()
	}

	if i.next != nil {
		return i.next.retrieve(root, nextNode, results)
	}

	if i.accessorMode {
		*results = append(*results, config.Accessor{
			Get: func() any { return currentMap[key] },
			Set: func(value any) { currentMap[key] = value },
		})
		return nil
	}

	*results = append(*results, nextNode)
	return nil
}

func (i *syntaxBasicNode) retrieveListNext(
	root any, currentList []any, index int, results *[]any) errors.ErrorRuntime {

	if i.next != nil {
		return i.next.retrieve(root, currentList[index], results)
	}

	if i.accessorMode {
		*results = append(*results, config.Accessor{
			Get: func() any { return currentList[index] },
			Set: func(value any) { currentList[index] = value },
		})
		return nil
	}

	*results = append(*results, currentList[index])
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
