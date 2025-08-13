package syntax

import "github.com/AsaiYusuke/jsonpath/errors"

type syntaxFilterFunction struct {
	*syntaxBasicNode

	function func(interface{}) (interface{}, error)
}

func (f *syntaxFilterFunction) retrieve(
	root, current interface{}, container *bufferContainer) errors.ErrorRuntime {

	filteredValue, err := f.function(current)
	if err != nil {
		return errors.NewErrorFunctionFailed(f.path, f.remainingPathLen, err)
	}

	return f.retrieveAnyValueNext(root, filteredValue, container)
}
