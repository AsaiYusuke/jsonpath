package syntax

type syntaxSliceNegativeStepSubscript struct {
	*syntaxBasicSubscript

	start *syntaxIndexSubscript
	end   *syntaxIndexSubscript
	step  *syntaxIndexSubscript
}

func (s *syntaxSliceNegativeStepSubscript) getLoopStart(srcLength int) int {
	if s.start.isOmitted {
		return s.getNormalizedValue(srcLength-1, srcLength)
	}
	return s.getNormalizedValue(s.start.number, srcLength)
}

func (s *syntaxSliceNegativeStepSubscript) getLoopEnd(srcLength int) int {
	if s.end.isOmitted {
		return s.getNormalizedValue(-srcLength-1, srcLength)
	}
	return s.getNormalizedValue(s.end.number, srcLength)
}

func (s *syntaxSliceNegativeStepSubscript) getNormalizedValue(value int, srcLength int) int {
	if value >= srcLength {
		return srcLength - 1
	}
	if value < -srcLength-1 {
		return -1
	}
	if value < 0 {
		return value + srcLength
	}
	return value
}

func (s *syntaxSliceNegativeStepSubscript) count(srcLength int) int {
	start, end := s.getLoopStart(srcLength), s.getLoopEnd(srcLength)
	if start <= end {
		return 0
	}
	step := -s.step.number
	return (start - end + step - 1) / step
}

func (s *syntaxSliceNegativeStepSubscript) indexAt(srcLength int, ordinal int) int {
	return s.getLoopStart(srcLength) + ordinal*s.step.number
}
