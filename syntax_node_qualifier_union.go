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
		foundType := `null`
		if current != nil {
			foundType = reflect.TypeOf(current).String()
		}
		return ErrorTypeUnmatched{
			errorBasicRuntime: &errorBasicRuntime{
				node: u.syntaxBasicNode,
			},
			expectedType: `array`,
			foundType:    foundType,
		}
	}

	if u.isValueGroup() {
		var deepestTextLen int
		deepestErrors := make([]errorRuntime, 0, 2)

		for _, subscript := range u.subscripts {
			for _, index := range subscript.getIndexes(srcArray) {
				if err := u.retrieveListNext(root, srcArray, index, container); err != nil {
					deepestTextLen, deepestErrors = u.addDeepestError(err, deepestTextLen, deepestErrors)
				}
			}
		}

		if len(container.result) == 0 {
			switch len(deepestErrors) {
			case 0:
				return ErrorIndexOutOfRange{
					errorBasicRuntime: &errorBasicRuntime{
						node: u.syntaxBasicNode,
					},
				}
			case 1:
				return deepestErrors[0]
			default:
				return ErrorNoneMatched{
					errorBasicRuntime: &errorBasicRuntime{
						node: deepestErrors[0].getSyntaxNode(),
					},
				}
			}
		}

		return nil
	}

	indexes := u.subscripts[0].getIndexes(srcArray)
	if len(indexes) == 0 {
		return ErrorIndexOutOfRange{
			errorBasicRuntime: &errorBasicRuntime{
				node: u.syntaxBasicNode,
			},
		}
	}

	return u.retrieveListNext(root, srcArray, indexes[0], container)
}

func (u *syntaxUnionQualifier) merge(union *syntaxUnionQualifier) {
	u.subscripts = append(u.subscripts, union.subscripts...)
}
