package syntax

import (
	"reflect"

	"github.com/AsaiYusuke/jsonpath/v2/errors"
)

type syntaxChildSingleIdentifier struct {
	*syntaxBasicNode

	identifier string
}

func (i *syntaxChildSingleIdentifier) retrieve(
	root, current any, container *bufferContainer) errors.ErrorRuntime {

	if srcMap, ok := current.(map[string]any); ok {
		return i.retrieveMapNext(root, srcMap, i.identifier, container)
	}

	if current != nil {
		return errors.NewErrorTypeUnmatched(
			i.path, i.remainingPathLen, msgTypeObject, reflect.TypeOf(current).String())
	}
	return errors.NewErrorTypeUnmatched(
		i.path, i.remainingPathLen, msgTypeObject, msgTypeNull)
}
