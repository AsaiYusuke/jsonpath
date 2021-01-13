package jsonpath

type syntaxRecursiveChildIdentifier struct {
	*syntaxBasicNode

	nextMapRequired  bool
	nextListRequired bool
}

func (i *syntaxRecursiveChildIdentifier) retrieve(
	root, current interface{}, container *bufferContainer) error {

	targetNodes := make([]interface{}, 1, 5)
	targetNodes[0] = current

	for len(targetNodes) > 0 {
		currentNode := targetNodes[len(targetNodes)-1]
		targetNodes = targetNodes[:len(targetNodes)-1]
		switch typedNodes := currentNode.(type) {
		case map[string]interface{}:
			if i.nextMapRequired {
				i.retrieveAnyValueNext(root, typedNodes, container)
			}

			sortKeys := container.getSortSlice(len(typedNodes))

			defer func() {
				container.putSortSlice(sortKeys)
			}()

			index := 0
			for key := range typedNodes {
				(*sortKeys)[index] = key
				index++
			}

			if len(*sortKeys) > 1 {
				sortKeys.Sort()
			}
			for index := len(typedNodes) - 1; index >= 0; index-- {
				targetNodes = append(targetNodes, typedNodes[(*sortKeys)[index]])
			}

		case []interface{}:
			if i.nextListRequired {
				i.retrieveAnyValueNext(root, typedNodes, container)
			}
			for index := len(typedNodes) - 1; index >= 0; index-- {
				targetNodes = append(targetNodes, typedNodes[index])
			}
		}
	}

	if len(container.result) == 0 {
		return ErrorNoneMatched{path: i.getConnectedText()}
	}

	return nil
}
