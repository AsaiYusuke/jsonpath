package syntax

import "github.com/AsaiYusuke/jsonpath/errors"

type syntaxBasicNode struct {
	text          string
	connectedText string
	valueGroup    bool
	next          syntaxNode
	accessorMode  bool
	errorRuntime  *errorBasicRuntime
}

func (i *syntaxBasicNode) setText(text string) {
	i.text = text
}

func (i *syntaxBasicNode) getText() string {
	return i.text
}

func (i *syntaxBasicNode) setValueGroup() {
	i.valueGroup = true
}

func (i *syntaxBasicNode) isValueGroup() bool {
	return i.valueGroup
}

func (i *syntaxBasicNode) setConnectedText(text string) {
	i.connectedText = text
}

func (i *syntaxBasicNode) getConnectedText() string {
	return i.connectedText
}

func (i *syntaxBasicNode) setNext(next syntaxNode) {
	if i.next != nil {
		i.next.setNext(next)
	} else {
		i.next = next
	}
}

func (i *syntaxBasicNode) getNext() syntaxNode {
	return i.next
}

func (i *syntaxBasicNode) retrieveAnyValueNext(
	root interface{}, nextSrc interface{}, container *bufferContainer) errorRuntime {

	if i.next != nil {
		return i.next.retrieve(root, nextSrc, container)
	}

	if i.accessorMode {
		container.result = append(container.result, Accessor{
			Get: func() interface{} { return nextSrc },
			Set: nil,
		})
	} else {
		container.result = append(container.result, nextSrc)
	}

	return nil
}

func (i *syntaxBasicNode) retrieveMapNext(
	root interface{}, currentMap map[string]interface{}, key string, container *bufferContainer) errorRuntime {

	nextNode, ok := currentMap[key]
	if !ok {
		return newErrorMemberNotExist(i.errorRuntime.node)
	}

	if i.next != nil {
		return i.next.retrieve(root, nextNode, container)
	}

	if i.accessorMode {
		container.result = append(container.result, Accessor{
			Get: func() interface{} { return currentMap[key] },
			Set: func(value interface{}) { currentMap[key] = value },
		})
	} else {
		container.result = append(container.result, nextNode)
	}

	return nil
}

func (i *syntaxBasicNode) retrieveListNext(
	root interface{}, currentList []interface{}, index int, container *bufferContainer) errorRuntime {

	if i.next != nil {
		return i.next.retrieve(root, currentList[index], container)
	}

	if i.accessorMode {
		container.result = append(container.result, Accessor{
			Get: func() interface{} { return currentList[index] },
			Set: func(value interface{}) { currentList[index] = value },
		})
	} else {
		container.result = append(container.result, currentList[index])
	}

	return nil
}

func (i *syntaxBasicNode) setAccessorMode(mode bool) {
	i.accessorMode = mode
}

func (i *syntaxBasicNode) addDeepestError(
	err errorRuntime, deepestTextLen int, deepestError errorRuntime) (int, errorRuntime) {

	textLength := len(err.getSyntaxNode().getConnectedText())

	if deepestTextLen == 0 || deepestTextLen > textLength {
		return textLength, err
	}

	if deepestTextLen == textLength {
		if adapter, ok := deepestError.(errorRuntimeAdapter); ok {
			if _, isTypeUnmatched := adapter.err.(errors.ErrorTypeUnmatched); isTypeUnmatched {
				return deepestTextLen, err
			}
		}
	}

	return deepestTextLen, deepestError
}
