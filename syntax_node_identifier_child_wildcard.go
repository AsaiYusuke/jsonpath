package jsonpath

import (
	"reflect"
)

type syntaxChildWildcardIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxChildWildcardIdentifier) retrieve(
	root, current interface{}, container *bufferContainer) errorRuntime {

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
		return ErrorTypeUnmatched{
			errorBasicRuntime: i.errorRuntime,
			expectedType:      msgTypeObjectOrArray,
			foundType:         foundType,
		}
	}
}

func (i *syntaxChildWildcardIdentifier) retrieveMap(
	root interface{}, srcMap map[string]interface{}, container *bufferContainer) errorRuntime {

	var deepestTextLen int
	var deepestError errorRuntime

	sortKeys := container.getSortedKeys(srcMap)

	for _, key := range *sortKeys {
		if err := i.retrieveMapNext(root, srcMap, key, container); err != nil {
			if len(container.result) == 0 {
				deepestTextLen, deepestError = i.addDeepestError(err, deepestTextLen, deepestError)
			}
		}
	}

	container.putSortSlice(sortKeys)

	if len(container.result) > 0 {
		return nil
	}

	if deepestError == nil {
		return ErrorMemberNotExist{
			errorBasicRuntime: i.errorRuntime,
		}
	}

	return deepestError
}

func (i *syntaxChildWildcardIdentifier) retrieveList(
	root interface{}, srcList []interface{}, container *bufferContainer) errorRuntime {

	var deepestTextLen int
	var deepestError errorRuntime

	for index := range srcList {
		if err := i.retrieveListNext(root, srcList, index, container); err != nil {
			if len(container.result) == 0 {
				deepestTextLen, deepestError = i.addDeepestError(err, deepestTextLen, deepestError)
			}
		}
	}

	if len(container.result) > 0 {
		return nil
	}

	if deepestError == nil {
		return ErrorMemberNotExist{
			errorBasicRuntime: i.errorRuntime,
		}
	}

	return deepestError
}
