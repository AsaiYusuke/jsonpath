package syntax

// errorBasicRuntime is a basic runtime error structure
type errorBasicRuntime struct {
	node *syntaxBasicNode
}

func (e *errorBasicRuntime) getSyntaxNode() *syntaxBasicNode {
	return e.node
}
