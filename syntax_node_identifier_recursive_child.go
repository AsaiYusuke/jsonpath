package jsonpath

import (
	"sort"
)

type syntaxRecursiveChildIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxRecursiveChildIdentifier) retrieve(current interface{}) error {
	switch current.(type) {
	case map[string]interface{}:
		srcMap := current.(map[string]interface{})
		i.retrieveNext(srcMap)

		index, keys := 0, make([]string, len(srcMap))
		for key := range srcMap {
			keys[index] = key
			index++
		}
		sort.Strings(keys)
		for _, key := range keys {
			i.retrieve(srcMap[key])
		}

	case []interface{}:
		srcArray := current.([]interface{})
		i.retrieveNext(srcArray)

		for _, child := range srcArray {
			i.retrieve(child)
		}
	}

	if len((**i.result)) == 0 {
		return ErrorNoneMatched{i.getConnectedText()}
	}
	return nil
}
