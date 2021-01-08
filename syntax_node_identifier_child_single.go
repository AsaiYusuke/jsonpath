package jsonpath

import "reflect"

type syntaxChildSingleIdentifier struct {
	*syntaxBasicNode

	identifier string
}

func (i *syntaxChildSingleIdentifier) retrieve(
	root, current interface{}, result *[]interface{}) error {

	srcMap, ok := current.(map[string]interface{})
	if !ok {
		foundType := `null`
		if current != nil {
			foundType = reflect.TypeOf(current).String()
		}
		return ErrorTypeUnmatched{
			expectedType: `object`,
			foundType:    foundType,
			path:         i.text,
		}
	}

	if _, ok := srcMap[i.identifier]; !ok {
		return ErrorMemberNotExist{path: i.text}
	}

	return i.retrieveNext(
		root, result,
		func() interface{} {
			return srcMap[i.identifier]
		},
		func(value interface{}) {
			srcMap[i.identifier] = value
		})
}
