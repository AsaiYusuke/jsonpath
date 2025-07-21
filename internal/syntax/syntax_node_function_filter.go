package syntax

type syntaxFilterFunction struct {
	*syntaxBasicNode

	function func(interface{}) (interface{}, error)
}

func (f *syntaxFilterFunction) retrieve(
	root, current interface{}, container *bufferContainer) errorRuntime {

	filteredValue, err := f.function(current)
	if err != nil {
		return newErrorFunctionFailed(f.errorRuntime.node, err.Error())
	}

	return f.retrieveAnyValueNext(root, filteredValue, container)
}
