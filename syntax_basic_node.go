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

func (i *syntaxBasicNode) retrieveNext(
	root interface{}, result *[]interface{}, getter func() interface{}, setter func(interface{})) error {

	if i.next != nil {
		return i.next.retrieve(root, getter(), result)
	}

	if i.accessorMode {
		*result = append(*result, Accessor{Get: getter, Set: setter})
	} else {
		*result = append(*result, getter())
	}

	return nil
}

func (i *syntaxBasicNode) setAccessorMode(mode bool) {
	i.accessorMode = mode
}
