package syntax

import (
	"reflect"

	"github.com/AsaiYusuke/jsonpath/errors"
)

type syntaxChildWildcardIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxChildWildcardIdentifier) retrieve(
	root, current interface{}, container *bufferContainer) errors.ErrorRuntime {

	switch typedNodes := current.(type) {
	case map[string]interface{}:
		return i.retrieveMap(root, typedNodes, container)

	case []interface{}:
		return i.retrieveList(root, typedNodes, container)

	default:
		foundType := msgTypeNull
		if current != nil {
			foundType = reflect.TypeOf(current).String()
		}
		return errors.NewErrorTypeUnmatched(i.path, i.remainingPathLen, msgTypeObjectOrArray, foundType)
	}
}

func (i *syntaxChildWildcardIdentifier) retrieveMap(
	root interface{}, srcMap map[string]interface{}, container *bufferContainer) errors.ErrorRuntime {

	var deepestError errors.ErrorRuntime

	sortKeys := getSortedKeys(srcMap)

	for _, key := range *sortKeys {
		if err := i.retrieveMapNext(root, srcMap, key, container); err != nil {
			if len(container.result) == 0 {
				deepestError = i.getMostResolvedError(err, deepestError)
			}
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
	root interface{}, srcList []interface{}, container *bufferContainer) errors.ErrorRuntime {

	var deepestError errors.ErrorRuntime

	for index := range srcList {
		if err := i.retrieveListNext(root, srcList, index, container); err != nil {
			if len(container.result) == 0 {
				deepestError = i.getMostResolvedError(err, deepestError)
			}
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
