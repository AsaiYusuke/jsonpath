package syntax

import (
	"reflect"

	"github.com/AsaiYusuke/jsonpath/errors"
)

type syntaxChildMultiIdentifier struct {
	*syntaxBasicNode

	identifiers    []syntaxNode
	isAllWildcard  bool
	unionQualifier syntaxUnionQualifier
}

func (i *syntaxChildMultiIdentifier) retrieve(
	root, current any, container *bufferContainer) errors.ErrorRuntime {

	if i.isAllWildcard {
		if _, ok := current.([]any); ok {
			// If the "current" variable points to the array structure
			// and only wildcards are specified for qualifier,
			// then switch to syntaxUnionQualifier.
			return i.unionQualifier.retrieve(root, current, container)
		}
	}

	if srcMap, ok := current.(map[string]any); ok {
		return i.retrieveMap(root, srcMap, container)
	}

	if current != nil {
		return errors.NewErrorTypeUnmatched(
			i.path, i.remainingPathLen, msgTypeObject, reflect.TypeOf(current).String())
	}
	return errors.NewErrorTypeUnmatched(
		i.path, i.remainingPathLen, msgTypeObject, msgTypeNull)
}

func (i *syntaxChildMultiIdentifier) retrieveMap(
	root any, srcMap map[string]any, container *bufferContainer) errors.ErrorRuntime {

	var deepestError errors.ErrorRuntime

	for _, identifier := range i.identifiers {
		if singleIdentifier, ok := identifier.(*syntaxChildSingleIdentifier); ok {
			if _, ok = srcMap[singleIdentifier.identifier]; !ok {
				continue
			}
		}

		if err := identifier.retrieve(root, srcMap, container); len(container.result) == 0 && err != nil {
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
