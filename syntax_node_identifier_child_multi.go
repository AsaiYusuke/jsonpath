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
		foundType := `null`
		if current != nil {
			foundType = reflect.TypeOf(current).String()
		}
		return ErrorTypeUnmatched{
			errorBasicRuntime: &errorBasicRuntime{
				node: i.syntaxBasicNode,
			},
			expectedType: `object`,
			foundType:    foundType,
		}
	}

	deepestErrors := make([]errorRuntime, 0, 2)

	deepestErrors = i.retrieveMap(root, srcMap, container, deepestErrors)

	if len(container.result) == 0 {
		switch len(deepestErrors) {
		case 0:
			return ErrorMemberNotExist{
				errorBasicRuntime: &errorBasicRuntime{
					node: i.syntaxBasicNode,
				},
			}
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

	return nil
}

func (i *syntaxChildMultiIdentifier) retrieveMap(
	root interface{}, srcMap map[string]interface{}, container *bufferContainer,
	deepestErrors []errorRuntime) []errorRuntime {

	deepestTextLen := -1

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
			deepestTextLen, deepestErrors = i.addDeepestError(err, deepestTextLen, deepestErrors)
		}
	}

	return deepestErrors
}
