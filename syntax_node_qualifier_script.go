package jsonpath

type syntaxScriptQualifier struct {
	*syntaxBasicNode

	command string
}

func (s *syntaxScriptQualifier) retrieve(_ interface{}) error {
	return ErrorNotSupported{`script`, s.text}
}
