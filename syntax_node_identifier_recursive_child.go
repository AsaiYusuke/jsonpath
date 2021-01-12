package jsonpath

import "sort"

type syntaxRecursiveChildIdentifier struct {
	*syntaxBasicNode

	nextMapRequired  bool
	nextListRequired bool
}

func (i *syntaxRecursiveChildIdentifier) retrieve(
	root, current interface{}, result *[]interface{}) error {

	targetNodes := make([]interface{}, 1, 5)
	targetNodes[0] = current

	keys := make(sort.StringSlice, 0, 2)

	for len(targetNodes) > 0 {
		currentNode := targetNodes[len(targetNodes)-1]
		targetNodes = targetNodes[:len(targetNodes)-1]
		switch typedNodes := currentNode.(type) {
		case map[string]interface{}:
			if i.nextMapRequired {
				i.retrieveAnyValueNext(root, typedNodes, result)
			}
			keys = keys[:0]
			for key := range typedNodes {
				keys = append(keys, key)
			}
			if len(keys) > 1 {
				keys.Sort()
			}
			for index := len(typedNodes) - 1; index >= 0; index-- {
				targetNodes = append(targetNodes, typedNodes[keys[index]])
			}
		case []interface{}:
			if i.nextListRequired {
				i.retrieveAnyValueNext(root, typedNodes, result)
			}
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
