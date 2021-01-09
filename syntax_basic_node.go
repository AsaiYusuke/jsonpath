package jsonpath

type syntaxBasicNode struct {
	text         string
	valueGroup   bool
	next         syntaxNode
	accessorMode bool
}

func (i *syntaxBasicNode) setText(text string) {
	i.text = text
}

func (i *syntaxBasicNode) setValueGroup() {
	i.valueGroup = true
}

func (i *syntaxBasicNode) isValueGroup() bool {
	return i.valueGroup
}

func (i *syntaxBasicNode) getConnectedText() string {
	if i.next != nil {
		return i.text + i.next.getConnectedText()
	}
	return i.text
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
	root interface{}, current interface{}, result *[]interface{}) error {

	if i.next != nil {
		return i.next.retrieve(root, current, result)
	}

	if i.accessorMode {
		*result = append(*result, Accessor{
			Get: func() interface{} { return current },
			Set: nil,
		})
	} else {
		*result = append(*result, current)
	}

	return nil
}

func (i *syntaxBasicNode) retrieveMapNext(
	root interface{}, currentMap map[string]interface{}, key string, result *[]interface{}) error {

	if i.next != nil {
		return i.next.retrieve(root, currentMap[key], result)
	}

	if i.accessorMode {
		*result = append(*result, Accessor{
			Get: func() interface{} { return currentMap[key] },
			Set: func(value interface{}) { currentMap[key] = value },
		})
	} else {
		*result = append(*result, currentMap[key])
	}

	return nil
}

func (i *syntaxBasicNode) retrieveListNext(
	root interface{}, currentList []interface{}, index int, result *[]interface{}) error {

	if i.next != nil {
		return i.next.retrieve(root, currentList[index], result)
	}

	if i.accessorMode {
		*result = append(*result, Accessor{
			Get: func() interface{} { return currentList[index] },
			Set: func(value interface{}) { currentList[index] = value },
		})
	} else {
		*result = append(*result, currentList[index])
	}

	return nil
}

func (i *syntaxBasicNode) setAccessorMode(mode bool) {
	i.accessorMode = mode
}
