package syntax

import (
	"reflect"

	"github.com/AsaiYusuke/jsonpath/errors"
)

type syntaxFilterQualifier struct {
	*syntaxBasicNode

	query syntaxQuery
}

func (f *syntaxFilterQualifier) retrieve(
	root, current any, container *bufferContainer) errors.ErrorRuntime {

	switch typedNodes := current.(type) {
	case map[string]any:
		return f.retrieveMap(root, typedNodes, container)

	case []any:
		return f.retrieveList(root, typedNodes, container)

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
	root any, srcMap map[string]any, container *bufferContainer) errors.ErrorRuntime {

	var deepestError errors.ErrorRuntime

	sortKeys, keyLength := getSortedKeys(srcMap)

	valueList := make([]any, keyLength)
	for index := range *sortKeys {
		valueList[index] = srcMap[(*sortKeys)[index]]
	}

	valueList = f.query.compute(root, valueList)

	isEachResult := len(valueList) == len(srcMap)

	if !isEachResult {
		if valueList[0] == emptyEntity {
			return f.newErrMemberNotExist()
		}
	}

	for index := range *sortKeys {
		if isEachResult {
			if valueList[index] == emptyEntity {
				continue
			}
		}
		if err := f.retrieveMapNext(root, srcMap, (*sortKeys)[index], container); len(container.result) == 0 && err != nil {
			deepestError = f.getMostResolvedError(err, deepestError)
		}
	}

	putSortSlice(sortKeys)

	if len(container.result) > 0 {
		return nil
	}

	if deepestError == nil {
		return f.newErrMemberNotExist()
	}

	return deepestError
}

func (f *syntaxFilterQualifier) retrieveList(
	root any, srcList []any, container *bufferContainer) errors.ErrorRuntime {

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
		if err := f.retrieveListNext(root, srcList, index, container); len(container.result) == 0 && err != nil {
			deepestError = f.getMostResolvedError(err, deepestError)
		}
	}

	if len(container.result) > 0 {
		return nil
	}

	if deepestError == nil {
		return f.newErrMemberNotExist()
	}

	return deepestError
}
