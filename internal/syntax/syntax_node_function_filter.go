package syntax

import "github.com/AsaiYusuke/jsonpath/v2/errors"

type syntaxFilterFunction struct {
	*syntaxBasicNode

	function func(any) (any, error)
}

func (f *syntaxFilterFunction) retrieve(
	root, current any, results *[]any) errors.ErrorRuntime {

	filteredValue, err := f.function(current)
	if err != nil {
		return errors.NewErrorFunctionFailed(f.path, f.remainingPathLen, err)
	}

	return f.retrieveAnyValueNext(root, filteredValue, results)
}
