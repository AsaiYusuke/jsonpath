package jsonpath

type errorRuntime interface {
	getSyntaxNode() *syntaxBasicNode
}
