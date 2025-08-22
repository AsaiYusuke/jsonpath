package syntax

import "github.com/AsaiYusuke/jsonpath/v2/errors"

type syntaxAggregateFunction struct {
	*syntaxBasicNode

	function func([]any) (any, error)
	param    syntaxNode
}

func (f *syntaxAggregateFunction) retrieve(
	root, current any, results *[]any) errors.ErrorRuntime {

	buf := getNodeSlice()
	defer func() { putNodeSlice(buf) }()

	if err := f.param.retrieve(root, current, buf); err != nil {
		return err
	}

	result := *buf
	if !f.param.isValueGroup() {
		if arrayParam, ok := (*buf)[0].([]any); ok {
			result = arrayParam
		}
	}

	filteredValue, err := f.function(result)
	if err != nil {
		return errors.NewErrorFunctionFailed(f.path, f.remainingPathLen, err)
	}

	return f.retrieveAnyValueNext(root, filteredValue, results)
}
