package jsonpath

import "sort"

type syntaxRecursiveChildIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxRecursiveChildIdentifier) retrieve(
	root, current interface{}, result *[]interface{}) error {

	switch typedNodes := current.(type) {
	case map[string]interface{}:
		i.retrieveAnyValueNext(root, typedNodes, result)

		index, keys := 0, make(sort.StringSlice, len(typedNodes))
		for key := range typedNodes {
			keys[index] = key
			index++
		}
		keys.Sort()
		for index := range keys {
			i.retrieve(root, typedNodes[keys[index]], result)
		}

	case []interface{}:
		i.retrieveAnyValueNext(root, typedNodes, result)

		for index := range typedNodes {
			i.retrieve(root, typedNodes[index], result)
		}
	}

	if len(*result) == 0 {
		return ErrorNoneMatched{path: i.getConnectedText()}
	}

	return nil
}
