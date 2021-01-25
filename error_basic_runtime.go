package jsonpath

type errorBasicRuntime struct {
	node *syntaxBasicNode
}

func (b errorBasicRuntime) getSyntaxNode() *syntaxBasicNode {
	return b.node
}
