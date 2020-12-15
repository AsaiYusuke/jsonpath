package jsonpath

import (
	"sort"
)

type syntaxChildAsteriskIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxChildAsteriskIdentifier) retrieve(current interface{}) error {

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
			i.retrieveNext(srcMap[key])
		}

	case []interface{}:
		for _, value := range current.([]interface{}) {
			i.retrieveNext(value)
		}
	}

	if len(**i.result) == 0 {
		return ErrorNoneMatched{i.getConnectedText()}
	}

	return nil
}
