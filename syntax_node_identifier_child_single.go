package jsonpath

import "reflect"

type syntaxChildSingleIdentifier struct {
	*syntaxBasicNode

	identifier string
}

func (i *syntaxChildSingleIdentifier) retrieve(
	root, current interface{}, result *[]interface{}) error {

	switch typedNodes := current.(type) {
	case map[string]interface{}:
		_, ok := typedNodes[i.identifier]
		if !ok {
			return ErrorMemberNotExist{path: i.text}
		}
		return i.retrieveNext(
			root, result,
			func() interface{} {
				return typedNodes[i.identifier]
			},
			func(value interface{}) {
				typedNodes[i.identifier] = value
			})

	case []interface{}:
		return ErrorTypeUnmatched{
			expectedType: `object`,
			foundType:    `array`,
			path:         i.text,
		}
	}

	foundType := `null`
	if current != nil {
		foundType = reflect.TypeOf(current).String()
	}

	return ErrorTypeUnmatched{
		expectedType: `object/array`,
		foundType:    foundType,
		path:         i.text,
	}
}
