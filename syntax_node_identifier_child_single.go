package jsonpath

import (
	"reflect"
)

type syntaxChildSingleIdentifier struct {
	*syntaxBasicNode

	identifier string
}

func (i *syntaxChildSingleIdentifier) retrieve(root, current interface{}) error {

	switch current.(type) {
	case map[string]interface{}:
		srcMap := current.(map[string]interface{})
		child, ok := srcMap[i.identifier]
		if !ok {
			return ErrorMemberNotExist{i.text}
		}
		return i.retrieveNext(root, child)

	case []interface{}:
		return ErrorTypeUnmatched{`object`, reflect.TypeOf(current).String(), i.text}
	}

	foundType := `null`
	if current != nil {
		foundType = reflect.TypeOf(current).String()
	}
	return ErrorTypeUnmatched{`object/array`, foundType, i.text}
}
