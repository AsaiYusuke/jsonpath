package syntax

import "reflect"

type syntaxChildSingleIdentifier struct {
	*syntaxBasicNode

	identifier string
}

func (i *syntaxChildSingleIdentifier) retrieve(
	root, current interface{}, container *bufferContainer) errorRuntime {

	if srcMap, ok := current.(map[string]interface{}); ok {
		return i.retrieveMapNext(root, srcMap, i.identifier, container)
	}

	foundType := msgTypeNull
	if current != nil {
		foundType = reflect.TypeOf(current).String()
	}
	return newErrorTypeUnmatched(i.errorRuntime.node, msgTypeObject, foundType)
}
