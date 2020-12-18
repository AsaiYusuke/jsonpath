package jsonpath

type syntaxAggregateFunction struct {
	*syntaxBasicNode

	function  func([]interface{}) (interface{}, error)
	param     syntaxNode
	resultPtr *[]interface{}
}

func (f *syntaxAggregateFunction) retrieve(current interface{}) error {
	values := make([]interface{}, 0)
	f.resultPtr = &values
	if err := f.param.retrieve(current); err != nil {
		return err
	}
	if !f.param.isMultiValue() {
		if arrayParam, ok := values[0].([]interface{}); ok {
			values = arrayParam
		}
	}
	filteredValue, err := f.function(values)
	if err != nil {
		return ErrorFunctionFailed{
			function: f.text,
			err:      err,
		}
	}
	return f.retrieveNext(filteredValue)
}
