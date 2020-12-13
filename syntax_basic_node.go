package jsonpath

type syntaxBasicNode struct {
	text       string
	multiValue bool
	next       syntaxNode
	result     **[]interface{}
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

func (i *syntaxBasicNode) retrieveNext(root, current interface{}) error {
	if i.next != nil {
		return i.next.retrieve(root, current)
	}
	(**i.result) = append((**i.result), current)
	return nil
}

func (i *syntaxBasicNode) setResultPtr(resultPtr **[]interface{}) {
	i.result = resultPtr
}
