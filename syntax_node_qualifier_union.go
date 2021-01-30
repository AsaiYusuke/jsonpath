package jsonpath

import "reflect"

type syntaxUnionQualifier struct {
	*syntaxBasicNode

	subscripts []syntaxSubscript
}

func (u *syntaxUnionQualifier) retrieve(
	root, current interface{}, container *bufferContainer) errorRuntime {

	srcArray, ok := current.([]interface{})
	if !ok {
		foundType := msgTypeNull
		if current != nil {
			foundType = reflect.TypeOf(current).String()
		}
		return ErrorTypeUnmatched{
			errorBasicRuntime: u.errorRuntime,
			expectedType:      msgTypeArray,
			foundType:         foundType,
		}
	}

	var deepestTextLen int
	var deepestError errorRuntime

	for _, subscript := range u.subscripts {
		for _, index := range subscript.getIndexes(srcArray) {
			if err := u.retrieveListNext(root, srcArray, index, container); err != nil {
				if len(container.result) == 0 {
					deepestTextLen, deepestError = u.addDeepestError(err, deepestTextLen, deepestError)
				}
			}
		}
	}

	if len(container.result) > 0 {
		return nil
	}

	if deepestError == nil {
		return ErrorMemberNotExist{
			errorBasicRuntime: u.errorRuntime,
		}
	}

	return deepestError
}

func (u *syntaxUnionQualifier) merge(union *syntaxUnionQualifier) {
	u.subscripts = append(u.subscripts, union.subscripts...)
}
