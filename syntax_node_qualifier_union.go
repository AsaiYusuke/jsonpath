package jsonpath

import "reflect"

type syntaxUnionQualifier struct {
	*syntaxBasicNode

	subscripts []syntaxSubscript
}

func (u *syntaxUnionQualifier) retrieve(root, current interface{}) error {
	if _, ok := current.(map[string]interface{}); ok {
		if len(u.subscripts) == 1 {
			if _, ok := u.subscripts[0].(*syntaxAsterisk); ok {
				// Switch to the all node analysis mode,
				// if "current" variable points the map structure and
				// specifying the Asterisk subscript
				asteriskIdentifier := syntaxChildAsteriskIdentifier{
					syntaxBasicNode: &syntaxBasicNode{
						text:   u.text,
						next:   u.next,
						result: u.result,
					},
				}
				return asteriskIdentifier.retrieve(root, current)
			}
		}
	}

	srcArray, ok := current.([]interface{})
	if !ok {
		foundType := `null`
		if current != nil {
			foundType = reflect.TypeOf(current).String()
		}
		return ErrorTypeUnmatched{`array`, foundType, u.text}
	}

	indexes := make([]int, 0)

	for _, subscript := range u.subscripts {
		indexes = append(indexes, subscript.getIndexes(srcArray)...)
	}

	if u.isMultiValue() {
		for _, index := range indexes {
			u.retrieveNext(root, srcArray[index])
		}

		if len((**u.result)) == 0 {
			return ErrorNoneMatched{u.getConnectedText()}
		}

		return nil
	}

	if len(indexes) == 0 {
		return ErrorIndexOutOfRange{u.text}
	}

	return u.retrieveNext(root, srcArray[indexes[0]])
}

func (u *syntaxUnionQualifier) merge(union *syntaxUnionQualifier) {
	u.subscripts = append(u.subscripts, union.subscripts...)
}
