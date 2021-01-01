package jsonpath

import "reflect"

type syntaxUnionQualifier struct {
	*syntaxBasicNode

	subscripts []syntaxSubscript
}

func (u *syntaxUnionQualifier) retrieve(
	root, current interface{}, result *[]interface{}) error {

	if _, ok := current.(map[string]interface{}); ok {
		if len(u.subscripts) == 1 {
			if _, ok := u.subscripts[0].(*syntaxWildcardSubscript); ok {
				// Switch to the all node analysis mode,
				// if "current" variable points the map structure and
				// specifying the Wildcard subscript
				wildcardIdentifier := syntaxChildWildcardIdentifier{
					syntaxBasicNode: &syntaxBasicNode{
						text: u.text,
						next: u.next,
					},
				}
				return wildcardIdentifier.retrieve(root, current, result)
			}
		}
	}

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
	for index := range u.subscripts {
		resultIndexes = append(resultIndexes, u.subscripts[index].getIndexes(srcArray)...)
	}

	if u.isValueGroup() {
		childErrorMap := make(map[error]struct{}, 1)
		var lastError error
		for index := range resultIndexes {
			localIndex := resultIndexes[index]
			err := u.retrieveNext(
				root, result,
				func() interface{} {
					return srcArray[localIndex]
				},
				func(value interface{}) {
					srcArray[localIndex] = value
				})
			if err != nil {
				childErrorMap[err] = struct{}{}
				lastError = err
			}
		}

		if len(*result) == 0 {
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

	return u.retrieveNext(
		root, result,
		func() interface{} {
			return srcArray[resultIndexes[0]]
		},
		func(value interface{}) {
			srcArray[resultIndexes[0]] = value
		})
}

func (u *syntaxUnionQualifier) merge(union *syntaxUnionQualifier) {
	u.subscripts = append(u.subscripts, union.subscripts...)
}
