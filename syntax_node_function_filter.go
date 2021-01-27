package jsonpath

type syntaxFilterFunction struct {
	*syntaxBasicNode

	function func(interface{}) (interface{}, error)
}

func (f *syntaxFilterFunction) retrieve(
	root, current interface{}, container *bufferContainer) errorRuntime {

	filteredValue, err := f.function(current)
	if err != nil {
		return ErrorFunctionFailed{
			errorBasicRuntime: f.errorRuntime,
			err:               err,
		}
	}

	return f.retrieveAnyValueNext(root, filteredValue, container)
}
