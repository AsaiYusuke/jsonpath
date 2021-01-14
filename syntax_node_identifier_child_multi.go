package jsonpath

import "reflect"

type syntaxChildMultiIdentifier struct {
	*syntaxBasicNode

	identifiers    []syntaxNode
	isAllWildcard  bool
	unionQualifier syntaxUnionQualifier
}

func (i *syntaxChildMultiIdentifier) retrieve(
	root, current interface{}, container *bufferContainer) error {

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
			expectedType: `object`,
			foundType:    foundType,
			path:         i.text,
		}
	}

	childErrorMap := make(map[error]struct{}, 1)

	lastError := i.retrieveMap(root, srcMap, container, childErrorMap)

	if len(container.result) == 0 {
		switch len(childErrorMap) {
		case 0:
			return ErrorMemberNotExist{path: i.text}
		case 1:
			return lastError
		default:
			return ErrorNoneMatched{path: i.next.getConnectedText()}
		}
	}

	return nil
}

func (i *syntaxChildMultiIdentifier) retrieveMap(
	root interface{}, srcMap map[string]interface{}, container *bufferContainer,
	childErrorMap map[error]struct{}) error {

	var lastError error

	for _, identifier := range i.identifiers {
		var found bool
		switch typedNode := identifier.(type) {
		case *syntaxChildWildcardIdentifier:
			found = true
		case *syntaxChildSingleIdentifier:
			if _, ok := srcMap[typedNode.identifier]; ok {
				found = true
			}
		}

		if found {
			if err := identifier.retrieve(root, srcMap, container); err != nil {
				childErrorMap[err] = struct{}{}
				lastError = err
				continue
			}
		}
	}

	return lastError
}
