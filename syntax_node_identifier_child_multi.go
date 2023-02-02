package jsonpath

import "reflect"

type syntaxChildMultiIdentifier struct {
	*syntaxBasicNode

	identifiers    []syntaxNode
	isAllWildcard  bool
	unionQualifier syntaxUnionQualifier
}

func (i *syntaxChildMultiIdentifier) retrieve(
	root, current interface{}, container *bufferContainer) errorRuntime {

	if i.isAllWildcard {
		if _, ok := current.([]interface{}); ok {
			// If the "current" variable points to the array structure
			// and only wildcards are specified for qualifier,
			// then switch to syntaxUnionQualifier.
			return i.unionQualifier.retrieve(root, current, container)
		}
	}

	srcMap, ok := current.(map[string]interface{})
	if !ok {
		foundType := msgTypeNull
		if current != nil {
			foundType = reflect.TypeOf(current).String()
		}
		return ErrorTypeUnmatched{
			errorBasicRuntime: i.errorRuntime,
			expectedType:      msgTypeObject,
			foundType:         foundType,
		}
	}

	return i.retrieveMap(root, srcMap, container)
}

func (i *syntaxChildMultiIdentifier) retrieveMap(
	root interface{}, srcMap map[string]interface{}, container *bufferContainer) errorRuntime {

	var deepestTextLen int
	var deepestError errorRuntime

	for _, identifier := range i.identifiers {
		if singleIdentifier, ok := identifier.(*syntaxChildSingleIdentifier); ok {
			if _, ok = srcMap[singleIdentifier.identifier]; !ok {
				continue
			}
		}

		if err := identifier.retrieve(root, srcMap, container); err != nil {
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
