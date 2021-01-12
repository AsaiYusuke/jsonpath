package jsonpath

import "sort"

type syntaxFilterQualifier struct {
	*syntaxBasicNode

	query syntaxQuery
}

func (f *syntaxFilterQualifier) retrieve(
	root, current interface{}, result *[]interface{}) error {

	childErrorMap := make(map[error]struct{}, 1)
	var lastError error

	switch typedNodes := current.(type) {
	case map[string]interface{}:
		lastError = f.retrieveMap(root, typedNodes, result, childErrorMap)

	case []interface{}:
		lastError = f.retrieveList(root, typedNodes, result, childErrorMap)

	}

	if len(*result) == 0 {
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
	valueList := make([]interface{}, len(keys))
	for index := range keys {
		valueList[index] = srcMap[keys[index]]
	}

	valueList = f.query.compute(root, valueList)

	for index := range keys {
		var nodeNotFound bool
		if len(valueList) == 1 {
			_, nodeNotFound = valueList[0].(struct{})
		} else {
			_, nodeNotFound = valueList[index].(struct{})
		}
		if !nodeNotFound {
			if err := f.retrieveMapNext(root, srcMap, keys[index], result); err != nil {
				childErrorMap[err] = struct{}{}
				lastError = err
			}
		}
	}

	return lastError
}

func (f *syntaxFilterQualifier) retrieveList(
	root interface{}, srcList []interface{}, result *[]interface{},
	childErrorMap map[error]struct{}) error {

	var lastError error

	valueList := f.query.compute(root, srcList)

	for index := range srcList {
		var nodeNotFound bool
		if len(valueList) == 1 {
			_, nodeNotFound = valueList[0].(struct{})
		} else {
			_, nodeNotFound = valueList[index].(struct{})
		}
		if !nodeNotFound {
			if err := f.retrieveListNext(root, srcList, index, result); err != nil {
				childErrorMap[err] = struct{}{}
				lastError = err
			}
		}
	}

	return lastError
}
