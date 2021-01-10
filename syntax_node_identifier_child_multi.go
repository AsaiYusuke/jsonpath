package jsonpath

import "reflect"

type syntaxChildMultiIdentifier struct {
	*syntaxBasicNode

	identifiers []syntaxNode
}

func (i *syntaxChildMultiIdentifier) retrieve(
	root, current interface{}, result *[]interface{}) error {

	if _, ok := current.([]interface{}); ok {
		wildcardSubscripts := make([]syntaxSubscript, 0, len(i.identifiers))
		isAllWildcard := true
		for index := range i.identifiers {
			if _, ok := i.identifiers[index].(*syntaxChildWildcardIdentifier); !ok {
				isAllWildcard = false
				break
			}
			wildcardSubscripts = append(wildcardSubscripts, &syntaxWildcardSubscript{})
		}
		if len(i.identifiers) > 0 && isAllWildcard {
			// If the "current" variable points to the array structure
			// and only wildcards are specified for qualifier,
			// then switch to syntaxUnionQualifier.
			unionQualifier := syntaxUnionQualifier{
				syntaxBasicNode: &syntaxBasicNode{
					text:         i.text,
					next:         i.next,
					valueGroup:   true,
					accessorMode: i.accessorMode,
				},
				subscripts: wildcardSubscripts,
			}
			return unionQualifier.retrieve(root, current, result)
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

	lastError := i.retrieveMap(root, srcMap, result, childErrorMap)

	if len(*result) == 0 {
		switch len(childErrorMap) {
		case 1:
			return lastError
		default:
			return ErrorNoneMatched{path: i.next.getConnectedText()}
		}
	}

	return nil
}

func (i *syntaxChildMultiIdentifier) retrieveMap(
	root interface{}, srcMap map[string]interface{}, result *[]interface{},
	childErrorMap map[error]struct{}) error {

	var lastError error
	var partialFound bool

	for index := range i.identifiers {
		var found bool
		switch typedNode := i.identifiers[index].(type) {
		case *syntaxChildWildcardIdentifier:
			found = true
		case *syntaxChildSingleIdentifier:
			if _, ok := srcMap[typedNode.identifier]; ok {
				found = true
			}
		}

		if found {
			partialFound = true
			if err := i.identifiers[index].retrieve(root, srcMap, result); err != nil {
				childErrorMap[err] = struct{}{}
				lastError = err
				continue
			}
		}
	}

	if !partialFound {
		err := ErrorMemberNotExist{path: i.text}
		childErrorMap[err] = struct{}{}
		lastError = err
	}

	return lastError
}
