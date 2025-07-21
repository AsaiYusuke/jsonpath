package syntax

import "github.com/AsaiYusuke/jsonpath/errors"

type errorRuntimeAdapter struct {
	err  error
	node *syntaxBasicNode
}

func (e errorRuntimeAdapter) getSyntaxNode() *syntaxBasicNode {
	return e.node
}

func newErrorTypeUnmatched(node *syntaxBasicNode, expectedType, foundType string) errorRuntime {
	return errorRuntimeAdapter{
		err:  errors.NewErrorTypeUnmatched(node.text, expectedType, foundType),
		node: node,
	}
}

func newErrorMemberNotExist(node *syntaxBasicNode) errorRuntime {
	return errorRuntimeAdapter{
		err:  errors.NewErrorMemberNotExist(node.text),
		node: node,
	}
}

func newErrorFunctionFailed(node *syntaxBasicNode, errorString string) errorRuntime {
	return errorRuntimeAdapter{
		err:  errors.NewErrorFunctionFailed(node.text, errorString),
		node: node,
	}
}
