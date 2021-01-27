package jsonpath

import "reflect"

type syntaxFilterQualifier struct {
	*syntaxBasicNode

	query syntaxQuery
}

func (f *syntaxFilterQualifier) retrieve(
	root, current interface{}, container *bufferContainer) errorRuntime {

	deepestErrors := make([]errorRuntime, 0, 2)

	switch typedNodes := current.(type) {
	case map[string]interface{}:
		deepestErrors = f.retrieveMap(root, typedNodes, container, deepestErrors)

	case []interface{}:
		deepestErrors = f.retrieveList(root, typedNodes, container, deepestErrors)

	default:
		foundType := `null`
		if current != nil {
			foundType = reflect.TypeOf(current).String()
		}
		return ErrorTypeUnmatched{
			errorBasicRuntime: &errorBasicRuntime{
				node: f.syntaxBasicNode,
			},
			expectedType: `object/array`,
			foundType:    foundType,
		}
	}

	switch len(deepestErrors) {
	case 0:
		return nil
	case 1:
		return deepestErrors[0]
	default:
		return ErrorNoneMatched{
			errorBasicRuntime: &errorBasicRuntime{
				node: deepestErrors[0].getSyntaxNode(),
			},
		}
	}
}

func (f *syntaxFilterQualifier) retrieveMap(
	root interface{}, srcMap map[string]interface{}, container *bufferContainer,
	deepestErrors []errorRuntime) []errorRuntime {

	var deepestTextLen int

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
			return append(deepestErrors, ErrorMemberNotExist{
				errorBasicRuntime: &errorBasicRuntime{
					node: f.syntaxBasicNode,
				},
			})
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
				deepestTextLen, deepestErrors = f.addDeepestError(err, deepestTextLen, deepestErrors)
			}
		}
	}

	container.putSortSlice(sortKeys)

	if len(container.result) > 0 {
		return nil
	}

	if len(deepestErrors) == 0 {
		return append(deepestErrors, ErrorMemberNotExist{
			errorBasicRuntime: &errorBasicRuntime{
				node: f.syntaxBasicNode,
			},
		})
	}

	return deepestErrors
}

func (f *syntaxFilterQualifier) retrieveList(
	root interface{}, srcList []interface{}, container *bufferContainer,
	deepestErrors []errorRuntime) []errorRuntime {

	var deepestTextLen int

	valueList := f.query.compute(root, srcList, container)

	isEachResult := len(valueList) == len(srcList)

	var nodeNotFound bool
	if !isEachResult {
		_, nodeNotFound = valueList[0].(struct{})
		if nodeNotFound {
			return append(deepestErrors, ErrorMemberNotExist{
				errorBasicRuntime: &errorBasicRuntime{
					node: f.syntaxBasicNode,
				},
			})
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
				deepestTextLen, deepestErrors = f.addDeepestError(err, deepestTextLen, deepestErrors)
			}
		}
	}

	if len(container.result) > 0 {
		return nil
	}

	if len(deepestErrors) == 0 {
		return append(deepestErrors, ErrorMemberNotExist{
			errorBasicRuntime: &errorBasicRuntime{
				node: f.syntaxBasicNode,
			},
		})
	}

	return deepestErrors
}
