package jsonpath

import "reflect"

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
		return ErrorMemberNotExist{
			errorBasicRuntime: i.errorRuntime,
		}
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
	err errorRuntime, deepestTextLen int, deepestErrors []errorRuntime) (int, []errorRuntime) {

	textLength := len(err.getSyntaxNode().getConnectedText())

	if deepestTextLen == 0 || deepestTextLen > textLength {
		deepestTextLen = textLength
		deepestErrors = deepestErrors[:0]
	}

	if deepestTextLen == textLength {
		switch len(deepestErrors) {
		case 0:
			return deepestTextLen, append(deepestErrors, err)
		case 1:
			if reflect.TypeOf(err) != reflect.TypeOf(deepestErrors[0]) {
				return deepestTextLen, append(deepestErrors, err)
			}
		}
	}

	return deepestTextLen, deepestErrors
}
