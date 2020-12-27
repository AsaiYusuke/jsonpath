package jsonpath

type syntaxBasicNode struct {
	text         string
	multiValue   bool
	next         syntaxNode
	accessorMode bool
	result       **[]interface{}
}

func (i *syntaxBasicNode) setText(text string) {
	i.text = text
}

func (i *syntaxBasicNode) setMultiValue() {
	i.multiValue = true
}

func (i *syntaxBasicNode) isMultiValue() bool {
	return i.multiValue
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

func (i *syntaxBasicNode) retrieveNext(getter func() interface{}, setter func(interface{})) error {
	if i.next != nil {
		return i.next.retrieve(getter())
	}
	if i.accessorMode {
		**i.result = append(**i.result, Accessor{Get: getter, Set: setter})
	} else {
		**i.result = append(**i.result, getter())
	}
	return nil
}

func (i *syntaxBasicNode) setResultPtr(resultPtr **[]interface{}) {
	i.result = resultPtr
}

func (i *syntaxBasicNode) setAccessorMode(mode bool) {
	i.accessorMode = mode
}
