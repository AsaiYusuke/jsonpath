package jsonpath

type syntaxBasicSubscript struct {
	multiValue bool
}

func (s *syntaxBasicSubscript) isMultiValue() bool {
	return s.multiValue
}
