package jsonpath

type syntaxChildWildcardIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxChildWildcardIdentifier) retrieve(
	root, current interface{}, container *bufferContainer) error {

	childErrorMap := make(map[error]struct{}, 1)
	var lastError error

	switch typedNodes := current.(type) {
	case map[string]interface{}:
		lastError = i.retrieveMap(root, typedNodes, container, childErrorMap)

	case []interface{}:
		lastError = i.retrieveList(root, typedNodes, container, childErrorMap)

	}

	if len(container.result) == 0 {
		switch len(childErrorMap) {
		case 0:
			return ErrorMemberNotExist{path: i.text}
		case 1:
			return lastError
		default:
			return ErrorNoneMatched{path: i.next.getConnectedText()}
		}
	}

	return nil
}

func (i *syntaxChildWildcardIdentifier) retrieveMap(
	root interface{}, srcMap map[string]interface{}, container *bufferContainer,
	childErrorMap map[error]struct{}) error {

	var lastError error

	container.expandSortSlice(len(srcMap))

	index := 0
	for key := range srcMap {
		(*container.sortKeys)[index] = key
		index++
	}
	if len(*container.sortKeys) > 1 {
		container.sortKeys.Sort()
	}
	for _, key := range *container.sortKeys {
		if err := i.retrieveMapNext(root, srcMap, key, container); err != nil {
			childErrorMap[err] = struct{}{}
			lastError = err
		}
	}

	return lastError
}

func (i *syntaxChildWildcardIdentifier) retrieveList(
	root interface{}, srcList []interface{}, container *bufferContainer,
	childErrorMap map[error]struct{}) error {

	var lastError error

	for index := range srcList {
		if err := i.retrieveListNext(root, srcList, index, container); err != nil {
			childErrorMap[err] = struct{}{}
			lastError = err
		}
	}

	return lastError
}
