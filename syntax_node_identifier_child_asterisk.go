package jsonpath

import (
	"sort"
)

type syntaxChildAsteriskIdentifier struct {
	*syntaxBasicNode
}

func (i syntaxChildAsteriskIdentifier) retrieve(
	root, current interface{}, result *resultContainer) error {

	switch current.(type) {
	case map[string]interface{}:
		srcMap := current.(map[string]interface{})
		keys := make([]string, 0, len(srcMap))
		for key := range srcMap {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			i.retrieveNext(root, srcMap[key], result)
		}

	case []interface{}:
		for _, value := range current.([]interface{}) {
			i.retrieveNext(root, value, result)
		}
	}

	if !result.hasResult() {
		return ErrorNoneMatched{i.getConnectedText()}
	}

	return nil
}
