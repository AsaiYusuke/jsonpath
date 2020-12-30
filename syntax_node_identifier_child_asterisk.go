package jsonpath

import (
	"sort"
)

type syntaxChildAsteriskIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxChildAsteriskIdentifier) retrieve(
	root, current interface{}, result *[]interface{}) error {

	childErrorMap := make(map[error]bool, 1)
	var lastError error

	switch current.(type) {
	case map[string]interface{}:
		lastError = i.retrieveMap(
			root, current.(map[string]interface{}), result, childErrorMap)

	case []interface{}:
		lastError = i.retrieveList(
			root, current.([]interface{}), result, childErrorMap)

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

func (i *syntaxChildAsteriskIdentifier) retrieveMap(
	root interface{}, srcMap map[string]interface{}, result *[]interface{},
	childErrorMap map[error]bool) error {

	var lastError error

	index, keys := 0, make([]string, len(srcMap))
	for key := range srcMap {
		keys[index] = key
		index++
	}
	sort.Strings(keys)
	for _, key := range keys {
		localKey := key
		err := i.retrieveNext(
			root, result,
			func() interface{} {
				return srcMap[localKey]
			},
			func(value interface{}) {
				srcMap[localKey] = value
			})
		if err != nil {
			childErrorMap[err] = true
			lastError = err
		}
	}

	return lastError
}

func (i *syntaxChildAsteriskIdentifier) retrieveList(
	root interface{}, srcList []interface{}, result *[]interface{},
	childErrorMap map[error]bool) error {

	var lastError error

	for index := range srcList {
		localIndex := index
		err := i.retrieveNext(
			root, result,
			func() interface{} {
				return srcList[localIndex]
			},
			func(value interface{}) {
				srcList[localIndex] = value
			})
		if err != nil {
			childErrorMap[err] = true
			lastError = err
		}
	}

	return lastError
}
