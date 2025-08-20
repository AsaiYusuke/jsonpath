package syntax

import (
	"reflect"
	"slices"

	"github.com/AsaiYusuke/jsonpath/v2/errors"
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
		if current != nil {
			return errors.NewErrorTypeUnmatched(
				i.path, i.remainingPathLen, msgTypeObjectOrArray, reflect.TypeOf(current).String())
		}
		return errors.NewErrorTypeUnmatched(
			i.path, i.remainingPathLen, msgTypeObjectOrArray, msgTypeNull)
	}

	var deepestError errors.ErrorRuntime

	pooledNodes := getNodeSlice()
	targetNodes := *pooledNodes
	targetNodes = append(targetNodes, current)

	for len(targetNodes) > 0 {
		currentTargetNode := targetNodes[len(targetNodes)-1]
		targetNodes = targetNodes[:len(targetNodes)-1]
		switch typedNodes := currentTargetNode.(type) {
		case map[string]any:
			if i.nextMapRequired {
				if err := i.next.retrieve(root, typedNodes, container); len(container.result) == 0 && err != nil {
					deepestError = i.getMostResolvedError(err, deepestError)
				}
			}

			sortKeys, keyLength := getSortedRecursiveKeys(typedNodes)
			if keyLength > 0 {
				oldLength := len(targetNodes)
				targetNodes = slices.Grow(targetNodes, keyLength)
				targetNodes = targetNodes[:oldLength+keyLength]

				appendIndex := oldLength
				for index := keyLength - 1; index >= 0; index-- {
					targetNodes[appendIndex] = typedNodes[(*sortKeys)[index]]
					appendIndex++
				}
			}

			putSortSlice(sortKeys)

		case []any:
			if i.nextListRequired {
				if err := i.next.retrieve(root, typedNodes, container); len(container.result) == 0 && err != nil {
					deepestError = i.getMostResolvedError(err, deepestError)
				}
			}

			keyLength := 0
			for index := range typedNodes {
				switch typedNodes[index].(type) {
				case map[string]any, []any:
					keyLength++
				}
			}
			if keyLength > 0 {
				oldLength := len(targetNodes)
				targetNodes = slices.Grow(targetNodes, keyLength)
				targetNodes = targetNodes[:oldLength+keyLength]

				appendIndex := oldLength
				for index := len(typedNodes) - 1; index >= 0; index-- {
					switch typedNodes[index].(type) {
					case map[string]any, []any:
						targetNodes[appendIndex] = typedNodes[index]
						appendIndex++
					}
				}
			}
		}
	}

	*pooledNodes = targetNodes
	putNodeSlice(pooledNodes)

	if len(container.result) > 0 {
		return nil
	}

	if deepestError == nil {
		return i.newErrMemberNotExist()
	}

	return deepestError
}
