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
			errorBasicRuntime: &errorBasicRuntime{
				node: i.syntaxBasicNode,
			},
			expectedType: msgTypeObject,
			foundType:    foundType,
		}
	}

	deepestErrors := make([]errorRuntime, 0, 2)

	deepestErrors = i.retrieveMap(root, srcMap, container, deepestErrors)

	switch len(deepestErrors) {
	case 0:
		return nil
	case 1:
		return deepestErrors[0]
	default:
		return ErrorNoneMatched{
			errorBasicRuntime: &errorBasicRuntime{
				node: deepestErrors[0].getSyntaxNode(),
			},
		}
	}

}

func (i *syntaxChildMultiIdentifier) retrieveMap(
	root interface{}, srcMap map[string]interface{}, container *bufferContainer,
	deepestErrors []errorRuntime) []errorRuntime {

	var deepestTextLen int

	for _, identifier := range i.identifiers {
		switch typedNode := identifier.(type) {
		case *syntaxChildWildcardIdentifier:
		case *syntaxChildSingleIdentifier:
			if _, ok := srcMap[typedNode.identifier]; !ok {
				continue
			}
		default:
			continue
		}

		if err := identifier.retrieve(root, srcMap, container); err != nil {
			if len(container.result) == 0 {
				deepestTextLen, deepestErrors = i.addDeepestError(err, deepestTextLen, deepestErrors)
			}
		}
	}

	if len(container.result) > 0 {
		return nil
	}

	if len(deepestErrors) == 0 {
		return append(deepestErrors, ErrorMemberNotExist{
			errorBasicRuntime: &errorBasicRuntime{
				node: i.syntaxBasicNode,
			},
		})
	}

	return deepestErrors
}
