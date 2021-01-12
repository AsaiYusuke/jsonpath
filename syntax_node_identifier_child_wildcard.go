package jsonpath

import "sort"

type syntaxChildWildcardIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxChildWildcardIdentifier) retrieve(
	root, current interface{}, result *[]interface{}) error {

	childErrorMap := make(map[error]struct{}, 1)
	var lastError error

	switch typedNodes := current.(type) {
	case map[string]interface{}:
		lastError = i.retrieveMap(root, typedNodes, result, childErrorMap)

	case []interface{}:
		lastError = i.retrieveList(root, typedNodes, result, childErrorMap)

	}

	if len(*result) == 0 {
		switch len(childErrorMap) {
		case 0:
			return ErrorNoneMatched{path: i.text}
		case 1:
			return lastError
		default:
			return ErrorNoneMatched{path: i.next.getConnectedText()}
		}
	}

	return nil
}

func (i *syntaxChildWildcardIdentifier) retrieveMap(
	root interface{}, srcMap map[string]interface{}, result *[]interface{},
	childErrorMap map[error]struct{}) error {

	var lastError error

	index, keys := 0, make(sort.StringSlice, len(srcMap))
	for key := range srcMap {
		keys[index] = key
		index++
	}
	if len(keys) > 1 {
		keys.Sort()
	}
	for _, key := range keys {
		if err := i.retrieveMapNext(root, srcMap, key, result); err != nil {
			childErrorMap[err] = struct{}{}
			lastError = err
		}
	}

	return lastError
}

func (i *syntaxChildWildcardIdentifier) retrieveList(
	root interface{}, srcList []interface{}, result *[]interface{},
	childErrorMap map[error]struct{}) error {

	var lastError error

	for index := range srcList {
		if err := i.retrieveListNext(root, srcList, index, result); err != nil {
			childErrorMap[err] = struct{}{}
			lastError = err
		}
	}

	return lastError
}
