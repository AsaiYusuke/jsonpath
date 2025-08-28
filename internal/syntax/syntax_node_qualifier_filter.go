package syntax

import (
	"reflect"

	"github.com/AsaiYusuke/jsonpath/v2/errors"
)

type syntaxFilterQualifier struct {
	*syntaxBasicNode

	query syntaxQuery
}

func (f *syntaxFilterQualifier) retrieve(
	root, current any, results *[]any) errors.ErrorRuntime {

	switch typedNodes := current.(type) {
	case map[string]any:
		return f.retrieveMap(root, typedNodes, results)

	case []any:
		return f.retrieveList(root, typedNodes, results)

	default:
		if current != nil {
			return errors.NewErrorTypeUnmatched(
				f.path, f.remainingPathLen, msgTypeObjectOrArray, reflect.TypeOf(current).String())
		}
		return errors.NewErrorTypeUnmatched(
			f.path, f.remainingPathLen, msgTypeObjectOrArray, msgTypeNull)
	}
}

func (f *syntaxFilterQualifier) retrieveMap(
	root any, srcMap map[string]any, results *[]any) errors.ErrorRuntime {

	if len(srcMap) == 0 {
		return f.newErrMemberNotExist()
	}

	var deepestError errors.ErrorRuntime

	sortKeys, keyLength := getSortedKeys(srcMap)

	buf := getNodeSlice()
	if cap(*buf) < keyLength {
		*buf = make([]any, keyLength)
	}
	*buf = (*buf)[:keyLength]
	for index := range *sortKeys {
		(*buf)[index] = srcMap[(*sortKeys)[index]]
	}

	valueList := f.query.compute(root, *buf)

	putNodeSlice(buf)

	isEachResult := len(valueList) == len(srcMap)

	if !isEachResult {
		if valueList[0] == emptyEntity {
			putSortSlice(sortKeys)
			return f.newErrMemberNotExist()
		}
	}

	for index := range *sortKeys {
		if isEachResult {
			if valueList[index] == emptyEntity {
				continue
			}
		}
		if err := f.retrieveMapNext(root, srcMap, (*sortKeys)[index], results); len(*results) == 0 && err != nil {
			deepestError = f.getMostResolvedError(err, deepestError)
		}
	}

	putSortSlice(sortKeys)

	if len(*results) > 0 {
		return nil
	}

	if deepestError == nil {
		return f.newErrMemberNotExist()
	}

	return deepestError
}

func (f *syntaxFilterQualifier) retrieveList(
	root any, srcList []any, results *[]any) errors.ErrorRuntime {

	if len(srcList) == 0 {
		return f.newErrMemberNotExist()
	}

	var deepestError errors.ErrorRuntime

	valueList := f.query.compute(root, srcList)

	isEachResult := len(valueList) == len(srcList)

	if !isEachResult {
		if valueList[0] == emptyEntity {
			return f.newErrMemberNotExist()
		}
	}

	for index := range srcList {
		if isEachResult {
			if valueList[index] == emptyEntity {
				continue
			}
		}
		if err := f.retrieveListNext(root, srcList, index, results); len(*results) == 0 && err != nil {
			deepestError = f.getMostResolvedError(err, deepestError)
		}
	}

	if len(*results) > 0 {
		return nil
	}

	if deepestError == nil {
		return f.newErrMemberNotExist()
	}

	return deepestError
}
