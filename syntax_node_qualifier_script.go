package jsonpath

type syntaxScriptQualifier struct {
	*syntaxBasicNode

	command string
}

func (s syntaxScriptQualifier) retrieve(root, current interface{}, result *resultContainer) error {
	return ErrorNotSupported{`script`, s.text}
}
