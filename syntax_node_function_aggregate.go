package jsonpath

type syntaxAggregateFunction struct {
	*syntaxBasicNode

	function func([]interface{}) (interface{}, error)
	param    syntaxNode
}

func (f *syntaxAggregateFunction) retrieve(
	root, current interface{}, container *bufferContainer) errorRuntime {

	values := getContainer()
	defer func() {
		putContainer(values)
	}()

	if err := f.param.retrieve(root, current, values); err != nil {
		return err
	}

	result := values.result
	if !f.param.isValueGroup() {
		if arrayParam, ok := values.result[0].([]interface{}); ok {
			result = arrayParam
		}
	}

	filteredValue, err := f.function(result)
	if err != nil {
		return ErrorFunctionFailed{
			errorBasicRuntime: f.errorRuntime,
			err:               err,
		}
	}

	return f.retrieveAnyValueNext(root, filteredValue, container)
}
