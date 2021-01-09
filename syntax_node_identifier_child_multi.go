package jsonpath

import "reflect"

type syntaxChildMultiIdentifier struct {
	*syntaxBasicNode

	identifiers []string
}

func (i *syntaxChildMultiIdentifier) retrieve(
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

	childErrorMap := make(map[error]struct{}, 1)

	lastError := i.retrieveMap(root, srcMap, result, childErrorMap)

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
	childErrorMap map[error]struct{}) error {

	var lastError error
	var partialFound bool

	for index := range i.identifiers {
		if _, ok := srcMap[i.identifiers[index]]; ok {
			partialFound = true
			if err := i.retrieveMapNext(root, srcMap, i.identifiers[index], result); err != nil {
				childErrorMap[err] = struct{}{}
				lastError = err
			}
		}
	}

	if !partialFound {
		err := ErrorMemberNotExist{path: i.text}
		childErrorMap[err] = struct{}{}
		lastError = err
	}

	return lastError
}
