package jsonpath

import (
	"reflect"
)

type syntaxChildWildcardIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxChildWildcardIdentifier) retrieve(
	root, current interface{}, container *bufferContainer) errorRuntime {

	deepestErrors := make([]errorRuntime, 0, 2)

	switch typedNodes := current.(type) {
	case map[string]interface{}:
		deepestErrors = i.retrieveMap(root, typedNodes, container, deepestErrors)

	case []interface{}:
		deepestErrors = i.retrieveList(root, typedNodes, container, deepestErrors)

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

	switch len(deepestErrors) {
	case 0:
		return nil
	case 1:
		return deepestErrors[0]
	default:
		return ErrorNoneMatched{
			errorBasicRuntime: deepestErrors[0].getSyntaxNode().errorRuntime,
		}
	}

}

func (i *syntaxChildWildcardIdentifier) retrieveMap(
	root interface{}, srcMap map[string]interface{}, container *bufferContainer,
	deepestErrors []errorRuntime) []errorRuntime {

	var deepestTextLen int

	sortKeys := container.getSortedKeys(srcMap)

	for _, key := range *sortKeys {
		if err := i.retrieveMapNext(root, srcMap, key, container); err != nil {
			if len(container.result) == 0 {
				deepestTextLen, deepestErrors = i.addDeepestError(err, deepestTextLen, deepestErrors)
			}
		}
	}

	container.putSortSlice(sortKeys)

	if len(container.result) > 0 {
		return nil
	}

	if len(deepestErrors) == 0 {
		return append(deepestErrors, ErrorMemberNotExist{
			errorBasicRuntime: i.errorRuntime,
		})
	}

	return deepestErrors
}

func (i *syntaxChildWildcardIdentifier) retrieveList(
	root interface{}, srcList []interface{}, container *bufferContainer,
	deepestErrors []errorRuntime) []errorRuntime {

	var deepestTextLen int

	for index := range srcList {
		if err := i.retrieveListNext(root, srcList, index, container); err != nil {
			if len(container.result) == 0 {
				deepestTextLen, deepestErrors = i.addDeepestError(err, deepestTextLen, deepestErrors)
			}
		}
	}

	if len(container.result) > 0 {
		return nil
	}

	if len(deepestErrors) == 0 {
		return append(deepestErrors, ErrorIndexOutOfRange{
			errorBasicRuntime: i.errorRuntime,
		})
	}

	return deepestErrors
}
