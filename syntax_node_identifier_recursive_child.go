package jsonpath

import "sort"

type syntaxRecursiveChildIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxRecursiveChildIdentifier) retrieve(
	root, current interface{}, result *[]interface{}) error {

	switch current.(type) {
	case map[string]interface{}:
		srcMap := current.(map[string]interface{})
		i.retrieveNext(
			root, result,
			func() interface{} {
				return srcMap
			},
			nil)

		index, keys := 0, make(sort.StringSlice, len(srcMap))
		for key := range srcMap {
			keys[index] = key
			index++
		}
		keys.Sort()
		for index := range keys {
			i.retrieve(root, srcMap[keys[index]], result)
		}

	case []interface{}:
		srcArray := current.([]interface{})
		i.retrieveNext(
			root, result,
			func() interface{} {
				return srcArray
			},
			nil)

		for index := range srcArray {
			i.retrieve(root, srcArray[index], result)
		}
	}

	if len(*result) == 0 {
		return ErrorNoneMatched{i.getConnectedText()}
	}

	return nil
}
