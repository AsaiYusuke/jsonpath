package syntax

import (
	"reflect"

	"github.com/AsaiYusuke/jsonpath/errors"
)

type syntaxChildSingleIdentifier struct {
	*syntaxBasicNode

	identifier string
}

func (i *syntaxChildSingleIdentifier) retrieve(
	root, current interface{}, container *bufferContainer) errors.ErrorRuntime {

	if srcMap, ok := current.(map[string]interface{}); ok {
		return i.retrieveMapNext(root, srcMap, i.identifier, container)
	}

	foundType := msgTypeNull
	if current != nil {
		foundType = reflect.TypeOf(current).String()
	}
	return errors.NewErrorTypeUnmatched(i.path, i.remainingPathLen, msgTypeObject, foundType)
}
