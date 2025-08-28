package syntax

type syntaxSlicePositiveStepSubscript struct {
	*syntaxBasicSubscript

	start *syntaxIndexSubscript
	end   *syntaxIndexSubscript
	step  *syntaxIndexSubscript
}

func (s *syntaxSlicePositiveStepSubscript) getLoopStart(srcLength int) int {
	if s.start.isOmitted {
		return s.getNormalizedValue(0, srcLength)
	}
	return s.getNormalizedValue(s.start.number, srcLength)
}

func (s *syntaxSlicePositiveStepSubscript) getLoopEnd(srcLength int) int {
	if s.end.isOmitted {
		return s.getNormalizedValue(srcLength, srcLength)
	}
	return s.getNormalizedValue(s.end.number, srcLength)
}

func (s *syntaxSlicePositiveStepSubscript) getNormalizedValue(value int, srcLength int) int {
	if value >= srcLength {
		return srcLength
	}
	if value < -srcLength {
		return 0
	}
	if value < 0 {
		return value + srcLength
	}
	return value
}

func (s *syntaxSlicePositiveStepSubscript) count(srcLength int) int {
	if s.step.number == 0 {
		return 0
	}
	start, end := s.getLoopStart(srcLength), s.getLoopEnd(srcLength)
	if end <= start {
		return 0
	}
	step := s.step.number
	return (end - start + step - 1) / step
}

func (s *syntaxSlicePositiveStepSubscript) indexAt(srcLength int, ordinal int) int {
	return s.getLoopStart(srcLength) + ordinal*s.step.number
}
