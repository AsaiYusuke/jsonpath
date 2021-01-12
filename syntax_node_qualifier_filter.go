package jsonpath

type syntaxFilterQualifier struct {
	*syntaxBasicNode

	query syntaxQuery
}

func (f *syntaxFilterQualifier) retrieve(
	root, current interface{}, container *bufferContainer) error {

	childErrorMap := make(map[error]struct{}, 1)
	var lastError error

	switch typedNodes := current.(type) {
	case map[string]interface{}:
		lastError = f.retrieveMap(root, typedNodes, container, childErrorMap)

	case []interface{}:
		lastError = f.retrieveList(root, typedNodes, container, childErrorMap)

	}

	if len(container.result) == 0 {
		switch len(childErrorMap) {
		case 0:
			return ErrorMemberNotExist{path: f.text}
		case 1:
			return lastError
		default:
			return ErrorNoneMatched{path: f.next.getConnectedText()}
		}
	}

	return nil
}

func (f *syntaxFilterQualifier) retrieveMap(
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
	valueList := make([]interface{}, len(*container.sortKeys))
	for index := range *container.sortKeys {
		valueList[index] = srcMap[(*container.sortKeys)[index]]
	}

	valueList = f.query.compute(root, valueList, container)

	for index := range *container.sortKeys {
		var nodeNotFound bool
		if len(valueList) == 1 {
			_, nodeNotFound = valueList[0].(struct{})
		} else {
			_, nodeNotFound = valueList[index].(struct{})
		}
		if !nodeNotFound {
			if err := f.retrieveMapNext(root, srcMap, (*container.sortKeys)[index], container); err != nil {
				childErrorMap[err] = struct{}{}
				lastError = err
			}
		}
	}

	return lastError
}

func (f *syntaxFilterQualifier) retrieveList(
	root interface{}, srcList []interface{}, container *bufferContainer,
	childErrorMap map[error]struct{}) error {

	var lastError error

	valueList := f.query.compute(root, srcList, container)

	for index := range srcList {
		var nodeNotFound bool
		if len(valueList) == 1 {
			_, nodeNotFound = valueList[0].(struct{})
		} else {
			_, nodeNotFound = valueList[index].(struct{})
		}
		if !nodeNotFound {
			if err := f.retrieveListNext(root, srcList, index, container); err != nil {
				childErrorMap[err] = struct{}{}
				lastError = err
			}
		}
	}

	return lastError
}
