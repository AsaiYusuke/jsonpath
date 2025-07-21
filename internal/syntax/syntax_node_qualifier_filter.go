package syntax

import "reflect"

type syntaxFilterQualifier struct {
	*syntaxBasicNode

	query syntaxQuery
}

func (f *syntaxFilterQualifier) retrieve(
	root, current interface{}, container *bufferContainer) errorRuntime {

	switch typedNodes := current.(type) {
	case map[string]interface{}:
		return f.retrieveMap(root, typedNodes, container)

	case []interface{}:
		return f.retrieveList(root, typedNodes, container)

	default:
		foundType := msgTypeNull
		if current != nil {
			foundType = reflect.TypeOf(current).String()
		}
		return newErrorTypeUnmatched(f.errorRuntime.node, msgTypeObjectOrArray, foundType)
	}
}

func (f *syntaxFilterQualifier) retrieveMap(
	root interface{}, srcMap map[string]interface{}, container *bufferContainer) errorRuntime {

	var deepestTextLen int
	var deepestError errorRuntime

	sortKeys := getSortedKeys(srcMap)

	valueList := make([]interface{}, len(*sortKeys))
	for index := range *sortKeys {
		valueList[index] = srcMap[(*sortKeys)[index]]
	}

	valueList = f.query.compute(root, valueList)

	isEachResult := len(valueList) == len(srcMap)

	if !isEachResult {
		if valueList[0] == emptyEntity {
			return newErrorMemberNotExist(f.errorRuntime.node)
		}
	}

	for index := range *sortKeys {
		if isEachResult {
			if valueList[index] == emptyEntity {
				continue
			}
		}
		if err := f.retrieveMapNext(root, srcMap, (*sortKeys)[index], container); err != nil {
			if len(container.result) == 0 {
				deepestTextLen, deepestError = f.addDeepestError(err, deepestTextLen, deepestError)
			}
		}
	}

	putSortSlice(sortKeys)

	if len(container.result) > 0 {
		return nil
	}

	if deepestError == nil {
		return newErrorMemberNotExist(f.errorRuntime.node)
	}

	return deepestError
}

func (f *syntaxFilterQualifier) retrieveList(
	root interface{}, srcList []interface{}, container *bufferContainer) errorRuntime {

	var deepestTextLen int
	var deepestError errorRuntime

	valueList := f.query.compute(root, srcList)

	isEachResult := len(valueList) == len(srcList)

	if !isEachResult {
		if valueList[0] == emptyEntity {
			return newErrorMemberNotExist(f.errorRuntime.node)
		}
	}

	for index := range srcList {
		if isEachResult {
			if valueList[index] == emptyEntity {
				continue
			}
		}
		if err := f.retrieveListNext(root, srcList, index, container); err != nil {
			if len(container.result) == 0 {
				deepestTextLen, deepestError = f.addDeepestError(err, deepestTextLen, deepestError)
			}
		}
	}

	if len(container.result) > 0 {
		return nil
	}

	if deepestError == nil {
		return newErrorMemberNotExist(f.errorRuntime.node)
	}

	return deepestError
}
