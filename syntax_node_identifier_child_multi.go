package jsonpath

import (
	"reflect"
)

type syntaxChildMultiIdentifier struct {
	*syntaxBasicNode

	identifiers []string
}

func (i *syntaxChildMultiIdentifier) retrieve(
	root, current interface{}, result *[]interface{}) error {

	childErrorMap := make(map[error]bool, 1)
	var lastError error

	switch current.(type) {
	case map[string]interface{}:
		lastError = i.retrieveMap(
			root, current.(map[string]interface{}), result, childErrorMap)

	case []interface{}:
		return ErrorTypeUnmatched{
			expectedType: `object`,
			foundType:    reflect.TypeOf(current).String(),
			path:         i.text,
		}
	}

	if len(*result) == 0 {
		switch len(childErrorMap) {
		case 0:
			return ErrorNoneMatched{path: i.text}
		case 1:
			return lastError
		default:
			return ErrorNoneMatched{path: i.next.getConnectedText()}
		}
	}

	return nil
}

func (i *syntaxChildMultiIdentifier) retrieveMap(
	root interface{}, srcMap map[string]interface{}, result *[]interface{},
	childErrorMap map[error]bool) error {

	var lastError error

	for _, key := range i.identifiers {
		if _, ok := srcMap[key]; ok {
			localKey := key
			err := i.retrieveNext(
				root, result,
				func() interface{} {
					return srcMap[localKey]
				},
				func(value interface{}) {
					srcMap[localKey] = value
				})
			if err != nil {
				childErrorMap[err] = true
				lastError = err
			}
		}
	}

	return lastError
}
