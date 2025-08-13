package syntax

import (
	"reflect"

	"github.com/AsaiYusuke/jsonpath/errors"
)

type syntaxUnionQualifier struct {
	*syntaxBasicNode

	subscripts []syntaxSubscript
}

func (u *syntaxUnionQualifier) retrieve(
	root, current any, container *bufferContainer) errors.ErrorRuntime {

	srcArray, ok := current.([]any)
	if !ok {
		if current != nil {
			return errors.NewErrorTypeUnmatched(
				u.path, u.remainingPathLen, msgTypeArray, reflect.TypeOf(current).String())
		}
		return errors.NewErrorTypeUnmatched(
			u.path, u.remainingPathLen, msgTypeArray, msgTypeNull)
	}

	var deepestError errors.ErrorRuntime

	for _, subscript := range u.subscripts {
		for _, index := range subscript.getIndexes(len(srcArray)) {
			if err := u.retrieveListNext(root, srcArray, index, container); err != nil {
				if len(container.result) == 0 {
					deepestError = u.getMostResolvedError(err, deepestError)
				}
			}
		}
	}

	if len(container.result) > 0 {
		return nil
	}

	if deepestError == nil {
		return u.newErrMemberNotExist()
	}

	return deepestError
}

func (u *syntaxUnionQualifier) merge(union *syntaxUnionQualifier) {
	u.subscripts = append(u.subscripts, union.subscripts...)
}
