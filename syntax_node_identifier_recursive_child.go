package jsonpath

import (
	"reflect"
)

type syntaxRecursiveChildIdentifier struct {
	*syntaxBasicNode

	nextMapRequired  bool
	nextListRequired bool
}

func (i *syntaxRecursiveChildIdentifier) retrieve(
	root, current interface{}, container *bufferContainer) errorRuntime {

	switch current.(type) {
	case map[string]interface{}, []interface{}:
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

	var deepestTextLen int
	deepestErrors := make([]errorRuntime, 0, 2)

	targetNodes := make([]interface{}, 1, 5)
	targetNodes[0] = current

	for len(targetNodes) > 0 {
		currentNode := targetNodes[len(targetNodes)-1]
		targetNodes = targetNodes[:len(targetNodes)-1]
		switch typedNodes := currentNode.(type) {
		case map[string]interface{}:
			if i.nextMapRequired {
				if err := i.next.retrieve(root, typedNodes, container); err != nil {
					if len(container.result) == 0 {
						deepestTextLen, deepestErrors = i.addDeepestError(err, deepestTextLen, deepestErrors)
					}
				}
			}

			sortKeys := container.getSortedKeys(typedNodes)
			for index := len(typedNodes) - 1; index >= 0; index-- {
				node := typedNodes[(*sortKeys)[index]]
				switch node.(type) {
				case map[string]interface{}, []interface{}:
					targetNodes = append(targetNodes, node)
				}
			}

			container.putSortSlice(sortKeys)

		case []interface{}:
			if i.nextListRequired {
				if err := i.next.retrieve(root, typedNodes, container); err != nil {
					if len(container.result) == 0 {
						deepestTextLen, deepestErrors = i.addDeepestError(err, deepestTextLen, deepestErrors)
					}
				}
			}

			for index := len(typedNodes) - 1; index >= 0; index-- {
				node := typedNodes[index]
				switch node.(type) {
				case map[string]interface{}, []interface{}:
					targetNodes = append(targetNodes, node)
				}
			}
		}
	}

	if len(container.result) > 0 {
		return nil
	}

	switch len(deepestErrors) {
	case 0:
		return ErrorMemberNotExist{
			errorBasicRuntime: i.errorRuntime,
		}
	case 1:
		return deepestErrors[0]
	default:
		return ErrorNoneMatched{
			errorBasicRuntime: deepestErrors[0].getSyntaxNode().errorRuntime,
		}
	}
}
