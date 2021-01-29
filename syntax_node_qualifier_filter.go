package jsonpath

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
		return ErrorTypeUnmatched{
			errorBasicRuntime: f.errorRuntime,
			expectedType:      msgTypeObjectOrArray,
			foundType:         foundType,
		}
	}
}

func (f *syntaxFilterQualifier) retrieveMap(
	root interface{}, srcMap map[string]interface{}, container *bufferContainer) errorRuntime {

	var deepestTextLen int
	var deepestError errorRuntime

	sortKeys := container.getSortedKeys(srcMap)

	valueList := make([]interface{}, len(*sortKeys))
	for index := range *sortKeys {
		valueList[index] = srcMap[(*sortKeys)[index]]
	}

	valueList = f.query.compute(root, valueList, container)

	isEachResult := len(valueList) == len(srcMap)

	var nodeNotFound bool
	if !isEachResult {
		_, nodeNotFound = valueList[0].(struct{})
		if nodeNotFound {
			return ErrorMemberNotExist{
				errorBasicRuntime: f.errorRuntime,
			}
		}
	}

	for index := range *sortKeys {
		if isEachResult {
			_, nodeNotFound = valueList[index].(struct{})
		}
		if nodeNotFound {
			continue
		}
		if err := f.retrieveMapNext(root, srcMap, (*sortKeys)[index], container); err != nil {
			if len(container.result) == 0 {
				deepestTextLen, deepestError = f.addDeepestError(err, deepestTextLen, deepestError)
			}
		}
	}

	container.putSortSlice(sortKeys)

	if len(container.result) > 0 {
		return nil
	}

	if deepestError == nil {
		return ErrorMemberNotExist{
			errorBasicRuntime: f.errorRuntime,
		}
	}

	return deepestError
}

func (f *syntaxFilterQualifier) retrieveList(
	root interface{}, srcList []interface{}, container *bufferContainer) errorRuntime {

	var deepestTextLen int
	var deepestError errorRuntime

	valueList := f.query.compute(root, srcList, container)

	isEachResult := len(valueList) == len(srcList)

	var nodeNotFound bool
	if !isEachResult {
		_, nodeNotFound = valueList[0].(struct{})
		if nodeNotFound {
			return ErrorMemberNotExist{
				errorBasicRuntime: f.errorRuntime,
			}
		}
	}

	for index := range srcList {
		if isEachResult {
			_, nodeNotFound = valueList[index].(struct{})
		}
		if nodeNotFound {
			continue
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
		return ErrorMemberNotExist{
			errorBasicRuntime: f.errorRuntime,
		}
	}

	return deepestError
}
