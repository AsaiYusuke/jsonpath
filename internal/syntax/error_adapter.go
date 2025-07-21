package syntax

import "github.com/AsaiYusuke/jsonpath/errors"

type errorRuntimeAdapter struct {
	*errorBasicRuntime
	err error
}

func newErrorTypeUnmatched(node *syntaxBasicNode, expectedType, foundType string) errorRuntime {
	return errorRuntimeAdapter{
		errorBasicRuntime: &errorBasicRuntime{node: node},
		err:               errors.NewErrorTypeUnmatched(node.text, expectedType, foundType),
	}
}

func newErrorMemberNotExist(node *syntaxBasicNode) errorRuntime {
	return errorRuntimeAdapter{
		errorBasicRuntime: &errorBasicRuntime{node: node},
		err:               errors.NewErrorMemberNotExist(node.text),
	}
}

func newErrorFunctionFailed(node *syntaxBasicNode, errorString string) errorRuntime {
	return errorRuntimeAdapter{
		errorBasicRuntime: &errorBasicRuntime{node: node},
		err:               errors.NewErrorFunctionFailed(node.text, errorString),
	}
}
