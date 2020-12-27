package jsonpath

type syntaxFilterFunction struct {
	*syntaxBasicNode

	function func(interface{}) (interface{}, error)
}

func (f *syntaxFilterFunction) retrieve(current interface{}) error {
	filteredValue, err := f.function(current)
	if err != nil {
		return ErrorFunctionFailed{
			function: f.text,
			err:      err,
		}
	}
	return f.retrieveNext(
		func() interface{} {
			return filteredValue
		},
		nil)
}
