package syntax

import (
	"reflect"

	"github.com/AsaiYusuke/jsonpath/v2/errors"
)

type syntaxUnionQualifier struct {
	*syntaxBasicNode

	subscripts []syntaxSubscript
}

func (u *syntaxUnionQualifier) retrieve(
	root, current any, results *[]any) errors.ErrorRuntime {

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
		srcLen := len(srcArray)
		for ord := range subscript.count(srcLen) {
			if err := u.retrieveListNext(root, srcArray, subscript.indexAt(srcLen, ord), results); len(*results) == 0 && err != nil {
				deepestError = u.getMostResolvedError(err, deepestError)
			}
		}
	}

	if len(*results) > 0 {
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
