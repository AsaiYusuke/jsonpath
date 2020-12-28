package jsonpath

import (
	"sort"
)

type syntaxChildAsteriskIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxChildAsteriskIdentifier) retrieve(
	root, current interface{}, result *[]interface{}) error {

	switch current.(type) {
	case map[string]interface{}:
		srcMap := current.(map[string]interface{})
		index, keys := 0, make([]string, len(srcMap))
		for key := range srcMap {
			keys[index] = key
			index++
		}
		sort.Strings(keys)
		for _, key := range keys {
			localKey := key
			i.retrieveNext(
				root, result,
				func() interface{} {
					return srcMap[localKey]
				},
				func(value interface{}) {
					srcMap[localKey] = value
				})
		}

	case []interface{}:
		srcArray := current.([]interface{})
		for index := range srcArray {
			localIndex := index
			i.retrieveNext(
				root, result,
				func() interface{} {
					return srcArray[localIndex]
				},
				func(value interface{}) {
					srcArray[localIndex] = value
				})
		}
	}

	if len(*result) == 0 {
		return ErrorNoneMatched{i.getConnectedText()}
	}

	return nil
}
