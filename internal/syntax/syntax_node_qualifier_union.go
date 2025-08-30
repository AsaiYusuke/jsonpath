package syntax

import (
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
		return u.newErrTypeUnmatched(msgTypeArray, current)
	}

	var deepestError errors.ErrorRuntime

	srcLen := len(srcArray)
	for _, subscript := range u.subscripts {
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
