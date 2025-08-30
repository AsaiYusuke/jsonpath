package syntax

import (
	"github.com/AsaiYusuke/jsonpath/v2/errors"
)

type syntaxChildWildcardIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxChildWildcardIdentifier) retrieve(
	root, current any, results *[]any) errors.ErrorRuntime {

	switch typedNodes := current.(type) {
	case map[string]any:
		return i.retrieveMap(root, typedNodes, results)

	case []any:
		return i.retrieveList(root, typedNodes, results)

	default:
		return i.newErrTypeUnmatched(msgTypeObjectOrArray, current)
	}
}

func (i *syntaxChildWildcardIdentifier) retrieveMap(
	root any, srcMap map[string]any, results *[]any) errors.ErrorRuntime {

	var deepestError errors.ErrorRuntime

	sortKeys, keyLength := getSortedKeys(srcMap)

	for index := range keyLength {
		if err := i.retrieveMapNext(root, srcMap, (*sortKeys)[index], results); len(*results) == 0 && err != nil {
			deepestError = i.getMostResolvedError(err, deepestError)
		}
	}

	putSortSlice(sortKeys)

	if len(*results) > 0 {
		return nil
	}

	if deepestError == nil {
		return i.newErrMemberNotExist()
	}

	return deepestError
}

func (i *syntaxChildWildcardIdentifier) retrieveList(
	root any, srcList []any, results *[]any) errors.ErrorRuntime {

	var deepestError errors.ErrorRuntime

	for index := range srcList {
		if err := i.retrieveListNext(root, srcList, index, results); len(*results) == 0 && err != nil {
			deepestError = i.getMostResolvedError(err, deepestError)
		}
	}

	if len(*results) > 0 {
		return nil
	}

	if deepestError == nil {
		return i.newErrMemberNotExist()
	}

	return deepestError
}
