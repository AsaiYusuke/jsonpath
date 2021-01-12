package jsonpath

import "reflect"

type syntaxUnionQualifier struct {
	*syntaxBasicNode

	subscripts []syntaxSubscript
}

func (u *syntaxUnionQualifier) retrieve(
	root, current interface{}, container *bufferContainer) error {

	srcArray, ok := current.([]interface{})
	if !ok {
		foundType := `null`
		if current != nil {
			foundType = reflect.TypeOf(current).String()
		}
		return ErrorTypeUnmatched{
			expectedType: `array`,
			foundType:    foundType,
			path:         u.text,
		}
	}

	var resultIndexes []int
	for _, subscript := range u.subscripts {
		resultIndexes = append(resultIndexes, subscript.getIndexes(srcArray)...)
	}

	if u.isValueGroup() {
		childErrorMap := make(map[error]struct{}, 1)
		var lastError error
		for _, indexes := range resultIndexes {
			if err := u.retrieveListNext(root, srcArray, indexes, container); err != nil {
				childErrorMap[err] = struct{}{}
				lastError = err
			}
		}

		if len(container.result) == 0 {
			switch len(childErrorMap) {
			case 0:
				return ErrorNoneMatched{path: u.text}
			case 1:
				return lastError
			default:
				return ErrorNoneMatched{path: u.next.getConnectedText()}
			}
		}

		return nil
	}

	if len(resultIndexes) == 0 {
		return ErrorIndexOutOfRange{path: u.text}
	}

	return u.retrieveListNext(root, srcArray, resultIndexes[0], container)
}

func (u *syntaxUnionQualifier) merge(union *syntaxUnionQualifier) {
	u.subscripts = append(u.subscripts, union.subscripts...)
}
