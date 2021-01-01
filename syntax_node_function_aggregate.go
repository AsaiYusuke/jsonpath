package jsonpath

type syntaxAggregateFunction struct {
	*syntaxBasicNode

	function func([]interface{}) (interface{}, error)
	param    syntaxNode
}

func (f *syntaxAggregateFunction) retrieve(
	root, current interface{}, result *[]interface{}) error {

	var values []interface{}
	if err := f.param.retrieve(root, current, &values); err != nil {
		return err
	}

	if !f.param.isValueGroup() {
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

	return f.retrieveNext(
		root, result,
		func() interface{} {
			return filteredValue
		},
		nil)
}
