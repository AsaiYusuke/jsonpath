package jsonpath

import (
	"sort"
)

type syntaxRecursiveChildIdentifier struct {
	*syntaxBasicNode
}

func (i syntaxRecursiveChildIdentifier) retrieve(root, current interface{}, result *[]interface{}) error {
	switch current.(type) {
	case map[string]interface{}:
		srcMap := current.(map[string]interface{})
		i.retrieveNext(root, srcMap, result)

		index, keys := 0, make([]string, len(srcMap))
		for key := range srcMap {
			keys[index] = key
			index++
		}
		sort.Strings(keys)
		for _, key := range keys {
			i.retrieve(root, srcMap[key], result)
		}

	case []interface{}:
		srcArray := current.([]interface{})
		i.retrieveNext(root, srcArray, result)

		for _, child := range srcArray {
			i.retrieve(root, child, result)
		}
	}

	if len(*result) == 0 {
		return ErrorNoneMatched{i.getConnectedText()}
	}
	return nil
}
