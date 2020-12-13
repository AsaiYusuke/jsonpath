package jsonpath

type syntaxScriptQualifier struct {
	*syntaxBasicNode

	command string
}

func (s *syntaxScriptQualifier) retrieve(current interface{}) error {
	return ErrorNotSupported{`script`, s.text}
}
