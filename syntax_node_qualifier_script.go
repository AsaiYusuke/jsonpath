package jsonpath

type syntaxScriptQualifier struct {
	*syntaxBasicNode

	command string
}

func (s *syntaxScriptQualifier) retrieve(root, current interface{}) error {
	return ErrorNotSupported{`script`, s.text}
}
