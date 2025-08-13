package syntax

import (
	"reflect"

	"github.com/AsaiYusuke/jsonpath/errors"
)

type syntaxRecursiveChildIdentifier struct {
	*syntaxBasicNode

	nextMapRequired  bool
	nextListRequired bool
}

func (i *syntaxRecursiveChildIdentifier) retrieve(
	root, current any, container *bufferContainer) errors.ErrorRuntime {

	switch current.(type) {
	case map[string]any, []any:
	default:
		foundType := msgTypeNull
		if current != nil {
			foundType = reflect.TypeOf(current).String()
		}
		return errors.NewErrorTypeUnmatched(i.path, i.remainingPathLen, msgTypeObjectOrArray, foundType)
	}

	var deepestError errors.ErrorRuntime

	targetNodes := make([]interface{}, 1, 5)
	targetNodes[0] = current

	for len(targetNodes) > 0 {
		currentTargetNode := targetNodes[len(targetNodes)-1]
		targetNodes = targetNodes[:len(targetNodes)-1]
		switch typedNodes := currentTargetNode.(type) {
		case map[string]any:
			if i.nextMapRequired {
				if err := i.next.retrieve(root, typedNodes, container); err != nil {
					if len(container.result) == 0 {
						deepestError = i.getMostResolvedError(err, deepestError)
					}
				}
			}

			sortKeys := getSortedKeys(typedNodes)
			for index := len(typedNodes) - 1; index >= 0; index-- {
				node := typedNodes[(*sortKeys)[index]]
				switch node.(type) {
				case map[string]interface{}, []interface{}:
					targetNodes = append(targetNodes, node)
				}
			}

			putSortSlice(sortKeys)

		case []any:
			if i.nextListRequired {
				if err := i.next.retrieve(root, typedNodes, container); err != nil {
					if len(container.result) == 0 {
						deepestError = i.getMostResolvedError(err, deepestError)
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

	if deepestError == nil {
		return i.newErrMemberNotExist()
	}

	return deepestError
}
