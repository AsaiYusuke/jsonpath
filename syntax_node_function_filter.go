package jsonpath

type syntaxFilterFunction struct {
	*syntaxBasicNode

	function func(interface{}) (interface{}, error)
}

func (f *syntaxFilterFunction) retrieve(
	root, current interface{}, container *bufferContainer) error {

	filteredValue, err := f.function(current)
	if err != nil {
		return ErrorFunctionFailed{
			function: f.text,
			err:      err,
		}
	}

	return f.retrieveAnyValueNext(root, filteredValue, container)
}
