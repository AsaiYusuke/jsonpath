package jsonpath

type syntaxChildWildcardIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxChildWildcardIdentifier) retrieve(
	root, current interface{}, container *bufferContainer) errorRuntime {

	deepestErrors := make([]errorRuntime, 0, 2)

	switch typedNodes := current.(type) {
	case map[string]interface{}:
		deepestErrors = i.retrieveMap(root, typedNodes, container, deepestErrors)

	case []interface{}:
		deepestErrors = i.retrieveList(root, typedNodes, container, deepestErrors)

	}

	if len(container.result) == 0 {
		switch len(deepestErrors) {
		case 0:
			return ErrorMemberNotExist{
				errorBasicRuntime: &errorBasicRuntime{
					node: i.syntaxBasicNode,
				},
			}
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

	return nil
}

func (i *syntaxChildWildcardIdentifier) retrieveMap(
	root interface{}, srcMap map[string]interface{}, container *bufferContainer,
	deepestErrors []errorRuntime) []errorRuntime {

	deepestTextLen := -1

	sortKeys := container.getSortedKeys(srcMap)

	for _, key := range *sortKeys {
		if err := i.retrieveMapNext(root, srcMap, key, container); err != nil {
			deepestTextLen, deepestErrors = i.addDeepestError(err, deepestTextLen, deepestErrors)
		}
	}

	container.putSortSlice(sortKeys)

	return deepestErrors
}

func (i *syntaxChildWildcardIdentifier) retrieveList(
	root interface{}, srcList []interface{}, container *bufferContainer,
	deepestErrors []errorRuntime) []errorRuntime {

	deepestTextLen := -1

	for index := range srcList {
		if err := i.retrieveListNext(root, srcList, index, container); err != nil {
			deepestTextLen, deepestErrors = i.addDeepestError(err, deepestTextLen, deepestErrors)
		}
	}

	return deepestErrors
}
