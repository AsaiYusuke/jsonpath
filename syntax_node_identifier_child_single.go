package jsonpath

import (
	"reflect"
)

type syntaxChildSingleIdentifier struct {
	*syntaxBasicNode

	identifier string
}

func (i syntaxChildSingleIdentifier) retrieve(
	root, current interface{}, result *resultContainer) error {

	switch current.(type) {
	case map[string]interface{}:
		srcMap := current.(map[string]interface{})
		identifier := i.identifier
		child, ok := srcMap[identifier]
		if !ok {
			return ErrorMemberNotExist{i.text}
		}
		return i.retrieveNext(root, child, result)

	case []interface{}:
		if len(i.identifier) > 0 {
			return ErrorTypeUnmatched{`object`, reflect.TypeOf(current).String(), i.text}
		}
		return i.retrieveNext(root, current, result)

	}

	foundType := `null`
	if current != nil {
		foundType = reflect.TypeOf(current).String()
	}
	return ErrorTypeUnmatched{`object/array`, foundType, i.text}
}
