package jsonpath

import (
	"reflect"
)

type syntaxChildMultiIdentifier struct {
	*syntaxBasicNode

	identifiers []string
}

func (i syntaxChildMultiIdentifier) retrieve(
	root interface{}, current interface{}, result *[]interface{}) error {

	switch current.(type) {
	case map[string]interface{}:
		srcMap := current.(map[string]interface{})
		for _, key := range i.identifiers {
			if _, ok := srcMap[key]; ok {
				i.retrieveNext(root, srcMap[key], result)
			}
		}

	case []interface{}:
		return ErrorTypeUnmatched{`object`, reflect.TypeOf(current).String(), i.text}
	}

	if len(*result) == 0 {
		return ErrorNoneMatched{i.getConnectedText()}
	}

	return nil
}
