package jsonpath

import (
	"sort"
)

type syntaxRecursiveChildIdentifier struct {
	syntaxChildIdentifier
}

func (i syntaxRecursiveChildIdentifier) retrieve(root, current interface{}, result *resultContainer) error {
	switch current.(type) {
	case map[string]interface{}:
		srcMap := current.(map[string]interface{})
		isRequiredAsteriskCheck := len(i.identifiers) == 1 && len(i.identifiers[0]) == 0
		_, ok := srcMap[i.identifiers[0]]
		if !ok && isRequiredAsteriskCheck {
			// If identifier is "", additionally checks whether this qualifier is [*].
			i.retrieveNext(root, srcMap, result)
		} else {
			i.syntaxChildIdentifier.retrieve(root, srcMap, result)
		}

		keys := make([]string, 0, len(srcMap))
		for key := range srcMap {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			child := srcMap[key]
			i.retrieve(root, child, result)
		}

	case []interface{}:
		srcArray := current.([]interface{})
		i.syntaxChildIdentifier.retrieve(root, srcArray, result)

		for _, child := range srcArray {
			i.retrieve(root, child, result)
		}
	}

	if !result.hasResult() {
		return ErrorNoneMatched{i.getConnectedText()}
	}
	return nil
}
