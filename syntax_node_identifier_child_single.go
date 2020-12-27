package jsonpath

import (
	"reflect"
)

type syntaxChildSingleIdentifier struct {
	*syntaxBasicNode

	identifier string
}

func (i *syntaxChildSingleIdentifier) retrieve(current interface{}) error {

	switch current.(type) {
	case map[string]interface{}:
		srcMap := current.(map[string]interface{})
		_, ok := srcMap[i.identifier]
		if !ok {
			return ErrorMemberNotExist{i.text}
		}
		return i.retrieveNext(
			func() interface{} {
				return srcMap[i.identifier]
			},
			func(value interface{}) {
				srcMap[i.identifier] = value
			})

	case []interface{}:
		return ErrorTypeUnmatched{`object`, reflect.TypeOf(current).String(), i.text}
	}

	foundType := `null`
	if current != nil {
		foundType = reflect.TypeOf(current).String()
	}
	return ErrorTypeUnmatched{`object/array`, foundType, i.text}
}
