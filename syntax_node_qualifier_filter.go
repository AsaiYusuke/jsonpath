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
			return ErrorNoneMatched{path: f.text}
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
	var partialFound bool

	index, keys := 0, make(sort.StringSlice, len(srcMap))
	for key := range srcMap {
		keys[index] = key
		index++
	}
	keys.Sort()
	argumentMap := make(map[int]interface{}, len(keys))
	for index := range keys {
		argumentMap[index] = srcMap[keys[index]]
	}

	computedMap := f.query.compute(root, argumentMap)

	if len(computedMap) > 0 {
		for index := range keys {
			if _, ok := computedMap[index]; ok {
				partialFound = true
				localKey := keys[index]
				err := f.retrieveNext(
					root, result,
					func() interface{} {
						return srcMap[localKey]
					},
					func(value interface{}) {
						srcMap[localKey] = value
					})
				if err != nil {
					childErrorMap[err] = struct{}{}
					lastError = err
				}
			}
		}
	}

	if !partialFound {
		err := ErrorMemberNotExist{path: f.text}
		childErrorMap[err] = struct{}{}
		lastError = err
	}

	return lastError
}

func (f *syntaxFilterQualifier) retrieveList(
	root interface{}, srcList []interface{}, result *[]interface{},
	childErrorMap map[error]struct{}) error {

	var lastError error
	var partialFound bool

	argumentMap := make(map[int]interface{}, len(srcList))
	for index := range srcList {
		argumentMap[index] = srcList[index]
	}

	computedMap := f.query.compute(root, argumentMap)

	if len(computedMap) > 0 {
		for index := range srcList {
			if _, ok := computedMap[index]; ok {
				partialFound = true
				localIndex := index
				err := f.retrieveNext(
					root, result,
					func() interface{} {
						return srcList[localIndex]
					},
					func(value interface{}) {
						srcList[localIndex] = value
					})
				if err != nil {
					childErrorMap[err] = struct{}{}
					lastError = err
				}
			}
		}
	}

	if !partialFound {
		err := ErrorMemberNotExist{path: f.text}
		childErrorMap[err] = struct{}{}
		lastError = err
	}

	return lastError
}
