package syntax

type syntaxBasicSubscript struct {
	valueGroup bool
}

func (s *syntaxBasicSubscript) isValueGroup() bool {
	return s.valueGroup
}
