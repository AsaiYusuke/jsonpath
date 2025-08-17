package syntax

import (
	"reflect"

	"github.com/AsaiYusuke/jsonpath/v2/errors"
)

type syntaxChildWildcardIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxChildWildcardIdentifier) retrieve(
	root, current any, container *bufferContainer) errors.ErrorRuntime {

	switch typedNodes := current.(type) {
	case map[string]any:
		return i.retrieveMap(root, typedNodes, container)

	case []any:
		return i.retrieveList(root, typedNodes, container)

	default:
		if current != nil {
			return errors.NewErrorTypeUnmatched(
				i.path, i.remainingPathLen, msgTypeObjectOrArray, reflect.TypeOf(current).String())
		}
		return errors.NewErrorTypeUnmatched(
			i.path, i.remainingPathLen, msgTypeObjectOrArray, msgTypeNull)
	}
}

func (i *syntaxChildWildcardIdentifier) retrieveMap(
	root any, srcMap map[string]any, container *bufferContainer) errors.ErrorRuntime {

	var deepestError errors.ErrorRuntime

	sortKeys, keyLength := getSortedKeys(srcMap)

	for index := range keyLength {
		if err := i.retrieveMapNext(root, srcMap, (*sortKeys)[index], container); len(container.result) == 0 && err != nil {
			deepestError = i.getMostResolvedError(err, deepestError)
		}
	}

	putSortSlice(sortKeys)

	if len(container.result) > 0 {
		return nil
	}

	if deepestError == nil {
		return i.newErrMemberNotExist()
	}

	return deepestError
}

func (i *syntaxChildWildcardIdentifier) retrieveList(
	root any, srcList []any, container *bufferContainer) errors.ErrorRuntime {

	var deepestError errors.ErrorRuntime

	for index := range srcList {
		if err := i.retrieveListNext(root, srcList, index, container); len(container.result) == 0 && err != nil {
			deepestError = i.getMostResolvedError(err, deepestError)
		}
	}

	if len(container.result) > 0 {
		return nil
	}

	if deepestError == nil {
		return i.newErrMemberNotExist()
	}

	return deepestError
}
