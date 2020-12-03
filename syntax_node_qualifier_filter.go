package jsonpath

import (
	"sort"
)

type syntaxFilter struct {
	*syntaxBasicNode

	query syntaxQuery
}

func (f syntaxFilter) retrieve(root, current interface{}, result *resultContainer) error {
	switch current.(type) {
	case map[string]interface{}:
		srcMap := current.(map[string]interface{})
		keys := make([]string, 0, len(srcMap))
		for key := range srcMap {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		argumentMap := make(map[int]interface{}, len(keys))
		for index, key := range keys {
			argumentMap[index] = srcMap[key]
		}

		computedMap := f.query.compute(root, argumentMap)

		if len(computedMap) > 0 {
			for index, key := range keys {
				if _, ok := computedMap[index]; ok {
					f.retrieveNext(root, srcMap[key], result)
				}
			}
		}

	case []interface{}:
		srcList := current.([]interface{})

		argumentMap := make(map[int]interface{}, len(srcList))
		for index, entity := range srcList {
			argumentMap[index] = entity
		}

		computedMap := f.query.compute(root, argumentMap)

		if len(computedMap) > 0 {
			for index, entity := range srcList {
				if _, ok := computedMap[index]; ok {
					f.retrieveNext(root, entity, result)
				}
			}
		}

	}

	if !result.hasResult() {
		return ErrorNoneMatched{f.getConnectedText()}
	}

	return nil
}
