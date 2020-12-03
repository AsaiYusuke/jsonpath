package jsonpath

type syntaxScript struct {
	*syntaxBasicNode

	command string
}

func (s syntaxScript) retrieve(root, current interface{}, result *resultContainer) error {
	return ErrorNotSupported{`script`, s.text}
}
