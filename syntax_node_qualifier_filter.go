package jsonpath

import (
	"sort"
)

type syntaxFilterQualifier struct {
	*syntaxBasicNode

	query syntaxQuery
}

func (f syntaxFilterQualifier) retrieve(root, current interface{}, result *[]interface{}) error {
	switch current.(type) {
	case map[string]interface{}:
		srcMap := current.(map[string]interface{})
		index, keys := 0, make([]string, len(srcMap))
		for key := range srcMap {
			keys[index] = key
			index++
		}
		sort.Strings(keys)
		argumentMap := make(map[int]interface{}, len(keys))
		for index, key := range keys {
			argumentMap[index] = srcMap[key]
		}

		computedMap := f.query.compute(root, argumentMap)

		if len(computedMap) > 0 {
			for index := range keys {
				if _, ok := computedMap[index]; ok {
					f.retrieveNext(root, argumentMap[index], result)
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

	if len(*result) == 0 {
		return ErrorNoneMatched{f.getConnectedText()}
	}

	return nil
}
