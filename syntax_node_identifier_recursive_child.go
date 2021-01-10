package jsonpath

import "sort"

type syntaxRecursiveChildIdentifier struct {
	*syntaxBasicNode
}

func (i *syntaxRecursiveChildIdentifier) retrieve(
	root, current interface{}, result *[]interface{}) error {

	targetNodes := []interface{}{current}

	for len(targetNodes) > 0 {
		currentNode := targetNodes[len(targetNodes)-1]
		targetNodes = targetNodes[:len(targetNodes)-1]

		switch typedNodes := currentNode.(type) {
		case map[string]interface{}:
			i.retrieveAnyValueNext(root, typedNodes, result)

			index, keys := 0, make(sort.StringSlice, len(typedNodes))
			for key := range typedNodes {
				keys[index] = key
				index++
			}
			sort.Sort(sort.Reverse(keys))
			for index := range keys {
				targetNodes = append(targetNodes, typedNodes[keys[index]])
			}

		case []interface{}:
			i.retrieveAnyValueNext(root, typedNodes, result)

			for index := len(typedNodes) - 1; index >= 0; index-- {
				targetNodes = append(targetNodes, typedNodes[index])
			}
		}
	}

	if len(*result) == 0 {
		return ErrorNoneMatched{path: i.getConnectedText()}
	}

	return nil
}
